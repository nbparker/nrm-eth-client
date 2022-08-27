package client

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/nbparker/nrm-eth-client/pkg/proto/nrm"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type NRMClient struct{}

// TODO replace updates with awaitable
func (c *NRMClient) RunStore(client nrm.NaturalResourceManagementClient) {
	// Connect to stream
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	stream, err := client.Store(ctx)
	if err != nil {
		log.Fatalf("client.Store failed: %v", err)
	}

	go c.handleSummaries(make(chan struct{}), stream)

	// Send protos to store
	updates := c.getUpdates()
	for i, update := range updates {
		if err := stream.Send(update); err != nil {
			log.Fatalf("client.Store: stream.Send(%v) failed: %v", update, err)
		}

		// TODO remove
		time.Sleep(time.Second * time.Duration(i+1))
	}
}

func (c *NRMClient) handleSummaries(waitc chan struct{}, stream nrm.NaturalResourceManagement_StoreClient) {
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

func (c *NRMClient) getUpdates() []*nrm.GenericUpdate {
	return []*nrm.GenericUpdate{
		{Units: 0, Start: timestamppb.Now(), Finish: timestamppb.Now()},
		{Units: 1, Start: timestamppb.Now(), Finish: timestamppb.Now()},
		{Units: 2, Start: timestamppb.Now(), Finish: timestamppb.Now()},
		{Units: 3, Start: timestamppb.Now(), Finish: timestamppb.Now()},
		{Units: 4, Start: timestamppb.Now(), Finish: timestamppb.Now()},
		{Units: 5, Start: timestamppb.Now(), Finish: timestamppb.Now()},
	}
}