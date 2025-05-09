package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"time"

	"google.golang.org/grpc"

	"VK/internal/config"
	"VK/internal/eventbus/consume/handle_message"
	"VK/internal/generated/service"
	"VK/internal/rpc"
	"VK/internal/service/stream_manager"
	"VK/internal/service/subpub"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		slog.Error("couldn't load config", slog.Any("error", err))
	}

	subPubService := subpub.NewSubPub()

	streamManager := stream_manager.New()
	go func() {
		for {
			time.Sleep(24 * time.Hour)
			streamManager.ClearIdpMap()
			slog.Info("Stream_manager: map IDP cleared")
		}
	}()

	consumeHandleMessages := handle_message.New(streamManager)

	rpcHandler := rpc.New(subPubService, streamManager, consumeHandleMessages.Consume)

	rpcServer := grpc.NewServer()
	service.RegisterPubSubServer(rpcServer, rpcHandler)

	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		slog.Error(
			"Couldn't start listening rpc address",
			slog.Any("error", err),
			slog.Any("address", addr),
		)
	}
	defer listener.Close()

	//graceful shutdown
	exit := make(chan os.Signal)
	signal.Notify(exit, os.Interrupt)

	go func() {
		slog.Info("Starting rpc server on", slog.Any("address", addr))

		if err = rpcServer.Serve(listener); err != nil {
			slog.Error("Couldn't start rpc server", slog.Any("error", err))
		}
	}()

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	//graceful shutdown
	<-exit
	exitCtx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	subPubService.Close(exitCtx)

	rpcServer.GracefulStop()

}
