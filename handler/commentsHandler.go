package handler

import (
	"encoding/json"
	"errors"
	"github.com/demos/kafka/topics"
	"github.com/gofiber/fiber/v2"
	"log"
	"os"
	"os/signal"
)

type Comment struct {
	Text string `form:"text" json:"text"`
}

type Comments struct {
	Texts []string `form:"texts" json:"texts"`
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

func (h commentsHandler) CreateComments(c *fiber.Ctx) error {

	// Instantiate new Message struct
	cmts := new(Comments)
	if err := c.BodyParser(cmts); err != nil {
		c.Status(400).JSON(&fiber.Map{
			"success": false,
			"message": err,
		})
		return err
	}

	// convert body into bytes and send it to kafka
	cmtsInBytes := make([][]byte, 0)
	for _, cmt := range cmts.Texts {
		cmtInBytes, err := json.Marshal(Comment{cmt})
		if err != nil {
			log.Printf("unable to seriallize comment: %v", cmt)
		}
		cmtsInBytes = append(cmtsInBytes, cmtInBytes)
	}

	if len(cmtsInBytes) == 0 {
		log.Println("no valid comments to send")
		return errors.New("no valid comments to send")
	}

	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	h.topic.PublishAsync(cmtsInBytes, h.brokers, signals)

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
