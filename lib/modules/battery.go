package modules

import (
	"fmt"
	"github.com/davidscholberg/go-i3barjson"
	"os"
)

// Battery represents the configuration for the battery block.
type Battery struct {
	BlockConfigBase `yaml:",inline"`
	BatteryNumber   int     `yaml:"battery_number"`
	CritBattery     float64 `yaml:"crit_battery"`
}

const (
	acFilePath = "/sys/class/power_supply/AC/online"
)

// UpdateBlock updates the battery status block.
func (c Battery) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	r, err := os.Open(acFilePath)
	if err != nil {
		b.Urgent = true
		b.FullText = err.Error()
		return
	}
	defer r.Close()
	var online int
	_, err = fmt.Fscanf(r, "%d", &online)
	if err != nil {
		b.Urgent = true
		b.FullText = err.Error()
		return
	}
	if online == 1 {
		c.Label = "🔌 "
	} else {
		c.Label = "🔋 "
	}

	fullTextFmt := fmt.Sprintf("%s%%d%%%%", c.Label)
	var capacity int
	batFilePath := fmt.Sprintf("/sys/class/power_supply/BAT%d/capacity", c.BatteryNumber)
	r, err = os.Open(batFilePath)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	defer r.Close()
	_, err = fmt.Fscanf(r, "%d", &capacity)
	if err != nil {
		b.Urgent = true
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	if float64(capacity) >= c.CritBattery {
		b.Urgent = false
	} else {
		b.Urgent = true
	}
	b.FullText = fmt.Sprintf(fullTextFmt, capacity)
}
