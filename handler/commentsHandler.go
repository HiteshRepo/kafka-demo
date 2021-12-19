package handler

import (
	"encoding/json"
	"github.com/demos/kafka/topics"
	"github.com/gofiber/fiber/v2"
	"log"
)

type Comment struct {
	Text string `form:"text" json:"text"`
}

type commentsHandler struct {
	topic topics.Topic
	brokers []string
}

func NewCommentsHandler(topic topics.Topic, brokers []string) commentsHandler {
	return commentsHandler{topic: topic, brokers: brokers}
}

func (h commentsHandler) CreateComment(c *fiber.Ctx) error {

	// Instantiate new Message struct
	cmt := new(Comment)
	if err := c.BodyParser(cmt); err != nil {
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	}

	// convert body into bytes and send it to kafka
	cmtInBytes, err := json.Marshal(cmt)
	err = h.topic.Publish(cmtInBytes, h.brokers)
	if err != nil {
		if err = sendResponse(c ,false, "error while publishing comment", cmt); err != nil {
			return err
		}
	}

	// Return Comment in JSON format
	if err = sendResponse(c ,true, "comment pushed successfully", cmt); err != nil {
		return err
	}

	return nil
}

func sendResponse(c *fiber.Ctx, success bool, message string, data interface{}) error {
	err := c.JSON(&fiber.Map{
		"success": success,
		"message": message,
		"comment": data,
	})

	if err != nil {
		c.Status(500).JSON(&fiber.Map{
			"success": false,
			"message": "error creating product",
		})
		log.Printf("error while sending response: %v", err)
	}

	return err
}
