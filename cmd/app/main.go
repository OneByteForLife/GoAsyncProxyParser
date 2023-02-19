package main

import (
	"GoAsyncProxyParser/internal/app"
	"GoAsyncProxyParser/pkg"
)

func init() {
	pkg.ConfigLog()
}

func main() {
	app.Run()
}
