package handler

import (
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"io/ioutil"
	"net/http"
	"gopkg.in/ini.v1"
)

func HandleRequest(c echo.Context) error {
	cfg, _ := ini.Load("config.ini")
	lat := c.QueryParam("lat")
	lng := c.QueryParam("lng")
	apiKey := cfg.Section("api").Key("key").String()
	url := fmt.Sprintf("http://webservice.recruit.co.jp/hotpepper/gourmet/v1/?key=%v&lat=%v&lng=%v&range=3&order=4&format=json&keyword=ラーメン",
		apiKey, lat, lng)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}

	if len(response.Results.Shop) <= 0 {
		return c.JSON(http.StatusOK, "近くにお店がないよ")
	}

	return c.JSON(http.StatusOK, response.Results)
}

type Response struct {
	Results struct {
		Shop []struct {
			Id          string `json:"id"`
			Name        string `json:"name"`
			LogoImage   string `json:"logo_image"`
			NameKana    string `json:"name_kana"`
			StationName string `json:"station_name"`
			Access      string `json:"access"`
			Address     string `json:"address"`
			Urls        struct {
				Pc string `json:"pc"`
			} `json:"urls"`
			Genre struct {
				Name  string `json:"name"`
				Catch string `json:"catch"`
			} `json:"genre"`
		} `json:"shop"`
	} `json:"results"`
}
