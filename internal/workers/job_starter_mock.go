package workers

import (
	"github.com/pkg/errors"

	"github.com/gocraft/work"
)

type JobStarterMockCallCheck struct {
	CalledArgs map[string][]map[string]interface{}
}

func (j *JobStarterMockCallCheck) EnqueueUnique(jobName string, args map[string]interface{}) (*work.Job, error) {
	if j.CalledArgs == nil {
		j.CalledArgs = map[string][]map[string]interface{}{}
	}
	if j.CalledArgs[jobName] == nil {
		j.CalledArgs[jobName] = []map[string]interface{}{}
	}
	j.CalledArgs[jobName] = append(j.CalledArgs[jobName], args)
	return &work.Job{}, nil
}

type JobStarterMockFailure struct{}

func (j *JobStarterMockFailure) Enqueue(jobName string, args map[string]interface{}) (*work.Job, error) {
	return nil, errors.New("unable to enqueue job")
}
