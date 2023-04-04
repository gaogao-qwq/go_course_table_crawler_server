package configs

import (
	"encoding/json"
	"os"
)

type Url struct {
	LoginUrl string `json:"loginUrl"`
	HomeUrl  string `json:"homeUrl"`
}

var URL Url

func init() {
	f, err := os.ReadFile("configs/url.json")
	if err != nil {
		panic(err)
		return
	}

	err = json.Unmarshal(f, &URL)
	if err != nil {
		panic(err)
		return
	}
}
