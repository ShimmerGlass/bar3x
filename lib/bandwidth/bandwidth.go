package bandwidth

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type BW struct {
	iface string

	lastAt       time.Time
	lastR, lastT int
}

func New(iface string) *BW {
	return &BW{
		iface: iface,
		lastR: -1,
		lastT: -1,
	}
}

func (b *BW) Read() (int, int, error) {
	now := time.Now()
	r, t, err := readBandwidth(b.iface)
	if err != nil {
		return 0, 0, err
	}

	if b.lastR == -1 && b.lastT == -1 {
		b.lastR = r
		b.lastT = t
		b.lastAt = now
		return 0, 0, nil
	}

	cr, ct := r-b.lastR, t-b.lastT
	secs := now.Sub(b.lastAt).Seconds()

	cr = int(float64(cr) / secs)
	ct = int(float64(ct) / secs)

	b.lastR, b.lastT = r, t
	b.lastAt = now

	return cr, ct, nil
}

func readLines(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return []string{""}, err
	}
	defer f.Close()

	var ret []string

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			break
		}
		ret = append(ret, strings.Trim(line, "\n"))
	}
	return ret, nil
}

func readBandwidth(dev string) (int, int, error) {
	lines, err := readLines("/proc/net/dev")
	if err != nil {
		return 0, 0, err
	}

	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		if key != dev {
			continue
		}

		value := strings.Fields(strings.TrimSpace(fields[1]))

		r, err := strconv.Atoi(value[0])
		if err != nil {
			return 0, 0, err
		}

		t, err := strconv.Atoi(value[8])
		if err != nil {
			return 0, 0, err
		}

		return r, t, nil
	}

	return 0, 0, fmt.Errorf("dev \"%s\" not found", dev)
}


func IsValidInterface (dev string) bool {
	lines, err := readLines("/proc/net/dev")
	if err != nil {
		return false
	}

	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) < 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		if key == dev {
			return true
		}
	}

	return false
}