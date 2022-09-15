package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/nbparker/nrm-eth-client/pkg/proto/nrm"
)

// GetUpdatesFunc defines a type for use with NRMClient for inputing updates
//
// Use of channels allows non-blocking input and potential for a more event-driven
// function than is implemented in GetUpdatesFromFolder
type GetUpdatesFunc func(updates chan *nrm.ResourceUpdate, errs chan error)

// GetUpdatesFromFolder returns GetUpdatesFunc
// Reads json from files and adds to updates chan. Closes channel once all files read.
func GetUpdatesFromFolder(path string) GetUpdatesFunc {
	return func(updates chan *nrm.ResourceUpdate, errs chan error) {
		defer close(updates)
		defer close(errs)

		// No updates if no folder name
		if path == "" {
			errs <- fmt.Errorf("no folder specified so no updates to send")
			return
		}

		files, err := os.ReadDir(path)
		if err != nil {
			errs <- err
			return
		}

		// Iterate files in updates folder, adding updates
		for _, file := range files {
			path := filepath.Join(path, file.Name())
			log.Printf("Reading file: %s", path)

			data, err := os.ReadFile(path)
			if err != nil {
				log.Fatalf("Failed to read updates file: %v", err)
			}

			// Attempt to marshall individual json object
			var _update *nrm.ResourceUpdate
			if err := json.Unmarshal(data, &_update); err == nil {
				updates <- _update

				time.Sleep(time.Millisecond * 200) // TODO remove
				continue
			}

			// Attempt to marshall list of json objects
			var _updates []*nrm.ResourceUpdate
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
}
