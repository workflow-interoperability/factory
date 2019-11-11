package main

import (
	"fmt"
	"log"
	"net"
	"os"

	flag "github.com/spf13/pflag"
	"google.golang.org/grpc"

	"github.com/workflow-interoperability/factory/controller"
	mygrpc "github.com/workflow-interoperability/factory/grpc"
)

func main() {
	var GRPCPORT = os.Getenv("GRPC_PORT")
	if len(GRPCPORT) == 0 {
		GRPCPORT = "8081"
	}
	var grpcPort = flag.StringP("grpc_port", "g", GRPCPORT, "Define the port where grpc service runs")
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", *grpcPort))
	if err != nil {
		log.Fatalf("failed to listen grpc: %v", err)
	}
	log.Printf("Listening on: %s", *grpcPort)
	gs := grpc.NewServer()
	mygrpc.RegisterInstanceServer(gs, &controller.Instance{})
	mygrpc.RegisterFactoryServer(gs, &controller.Factory{})
	mygrpc.RegisterServiceRegistryServer(gs, &controller.ServiceRegistry{})
	gs.Serve(lis)
}
