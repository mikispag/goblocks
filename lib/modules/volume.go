package modules

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/davidscholberg/go-i3barjson"
)

// Volume represents the configuration for the volume display block.
type Volume struct {
	BlockConfigBase `yaml:",inline"`
	MixerDevice     string `yaml:"mixer_device"`
	Channel         string `yaml:"channel"`
}

// UpdateBlock updates the volume display block.
// Currently, only the ALSA master channel volume is supported.
func (c Volume) UpdateBlock(b *i3barjson.Block) {
	b.Color = c.Color
	amixerCmd := "amixer"
	if c.MixerDevice == "" {
		c.MixerDevice = "default"
	}
	if c.Channel == "" {
		c.Channel = "Master"
	}
	amixerArgs := []string{"-D", c.MixerDevice, "get", c.Channel}
	out, err := exec.Command(amixerCmd, amixerArgs...).Output()
	if err != nil {
		b.FullText = err.Error()
		return
	}
	outStr := string(out)
	iBegin := strings.Index(outStr, "[")
	if iBegin == -1 {
		b.FullText = "cannot parse amixer output"
		return
	}
	iEnd := strings.Index(outStr, "]")
	if iEnd == -1 {
		b.FullText = "cannot parse amixer output"
		return
	}

	v, err := strconv.Atoi(outStr[iBegin+1 : iEnd-1])
	if err != nil {
		b.FullText = "cannot parse amixer output"
		return
	}

	// Update label depending on the volume label.
	c.Label = "🔊 "
	if v < 20 {
		c.Label = "🔈 "
	} else if v < 75 {
		c.Label = "🔉 "
	}

	// If the device is "off", gray out the indicator and update the label.
	if len(outStr) >= iEnd+5 && outStr[iEnd+4] == 'f' {
		c.Label = "🔇 "
		b.Color = "#9e9e9e"
	}

	b.FullText = fmt.Sprintf("%s%d%%",c.Label, v)
}
