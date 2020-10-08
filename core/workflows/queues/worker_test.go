package queues

import (
	"context"
	"syscall"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/9count/go-services/core/clock"
	commonSQS "github.com/9count/go-services/core/sqs"
	"github.com/9count/go-services/core/workflows/types"
)

type mockTask struct {
	mock.Mock
}

func (m *mockTask) NoopFunc(ctx context.Context, input string) (string, error) {
	args := m.Called(ctx, input)

	if args.Error(1) != nil {
		return "", args.Error(1)
	}

	return args.Get(0).(string), nil
}

func setupSQSWorkerFixture(f types.TaskFunc) (*SQSWorker, error) {
	worker, err := NewSQSWorker(&SQSWorkerConfig{
		SqsClient: &commonSQS.MockSQSClient{},
		QueueURL:  "https://sqs.us-east-1.amazonaws.com/1234567890/foo-queue-ABCD",
		TaskFunc:  f,
	})

	if err != nil {
		return nil, err
	}

	worker.timeClock = clock.NewMock()

	return worker, nil
}

func TestNewTaskWorker(t *testing.T) {
	queueURL := "https://sqs.us-east-1.amazonaws.com/1234567890/foo-queue-ABCD"

	t.Run("returns an error with nil config", func(t *testing.T) {
		_, err := NewSQSWorker(nil)

		assert.Error(t, err)
	})

	t.Run("returns an error with no QueueUrl", func(t *testing.T) {
		_, err := NewSQSWorker(&SQSWorkerConfig{})

		assert.Error(t, err)
	})

	t.Run("returns an error with no TaskFunc", func(t *testing.T) {
		_, err := NewSQSWorker(&SQSWorkerConfig{QueueURL: queueURL})

		assert.Error(t, err)
	})

	t.Run("returns a valid object with default options", func(t *testing.T) {
		w, err := NewSQSWorker(&SQSWorkerConfig{
			AwsEndpoint: "http://localhost:8083",
			AwsRegion:   "us-west-1",
			QueueURL:    queueURL,
			TaskFunc: func(ctx context.Context, input string) (s string, e error) {
				return "", nil
			},
		})

		assert.Nil(t, err)
		assert.NotNil(t, w)

		assert.Equal(t, false, w.options.IgnoreTaskErrors)
		assert.Equal(t, false, w.options.LogTaskOutput)
		assert.Equal(t, defaultShutdownTimeout, w.options.ShutdownTimeout)
		assert.Equal(t, defaultTaskTimeout, w.options.TaskTimeout)
		assert.Equal(t, defaultWorkerCount, w.options.WorkerCount)
	})

	t.Run("sets custom options", func(t *testing.T) {
		expShutdownTimeout := 88 * time.Second
		expTaskTimeout := 888 * time.Second
		var expWorkerCount uint = 11

		w, err := NewSQSWorker(
			&SQSWorkerConfig{
				QueueURL: queueURL,
				TaskFunc: func(ctx context.Context, input string) (s string, e error) {
					return "", nil
				},
			},
			IgnoreTaskErrors(true),
			LogTaskOutput(true),
			ShutdownTimeout(expShutdownTimeout),
			TaskTimeout(expTaskTimeout),
			WorkerCount(expWorkerCount),
		)

		assert.Nil(t, err)
		assert.NotNil(t, w)

		assert.Equal(t, true, w.options.IgnoreTaskErrors)
		assert.Equal(t, true, w.options.LogTaskOutput)
		assert.Equal(t, expShutdownTimeout, w.options.ShutdownTimeout)
		assert.Equal(t, expTaskTimeout, w.options.TaskTimeout)
		assert.Equal(t, expWorkerCount, w.options.WorkerCount)
	})
}

