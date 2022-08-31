package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"testing"

	"github.com/nbparker/nrm-eth-client/pkg/proto/nrm"
	"google.golang.org/protobuf/proto"
)

type pbTimestamp struct {
	Seconds, Nanos int
}

type resourceUpdateJSON struct {
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
		updates := make(chan *nrm.ResourceUpdate, 10) // buffered channel to allow len check
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
	cases := [][]interface{}{
		// -- single files
		// empty object
		{resourceUpdateJSON{}},
		// single object
		{
			resourceUpdateJSON{
				Start:  pbTimestamp{Seconds: 0, Nanos: 857632152},
				Finish: pbTimestamp{Seconds: 0, Nanos: 857633064},
				Units:  1111111,
			},
		},
		// single object in array
		{
			[]resourceUpdateJSON{
				{
					Start:  pbTimestamp{Seconds: 0, Nanos: 857632152},
					Finish: pbTimestamp{Seconds: 0, Nanos: 857633064},
					Units:  1111111,
				},
			},
		},
		// multiple objects in array
		{
			[]resourceUpdateJSON{
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
		},
		// -- multiple files
		// single objects
		{
			resourceUpdateJSON{
				Start:  pbTimestamp{Seconds: 1111, Nanos: 1111},
				Finish: pbTimestamp{Seconds: 1111, Nanos: 1111},
				Units:  1111,
			},
			resourceUpdateJSON{
				Start:  pbTimestamp{Seconds: 2222, Nanos: 2222},
				Finish: pbTimestamp{Seconds: 2222, Nanos: 2222},
				Units:  2222,
			},
		},
		// mixed objects
		{
			resourceUpdateJSON{
				Start:  pbTimestamp{Seconds: 1111, Nanos: 1111},
				Finish: pbTimestamp{Seconds: 1111, Nanos: 1111},
				Units:  1111,
			},
			[]resourceUpdateJSON{
				{
					Start:  pbTimestamp{Seconds: 1111, Nanos: 1111},
					Finish: pbTimestamp{Seconds: 1111, Nanos: 1111},
					Units:  1111,
				},
				{
					Start:  pbTimestamp{Seconds: 2222, Nanos: 2222},
					Finish: pbTimestamp{Seconds: 2222, Nanos: 2222},
					Units:  2222,
				},
				{
					Start:  pbTimestamp{Seconds: 3333, Nanos: 3333},
					Finish: pbTimestamp{Seconds: 3333, Nanos: 3333},
					Units:  3333,
				},
			},
			[]resourceUpdateJSON{
				{
					Start:  pbTimestamp{Seconds: 7777, Nanos: 7777},
					Finish: pbTimestamp{Seconds: 7777, Nanos: 7777},
					Units:  7777,
				},
				{
					Start:  pbTimestamp{Seconds: 6666, Nanos: 6666},
					Finish: pbTimestamp{Seconds: 6666, Nanos: 6666},
					Units:  6666,
				},
				{
					Start:  pbTimestamp{Seconds: 4444, Nanos: 4444},
					Finish: pbTimestamp{Seconds: 4444, Nanos: 4444},
					Units:  4444,
				},
			},
		},
	}

	const folderPath = "new_folder"
	for _, jsons := range cases {
		setup(folderPath)

		var stored []*nrm.ResourceUpdate
		for i, _json := range jsons {
			// Specify filename to ensure read in order
			writeJsonToFile(folderPath, fmt.Sprintf("%d.json", i), _json)
			_stored := jsonToResourceUpdates(_json)
			stored = append(stored, _stored...)
		}

		updates := make(chan *nrm.ResourceUpdate)
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

func writeJsonToFile(folderPath, fileName string, _json interface{}) {
	// Write json to file
	file, err := os.CreateTemp(folderPath, fileName)
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	enc := json.NewEncoder(file)
	if err := enc.Encode(_json); err != nil {
		log.Fatalf("Failed to encode: %v", err)
	}
}

func jsonToResourceUpdates(_json interface{}) []*nrm.ResourceUpdate {
	// Marshall json to ResourceUpdate(s)
	var stored []*nrm.ResourceUpdate
	b, err := json.Marshal(_json)
	if err != nil {
		log.Fatalf("Failed to marshal: %v", err)
	}

	switch reflect.TypeOf(_json).Kind() {
	case reflect.Slice, reflect.Array:
		json.Unmarshal(b, &stored)
	default:
		var stored_ *nrm.ResourceUpdate
		json.Unmarshal(b, &stored_)
		stored = []*nrm.ResourceUpdate{stored_}
	}
	return stored
}
