package cpu

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

func readSpec() (CPUSpec, error) {
	var spec CPUSpec

	file, err := os.Open("/proc/cpuinfo")
	if err != nil {
		return spec, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "vendor_id":
			spec.Vendor = value
		case "model name":
			spec.ModelName = value
		case "cpu cores":
			if n, err := strconv.Atoi(value); err == nil {
				spec.Cores = n
			}
		case "siblings":
			if n, err := strconv.Atoi(value); err == nil {
				spec.Threads = n
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return spec, err
	}

	return spec, nil
}
