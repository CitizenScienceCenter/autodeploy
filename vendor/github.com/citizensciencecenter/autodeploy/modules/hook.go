package modules

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/viper"
)

// HookBody webhook payload for the webhook
type HookBody struct {
	Source string
	Status string
	Stage  string
	Msg    string
}

// AutoDeploy to hold the loaded config and the hook data that is updated at each stage
type AutoDeploy struct {
	Config   *viper.Viper
	HookBody HookBody
	Travis   TravisResp
	Hash     string
	Dir string
}

// Notify sends hook details to your endpoint of choice
func Notify(ad AutoDeploy) {
	body, err := json.Marshal(ad.HookBody)
	fmt.Println(string(body))
	ErrHandler(err)
	url := ad.Config.GetString("webhook.url")
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	ErrHandler(err)
	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	ErrHandler(err)
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	ErrHandler(err)
	fmt.Printf("%s\n", string(body))
}
