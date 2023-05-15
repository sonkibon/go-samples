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

// marshalNoFieldTag is a struct to marshal YAML documents with exported struct fields.
type marshalNoFieldTag struct {
	B           bool
	I           int
	S           string
	IntSlice    []int
	StringSlice []string
}

// marshalFieldTag is a struct for exporting and marshaling into a YAML document with a struct fields with a "yaml" field tag.
type marshalFieldTag struct {
	B           bool     `yaml:"bool"`
	I           int      `yaml:"int"`
	S           string   `yaml:"string"`
	IntSlice    []int    `yaml:"int_array"`
	StringSlice []string `yaml:"string_array"`
}

// marshalOmitempty is a struct to marshal a YAML document with struct fields that have a "yaml" field tag that has been exported and has the "omitempty" flag.
type marshalOmitempty struct {
	B           bool     `yaml:"bool,omitempty"`
	I           int      `yaml:"int,omitempty"`
	S           string   `yaml:"string,omitempty"`
	IntSlice    []int    `yaml:"int_array,omitempty"`
	StringSlice []string `yaml:"string_array,omitempty"`
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

	mnft := marshalNoFieldTag{
		B: true,
		I: 10,
		S: "hoge",
		IntSlice: []int{
			1,
			2,
			3,
		},
		StringSlice: []string{
			"aaa",
			"bbb",
			"ccc",
		},
	}

	mnftOut, err := marshal(&mnft)
	if err != nil {
		log.Fatalf("encode: %v", err)
	}
	fmt.Printf("Result of marshaling struct without field tag\n%s\n", string(mnftOut))

	mft := marshalFieldTag{
		B: true,
		I: 10,
		S: "hoge",
		IntSlice: []int{
			1,
			2,
			3,
		},
		StringSlice: []string{
			"aaa",
			"bbb",
			"ccc",
		},
	}

	mftOut, err := marshal(&mft)
	if err != nil {
		log.Fatalf("encode: %v", err)
	}
	fmt.Printf("Result of marshaling struct with field tag\n%s\n", string(mftOut))

	meft := marshalFieldTag{}

	meftOut, err := marshal(&meft)
	if err != nil {
		log.Fatalf("encode: %v", err)
	}
	fmt.Printf("Result of marshaling struct without omitempty flag\n%s\n", string(meftOut))

	moe := marshalOmitempty{}
	moeOut, err := marshal(&moe)
	if err != nil {
		log.Fatalf("encode: %v", err)
	}
	fmt.Printf("Result of marshaling struct with omitempty flag\n%s\n", string(moeOut))
}
