package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/spacemeshos/go-spacemesh/log"
	"github.com/spacemeshos/node-mock/services"
)

var (
	//Options
	flagServer = flag.Bool("server", false, "(option) start server")

	//Params
	flagConfig = flag.String("config", "", "(param) config")

	//Debug
)

// GrpcService -
type GrpcService struct {
	Server *grpc.Server
	Port   uint
}

func (s GrpcService) startServices() error {
	addr := ":" + strconv.Itoa(int(s.Port))

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	log.Info("grpc API listening on port %d", s.Port)

	// start serving - this blocks until err or server is stopped
	err = s.Server.Serve(listen)
	if err != nil {
		return err
	}

	return nil
}

// NewGrpcService create a new grpc service using config data.
func NewGrpcService(port uint) *GrpcService {
	options := []grpc.ServerOption{
		// XXX: this is done to prevent routers from cleaning up our connections (e.g aws load balances..)
		// TODO: these parameters work for now but we might need to revisit or add them as configuration
		// TODO: Configure maxconns, maxconcurrentcons ..
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle:     time.Minute * 120,
			MaxConnectionAge:      time.Minute * 180,
			MaxConnectionAgeGrace: time.Minute * 10,
			Time:                  time.Minute,
			Timeout:               time.Minute * 3,
		}),
	}

	server := grpc.NewServer(options...)
	return &GrpcService{
		Server: server,
		Port:   uint(port),
	}
}

func startServer(port uint) *GrpcService {
	grpcService := NewGrpcService(port)

	services.InitMocker(grpcService.Server)

	grpcService.startServices()

	return grpcService
}

func main() {
	flag.Parse()

	var err error

	if len(*flagConfig) == 0 {
		fmt.Println("ERROR: -config is mandatory")
		os.Exit(1)
	}

	services.Config, err = parseConfig(*flagConfig)
	if err != nil {
		fmt.Printf("ERROR: %v\n", err)
		os.Exit(1)
	}

	if *flagServer {
		startServer(services.Config.RPCPort)
	}
}
