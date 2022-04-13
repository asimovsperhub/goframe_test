package main

import (
	_ "firstproject/internal/packed"

	"github.com/gogf/gf/v2/os/gctx"

	"firstproject/internal/cmd"
)

func main() {
	cmd.Main.Run(gctx.New())
}
