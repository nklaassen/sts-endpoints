package main

import (
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gravitational/teleport/lib/utils/golden"
)

// The point of this test is to cause CI to fail in dependabot PRs so I know
// that consumers need to be updated.
func TestEndpointsChanged(t *testing.T) {
	var buf bytes.Buffer
	for _, e := range stsEndpoints() {
		if _, err := buf.WriteString(e); err != nil {
			t.Fatal(err)
		}
		if err := buf.WriteByte('\n'); err != nil {
			t.Fatal(err)
		}
	}

	if golden.ShouldSet() {
		golden.Set(t, buf.Bytes())
		return
	}

	expected := golden.Get(t)
	diff := cmp.Diff(string(expected), buf.String())
	if diff != "" {
		t.Log("Endpoints list has changed (-old, +new)")
		t.Log(diff)
		t.Log("Run the following command to update the expected endpoints list")
		t.Log("> GOLDEN_UPDATE=1 go test")
		t.Fail()
	}
}
