package modprobe

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// Status of a loaded Kernel module.
type Status struct {
	//
	Name string

	//
	Size int

	//
	Instances int

	//
	Dependencies []string

	//
	State string
}

// List all the currently loaded modules.
func List() ([]Status, error) {
	file, err := os.Open("/proc/modules")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	ret := []Status{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		status, err := parseStatus(scanner.Text())
		if err != nil {
			return nil, err
		}
		ret = append(ret, *status)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}

// Parse the proc modules status line into a Status struct.
func parseStatus(line string) (*Status, error) {
	chunks := strings.SplitN(line, " ", 6)

	size, err := strconv.ParseInt(chunks[1], 10, 32)
	if err != nil {
		return nil, err
	}

	instances, err := strconv.ParseInt(chunks[2], 10, 32)
	if err != nil {
		return nil, err
	}

	dependencies := []string{}
	if strings.Compare(chunks[3], "-") != 0 {
		dependencies = strings.Split(chunks[3], ",")
	}

	return &Status{
		Name:         chunks[0],
		Size:         int(size),
		Instances:    int(instances),
		Dependencies: dependencies,
		State:        chunks[4],
	}, nil
}
