package main

import (
	"encoding/json"
	"fmt"
)

type Parser interface {
	Parse() (GroupedStratagems, error)
}

const (
	BaseUrl  = "https://helldivers.fandom.com/wiki/Stratagem_Codes_(Helldivers_2)"
	Filename = "stratagems.json"
	Success  = "✅ Json file has been created successfully!"
)

func main() {
	helldiversParser := NewParser(BaseUrl)
	stratagems, err := runParser(helldiversParser)
	if err != nil {
		panic(err)
	}

	b, err := json.Marshal(stratagems)
	if err != nil {
		panic(err)
	}

	err = writeDataToFile(b, Filename)
	if err != nil {
		panic(err)
	}

	fmt.Println(Success)
}

func runParser(p Parser) (GroupedStratagems, error) {
	return p.Parse()
}
