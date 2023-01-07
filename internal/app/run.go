package app

import (
	"GoAsyncProxyParser/internal/config"
	"GoAsyncProxyParser/internal/models"
	"sync"
)

func Run() {
	conf := config.ReadConfig()
	var wg sync.WaitGroup

	wg.Add(3)
	go func() {
		models.CollectingProxies("https://proxylist.geonode.com/api/proxy-list?limit=500&page=1")
		wg.Done()
	}()

	go func() {
		models.CollectingProxies("https://proxylist.geonode.com/api/proxy-list?limit=500&page=2")
		wg.Done()
	}()

	go func() {
		models.CollectingProxies("https://proxylist.geonode.com/api/proxy-list?limit=500&page=3")
		wg.Done()
	}()
	wg.Wait()
	models.SendingData(conf.OutServiceURL, conf.JwtToken)
}
