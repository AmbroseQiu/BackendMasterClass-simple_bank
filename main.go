package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/backendmaster/simple_bank/api"
	"github.com/backendmaster/simple_bank/db/gorm"
	db "github.com/backendmaster/simple_bank/db/sqlc"
	"github.com/backendmaster/simple_bank/delivery"
	"github.com/backendmaster/simple_bank/gapi"
	"github.com/backendmaster/simple_bank/pb"
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

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("can't not connect to database ")
	}

	// store := db.NewStore(conn)

	runGormHttpServer(config, conn)
	// go runGateWayServer(config, store)
	// runGrpcServer(config, store)
	// runGinServer(config, store)
}

// type User struct {
// 	Name  string `json:"name"`
// 	Email string `json:"email"`
// }

// type Response struct {
// 	User
// 	UserAgent string
// 	ClientIP  string
// }

// func main() {
// 	http.HandleFunc("/", handler)
// 	http.ListenAndServe(":8080", nil)
// }

// func handler(w http.ResponseWriter, r *http.Request) {
// 	// Read the request body
// 	body, err := ioutil.ReadAll(r.Body)
// 	if err != nil {
// 		http.Error(w, "Unable to read request body", http.StatusBadRequest)
// 		return
// 	}

// 	var user User

// 	err = json.Unmarshal(body, &user)
// 	if err != nil {
// 		http.Error(w, "Bad Request", http.StatusBadRequest)
// 		return
// 	}

// 	rsp := Response{
// 		User:      user,
// 		UserAgent: r.Header.Get("User-Agent"),
// 		ClientIP:  r.Header.Get("X-Forwarded-For"),
// 	}
// 	if rsp.ClientIP == "" {
// 		rsp.ClientIP = r.RemoteAddr
// 	}
// 	// Close the request body
// 	defer r.Body.Close()

// 	// Print the request body
// 	// fmt.Fprintf(w, "Request body: %s", body)
// 	// fmt.Fprintf(w, "Request user: %v", user)
// 	json.NewEncoder(w).Encode(rsp)
// }

func runGormHttpServer(config util.Config, conn *sql.DB) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	server, err := delivery.NewGormServer(config, conn)
	if err != nil {
		log.Fatal().Err(err).Msg("Can't not create gapi server ")
	}
	err = server.Start(config.HTTPServerAddress)

	if err != nil {
		log.Fatal().Err(err).Msg("can't not start server ")
	}
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
