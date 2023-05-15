package main

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

// marshal is a function that wraps the Marshal function from "gopkg.in/yaml.v3
func marshal(in interface{}) ([]byte, error) {
	out, err := yaml.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("yaml.Marshal: %w", err)
	}

	return out, nil
}

func main() {}
