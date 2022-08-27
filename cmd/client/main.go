package main

import (
	"flag"
	"log"
	"time"

	"github.com/nbparker/nrm-eth-client/pkg/client"
	"github.com/nbparker/nrm-eth-client/pkg/proto/nrm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	serverAddr = flag.String("addr", "localhost:50051", "The server address in the format of host:port")
)

func main() {
	flag.Parse()

	var opts []grpc.DialOption
	// Run without transport creds
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to dial: %v", err)
	}
	defer conn.Close()

	cli := nrm.NewNaturalResourceManagementClient(conn)

	log.Printf("Starting storing: %v", time.Now())
	nrmClient := &client.NRMClient{}
	nrmClient.RunStore(cli)
	log.Printf("Finished storing: %v", time.Now())
}
