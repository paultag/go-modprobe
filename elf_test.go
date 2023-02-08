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

func TestNotFound(t *testing.T) {
	_, err := modprobe.ResolveName("not-found")
	if err == nil {
		t.Fail()
	}

	if !strings.Contains(err.Error(), "not-found") {
		t.Fail()
	}
}
