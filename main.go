package main

import (
	"context"
	"net"
	"net/http"
	"os"

	"github.com/backendmaster/simple_bank/api"
	"github.com/backendmaster/simple_bank/db/gorm"
	db "github.com/backendmaster/simple_bank/db/sqlc"
	"github.com/backendmaster/simple_bank/delivery"
	"github.com/backendmaster/simple_bank/gapi"
	"github.com/backendmaster/simple_bank/pb"
	repository "github.com/backendmaster/simple_bank/repository/postgresql"
	"github.com/backendmaster/simple_bank/token"
	"github.com/backendmaster/simple_bank/usercase"
	"github.com/backendmaster/simple_bank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("Load Config Failed: ")
	}

	// conn, err := sql.Open(config.DBDriver, config.DBSource)
	// if err != nil {
	// 	log.Fatal().Err(err).Msg("can't not connect to database ")
	// }

	// store := db.NewStore(conn)

	runCleanHttpServer(config)
	// go runGateWayServer(config, store)
	// runGrpcServer(config, store)
	// runGinServer(config, store)
}

func runCleanHttpServer(config util.Config) {
	dbClinet := gorm.DBClient{}
	dbClinet.Connect(config)
	userRepo := repository.NewpostgresqlUserRepository(dbClinet.Client)
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't not create token ")
	}
	usersTableUseCase := usercase.NewusersTableUserCase(config, userRepo, tokenMaker)
	handler := delivery.NewUsersHandlerDelivery(usersTableUseCase)
	handler.SetupRouter()
	handler.Start(config.HTTPServerAddress)
}
func runGrpcServer(config util.Config, store db.Store) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't not create gapi server ")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("can't not create listener ")
	}

	log.Info().Msgf("start grpc server at %s", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal().Err(err).Msg("can't not create grpc server ")
	}

}
func runGateWayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't not create gapi server ")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})
	grpcMux := runtime.NewServeMux(jsonOption)

	contex, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = pb.RegisterSimpleBankHandlerServer(contex, grpcMux, server)
	if err != nil {
		log.Fatal().Err(err).Msg("can not register handler server ")
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("can't not create listener ")
	}

	log.Info().Msgf("start http gateway server at %s", listener.Addr().String())
	hander := gapi.HttpLogger(mux)
	err = http.Serve(listener, hander)
	if err != nil {
		log.Fatal().Err(err).Msg("can't not create http gateway server ")
	}

}
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't not create gin server ")
	}

	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal().Err(err).Msg("can't not start server ")
	}
}

func runGormServer(config util.Config, client gorm.DBClient) {
	server, err := gorm.NewHttpServer(config, client.Client)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't not create gorm server ")
	}

	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal().Err(err).Msg("can't not start gorm server ")
	}
}
