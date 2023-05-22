package workers

// This is a wrapper class. It only exists to enable mocking of the gocraft/work library's enqueuer
// TODO: Find a better solution and replace all of this when possible.

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
)

type JobStarter interface {
	EnqueueUnique(jobName string, args map[string]interface{}) (*work.Job, error)
}

type jobStarter struct {
	enqueuer *work.Enqueuer
}

func NewJobStarter(namespace string, redisPool *redis.Pool) JobStarter {
	return &jobStarter{
		enqueuer: work.NewEnqueuer(namespace, redisPool),
	}
}

func (j *jobStarter) EnqueueUnique(jobName string, args map[string]interface{}) (*work.Job, error) {
	return j.enqueuer.EnqueueUnique(jobName, args)
}
