package app

import (
	"GoAsyncProxyParser/internal/config"
	"GoAsyncProxyParser/internal/models"
	"sync"
)

func Run() {
	conf := config.ReadConfig()
	var wg sync.WaitGroup

	wg.Add(5)
	go func() {
		models.CollectingProxies("https://proxylist.geonode.com/api/proxy-list?limit=500&page=1&sort_by=lastChecked")
		wg.Done()
	}()

	go func() {
		models.CollectingProxies("https://proxylist.geonode.com/api/proxy-list?limit=500&page=2&sort_by=lastChecked")
		defer wg.Done()
	}()

	go func() {
		models.CollectingProxies("https://proxylist.geonode.com/api/proxy-list?limit=500&page=3&sort_by=lastChecked")
		defer wg.Done()
	}()
	go func() {
		models.CollectingProxies("https://proxylist.geonode.com/api/proxy-list?limit=500&page=4&sort_by=lastChecked")
		defer wg.Done()
	}()

	go func() {
		models.CollectingProxies("https://proxylist.geonode.com/api/proxy-list?limit=500&page=5&sort_by=lastChecked")
		defer wg.Done()
	}()
	wg.Wait()

	models.SendingData(conf.OutServiceURL, conf.JwtToken)
}
