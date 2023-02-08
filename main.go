package main

import (
	"database/sql"
	"log"
	"net"

	"github.com/backendmaster/simple_bank/api"
	db "github.com/backendmaster/simple_bank/db/sqlc"
	"github.com/backendmaster/simple_bank/gapi"
	"github.com/backendmaster/simple_bank/pb"
	"github.com/backendmaster/simple_bank/util"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("Load Config Failed: %v", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("can't not connect to database")
	}

	store := db.NewStore(conn)
	runGrpcServer(config, store)
}
func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("Can't not create gapi server %v", err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal("can't not create listener")
	}

	log.Printf("start grpc server at %v", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatal("can't not create grpc server")
	}

}
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("Can't not create gin server %v", err)
	}

	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatalf("can't not start server")
	}
}
