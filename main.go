package main

import (
	"flag"
	"github.com/demos/kafka/config"
	"github.com/demos/kafka/router"
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

	// serverConfig := appConfig.GetServerConfig()
	// topicConfig := appConfig.GetTopicsConfig()
	// brokers := strings.Split(serverConfig.Brokers, ";")
	// go consumer.StartConsumer(brokers, topicConfig.Name, appConfig.GetConsumerConfig())
	// go consumerGroup.StartConsumerGroup(brokers, topicConfig.Name, appConfig.GetConsumerGroupConfig())

	app, err := router.InitRouter(appConfig)
	app.Listen(":3000")
}



