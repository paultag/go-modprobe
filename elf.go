package modprobe

import (
	"os"
	"strings"

	"debug/elf"
)

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
