package queues

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/jpillora/backoff"

	"github.com/9count/go-services/core/clock"
	"github.com/9count/go-services/core/log"
	"github.com/9count/go-services/core/workflows/types"
)

const (
	defaultReceiveWaitTimeout      = 20 * time.Second
	defaultShutdownTimeout         = 30 * time.Second
	defaultTaskTimeout             = 60 * time.Second
	defaultWorkerCount        uint = 10
)

type OutputQueueFunc func(taskOutput string) []string

type WorkerOptions struct {
	IgnoreTaskErrors bool
	LogTaskOutput    bool
	ShutdownTimeout  time.Duration
	TaskTimeout      time.Duration
	WorkerCount      uint
}

type WorkerOption func(*WorkerOptions)

func IgnoreTaskErrors(ignore bool) func(*WorkerOptions) {
	return func(args *WorkerOptions) {
		args.IgnoreTaskErrors = ignore
	}
}

func LogTaskOutput(b bool) func(*WorkerOptions) {
	return func(args *WorkerOptions) {
		args.LogTaskOutput = b
	}
}

func TaskTimeout(d time.Duration) func(*WorkerOptions) {
	return func(args *WorkerOptions) {
		args.TaskTimeout = d
	}
}

func ShutdownTimeout(d time.Duration) func(*WorkerOptions) {
	return func(args *WorkerOptions) {
		args.ShutdownTimeout = d
	}
}

func WorkerCount(count uint) func(*WorkerOptions) {
	return func(args *WorkerOptions) {
		args.WorkerCount = count
	}
}

type SQSWorker struct {
	pending   chan bool
	quit      chan bool
	sqsClient sqsiface.SQSAPI

	options  *WorkerOptions
	queueURL string
	taskFunc types.TaskFunc

	outputQueueFunc OutputQueueFunc
	outputQueueUrls []string

	// for unit testing
	timeClock clock.Clock
}

type SQSWorkerConfig struct {
	AwsRegion   string
	AwsEndpoint string

	SqsClient sqsiface.SQSAPI

	QueueURL string
	TaskFunc types.TaskFunc

	// OutputQueueFunc takes precedence
	OutputQueueFunc OutputQueueFunc
	OutputQueueUrls []string
}

func NewSQSWorker(cfg *SQSWorkerConfig, opts ...WorkerOption) (*SQSWorker, error) {
	if cfg == nil {
		return nil, errors.New("invalid: cfg")
	}

	args := &WorkerOptions{
		IgnoreTaskErrors: false,
		LogTaskOutput:    false,
		ShutdownTimeout:  defaultShutdownTimeout,
		TaskTimeout:      defaultTaskTimeout,
		WorkerCount:      defaultWorkerCount,
	}

	for _, opt := range opts {
		opt(args)
	}

	awsConf := aws.NewConfig()

	if cfg.AwsEndpoint != "" {
		awsConf = awsConf.WithEndpoint(cfg.AwsEndpoint)
	}

	if cfg.AwsRegion != "" {
		awsConf = awsConf.WithRegion(cfg.AwsRegion)
	}

	sqsClient := cfg.SqsClient

	if sqsClient == nil {
		sqsClient = sqs.New(session.Must(session.NewSession()), awsConf)
	}

	if cfg.QueueURL == "" {
		return nil, errors.New("invalid: QueueURL")
	}

	if cfg.TaskFunc == nil {
		return nil, errors.New("invalid: TaskFunc")
	}

	return &SQSWorker{
		options: args,
		pending: make(chan bool, 1),
		quit:    make(chan bool, 1),

		outputQueueFunc: cfg.OutputQueueFunc,
		outputQueueUrls: cfg.OutputQueueUrls,
		queueURL:        cfg.QueueURL,
		sqsClient:       sqsClient,
		taskFunc:        cfg.TaskFunc,

		timeClock: clock.New(),
	}, nil
}

type taskInput struct {
	Input         string
	ReceiptHandle string
}

func (s *SQSWorker) stubbornDeleteMessage(taskStart time.Time, receiptHandle string) error {
	delay := &backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    10 * time.Second,
		Factor: 2,
		Jitter: true,
	}

	// This operation is critical in that there is a potential for a successful processing of a message
	// to get re-queued if the DeleteMessage fails, such as in a networking error incident
	for s.timeClock.Now().Sub(taskStart) < s.options.TaskTimeout {
		_, err := s.sqsClient.DeleteMessage(&sqs.DeleteMessageInput{
			QueueUrl:      aws.String(s.queueURL),
			ReceiptHandle: aws.String(receiptHandle),
		})

		if err == nil {
			return nil
		}

		s.timeClock.Sleep(delay.Duration())
	}

	return fmt.Errorf("error deleting from queue: %s", s.queueURL)
}

