package main

import (
	"golang-exercise/api"
	"golang-exercise/storage/pg_storage"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net"
)

func main() {

	dbURL := "host=localhost user=myuser password=mypassword dbname=postgres port=5432 sslmode=disable"

	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{})

	if err != nil {
		log.Fatal("gorm open error:", err)
	}
	err = pg_storage.Init(db)
	if err != nil {
		log.Fatal("pg_storage init error:", err)
	}

	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	api.Init(grpcServer)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to run grpc: %v", err)
	}
}
