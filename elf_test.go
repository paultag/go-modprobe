package modprobe_test

import (
	"strings"
	"testing"

	"pault.ag/go/modprobe"
)

func TestResolve(t *testing.T) {
	path, err := modprobe.ResolveName("snd")
	if err != nil {
		t.Errorf("%s", err)
	}

	if !strings.Contains(path, "snd") {
		t.Fail()
	}

	if !strings.Contains(path, "/usr/modules") {
		t.Fail()
	}
}
