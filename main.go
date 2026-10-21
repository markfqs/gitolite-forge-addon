package main

import (
	"go.woodpecker-ci.org/woodpecker/v3/server/forge/addon"
)

func main() {
	addon.Serve(&Generic{})
}
