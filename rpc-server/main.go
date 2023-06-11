package main

import (
	"database/sql"
	"os"
	"os/signal"
	"syscall"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/rs/zerolog/log"

	"github.com/jh-chee/kitewave/rpc-server/database"
	rpc "github.com/jh-chee/kitewave/rpc-server/kitex_gen/rpc/imservice"

	"github.com/jh-chee/kitewave/rpc-server/handler"
	"github.com/jh-chee/kitewave/rpc-server/repository"
	"github.com/jh-chee/kitewave/rpc-server/service"
)

func main() {
	db, err := database.InitDB()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to db")
		os.Exit(1)
	}
	registerCleanup(db)

	msgRepo := repository.NewMessageRepository(db)
	msgSvc := service.NewMessageRepository(msgRepo)
	msgHandler := handler.NewMessageHandler(msgSvc)

	r, err := etcd.NewEtcdRegistry([]string{"kitewave-etcd:2379"})
	if err != nil {
		log.Fatal().Err(err)
		os.Exit(1)
	}

	svr := rpc.NewServer(
		msgHandler,
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "kitewave.rpc.server",
			},
		),
	)

	if err = svr.Run(); err != nil {
		log.Fatal().Err(err)
		os.Exit(1)
	}

	waitForTerminationSignal()
}

func registerCleanup(db *sql.DB) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		db.Close()
		os.Exit(0)
	}()
}

func waitForTerminationSignal() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
}
