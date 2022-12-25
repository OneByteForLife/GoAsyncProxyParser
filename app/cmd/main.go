package main

import (
	"GoProxyService/internal/app"
	"GoProxyService/pkg"
)

func init() {
	pkg.ConfigLog()
}

func main() {
	app.Run()
}
