package modprobe

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// Given a short module name (such as `g_ether`), determine where the kernel
// module is located, determine any dependencies, and load all required modules.
func Load(module string) error {
	path, err := ResolveName(module)
	if err != nil {
		return err
	}

	order, err := Dependencies(path)
	if err != nil {
		return err
	}

	fmt.Printf("%s\n", order)

	for _, module := range order {
		fd, err := os.Open(module)
		if err != nil {
			return err
		}
		/* not doing a defer since we're in a loop */
		if err := Init(fd, ""); err != nil && err != unix.EEXIST {
			fd.Close()
			return err
		}
		fd.Close()
	}

	return nil
}
