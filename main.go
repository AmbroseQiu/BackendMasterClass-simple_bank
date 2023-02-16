package main

import (
	"context"
	"database/sql"
	"log"
	"net"
	"net/http"

	"github.com/backendmaster/simple_bank/api"
	db "github.com/backendmaster/simple_bank/db/sqlc"
	"github.com/backendmaster/simple_bank/gapi"
	"github.com/backendmaster/simple_bank/pb"
	"github.com/backendmaster/simple_bank/util"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("Load Config Failed: %v", err)
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatalf("can't not connect to database %v", err)
	}

	store := db.NewStore(conn)
	go runGateWayServer(config, store)
	runGrpcServer(config, store)
	// runGinServer(config, store)
}
func runGrpcServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("Can't not create gapi server %v", err)
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatalf("can't not create listener %v", err)
	}

	log.Printf("start grpc server at %v", listener.Addr().String())
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("can't not create grpc server %v", err)
	}

}
func runGateWayServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
	if err != nil {
		log.Fatalf("Can't not create gapi server %v", err)
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
		log.Fatalf("can not register handler server %v", err)
	}

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatalf("can't not create listener %v", err)
	}

	log.Printf("start http gateway server at %v", listener.Addr().String())
	err = http.Serve(listener, mux)
	if err != nil {
		log.Fatalf("can't not create http gateway server %v", err)
	}

}
func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatalf("Can't not create gin server %v", err)
	}

	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatalf("can't not start server %v", err)
	}
}
