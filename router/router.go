package router

import (
	"errors"
	"fmt"
	"github.com/demos/kafka/config"
	"github.com/demos/kafka/handler"
	"github.com/demos/kafka/topics"
	"github.com/gofiber/fiber/v2"
	"log"
)

func InitRouter(topicConfig *config.TopicsConfig, producerConfig *config.ProducerConfig) (*fiber.App, error) {
	app := fiber.New()
	api := app.Group("/api/v1")
	commentsHandler, err := getCommentsHandler(topicConfig, producerConfig)
	if err != nil {
		log.Println("invalid topic")
		return nil, err
	}

	api.Post("/comment", commentsHandler)

	return app, nil
}

func getCommentsHandler(topicConfig *config.TopicsConfig, producerConfig *config.ProducerConfig) (fiber.Handler, error) {
	topic := topics.GetNewTopic(topicConfig, producerConfig)

	if !topic.IsTopicAvailable() {
		return nil, errors.New(fmt.Sprintf("invalid topic: %v", topicConfig.Name))
	}

	return handler.NewCommentsHandler(topic).CreateComment, nil
}
