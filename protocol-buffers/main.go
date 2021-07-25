package main

import (
	"fmt"
	"io/ioutil"

	"google.golang.org/protobuf/proto"
)

func write(path string) {
	d := &Dog{Name: "sophia", Age: 7}
	b, err := proto.Marshal(d)
	if err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path, b, 0644); err != nil {
		panic(err)
	}
}

func read(path string) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	d := &Dog{}
	err = proto.Unmarshal(b, d)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", d)
}

func main() {
	path := "dog"
	write(path)
	read(path)
}
