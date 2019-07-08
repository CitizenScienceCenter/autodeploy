package modules

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func notify() {
	req, err := http.NewRequest("POST", "", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		//handle response error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		//handle read response error
	}

	fmt.Printf("%s\n", string(body))
}
