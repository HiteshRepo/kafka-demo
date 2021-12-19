package router

import (
	"errors"
	"fmt"
	"github.com/demos/kafka/config"
	"github.com/demos/kafka/handler"
	"github.com/demos/kafka/topics"
	"github.com/gofiber/fiber/v2"
	"log"
	"strings"
)

func InitRouter(appConfig config.AppConfig) (*fiber.App, error) {
	app := fiber.New()
	api := app.Group("/api/v1")
	brokers := strings.Split(appConfig.GetServerConfig().Brokers, ";")
	commentsHandler, err := getCommentsHandler(appConfig.GetTopicsConfig(), appConfig.GetProducerConfig(), brokers)
	if err != nil {
		log.Println("invalid topic")
		return nil, err
	}

	api.Post("/comment", commentsHandler)

	return app, nil
}

func getCommentsHandler(topicConfig *config.TopicsConfig, producerConfig *config.ProducerConfig, brokers []string) (fiber.Handler, error) {
	topic := topics.GetNewTopic(topicConfig, producerConfig)

	if !topic.IsTopicAvailable(brokers) {
		return nil, errors.New(fmt.Sprintf("invalid topic: %v", topicConfig.Name))
	}

	return handler.NewCommentsHandler(topic, brokers).CreateComment, nil
}
