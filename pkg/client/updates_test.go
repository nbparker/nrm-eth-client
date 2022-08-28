package client

import (
	"fmt"
	"testing"

	"github.com/nbparker/nrm-eth-client/pkg/proto/nrm"
)

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

		got := GetUpdates(c.inPath, updates)
		if got.Error() != c.want.Error() {
			t.Errorf("Incorrect error: got '%v', want '%v'", got, c.want)
		}

		// Test for empty channel
		if len(updates) != 0 {
			t.Errorf("Channel not empty: got %d, want 0", len(updates))
		}
	}
}
