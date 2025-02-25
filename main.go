package main

import (
	"fmt"
	"github.com/webitel/webitel-fts/cmd"
)

//go:generate go run github.com/bufbuild/buf/cmd/buf@latest generate --template buf.gen.yaml
//go:generate go run github.com/bufbuild/buf/cmd/buf@latest generate --template buf.gen.webitel.yaml
//go:generate go mod tidy
//go:generate go run github.com/google/wire/cmd/wire@latest gen ./cmd

func main() {
	if err := cmd.Run(); err != nil {
		fmt.Println(err.Error())
		return
	}
}
