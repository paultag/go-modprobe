package modprobe

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInit(t *testing.T) {
	if os.Getenv("TEST_MODULE_INIT") == "" {
		t.Skipf("Skipping module init testing")
	}

	modulePath := filepath.Join("testing", "test_kernel_module", "test.ko.xz")

	f, err := os.Open(modulePath)
	if err != nil {
		t.Fatalf("failed to open test module file: %s", err)
	}

	err = Init(f, "")
	if err != nil {
		t.Fatalf("failed to init test module: %s", err)
	}

	t.Cleanup(func() {
		err := Remove("test")
		if err != nil {
			t.Errorf("failed to remove test module: %s", err)
		}
	})
}
