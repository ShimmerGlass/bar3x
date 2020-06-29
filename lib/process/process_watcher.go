package process

import (
	"bytes"
	"io"
	"os"
	"strconv"
	"time"
)

type ProcessWatcher struct {
	statsBuf *bytes.Buffer
	MaxCPU   string
	MaxRAM   string
}

func New() *ProcessWatcher {
	p := &ProcessWatcher{
		statsBuf: &bytes.Buffer{},
	}
	go p.watch()
	return p
}

func (w *ProcessWatcher) watch() {
	var lastProcessTimes map[string]int

	for range time.Tick(2 * time.Second) {
		pids, err := w.listProcesses()
		if err != nil {
			continue
		}

		cpuVals := map[string]int{}
		ramPerName := map[string]int{}
		cpuPerName := map[string]int{}
		for _, pid := range pids {
			name, cpu, ram, err := w.processStats(pid)
			if err != nil {
				continue
			}

			cpuVals[pid] = cpu

			lastTime, ok := lastProcessTimes[pid]
			if !ok {
				continue
			}

			timeD := cpu - lastTime
			cpuPerName[name] += timeD
			ramPerName[name] += ram
		}

		maxTime := 0
		maxRAM := 0
		maxCPUName := ""
		maxRAMName := ""

		for name, val := range cpuPerName {
			if val > maxTime {
				maxTime = val
				maxCPUName = name
			}
		}
		for name, val := range ramPerName {
			if val > maxRAM {
				maxRAM = val
				maxRAMName = name
			}
		}

		lastProcessTimes = cpuVals
		w.MaxCPU = maxCPUName
		w.MaxRAM = maxRAMName
	}
}

func (w *ProcessWatcher) listProcesses() ([]string, error) {
	d, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer d.Close()

	results := make([]string, 0, 50)
	for {
		names, err := d.Readdirnames(10)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		for _, name := range names {
			// We only care if the name starts with a numeric
			if name[0] < '0' || name[0] > '9' {
				continue
			}

			// From this point forward, any errors we just ignore, because
			// it might simply be that the process doesn't exist anymore.
			results = append(results, name)
		}
	}

	return results, nil
}

func (w *ProcessWatcher) processStats(pid string) (string, int, int, error) {
	f, err := os.Open("/proc/" + pid + "/stat")
	if err != nil {
		return "", 0, 0, err
	}
	defer f.Close()

	w.statsBuf.ReadFrom(f)
	defer w.statsBuf.Reset()

	stats := w.statsBuf.Bytes()

	cpu := 0
	ram := 0

	// extract the name
	p := bytes.IndexByte(stats, '(')
	p2 := bytes.IndexByte(stats, ')')
	name := string(stats[p+1 : p2])

	adv := p2
	i := 1

	for {
		p := bytes.IndexByte(stats[adv:], ' ')
		if i == 13 || i == 14 {
			s := stats[adv : p+adv]
			t, err := strconv.Atoi(string(s))
			if err != nil {
				panic(err)
			}

			cpu += t
		}
		if i == 23 {
			s := stats[adv : p+adv]
			t, err := strconv.Atoi(string(s))
			if err != nil {
				panic(err)
			}

			ram += t
		}
		if i == 23 {
			break
		}
		adv = adv + p + 1
		i++
	}

	return name, cpu, ram, nil
}
