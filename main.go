package main

import (
	"net"
	"time"

	"github.com/ajdinahmetovic/item-service/db"
	"github.com/ajdinahmetovic/item-service/logger"
	"github.com/ajdinahmetovic/item-service/proto/v1"
	"google.golang.org/grpc"
)

//Server struct
type Server struct{}

func main() {
	logger.InitLogger()

	err := db.ConnectDB()
	if err != nil {
		logger.Error("Failed to connect to daabase", "time", time.Now(), "err", err)
		return
	}

	listener, err := net.Listen("tcp", ":4040")
	if err != nil {
		logger.Error("Server failed to start", "time", time.Now(), "err", err)
		return
	}
	logger.Info("Server started", "time", time.Now())
	srv := grpc.NewServer()
	proto.RegisterUserServiceServer(srv, &Server{})
	srv.Serve(listener)
}
