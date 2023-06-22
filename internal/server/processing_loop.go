package server

import (
	"context"
	"sync"
	"time"

	"github.com/gocraft/work"
	"github.com/vipulvpatil/candidate-tracker-go/internal/workers"
)

func (s *CandidateTrackerGoService) ProcessingLoop(ctx context.Context, tickerDuration time.Duration, wg *sync.WaitGroup, jobStarter workers.JobStarter) {
	wg.Add(1)
	defer wg.Done()

	ticker := time.NewTicker(tickerDuration)
	for {
		select {
		case <-ticker.C:
			s.processFileUpload(jobStarter)
		case <-ctx.Done():
			return
		}
	}
}

func (s *CandidateTrackerGoService) processFileUpload(jobStarter workers.JobStarter) {
	fileUploadIds, err := s.storage.GetAllProcessingNotStartedFileUploadIds()
	if err != nil {
		s.logger.LogError(err)
		return
	}
	for _, fileUploadId := range fileUploadIds {
		_, err := jobStarter.EnqueueUnique(workers.PROCESS_FILE_UPLOAD, work.Q{"fileUploadId": fileUploadId})
		if err != nil {
			s.logger.LogError(err)
		}
	}
}
