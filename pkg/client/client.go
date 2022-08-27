package client

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"time"

	"github.com/nbparker/nrm-eth-client/pkg/proto/nrm"
)

type NRMClient struct {
	Client          nrm.NaturalResourceManagementClient
	JsonUpdatesFile string

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
	updates := c.getUpdates()
	for i, update := range updates {
		if err := stream.Send(update); err != nil {
			log.Fatalf("client.Store: stream.Send(%v) failed: %v", update, err)
		}

		// TODO remove
		time.Sleep(time.Second * time.Duration(i+1))
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

func (c *NRMClient) getUpdates() []*nrm.GenericUpdate {
	var updates []*nrm.GenericUpdate
	// No updates if no file
	if c.JsonUpdatesFile == "" {
		return updates
	}

	data, err := ioutil.ReadFile(c.JsonUpdatesFile)
	if err != nil {
		log.Fatalf("Failed to read updates file: %v", err)
	}
	if err := json.Unmarshal(data, &updates); err != nil {
		log.Fatalf("Failed to load updates from json: %v", err)
	}
	return updates
}
