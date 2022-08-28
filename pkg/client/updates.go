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

func GetUpdates(folderPath string, updates chan *nrm.GenericUpdate) error {
	// No updates if no folder name
	if folderPath == "" {
		return fmt.Errorf("no folder specified so no updates to send")
	}

	files, err := os.ReadDir(folderPath)
	if err != nil {
		return err
	}

	// Iterate files in updates folder, adding updates
	for _, file := range files {
		path := filepath.Join(folderPath, file.Name())
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

	return nil

}
