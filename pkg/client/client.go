package client

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/nbparker/nrm-eth-client/pkg/proto/nrm"
)

type NRMClient struct {
	Client     nrm.NaturalResourceManagementClient
	GetUpdates GetUpdatesFunc

	stream  nrm.NaturalResourceManagement_StoreClient
	errors  chan error
	updates chan *nrm.ResourceUpdate
	waitc   chan struct{}
}

// RunStore finds updates and sends to server
//
// Spawns new Goroutines for handleSummaries and getUpdates to allow
// each method to run concurrently. These share the same error channel
// which is monitored by sendUpdates to allow effective error handling
// without the need for blocking functions waiting for an error to be
// returned
func (c *NRMClient) RunStore() error {
	// Connect to stream
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	stream, err := c.Client.Store(ctx)
	if err != nil {
		return fmt.Errorf("client.Store failed: %w", err)
	}

	// Initialise client requirements
	c.stream = stream
	c.updates = make(chan *nrm.ResourceUpdate)
	c.errors = make(chan error)
	c.waitc = make(chan struct{})

	// Spawn as Goroutines
	go c.handleSummaries()
	go c.GetUpdates(c.updates, c.errors)
	err = c.sendUpdates()
	c.stream.CloseSend()

	<-c.waitc
	return err
}

// sendUpdates to server
// iterates update and error channels until closed or error
func (c *NRMClient) sendUpdates() error {
	for {
		select {
		case update, ok := <-c.updates:
			if !ok {
				return nil
			}
			if err := c.stream.Send(update); err != nil {
				return fmt.Errorf("client.Store: stream.Send(%v) failed: %w", update, err)
			}
		case err, ok := <-c.errors:
			if !ok {
				return nil
			}
			return err
		}
	}
}

func (c *NRMClient) handleSummaries() {
	for {
		in, err := c.stream.Recv()
		if err == io.EOF {
			// Stream closed
			close(c.waitc)
			return
		}
		if err != nil {
			c.errors <- fmt.Errorf("client.Store recv failed: %w", err)
		}
		log.Printf(
			"Received summary. Success: %v, Attempts: %d, Last Attempt: %v (error: %s)",
			in.Success, in.Attempts, in.LastAttemptedAt, in.FailureMessage,
		)
	}
}
