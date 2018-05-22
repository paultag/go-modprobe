package modprobe

import (
	"os"
	"path/filepath"
	"strings"

	"debug/elf"

	"golang.org/x/sys/unix"
)

// Open every single kernel module under the kernel module directory
// (/lib/modules/$(uname -r)/), and parse the ELF headers to extract the
// module name.
func Map() (map[string]string, error) {
	uname := unix.Utsname{}
	if err := unix.Uname(&uname); err != nil {
		return nil, err
	}

	i := 0
	for ; uname.Release[i] != 0; i++ {
	}

	return elfMap(filepath.Join("/lib/modules", string(uname.Release[:i])))
}

// Open every single kernel module under the root, and parse the ELF headers to
// extract the module name.
func elfMap(root string) (map[string]string, error) {
	ret := map[string]string{}

	err := filepath.Walk(
		root,
		func(path string, info os.FileInfo, err error) error {
			if !info.Mode().IsRegular() {
				return nil
			}
			fd, err := os.Open(path)
			if err != nil {
				return err
			}
			defer fd.Close()
			name, err := Name(fd)
			if err != nil {
				/* For now, let's just ignore that and avoid adding to it */
				return nil
			}

			ret[name] = path
			return nil
		})

	if err != nil {
		return nil, err
	}

	return ret, nil
}

// Given a file descriptor, go ahead and parse out the module name from the
// Symbols.
func Name(file *os.File) (string, error) {
	f, err := elf.NewFile(file)
	if err != nil {
		return "", err
	}

	syms, err := f.Symbols()
	if err != nil {
		return "", err
	}

	for _, sym := range syms {
		if strings.Compare(sym.Name, "__this_module") == 0 {
			section := f.Sections[sym.Section]
			data, err := section.Data()
			if err != nil {
				return "", err
			}

			data = data[24:]
			i := 0
			for ; data[i] != 0x00; i++ {
			}
			return string(data[:i]), nil
		}
	}

	return "", nil
	// .gnu.linkonce.this_module
}
