package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/citizensciencecenter/autodeploy/modules"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func loadConfig() {
	viper.SetConfigType("json")
	viper.SetConfigFile("./config/conf.json")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	fmt.Println(viper.Get("repo_dir"))
}

func runHookServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", hookHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":9898", r))
}

func hookHandler(w http.ResponseWriter, r *http.Request) {
	// TODO check token is matching
	body, err := ioutil.ReadAll(r.Body)
	modules.ErrHandler(err)
	hook := modules.TravisResp{}
	json.Unmarshal(body, &hook)
	w.Header().Set("Content-Type", "application/json")

	src := fmt.Sprintf("%s:%s", hook.Repository.Name, hook.Branch)
	rc := modules.HookBody{Source: src, Status: "SUCCESS", Stage: "Hook Triggered", Msg: "Hook started"}
	ad := modules.AutoDeploy{Config: viper.GetViper(), HookBody: rc, Travis: hook}
	if strings.Compare(hook.State, "passed") == 0 {
		modules.Notify(ad)
		w.WriteHeader(200)
		w.Write([]byte("{data: Hook started}"))
		go modules.InitRepo(hook.Repository.Name, hook.Branch, ad)
	} else {
		ad.HookBody.Stage = "Travis"
		ad.HookBody.Status = "FAILED"
		ad.HookBody.Msg = "Tests Failed"
		modules.Notify(ad)
		w.WriteHeader(500)
		w.Write([]byte("{data: Hook Failed}"))
	}
}

func main() {
	loadConfig()
	runHookServer()
}
