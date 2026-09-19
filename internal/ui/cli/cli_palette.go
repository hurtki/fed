package cli_ui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

type CLIPalette struct {
	Status          *color.Color
	Log             *color.Color
	Plan            *color.Color
	ResultSuccess   *color.Color
	ResultFailure   *color.Color
	Approve         *color.Color
	RequestMulLines *color.Color
	DiffRemoved     *color.Color
	DiffAdded       *color.Color
}

func NewChillCLIPalette() CLIPalette {
	return CLIPalette{
		Status: color.RGB(71, 120, 144), // cerulean
		Log:    color.RGB(76, 134, 168), // air-force-blue
		Plan:   color.RGB(76, 134, 168), // air-force-blue

		ResultSuccess: color.RGB(225, 221, 143), // light-gold
		Approve:       color.RGB(225, 221, 143), // light-gold

		ResultFailure: color.RGB(142, 59, 70), // burnt-rose

		RequestMulLines: color.RGB(224, 119, 125), // light-coral
		DiffRemoved:     color.RGB(200, 70, 70),   // soft red
		DiffAdded:       color.RGB(80, 160, 100),  // soft green
	}
}

func NewWarmCLIPalette() CLIPalette {
	return CLIPalette{
		Status: color.RGB(255, 140, 0), // dark orange
		Log:    color.RGB(255, 165, 0), // orange
		Plan:   color.RGB(255, 165, 0), // orange

		ResultSuccess: color.RGB(255, 215, 0), // gold
		Approve:       color.RGB(255, 215, 0), // gold

		ResultFailure: color.RGB(178, 34, 34), // firebrick

		RequestMulLines: color.RGB(255, 99, 71), // tomato
		DiffRemoved:     color.RGB(220, 20, 60),  // crimson
		DiffAdded:       color.RGB(50, 205, 50),  // lime green
	}
}

func NewNeonCLIPalette() CLIPalette {
	return CLIPalette{
		Status: color.RGB(0, 255, 255), // cyan
		Log:    color.RGB(255, 0, 255), // magenta
		Plan:   color.RGB(255, 0, 255), // magenta

		ResultSuccess: color.RGB(57, 255, 20), // neon green
		Approve:       color.RGB(255, 255, 0), // yellow

		ResultFailure: color.RGB(255, 7, 58), // neon red

		RequestMulLines: color.RGB(255, 20, 147), // deep pink
		DiffRemoved:     color.RGB(255, 0, 0),    // neon red
		DiffAdded:       color.RGB(0, 255, 0),    // neon green
	}
}

func NewMonochromeCLIPalette() CLIPalette {
	return CLIPalette{
		Status: color.RGB(169, 169, 169), // dark gray
		Log:    color.RGB(211, 211, 211), // light gray
		Plan:   color.RGB(211, 211, 211), // light gray

		ResultSuccess: color.RGB(255, 255, 255), // white
		Approve:       color.RGB(255, 255, 255), // white

		ResultFailure: color.RGB(105, 105, 105), // dim gray

		RequestMulLines: color.RGB(128, 128, 128), // gray
		DiffRemoved:     color.RGB(160, 160, 160), // gray
		DiffAdded:       color.RGB(255, 255, 255), // white
	}
}

func (p CLIPalette) FormatDiff(path, find, replace string) string {
	var b strings.Builder
	fmt.Fprintf(&b, ">>>>>> %s\n", path)
	for _, line := range strings.Split(find, "\n") {
		if p.DiffRemoved != nil {
			b.WriteString(p.DiffRemoved.Sprintf("- %s\n", line))
		} else {
			fmt.Fprintf(&b, "- %s\n", line)
		}
	}
	fmt.Fprintf(&b, "====== %s\n", path)
	for _, line := range strings.Split(replace, "\n") {
		if p.DiffAdded != nil {
			b.WriteString(p.DiffAdded.Sprintf("+ %s\n", line))
		} else {
			fmt.Fprintf(&b, "+ %s\n", line)
		}
	}
	return b.String()
}
