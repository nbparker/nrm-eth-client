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
	Client            nrm.NaturalResourceManagementClient
	UpdatesFolderPath string

	stream nrm.NaturalResourceManagement_StoreClient
}

// RunStore finds updates and sends to server
func (c *NRMClient) RunStore() error {
	// Connect to stream
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	stream, err := c.Client.Store(ctx)
	if err != nil {
		return fmt.Errorf("client.Store failed: %w", err)
	}
	c.stream = stream

	updates := make(chan *nrm.GenericUpdate)
	errs := make(chan error)

	go c.handleSummaries(errs)
	go GetUpdates(c.UpdatesFolderPath, updates, errs)
	return c.sendUpdates(updates, errs)
}

// sendUpdates to server
// iterates update and error channels until closed or error
func (c *NRMClient) sendUpdates(updates chan *nrm.GenericUpdate, errs chan error) error {
	for {
		select {
		case update, ok := <-updates:
			if !ok {
				return nil
			}
			if err := c.stream.Send(update); err != nil {
				return fmt.Errorf("client.Store: stream.Send(%v) failed: %w", update, err)
			}
		case err, ok := <-errs:
			if !ok {
				return nil
			}
			return err
		}
	}
}

func (c *NRMClient) handleSummaries(errs chan error) {
	waitc := make(chan struct{})
	for {
		in, err := c.stream.Recv()
		if err == io.EOF {
			// Stream closed
			close(waitc)
			return
		}
		if err != nil {
			errs <- fmt.Errorf("client.Store recv failed: %w", err)
		}
		log.Printf(
			"Received summary. Success: %v, Attempts: %d, Last Attempt: %v (error: %s)",
			in.Success, in.Attempts, in.LastAttemptedAt, in.FailureMessage,
		)
	}
}
