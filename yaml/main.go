package main

import (
	"fmt"
	"log"
	"os"

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

// unmarshalUnexported is a struct to unmarshal YAML documents that do not have exported struct fields.
type unmarshalUnexported struct {
	boolean bool     `yaml:"boolean"`
	integer int      `yaml:"integer"`
	str     string   `yaml:"str"`
	array   []string `yaml:"array"`
}

// unmarshalNoFieldTag is a struct to unmarshal YAML documents with exported struct fields.
type unmarshalNoFieldTag struct {
	Boolean bool
	Integer int
	Str     string
	Array   []string
}

// unmarshalFieldTag is a struct for exporting and unmarshaling into a YAML document with a struct fields with a "yaml" field tag.
type unmarshalFieldTag struct {
	Boolean bool     `yaml:"boolean"`
	Integer int      `yaml:"integer"`
	Str     string   `yaml:"str"`
	Array   []string `yaml:"array"`
}

// unmarshalFieldTagNoFieldTag is a struct to unmarshal YAML documents with exported struct fields with different field name.
type unmarshalNoFieldTagWithDifferentFieldName struct {
	Hoge bool
	Fuga int
	Bar  string
	Foo  []string
}

// marshal is a function that wraps the Marshal function from "gopkg.in/yaml.v3
func marshal(in interface{}) ([]byte, error) {
	out, err := yaml.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("yaml.Marshal: %w", err)
	}

	return out, nil
}

// unmarshal is a function that wraps the Unmarshal function from "gopkg.in/yaml.v3
func unmarshal(in []byte, out interface{}) error {
	if err := yaml.Unmarshal(in, out); err != nil {
		return fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	return nil
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

	const (
		fileName = "yaml/sample.yaml"
	)

	b, err := os.ReadFile(fileName)
	if err != nil {
		log.Fatalf("os.ReadFile: %v", err)
	}

	var (
		uu   unmarshalUnexported
		unft unmarshalNoFieldTag
		uft  unmarshalFieldTag
		udfn unmarshalNoFieldTagWithDifferentFieldName
	)

	if err := unmarshal(b, &uu); err != nil {
		log.Fatalf("decode: %v", err)
	}
	fmt.Printf("Result of unmarshal of yaml document to unmarshalUnexported struct\n%#v\n", uu)

	if err := unmarshal(b, &unft); err != nil {
		log.Fatalf("decode: %v", err)
	}
	fmt.Printf("Result of unmarshal of yaml document to unmarshalNoFieldTag struct\n%#v\n", unft)

	if err := unmarshal(b, &uft); err != nil {
		log.Fatalf("decode: %v", err)
	}
	fmt.Printf("Result of unmarshal of yaml document to unmarshalFieldTagNoFieldTag struct\n%#v\n", uft)

	if err := unmarshal(b, &udfn); err != nil {
		log.Fatalf("decode: %v", err)
	}
	fmt.Printf("Result of unmarshal of yaml document to unmarshalNoFieldTagWithDifferentFieldName struct\n%#v\n", udfn)
}
