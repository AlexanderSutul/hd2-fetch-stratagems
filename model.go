package main

type InputCode = string

const (
	UP    InputCode = "up"
	DOWN  InputCode = "down"
	LEFT  InputCode = "left"
	RIGHT InputCode = "right"
)

type Stratagem struct {
	Group      string      `json:"group"`
	Name       string      `json:"name"`
	InputCode  []InputCode `json:"input_code"`
	Cooldown   string      `json:"cooldown"`
	Uses       string      `json:"uses"`
	Activation string      `json:"activation"`
	IconUrl    string      `json:"icon_url"`
}

type GroupedStratagems map[string][]*Stratagem
