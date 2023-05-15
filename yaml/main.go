package main

import (
	"fmt"
	"log"

	"gopkg.in/yaml.v3"
)

// marshalUnexported is a struct to marshal YAML documents that do not have exported struct fields.
type marshalUnexported struct {
	b           bool
	i           int
	s           string
	intSlice    []int
	stringSlice []string
}

// marshal is a function that wraps the Marshal function from "gopkg.in/yaml.v3
func marshal(in interface{}) ([]byte, error) {
	out, err := yaml.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("yaml.Marshal: %w", err)
	}

	return out, nil
}

func main() {
	mu := marshalUnexported{
		b: true,
		i: 10,
		s: "hoge",
		intSlice: []int{
			1,
			2,
			3,
		},
		stringSlice: []string{
			"aaa",
			"bbb",
			"ccc",
		},
	}

	muOut, err := marshal(&mu)
	if err != nil {
		log.Fatalf("encode: %v", err)
	}
	fmt.Printf("Result of marshaling a struct that field is not exported\n%s\n", string(muOut))
}
