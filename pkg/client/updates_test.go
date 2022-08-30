package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
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

		go GetUpdatesFromFolder(c.inPath)(updates, errs)
		if err := <-errs; err.Error() != c.want.Error() {
			t.Errorf("Incorrect error: got '%v', want '%v'", err, c.want)
		}

		// Test for empty channel
		if len(updates) != 0 {
			t.Errorf("Channel not empty: got %d, want 0", len(updates))
		}
	}
}

func TestGetUpdatesFromJSONArray(t *testing.T) {
	// Create json
	cases := []interface{}{
		// empty object
		genericUpdateJSON{},
		// single object
		genericUpdateJSON{
			Start:  pbTimestamp{Seconds: 0, Nanos: 857632152},
			Finish: pbTimestamp{Seconds: 0, Nanos: 857633064},
			Units:  1111111,
		},
		// single object in array
		[]*genericUpdateJSON{
			{
				Start:  pbTimestamp{Seconds: 0, Nanos: 857632152},
				Finish: pbTimestamp{Seconds: 0, Nanos: 857633064},
				Units:  1111111,
			},
		},
		// multiple objects in array
		[]*genericUpdateJSON{
			{
				Start:  pbTimestamp{Seconds: 1111, Nanos: 1111},
				Finish: pbTimestamp{Seconds: 11111, Nanos: 11111},
				Units:  1111,
			},
			{
				Start:  pbTimestamp{Seconds: 2222, Nanos: 2222},
				Finish: pbTimestamp{Seconds: 22222, Nanos: 22222},
				Units:  2222,
			},
			{
				Start:  pbTimestamp{Seconds: 3333, Nanos: 3333},
				Finish: pbTimestamp{Seconds: 33333, Nanos: 33333},
				Units:  3333,
			},
		},
	}

	const folderPath = "new_folder"
	for _, _json := range cases {
		setup(folderPath)
		writeJsonToFile(folderPath, _json)
		stored := jsonToGenericUpdates(_json)

		updates := make(chan *nrm.GenericUpdate)
		errs := make(chan error)
		go GetUpdatesFromFolder(folderPath)(updates, errs)

		// Test updates
		for _, want := range stored {
			update := <-updates
			if !proto.Equal(update, want) {
				t.Errorf("Stored json not matched: got '%v', want '%v'", update, want)
			}
		}

		// Test errors
		if err, ok := <-errs; err != nil || ok {
			t.Errorf("Unexpected error: %v", err)
		}

		cleanup(folderPath)
	}
}

func cleanup(folderPath string) {
	if err := os.RemoveAll(folderPath); err != nil {
		log.Fatalf("Failed to cleanup: %s", err)
	}
}

func setup(folderPath string) {
	cleanup(folderPath) // in case folder remaining

	// Create folder to write test files to
	err := os.Mkdir(folderPath, 0777)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Failed to create folder: %s", err)
	}
}

func writeJsonToFile(folderPath string, _json interface{}) {
	// Write json to file
	file, err := os.CreateTemp(folderPath, "*.json")
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	enc := json.NewEncoder(file)
	if err := enc.Encode(_json); err != nil {
		log.Fatalf("Failed to encode: %v", err)
	}
}

func jsonToGenericUpdates(_json interface{}) []*nrm.GenericUpdate {
	// Marshall json to GenericUpdate(s)
	var stored []*nrm.GenericUpdate
	b, err := json.Marshal(_json)
	if err != nil {
		log.Fatalf("Failed to marshal: %v", err)
	}

	switch reflect.TypeOf(_json).Kind() {
	case reflect.Slice, reflect.Array:
		json.Unmarshal(b, &stored)
	default:
		var stored_ *nrm.GenericUpdate
		json.Unmarshal(b, &stored_)
		stored = []*nrm.GenericUpdate{stored_}
	}
	return stored
}
