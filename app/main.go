package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"go.uber.org/zap"
	"golang-exercise/api"
	"golang-exercise/auth"
	"golang-exercise/config"
	"golang-exercise/logger"
	"golang-exercise/storage/pg_storage"
	"golang-exercise/writer"
	"google.golang.org/grpc"
	"log"
	"net"
	"strings"
)

func main() {
	config.Init()
	if err := logger.Init(); err != nil {
		log.Fatalf("logger fail to init, error: %v", err)
	}
	var urlBuilder strings.Builder
	urlBuilder.WriteString(fmt.Sprintf("host=%s ", config.DBHost()))
	urlBuilder.WriteString(fmt.Sprintf("user=%s ", config.DBUser()))
	urlBuilder.WriteString(fmt.Sprintf("password=%s ", config.DBPassword()))
	urlBuilder.WriteString(fmt.Sprintf("dbname=%s ", config.DBName()))
	urlBuilder.WriteString(fmt.Sprintf("port=%d ", config.DBPort()))
	urlBuilder.WriteString("sslmode=disable")

	db, err := gorm.Open("postgres", urlBuilder.String())
	if err != nil {
		logger.Fatal("gorm open",
			zap.Error(err),
		)
	}
	err = pgstorage.Init(db)
	if err != nil {
		logger.Fatal("pg_storage init fail",
			zap.Error(err),
		)
	}

	w, err := writer.NewWriter(
		writer.Brokers(config.Brokers()...),
		writer.DestinationTopic(config.DestinationTopic()),
	)
	if err != nil {
		logger.Fatal("new writer fail",
			zap.Error(err),
		)
	}

	lis, err := net.Listen("tcp", config.ServerAddress())
	if err != nil {
		logger.Fatal("failed to listen",
			zap.Error(err),
		)
	}
	byteSecretKey := []byte(config.SecretKey())
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(auth.JWTInterceptor(byteSecretKey)))
	api.Init(grpcServer,
		api.Writer(w),
	)
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatal("failed to run grpc",
			zap.Error(err),
		)
	}
}
