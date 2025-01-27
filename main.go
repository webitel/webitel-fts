package main

import (
	"github.com/webitel/webitel-fts/cmd"
)

//go:generate go run github.com/bufbuild/buf/cmd/buf@latest generate --template buf.gen.yaml
//go:generate go run github.com/bufbuild/buf/cmd/buf@latest generate --template buf.gen.webitel.yaml
//go:generate go run github.com/google/wire/cmd/wire@latest gen ./cmd

func main() {
	if err := cmd.Run(); err != nil {
		return
	}
}
