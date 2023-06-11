package workers

import (
	"github.com/gocraft/work"
	"github.com/gomodule/redigo/redis"
	"github.com/vipulvpatil/candidate-tracker-go/internal/clients/openai"
	"github.com/vipulvpatil/candidate-tracker-go/internal/storage"
	"github.com/vipulvpatil/candidate-tracker-go/internal/utilities"
)

const START_GAME_ONCE_PLAYERS_HAVE_JOINED = "start_game_once_players_have_joined"
const ASK_QUESTION_ON_BEHALF_OF_BOT = "ask_question_on_behalf_of_bot"
const ANSWER_QUESTION_ON_BEHALF_OF_BOT = "answer_question_on_behalf_of_bot"
const DELETE_EXPIRED_GAMES = "delete_expired_games"

var workerStorage storage.StorageAccessor
var openAiClient openai.Client
var minDelayAfterAIResponse int
var maxDelayAfterAIResponse int
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

	// TODO: Not sure if this is the best way to do this. But using Package variables for all dependencies required inside any of the jobs.
	workerStorage = deps.Storage
	logger = deps.Logger
	openAiClient = openai.NewClient(openai.ClientOptions{ApiKey: deps.OpenAiApiKey}, logger)
	minDelayAfterAIResponse = 8
	maxDelayAfterAIResponse = 15
	return pool
}
