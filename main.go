package main

import (
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	v := &validate{}
	protogen.Options{}.Run(v.Generate)
}