func (s *SQSWorker) doTask(ctx context.Context, in *taskInput) {
	if in == nil {
		log.Error(ctx, "invalid nil task in doTask")
		return
	}

	taskStart := s.timeClock.Now()

	taskInput := "{}"

	if in.Input != "" {
		taskInput = in.Input
	}

	// TODO: implement cancelFunc
	taskCtx, taskCancel := context.WithTimeout(ctx, s.options.TaskTimeout)

	defer taskCancel()

	taskOutput, taskErr := s.taskFunc(taskCtx, taskInput)

	if taskErr != nil {
		log.Error(ctx, "error in workerFunc", log.WithError(taskErr))

		// an early return here will implicitly allow the Message to be re-queued
		if !s.options.IgnoreTaskErrors {
			return
		}
	}

	if s.options.LogTaskOutput {
		log.Info(ctx, "task success", log.WithValue("output", taskOutput))
	}

	deleteErr := s.stubbornDeleteMessage(taskStart, in.ReceiptHandle)

	// this is very bad, and we should probably even notify bugsnag because either there is a permissions issue
	// on the queue or aws itself is having problems
	if deleteErr != nil {
		log.Error(ctx, "DeleteMessage error: exceeded timeout", log.WithValue("message", in.Input))

		return
	}

	// an empty string implies do not pass to output queues
	if taskOutput == "" {
		return
	}

	var outputQueues []string

	if s.outputQueueFunc != nil {
		outputQueues = s.outputQueueFunc(taskOutput)
	} else {
		outputQueues = s.outputQueueUrls
	}

	// only send to output queues if everything has gone well to this point as we would not want to double-send
	// if for example, the DeleteMessage failed
	for _, outputUrl := range outputQueues {
		publisher, err := NewSQSPublisher(&SQSPublisherConfig{
			SqsClient: s.sqsClient,
			QueueURL:  outputUrl,
		})

		if err != nil {
			log.Error(ctx, "NewSQSPublisher error",
				log.WithError(err),
				log.WithValue("outputUrl", outputUrl),
			)

			continue
		}

		sendErr := publisher.SendMessage(context.Background(), taskOutput)

		if sendErr != nil {
			log.Error(ctx, "error sending to output queue",
				log.WithValue("outputQueueUrl", outputUrl),
				log.WithValue("messageBody", taskOutput))
		}
	}
}

func (s *SQSWorker) worker(ctx context.Context) {
	delay := &backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    20 * time.Second,
		Factor: 2,
		Jitter: true,
	}

	active := true

	for active {
		select {
		case <-ctx.Done():
			active = false

		default:
			res, err := s.sqsClient.ReceiveMessage(&sqs.ReceiveMessageInput{
				MaxNumberOfMessages: aws.Int64(1),
				QueueUrl:            aws.String(s.queueURL),
				WaitTimeSeconds:     aws.Int64(int64(defaultReceiveWaitTimeout.Seconds())),
				VisibilityTimeout:   aws.Int64(int64(s.options.TaskTimeout.Seconds())),
			})

			// this is unlikely to happen, but would likely be a networking issue, over-limit, or permission error
			// TODO: implement exponential back-off
			if err != nil {
				log.Error(ctx, "error with ReceiveMessage", log.WithError(err))

				s.timeClock.Sleep(delay.Duration())
			} else {
				delay.Reset()

				// at the moment we set a limit of 1 message per ReceiveMessage, but let's code for the future
				for _, message := range res.Messages {
					s.doTask(ctx, &taskInput{Input: aws.StringValue(message.Body), ReceiptHandle: aws.StringValue(message.ReceiptHandle)})
				}
			}
		}
	}
}

func (s *SQSWorker) startWorkers(ctx context.Context) {
	for i := 0; i < int(s.options.WorkerCount); i++ {
		go s.worker(ctx)
	}
}

func (s *SQSWorker) Run(ctx context.Context) {
	if s.isRunning() {
		log.Warn(ctx, "worker is already running")
		return
	}

	runCtx, runCancel := context.WithCancel(ctx)

	s.startWorkers(runCtx)

	s.pending <- true

	sigKill := make(chan os.Signal, 2)
	signal.Notify(sigKill, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Info(ctx, "received ctx cancel")
		runCancel()

	case <-s.quit:
		log.Info(ctx, "received quit")
		runCancel()

	case <-sigKill:
		log.Info(ctx, "received kill/term signal")
		runCancel()
	}

	log.Info(ctx, fmt.Sprintf("Shutting down Worker: %v", s.options.ShutdownTimeout))

	s.timeClock.Sleep(s.options.ShutdownTimeout)

	<-s.pending
}

func (s *SQSWorker) isRunning() bool {
	return len(s.pending) > 0
}

func (s *SQSWorker) Stop() {
	if s.isRunning() {
		s.quit <- true
	}
}
