package models

import (
	"fmt"

	gojson "github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
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
	agent := fiber.Get(source)

	status, body, errs := agent.Bytes()
	if errs != nil {
		logrus.Errorf("Err request to source - %s", errs[0])
		return
	}

	if status != fiber.StatusOK {
		logrus.Warnf("Problems on the resource side - %d", status)
		return
	}

	var s SourceData
	if err := gojson.Unmarshal(body, &s); err != nil {
		logrus.Errorf("Err unmarshal source data to struct - %s", err)
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

func SendingData(service string, token string) bool {
	agent := fiber.AcquireAgent()
	req := agent.Request()

	req.Header.SetMethod(fiber.MethodPost)
	req.SetRequestURI(service)

	req.Header.SetMethod(fiber.MethodPost)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	agent.JSON(data)

	status, _, errs := agent.Bytes()
	if errs != nil {
		logrus.Errorf("Err request to service - %s", errs[0])
		return false
	}

	if status != fiber.StatusOK {
		logrus.Warnf("Something went wrong - %d", status)
		return false
	}
	return true
}
