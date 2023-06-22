package workers

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/vipulvpatil/candidate-tracker-go/internal/clients/openai"
	"github.com/vipulvpatil/candidate-tracker-go/internal/storage"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

const PROCESS_FILE_UPLOADS = "process_file_uploads"

var workerStorage storage.StorageAccessor
var openAiClient openai.Client
var logger utilities.Logger

type PoolDependencies struct {
	Namespace    string
	RedisPool    *redis.Pool
	Storage      storage.StorageAccessor
	OpenAiApiKey string
	Logger       utilities.Logger
}

func NewPool(deps PoolDependencies) *work.WorkerPool {
	pool := work.NewWorkerPool(jobContext{}, 10, deps.Namespace, deps.RedisPool)

	pool.Job(PROCESS_FILE_UPLOADS, (*jobContext).processFileUploads)

	// TODO: Not sure if this is the best way to do this. But using Package variables for all dependencies required inside any of the jobs.
	workerStorage = deps.Storage
	logger = deps.Logger
	openAiClient = openai.NewClient(openai.ClientOptions{ApiKey: deps.OpenAiApiKey}, logger)
	return pool
}
