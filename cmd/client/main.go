package main

import (
	"context"
	"flag"
	"io"
	"log"
	"time"

	"github.com/nbparker/nrm-eth-client/pkg/proto/nrm"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
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
	client := nrm.NewNaturalResourceManagementClient(conn)

	log.Printf("Starting storing: %v", time.Now())
	runNRMStore(client)
	log.Printf("Finished storing: %v", time.Now())
}

func runNRMStore(client nrm.NaturalResourceManagementClient) {
	updates := []*nrm.GenericUpdate{
		{Units: 0, Start: timestamppb.Now(), Finish: timestamppb.Now()},
		{Units: 1, Start: timestamppb.Now(), Finish: timestamppb.Now()},
		{Units: 2, Start: timestamppb.Now(), Finish: timestamppb.Now()},
		{Units: 3, Start: timestamppb.Now(), Finish: timestamppb.Now()},
		{Units: 4, Start: timestamppb.Now(), Finish: timestamppb.Now()},
		{Units: 5, Start: timestamppb.Now(), Finish: timestamppb.Now()},
	}

	// Connect to stream
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	stream, err := client.Store(ctx)
	if err != nil {
		log.Fatalf("client.Store failed: %v", err)
	}

	go handleNRMSummaries(make(chan struct{}), stream)

	// Send protos to store
	// TODO take as input
	for i, update := range updates {
		if err := stream.Send(update); err != nil {
			log.Fatalf("client.Store: stream.Send(%v) failed: %v", update, err)
		}

		// TODO remove
		time.Sleep(time.Second * time.Duration(i+1))
	}
}

func handleNRMSummaries(waitc chan struct{}, stream nrm.NaturalResourceManagement_StoreClient) {
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			// Stream closed
			close(waitc)
			return
		}
		if err != nil {
			log.Fatalf("client.Store recv failed: %v", err)
		}
		log.Printf(
			"Received summary. Success: %v, Attempts: %d, Last Attempt: %v (error: %s)",
			in.Success, in.Attempts, in.LastAttemptedAt, in.FailureMessage,
		)
	}
}
