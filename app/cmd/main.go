package main

import (
	"GoAsyncProxyParser/internal/app"
	"GoAsyncProxyParser/pkg"
	"time"
)

func init() {
	pkg.ConfigLog()
}

func main() {
	for {
		timer1 := time.NewTimer(3 * time.Hour)
		<-timer1.C
		app.Run()
	}
}
