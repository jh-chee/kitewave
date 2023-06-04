package main

import (
	"os"

	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	etcd "github.com/kitex-contrib/registry-etcd"
	"github.com/rs/zerolog/log"

	db "github.com/jh-chee/kitewave/rpc-server/database"
	rpc "github.com/jh-chee/kitewave/rpc-server/kitex_gen/rpc/imservice"

	"github.com/jh-chee/kitewave/rpc-server/handler"
	"github.com/jh-chee/kitewave/rpc-server/repository"
	"github.com/jh-chee/kitewave/rpc-server/service"
)

func main() {
	db, err := db.InitDB()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to connect to db")
		os.Exit(1)
	}

	msgRepo := repository.NewMessageRepository(db)
	msgSvc := service.NewMessageRepository(msgRepo)
	msgHandler := handler.NewMessageHandler(msgSvc)

	r, err := etcd.NewEtcdRegistry([]string{"etcd:2379"})
	if err != nil {
		log.Fatal().Err(err)
		os.Exit(1)
	}

	svr := rpc.NewServer(
		msgHandler,
		server.WithRegistry(r),
		server.WithServerBasicInfo(
			&rpcinfo.EndpointBasicInfo{
				ServiceName: "demo.rpc.server",
			},
		),
	)

	if err = svr.Run(); err != nil {
		log.Fatal().Err(err)
		os.Exit(1)
	}
}
