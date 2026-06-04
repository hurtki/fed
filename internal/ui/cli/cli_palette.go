package cli_ui

import "github.com/fatih/color"

type CLIPalette struct {
	Status          *color.Color
	Log             *color.Color
	Plan            *color.Color
	ResultSuccess   *color.Color
	ResultFailure   *color.Color
	Approve         *color.Color
	RequestMulLines *color.Color
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
	}
}
