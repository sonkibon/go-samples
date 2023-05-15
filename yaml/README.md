# About The Project
This will help you understand the Marshal and Unmarshal function specifications of the "gopkg.in/yaml.v3" package with many examples.

# Getting Started

1. Go to the root of the repository
```
cd <root of repository>
```
2. Run main.go
```
go run yaml/main.go
```

# Example
```
Result of marshaling a struct that field is not exported
{}

Result of marshaling struct without field tag
b: true
i: 10
s: hoge
intslice:
    - 1
    - 2
    - 3
stringslice:
    - aaa
    - bbb
    - ccc

Result of marshaling struct with field tag
bool: true
int: 10
string: hoge
int_array:
    - 1
    - 2
    - 3
string_array:
    - aaa
    - bbb
    - ccc

Result of marshaling struct without omitempty flag
bool: false
int: 0
string: ""
int_array: []
string_array: []

Result of marshaling struct with omitempty flag
{}

Result of unmarshal of yaml document to unmarshalUnexported struct
main.unmarshalUnexported{boolean:false, integer:0, str:"", array:[]string(nil)}
Result of unmarshal of yaml document to unmarshalNoFieldTag struct
main.unmarshalNoFieldTag{Boolean:true, Integer:1, Str:"hoge", Array:[]string{"fizz", "buzz"}}
Result of unmarshal of yaml document to unmarshalFieldTagNoFieldTag struct
main.unmarshalFieldTag{Boolean:true, Integer:1, Str:"hoge", Array:[]string{"fizz", "buzz"}}
Result of unmarshal of yaml document to unmarshalNoFieldTagWithDifferentFieldName struct
main.unmarshalNoFieldTagWithDifferentFieldName{Hoge:false, Fuga:0, Bar:"", Foo:[]string(nil)}
```
