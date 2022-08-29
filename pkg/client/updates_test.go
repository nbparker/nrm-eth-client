package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/nbparker/nrm-eth-client/pkg/proto/nrm"
)

type pbTimestamp struct {
	Seconds, Nanos int
}

type genericUpdateJSON struct {
	Start, Finish pbTimestamp
	Units         int
}

func TestGetUpdatesErrors(t *testing.T) {
	missingPath := "this/will/not/exist"
	cases := []struct {
		inPath string
		want   error
	}{
		{"", fmt.Errorf("no folder specified so no updates to send")},
		{missingPath, fmt.Errorf("open %s: no such file or directory", missingPath)},
	}

	for _, c := range cases {
		updates := make(chan *nrm.GenericUpdate, 10) // buffered channel to allow len check
		errs := make(chan error)

		go GetUpdates(c.inPath, updates, errs)
		if err := <-errs; err.Error() != c.want.Error() {
			t.Errorf("Incorrect error: got '%v', want '%v'", err, c.want)
		}

		// Test for empty channel
		if len(updates) != 0 {
			t.Errorf("Channel not empty: got %d, want 0", len(updates))
		}
	}
}

func TestGetUpdatesFromJSON(t *testing.T) {
	// Create folder to write test files to
	folderPath := "new_folder"
	err := os.Mkdir(folderPath, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Failed to create folder: %s", err)
	}

	// Create json
	j := genericUpdateJSON{
		Start:  pbTimestamp{Seconds: 0, Nanos: 857632152},
		Finish: pbTimestamp{Seconds: 0, Nanos: 857633064},
		Units:  1111111,
	}

	// Write json to file
	file, err := os.CreateTemp(folderPath, "*.json")
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	enc := json.NewEncoder(file)
	enc.Encode(j)

	// Marshall json to GenericUpdate
	b, _ := json.Marshal(j)
	var stored *nrm.GenericUpdate
	json.Unmarshal(b, &stored)

	updates := make(chan *nrm.GenericUpdate)
	errs := make(chan error)
	go GetUpdates(folderPath, updates, errs)

	// Test update
	update := <-updates
	if !proto.Equal(update, stored) {
		t.Errorf("Stored json not matched: got '%v', want '%v'", update, stored)
	}

	// Test errors
	if err, ok := <-errs; err != nil || ok {
		t.Errorf("Unexpected error: %v", err)
	}

	// Cleanup
	err = os.RemoveAll(folderPath)
	if err != nil {
		log.Fatalf("Failed to cleanup: %s", err)
	}

	// }
}
