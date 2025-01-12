package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	brokers          []string
	destinationTopic string

	serverAddress string

	dbName     string
	dbHost     string
	dbPort     int
	dbUser     string
	dbPassword string

	secretKey string
}

var c *Config

func Init() {
	c = new(Config)

	c.serverAddress = os.Getenv("SERVER_ADDRESS")
	c.secretKey = os.Getenv("SECRET_KEY")
	c.dbPort = 5432
	c.dbHost = os.Getenv("DB_HOST")

	dbPortEnv := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(dbPortEnv)
	if err != nil || dbPortEnv == "" {
		c.dbPort = dbPort
	}

	c.dbUser = os.Getenv("DB_USER")
	c.dbPassword = os.Getenv("DB_PASSWORD")
	c.dbName = os.Getenv("DB_NAME")
	c.destinationTopic = os.Getenv("DESTINATION_TOPIC")
	c.brokers = strings.Split(os.Getenv("KAFKA_BROKERS"), ",")
}

func ServerAddress() string {
	return c.serverAddress
}
func DBName() string {
	return c.dbName
}

func DBHost() string {
	return c.dbHost
}

func DBPort() int {
	return c.dbPort
}

func DBUser() string {
	return c.dbUser
}

func DBPassword() string {
	return c.dbPassword
}

func Brokers() []string {
	return c.brokers
}

func DestinationTopic() string {
	return c.destinationTopic
}

func SecretKey() string {
	return c.secretKey
}
