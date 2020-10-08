package queues

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
	"github.com/jpillora/backoff"

	"github.com/9count/go-services/core/clock"
)

const (
	defaultSendTimeout = 60 * time.Second
)

type SQSPublisher struct {
	sqsClient sqsiface.SQSAPI
	queueURL  string

	// for unit testing
	timeClock clock.Clock
}

type SQSPublisherConfig struct {
	AwsRegion   string
	AwsEndpoint string

	SqsClient sqsiface.SQSAPI
	QueueURL  string
}

func NewSQSPublisher(cfg *SQSPublisherConfig) (*SQSPublisher, error) {
	if cfg == nil {
		return nil, errors.New("invalid: cfg")
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

	return &SQSPublisher{
		queueURL:  cfg.QueueURL,
		sqsClient: sqsClient,

		timeClock: clock.New(),
	}, nil
}

func (s *SQSPublisher) stubbornSendMessage(queueUrl, messageBody string) error {
	delay := &backoff.Backoff{
		Min:    100 * time.Millisecond,
		Max:    10 * time.Second,
		Factor: 2,
		Jitter: true,
	}

	sendStart := time.Now()

	for s.timeClock.Now().Sub(sendStart) < defaultSendTimeout {
		_, err := s.sqsClient.SendMessage(&sqs.SendMessageInput{
			MessageBody: aws.String(messageBody),
			QueueUrl:    aws.String(queueUrl),
		})

		if err == nil {
			return nil
		}

		s.timeClock.Sleep(delay.Duration())
	}

	return fmt.Errorf("error sending to queue: %s", queueUrl)
}

// SendMessage takes a JSON serializable type or a raw String and performs a SendMessage
// to the QueueUrl attached to this SQSWorker
func (s *SQSPublisher) SendMessage(ctx context.Context, msg interface{}) error {
	if msg == nil {
		return errors.New("invalid: msg")
	}

	// check to see if the input is a string, and if so, just pass it straight through
	if msgStr, ok := msg.(string); ok {
		return s.stubbornSendMessage(s.queueURL, msgStr)
	}

	msgBuf, err := json.Marshal(msg)

	if err != nil {
		return err
	}

	return s.stubbornSendMessage(s.queueURL, string(msgBuf))
}
