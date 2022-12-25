package app

import (
	"GoProxyService/internal/config"
	"GoProxyService/internal/models"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

func Run() {
	conf := config.ReadConfig()
	for i := 1; i <= 2; i++ {
		models.CollectingProxies(fmt.Sprintf("https://proxylist.geonode.com/api/proxy-list?limit=500&page=%s", strconv.Itoa(i)))
	}

	if status := models.SendingData(conf.OutServiceURL, conf.JwtToken); !status {
		logrus.Errorf("Err sending data to %s", conf.OutServiceURL)
	}
}
