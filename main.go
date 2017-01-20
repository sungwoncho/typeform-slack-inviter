package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type payload struct {
	responses []response
}

type response struct {
	answers map[string]string
}

var apiKeyFlag = flag.String("apiKey", "", "TypeForm API key")
var formUIDFlag = flag.String("formUID", "", "TypeForm form UID")
var slackAPITokenFlag = flag.String("slackAPIToken", "", "SlackAPIToken")
var intervalFlag = flag.Int("interval", 1, "Interval for checking TypeForm responses (in minutes)")

func fetchTypeformData(apiKey string, formUID string, since time.Time) (payload, error) {
	sinceTimestamp := since.Unix()
	endpoint := fmt.Sprintf("https://api.typeform.com/v1/form/%s?key=%s&since=%d", formUID, apiKey, sinceTimestamp)
	var ret payload
	fmt.Println(endpoint)

	res, err := http.Get(endpoint)
	if err != nil {
		fmt.Println("err", err)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ret, err
	}

	var data payload
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("error unmarshalling")
		return ret, err
	}

	return data, nil
}

func sendSlackInvitation(email string) error {
	fmt.Println("Inviting", email)
	endpoint := fmt.Sprintf("https://slack.com/api/users.admin.invite?token=%s&email=%s&resend=false", *slackAPITokenFlag, email)

	req, err := http.NewRequest("POST", endpoint, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response body:", string(body))
	return nil
}

func run() {
	now := time.Now().UTC()
	targetTime := now.Add(-time.Duration(*intervalFlag) * time.Minute)

	data, err := fetchTypeformData(*apiKeyFlag, *formUIDFlag, targetTime)
	if err != nil {
		panic(err)
	}

	fmt.Println("Need to invite", len(data.responses))

	for _, response := range data.responses {
		email := response.answers["email_40622900"]
		if email != "" {
			sendSlackInvitation(email)
			fmt.Println("Inviting", email)
		}
	}
}

func main() {
	flag.Parse()

	t := time.NewTicker(time.Duration(*intervalFlag) * time.Minute)

	for {
		run()
		<-t.C
	}
}
