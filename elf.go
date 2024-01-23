package modprobe

import (
	"bytes"
	"compress/gzip"
	"debug/elf"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/klauspost/compress/zstd"
	"github.com/pierrec/lz4"
	"github.com/xi2/xz"
	"golang.org/x/sys/unix"
)

var (
	// get the root directory for the kernel modules. If this line panics,
	// it's because getModuleRoot has failed to get the uname of the running
	// kernel (likely a non-POSIX system, but maybe a broken kernel?)
	moduleRoot = getModuleRoot()
)

// Get the module root (/lib/modules/$(uname -r)/)
func getModuleRoot() string {
	uname := unix.Utsname{}
	if err := unix.Uname(&uname); err != nil {
		panic(err)
	}

	i := 0
	for ; uname.Release[i] != 0; i++ {
	}

	return filepath.Join(
		"/lib/modules",
		string(uname.Release[:i]),
	)
}

// Get a path relitive to the module root directory.
func modulePath(path string) string {
	return filepath.Join(moduleRoot, path)
}

// ResolveName will, given a module name (such as `g_ether`) return an absolute
// path to the .ko that provides that module.
func ResolveName(name string) (string, error) {
	paths, err := generateMap()
	if err != nil {
		return "", err
	}

	fsPath := paths[name]
	if !strings.HasPrefix(fsPath, moduleRoot) {
		return "", fmt.Errorf("Module '%s' isn't in the module directory", name)
	}

	return fsPath, nil
}

// Open every single kernel module under the kernel module directory
// (/lib/modules/$(uname -r)/), and parse the ELF headers to extract the
// module name.
func generateMap() (map[string]string, error) {
	return elfMap(moduleRoot)
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

			// switch to regex probably idk
			if !strings.Contains(path, ".ko") {
				return nil
			}

			if filepath.Base(path)[0] == '.' {
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

func ModInfo(file *os.File) (map[string]string, error) {
	ext := filepath.Ext(file.Name())
	var r io.Reader
	var err error

	switch ext {
	case ".ko":
		r = file
	case ".zst":
		r, err = zstd.NewReader(file)
	case ".xz":
		r, err = xz.NewReader(file, 0)
	case ".lz4":
		r = lz4.NewReader(file)
	case ".gz":
		r, err = gzip.NewReader(file)

	default:
		err = fmt.Errorf("unknown module compression format: %s", ext)
	}

	content, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	f, err := elf.NewFile(bytes.NewReader(content))
	if err != nil {
		return nil, err
	}

	attrs := map[string]string{}

	sec := f.Section(".modinfo")
	if sec == nil {
		return nil, errors.New("missing modinfo section")
	}

	data, err := sec.Data()
	if err != nil {
		return nil, fmt.Errorf("failed to get section data: %w", err)
	}

	for _, info := range bytes.Split(data, []byte{0}) {
		if parts := strings.SplitN(string(info), "=", 2); len(parts) == 2 {
			attrs[parts[0]] = parts[1]
		}
	}

	return attrs, nil
}

func unzstd(w io.Writer, r io.Reader) error {
	zstdReader, err := zstd.NewReader(r)
	if err != nil {
		return fmt.Errorf("failed to create new reader: %v", err)
	}
	defer zstdReader.Close()

	if _, err := io.Copy(w, zstdReader); err != nil {
		return fmt.Errorf("failed writing decompressed bytes to writer: %v", err)
	}
	return nil
}

// Name will, given a file descriptor to a Kernel Module (.ko file), parse the
// binary to get the module name. For instance, given a handle to the file at
// `kernel/drivers/usb/gadget/legacy/g_ether.ko`, return `g_ether`.
func Name(file *os.File) (string, error) {
	mi, err := ModInfo(file)
	if err != nil {
		return "", fmt.Errorf("failed to get module information: %w", err)
	}

	if name, ok := mi["name"]; !ok {
		return "", errors.New("module information is missing name")
	} else {
		return name, nil
	}
}
