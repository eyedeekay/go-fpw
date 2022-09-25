package fcw

import (
	"testing"
)

func TestFPWRun(t *testing.T) {
	err := Run()
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Success")
}
