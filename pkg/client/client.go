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
	updates chan *nrm.GenericUpdate
	waitc   chan struct{}
	cancel  context.CancelFunc
}

// RunStore finds updates and sends to server
func (c *NRMClient) RunStore() error {
	c.initialise()
	defer c.cancel()

	go c.handleSummaries()
	go c.GetUpdates(c.updates, c.errors)
	err := c.sendUpdates()
	c.stream.CloseSend()

	<-c.waitc
	return err
}

func (c *NRMClient) initialise() error {
	// Connect to stream
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	stream, err := c.Client.Store(ctx)
	if err != nil {
		return fmt.Errorf("client.Store failed: %w", err)
	}

	c.stream = stream
	c.updates = make(chan *nrm.GenericUpdate)
	c.errors = make(chan error)
	c.waitc = make(chan struct{})
	c.cancel = cancel
	return nil
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
