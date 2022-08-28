package client

import (
	"context"
	"io"
	"log"
	"time"

	"github.com/nbparker/nrm-eth-client/pkg/proto/nrm"
)

type NRMClient struct {
	Client            nrm.NaturalResourceManagementClient
	UpdatesFolderPath string

	stream nrm.NaturalResourceManagement_StoreClient
}

// TODO replace updates with awaitable
func (c *NRMClient) RunStore() {
	// Connect to stream
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	stream, err := c.Client.Store(ctx)
	if err != nil {
		log.Fatalf("client.Store failed: %v", err)
	}
	c.stream = stream

	go c.handleSummaries()

	// Send protos to store
	updates := make(chan *nrm.GenericUpdate)
	go GetUpdates(c.UpdatesFolderPath, updates)
	for update := range updates {
		if err := stream.Send(update); err != nil {
			log.Fatalf("client.Store: stream.Send(%v) failed: %v", update, err)
		}
	}
}

func (c *NRMClient) handleSummaries() {
	waitc := make(chan struct{})
	for {
		in, err := c.stream.Recv()
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
