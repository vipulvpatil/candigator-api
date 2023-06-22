package server

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vipulvpatil/candidate-tracker-go/internal/storage"
	"github.com/vipulvpatil/candidate-tracker-go/internal/workers"
)

func Test_Loop(t *testing.T) {
	t.Run("looks for file uploads that are not already being processed and calls job starter to process them, until canceled", func(t *testing.T) {
		jobStarterMock := &workers.JobStarterMockCallCheck{}
		tickerDuration := 10 * time.Millisecond

		fileUploadAccessorMock1 := fileUploadAccessorCallerInspectableMock{
			&functionCallInspectableMock{
				ReturnData:  [][]string{{"fp_id1", "fp_id2"}, {"fp_id3"}},
				ReturnCount: 2,
			},
		}

		functionsToCheck := []struct {
			name              string
			functionCall      functionCallInspectable
			expectedCallCount int
		}{
			{
				name:              "getAllUnprocessed, %s",
				functionCall:      fileUploadAccessorMock1,
				expectedCallCount: 4,
			},
		}

		jobStartedCallsToVerify := []struct {
			jobName string
			jobArgs []map[string]any
		}{
			{
				jobName: workers.PROCESS_FILE_UPLOAD,
				jobArgs: []map[string]any{
					{"fileUploadId": "fp_id1"},
					{"fileUploadId": "fp_id2"},
					{"fileUploadId": "fp_id3"},
				},
			},
		}

		server, _ := NewServer(
			ServerDependencies{
				Storage: storage.NewStorageAccessorMock(
					storage.WithFileUploadAccessorMock(
						&storage.FileUploadAccessorConfigurableMock{
							GetAllProcessingNotStartedFileUploadIdsInternal: fileUploadAccessorMock1.getAllUnprocessed,
						},
					),
				),
			},
		)

		var wg sync.WaitGroup
		loopCtx, cancelLoop := context.WithCancel(context.Background())
		go server.Loop(loopCtx, tickerDuration, &wg, jobStarterMock)
		time.Sleep(45 * time.Millisecond)

		for _, jobsStarted := range jobStartedCallsToVerify {
			assertJobStarterCalledWithArgsForJob(
				t,
				jobsStarted.jobArgs,
				jobStarterMock,
				jobsStarted.jobName,
			)
		}

		for _, f := range functionsToCheck {
			assertCallCount(t, f.expectedCallCount, f.functionCall, f.name, "loop should run continuously until canceled")
		}
		cancelLoop()
		time.Sleep(45 * time.Millisecond)
		for _, f := range functionsToCheck {
			assertCallCount(t, f.expectedCallCount, f.functionCall, f.name, "function call count should not change once loop is canceled")
		}
	})
}

type functionCallInspectable interface {
	FunctionCalledCount() int
}

type functionCallInspectableMock struct {
	ReturnData  [][]string
	ReturnCount int
	callCount   int
}

func (f *functionCallInspectableMock) FunctionCalledCount() int {
	return f.callCount
}

func assertCallCount(t *testing.T, expectedCallCount int, functionCall functionCallInspectable, msgAndArgs ...any) bool {
	return assert.Equal(t, expectedCallCount, functionCall.FunctionCalledCount(), msgAndArgs...)
}

type fileUploadAccessorCallerInspectableMock struct {
	*functionCallInspectableMock
}

func (m *fileUploadAccessorCallerInspectableMock) getAllUnprocessed() ([]string, error) {
	m.callCount++
	if m.ReturnCount >= m.callCount {
		return m.ReturnData[m.callCount-1], nil
	}
	return nil, nil
}

func assertJobStarterCalledWithArgsForJob(t *testing.T, expectedCalledArgs []map[string]any, jobStarter *workers.JobStarterMockCallCheck, jobName string) bool {
	return assert.EqualValues(
		t,
		expectedCalledArgs,
		jobStarter.CalledArgs[jobName],
		"appropriate jobs should be started from the loop",
	)
}
