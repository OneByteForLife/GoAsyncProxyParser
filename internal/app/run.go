package app

import (
	"GoProxyService/internal/config"
	"GoProxyService/internal/models"
	"fmt"
	"strconv"
)

func Run() {
	conf := config.ReadConfig()
	for i := 1; i <= 15; i++ {
		models.CollectingProxies(fmt.Sprintf("https://proxylist.geonode.com/api/proxy-list?limit=500&page=%s", strconv.Itoa(i)))
	}

	models.SendingData(conf.OutServiceURL, conf.JwtToken)
}
