package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Webhook struct {
	Destination string           `json:"destination"`
	Events      []*linebot.Event `json:"events"`
}

func HandleRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	bot, err := linebot.New(
		os.Getenv("secret"),
		os.Getenv("token"),
	)
	if err != nil {
		log.Fatal(err)
	}
	var webhook Webhook
	if err := json.Unmarshal([]byte(request.Body), &webhook); err != nil {
		log.Print(err)
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusBadRequest,
			Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", http.StatusText(http.StatusBadRequest)),
		}, nil
	}
	for _, event := range webhook.Events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.LocationMessage:
				lat := message.Latitude
				lng := message.Longitude
				shops, err := getShops(lat, lng)
				if err != nil {
					return events.APIGatewayProxyResponse{
						StatusCode: http.StatusInternalServerError,
						Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", http.StatusText(http.StatusInternalServerError)),
					}, err
				}
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(*shops)).Do(); err != nil {
					return events.APIGatewayProxyResponse{
						StatusCode: http.StatusInternalServerError,
						Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", http.StatusText(http.StatusInternalServerError)),
					}, err
				}
			default:
				m := "入力フォーム左の+ボタン　→　位置情報　から位置情報を送ってくれよなっ😎"
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(m)).Do(); err != nil {
					return events.APIGatewayProxyResponse{
						StatusCode: http.StatusInternalServerError,
						Body:       fmt.Sprintf(`{"message":"%s"}`+"\n", http.StatusText(http.StatusInternalServerError)),
					}, err
				}
			}
		}
	}
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func getShops(lat, lng float64) (*string, error) {
	apiKey := os.Getenv("key")
	url := fmt.Sprintf("http://webservice.recruit.co.jp/hotpepper/gourmet/v1/?key=%v&lat=%v&lng=%v&range=3&order=4&format=json&keyword=ラーメン",
		apiKey, lat, lng)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	var response struct {
		Results struct {
			Shop []struct {
				Name string `json:"name"`
				Urls struct {
					Pc string `json:"pc"`
				} `json:"urls"`
			} `json:"shop"`
		} `json:"results"`
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	var shops string
	if len(response.Results.Shop) == 0 {
		msg := "近くにラーメン屋がないぜ"
		return &msg, nil
	}
	for _, shop := range response.Results.Shop {
		info := shop.Name + "\n" + shop.Urls.Pc + "\n"
		shops += info + "\n"
	}
	return &shops, nil
}

func main() {
	lambda.Start(HandleRequest)
}
