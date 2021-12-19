package main

import (
	"flag"
	"github.com/demos/kafka/config"
	"github.com/demos/kafka/topics"
	"log"
	"os"
)

const (
	configFileKey     = "configFile"
	defaultConfigFile = ""
	configFileUsage   = "this is config file path"
)

func main() {
	var configFile string
	flag.StringVar(&configFile, configFileKey, defaultConfigFile, configFileUsage)
	flag.Parse()

	configReader, err := os.Open(configFile)
	if err != nil {
		log.Printf("error while opening config file: %v", err)
		return
	}

	appConfig, err := config.LoadConfig(configReader)
	if err != nil {
		log.Printf("error while loading config file: %v", err)
		return
	}

	producerConfig := appConfig.GetProducerConfig()
	topicConfig := appConfig.GetTopicsConfig()

	topic := topics.GetNewTopic(topicConfig, producerConfig)

	if !topic.IsTopicAvailable() {
		log.Println("invalid topic")
		return
	}

	err = topic.Publish([]byte("hello"))
	if err != nil {
		log.Printf("error while publishing topic to kafka: %v", err)
		return
	}
}
