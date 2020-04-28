package modules

import (
	"fmt"
	"os/exec"
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
	fullTextFmt := fmt.Sprintf("%s%%s", c.Label)
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
		b.FullText = fmt.Sprintf(fullTextFmt, err.Error())
		return
	}
	outStr := string(out)
	iBegin := strings.Index(outStr, "[")
	if iBegin == -1 {
		b.FullText = fmt.Sprintf(fullTextFmt, "cannot parse amixer output")
		return
	}
	iEnd := strings.Index(outStr, "]")
	if iEnd == -1 {
		b.FullText = fmt.Sprintf(fullTextFmt, "cannot parse amixer output")
		return
	}
	v := outStr[iBegin+1 : iEnd]

	// If the device is "off", gray out the indicator.
	if len(outStr) >= iEnd+4 && outStr[iEnd+3] == 'f' {
		b.Color = "#222222"
	}

	b.FullText = fmt.Sprintf(fullTextFmt, v)
}
