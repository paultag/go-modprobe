package modprobe

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolve(t *testing.T) {
	path, err := ResolveName("snd")
	if err != nil {
		t.Errorf("%s", err)
	}

	if !strings.Contains(path, "snd") {
		t.Fail()
	}

	_, err = os.Stat(path)
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestResolveCompressed(t *testing.T) {
	moduleRoot = filepath.Join("testing", "test_kernel_module")
	t.Cleanup(func() {
		moduleRoot = getModuleRoot()
	})

	path, err := ResolveName("test")
	if err != nil {
		t.Fatalf("%s", err)
	}

	if !strings.Contains(path, "test") {
		t.Fatalf("expected response path to contain 'test', got %s", path)
	}

	_, err = os.Stat(path)
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestNotFound(t *testing.T) {
	_, err := ResolveName("not-found")
	if err == nil {
		t.Fail()
	}

	if !strings.Contains(err.Error(), "not-found") {
		t.Fail()
	}
}
