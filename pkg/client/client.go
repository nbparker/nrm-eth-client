package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
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
	go c.getUpdates(updates)
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

func (c *NRMClient) getUpdates(updates chan *nrm.GenericUpdate) {
	// No updates if no folder name
	if c.UpdatesFolderPath == "" {
		log.Fatal("No folder specified so no updates to send")
	}

	files, err := os.ReadDir(c.UpdatesFolderPath)
	if err != nil {
		log.Fatalf("Failed to read dir: %v", err)
	}

	// Iterate files in updates folder, adding updates
	for _, file := range files {
		path := filepath.Join(c.UpdatesFolderPath, file.Name())
		log.Printf("Reading file: %s", path)

		data, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read updates file: %v", err)
		}

		// Attempt to marshall individual json object
		var _update *nrm.GenericUpdate
		if err := json.Unmarshal(data, &_update); err == nil {
			updates <- _update

			time.Sleep(time.Millisecond * 200) // TODO remove
			continue
		}

		// Attempt to marshall list of json objects
		var _updates []*nrm.GenericUpdate
		if err := json.Unmarshal(data, &_updates); err == nil {
			for _, _update := range _updates {
				updates <- _update

				time.Sleep(time.Millisecond * 200) // TODO remove
			}
			continue
		}

		// Log but continue to next file
		fmt.Printf("Failed to read json from file: %s", path)
	}
}
