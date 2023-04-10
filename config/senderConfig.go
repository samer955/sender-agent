package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

type SenderConfig struct {
	topics       []string
	discoveryTag string
	frequency    int
}

var config SenderConfig

func init() {

	if err := godotenv.Load("config.env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	topics := os.Getenv("TOPICS")
	if topics == "" {
		log.Fatal("TOPICS ENV VARIABLE NOT FOUND")
	}
	config.topics = strings.Split(topics, ",")

	discovery := os.Getenv("DISCOVERY_TAG")
	if discovery == "" {
		log.Fatal("DISCOVERY_TAG ENV VARIABLE NOT FOUND")
	}
	config.discoveryTag = discovery

	frequencyString := os.Getenv("FREQUENCY")
	if frequencyString == "" {
		log.Println("FREQUENCY ENV NOT FOUND. SET AUTO-FREQUENCY TO 20 SEC")
		config.frequency = 20
	} else {
		config.frequency, _ = strconv.Atoi(frequencyString)
	}

}

func GetConfig() SenderConfig {
	return config
}

func (c *SenderConfig) Topics() []string {
	return c.topics
}

func (c *SenderConfig) DiscoveryTag() string {
	return c.discoveryTag
}

func (c *SenderConfig) Frequency() int {
	return c.frequency
}