func TestTaskWorker_Run(t *testing.T) {
	setup := func(t *testing.T) (*SQSWorker, error) {
		w, err := setupSQSWorkerFixture((&mockTask{}).NoopFunc)

		if err != nil {
			return nil, err
		}

		mockSqs := w.sqsClient.(*commonSQS.MockSQSClient)

		mockSqs.On(
			"ReceiveMessage",
			mock.Anything,
		).Maybe().Return(&sqs.ReceiveMessageOutput{}, nil)

		mockClock := w.timeClock.(*clock.MockClock)

		mockClock.On(
			"Sleep",
			defaultShutdownTimeout,
		).Once().Return()

		return w, nil
	}

	t.Run("runs and stops correctly", func(t *testing.T) {
		w, err := setup(t)

		if err != nil {
			t.Error(err)
			return
		}

		// stop the worker after it starts running
		go func() {
			for !w.isRunning() {
				time.Sleep(10 * time.Millisecond)
			}
			w.Stop()
		}()

		w.Run(context.Background())

		w.sqsClient.(*commonSQS.MockSQSClient).AssertExpectations(t)
		w.timeClock.(*clock.MockClock).AssertExpectations(t)
	})

	t.Run("runs and stops correctly (ctx cancel)", func(t *testing.T) {
		w, err := setup(t)

		if err != nil {
			t.Error(err)
			return
		}

		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			for !w.isRunning() {
				time.Sleep(10 * time.Millisecond)
			}
			cancel()
		}()

		w.Run(ctx)

		w.sqsClient.(*commonSQS.MockSQSClient).AssertExpectations(t)
		w.timeClock.(*clock.MockClock).AssertExpectations(t)
	})

	t.Run("runs and stops correctly (SIGKILL)", func(t *testing.T) {
		w, err := setup(t)

		if err != nil {
			t.Error(err)
			return
		}

		go func() {
			for !w.isRunning() {
				time.Sleep(10 * time.Millisecond)
			}
			_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
		}()

		w.Run(context.Background())

		w.sqsClient.(*commonSQS.MockSQSClient).AssertExpectations(t)
		w.timeClock.(*clock.MockClock).AssertExpectations(t)
	})
}

// easier to test this private method in isolation
func TestSQSWorker_DoTask(t *testing.T) {
	t.Run("returns early with nil task input", func(t *testing.T) {
		task := &mockTask{}

		w, err := setupSQSWorkerFixture(task.NoopFunc)

		if err != nil {
			t.Error(err)
			return
		}

		mockSqs := w.sqsClient.(*commonSQS.MockSQSClient)

		w.doTask(context.Background(), nil)

		mockSqs.AssertExpectations(t)
		task.AssertExpectations(t)
	})

	t.Run("ignores output queues with empty TaskFunc output", func(t *testing.T) {
		task := &mockTask{}

		w, err := setupSQSWorkerFixture(task.NoopFunc)

		if err != nil {
			t.Error(err)
			return
		}

		w.outputQueueUrls = []string{"http://foo", "http://bar"}

		input := &taskInput{
			Input:         `{"hello":"world"}`,
			ReceiptHandle: "foo-1234",
		}

		timeClock := w.timeClock.(*clock.MockClock)

		timeClock.On("Now").Return(time.Now())

		sqsMock := w.sqsClient.(*commonSQS.MockSQSClient)

		sqsMock.On(
			"DeleteMessage",
			&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(w.queueURL),
				ReceiptHandle: aws.String(input.ReceiptHandle),
			},
		).Once().Return(&sqs.DeleteMessageOutput{}, nil)

		taskOutput := ""

		task.On(
			"NoopFunc",
			mock.Anything,
			input.Input,
		).Once().Return(taskOutput, nil)

		w.doTask(context.Background(), input)

		sqsMock.AssertExpectations(t)
		timeClock.AssertExpectations(t)
		task.AssertExpectations(t)
	})

	t.Run("send task output to output queues", func(t *testing.T) {
		task := &mockTask{}

		w, err := setupSQSWorkerFixture(task.NoopFunc)

		if err != nil {
			t.Error(err)
			return
		}

		w.outputQueueUrls = []string{"http://foo", "http://bar"}

		input := &taskInput{
			Input:         `{"hello":"world"}`,
			ReceiptHandle: "foo-1234",
		}

		mockClock := w.timeClock.(*clock.MockClock)

		mockClock.On("Now").Return(time.Now())

		sqsMock := w.sqsClient.(*commonSQS.MockSQSClient)

		sqsMock.On(
			"DeleteMessage",
			&sqs.DeleteMessageInput{
				QueueUrl:      aws.String(w.queueURL),
				ReceiptHandle: aws.String(input.ReceiptHandle),
			},
		).Once().Return(&sqs.DeleteMessageOutput{}, nil)

		taskOutput := "foo-output"

		task.On(
			"NoopFunc",
			mock.Anything,
			input.Input,
		).Once().Return(taskOutput, nil)

		for _, outputUrl := range w.outputQueueUrls {
			sqsMock.On(
				"SendMessage",
				&sqs.SendMessageInput{
					MessageBody: aws.String(taskOutput),
					QueueUrl:    aws.String(outputUrl),
				},
			).Once().Return(&sqs.SendMessageOutput{}, nil)
		}

		w.doTask(context.Background(), input)

		task.AssertExpectations(t)
		sqsMock.AssertExpectations(t)
	})
}
