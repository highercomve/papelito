package helperrepo

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/highercomve/papelito/modules/helpers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

var (
	storage StorageMongo
)

type Storage interface {
	GetDatabase() *mongo.Database
	GetCollection(name string) *mongo.Collection
}

// StorageMongo define all storage action and methods
type StorageMongo struct {
	Database         string
	CollectionPrefix string
	client           *mongo.Client
	timeoutDuration  time.Duration
}

func (s *StorageMongo) GetDatabase() *mongo.Database {
	return s.client.Database(s.Database)
}

func (s *StorageMongo) GetCollection(name string) *mongo.Collection {
	name = s.CollectionPrefix + name
	return s.client.Database(s.Database).Collection(name)
}

// IsNotFound resource not found
func IsNotFound(err error) bool {
	return err == mongo.ErrNoDocuments
}

// IsKeyDuplicated test if a key already exist on storage
func IsKeyDuplicated(err error) bool {
	return strings.Contains(err.Error(), "duplicate key error collection")
}

// IsDuplicateKey test if a key already exist on storage
func IsDuplicateKey(key string, err error) bool {
	return strings.Contains(err.Error(), "duplicate key error collection") &&
		strings.Contains(err.Error(), "index: "+key)

}

// New create new Storage Struct
func NewStorage(prefix string) (Storage, error) {
	if storage.client != nil {
		return &storage, nil
	}

	client, err := GetMongoClient()
	if err != nil {
		return nil, err
	}

	timeout, err := time.ParseDuration(helpers.GetEnv("MONGO_TIMEOUT_DURATION", "30m"))
	if err != nil {
		return nil, err
	}

	mongoDb := helpers.GetEnv("MONGO_DB", "")
	storage = StorageMongo{
		client:           client,
		Database:         mongoDb,
		timeoutDuration:  timeout,
		CollectionPrefix: prefix,
	}
	return &storage, nil
}

// GetMongoClient : To Get Mongo Client Object
func GetMongoClient() (*mongo.Client, error) {
	mongoDb := helpers.GetEnv("MONGO_DB", "")
	user := helpers.GetEnv("MONGO_USER", "")
	pass := helpers.GetEnv("MONGO_PASS", "")
	host := helpers.GetEnv("MONGO_HOST", "localhost")
	port := helpers.GetEnv("MONGO_PORT", "27017")
	mongoRs := helpers.GetEnv("MONGO_RS", "")

	//Setting Client Options
	clientOptions := options.Client()
	mongoConnect := "mongodb://"
	if user != "" {
		mongoConnect += user
		if pass != "" {
			mongoConnect += ":"
			mongoConnect += pass
		}
		mongoConnect += "@"
	}
	mongoConnect += host

	if port != "" {
		mongoConnect += ":"
		mongoConnect += port
	}

	mongoConnect += "/" + mongoDb + "?"

	if user != "" {
		mongoConnect += "authSource=" + mongoDb
		mongoConnect += "&authMechanism=SCRAM-SHA-1"
	}

	if mongoRs != "" {
		mongoConnect += "&replicaSet=" + mongoRs
	}

	clientOptions = clientOptions.ApplyURI(mongoConnect)
	if mongoRs != "" {
		clientOptions = clientOptions.SetReplicaSet(mongoRs)
	}
	clientOptions.SetMaxPoolSize(6)
	if os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT") != "" {
		clientOptions.SetMonitor(otelmongo.NewMonitor())
	}

	timeoutEnv := helpers.GetEnv("MONGO_TIMEOUT_DURATION", "30m")
	timeout, err := time.ParseDuration(timeoutEnv)
	if err != nil {
		timeout = 10 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	log.Println("Will connect to mongodb with: " + mongoConnect)
	client, err := mongo.Connect(ctx, clientOptions)
	return client, err
}
