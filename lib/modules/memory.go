package modules

import (
	"fmt"
	"os"

	"github.com/davidscholberg/go-i3barjson"
)

// Memory represents the configuration for the memory block.
type Memory struct {
	BlockConfigBase `yaml:",inline"`
	CritMem         float64 `yaml:"crit_mem"`
}

// UpdateBlock updates the system memory block status.
func (c Memory) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)
	var memAvail, memFree, memTotal int64
	r, err := os.Open("/proc/meminfo")
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	defer r.Close()
	_, err = fmt.Fscanf(
		r,
		"MemTotal: %d kB\nMemFree: %d kB\nMemAvailable: %d ",
		&memTotal, &memFree, &memAvail)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	memTotalG := float64(memTotal) / 1048576
	memUsedG := float64(memTotal-memAvail) / 1048576
	memAvailG := float64(memAvail) / 1048576
	if memAvailG < c.CritMem {
		b.Urgent = true
	} else {
		b.Urgent = false
	}
	b.FullText = fmt.Sprintf(fullTextFmt, fmt.Sprintf("%.2fG / %.2fG", memUsedG, memTotalG))
}
