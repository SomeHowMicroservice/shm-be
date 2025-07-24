package initialization

import (
	"fmt"

	"github.com/SomeHowMicroservice/shm-be/services/product/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func InitDB(cfg *config.Config) (*mongo.Client, error) {
	mongoURI := fmt.Sprintf("mongodb+srv://%s:%s@%s/%s?retryWrites=%t&w=%s&appName=%s", 
		cfg.Database.DBUser,
		cfg.Database.DBPassword,
		cfg.Database.DBHost,
		cfg.Database.DBName,
		cfg.Database.DBRetryWrites,
		cfg.Database.DBW,
		cfg.Database.DBAppName,
	)
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	clientOptions := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil, fmt.Errorf("kết nối MongoDB thất bại: %w", err)
	}

	return client, nil
}
