package models

import (
	"bytes"
	"fmt"
	"io"
	"net/http"

	gojson "github.com/goccy/go-json"
	"github.com/sirupsen/logrus"
)

/*
	Алгоритм прост
	1 - Собираем прокси
	2 - Проверяем их на валидность
	3 - Отправляем в другой сервис
*/

type SourceData struct {
	Data []struct {
		IP             string   `json:"ip"`
		AnonymityLevel string   `json:"anonymityLevel"`
		City           string   `json:"city"`
		Country        string   `json:"country"`
		Port           string   `json:"port"`
		Protocols      []string `json:"protocols"`
		Speed          int      `json:"speed"`
	} `json:"data"`
}

type ModedData struct {
	Types []string `json:"protocols"`
	Data  struct {
		IP      string `json:"ip"`
		Port    string `json:"port"`
		Speed   int    `json:"speed"`
		AnonLvL string `json:"anon_lvl"`
		Geo     struct {
			City    string `json:"city"`
			Country string `json:"country"`
		} `json:"geo"`
	} `json:"data"`
}

var data []ModedData

// Сбор прокси
func CollectingProxies(source string) {
	resp, err := http.Get(source)
	if err != nil {
		logrus.Errorf("Err request to source - %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Warnf("Problems on the resource side - %d", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Err read body - %s", err)
		return
	}

	var s SourceData
	if err := gojson.Unmarshal(body, &s); err != nil {
		logrus.Errorf("Err unmarshal source data to struct - %s", err)
		return
	}

	for _, val := range s.Data {
		data = append(data, ModedData{
			Types: val.Protocols,
			Data: struct {
				IP      string "json:\"ip\""
				Port    string "json:\"port\""
				Speed   int    "json:\"speed\""
				AnonLvL string "json:\"anon_lvl\""
				Geo     struct {
					City    string "json:\"city\""
					Country string "json:\"country\""
				} "json:\"geo\""
			}{
				IP:      val.IP,
				Port:    val.Port,
				Speed:   val.Speed,
				AnonLvL: val.AnonymityLevel,
				Geo: struct {
					City    string "json:\"city\""
					Country string "json:\"country\""
				}{
					City:    val.City,
					Country: val.Country,
				},
			},
		})
	}
}

func SendingData(service string, token string) {
	client := &http.Client{}

	out, err := gojson.Marshal(&data)
	if err != nil {
		logrus.Errorf("Err marshal data from payload - %s", err)
		return
	}

	req, err := http.NewRequest(http.MethodPost, service, bytes.NewReader(out))
	if err != nil {
		logrus.Errorf("Err generation request from proxy service - %s", err)
		return
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Err sending data to proxy service - %s", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logrus.Warnf("Proxy service returning status - %d", resp.StatusCode)
	}

	if resp.StatusCode == http.StatusUnauthorized {
		logrus.Warnf("Check jwt token - %d", resp.StatusCode)
	}

	logrus.Info("Success sending data")
}
