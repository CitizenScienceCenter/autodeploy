package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/citizensciencecenter/autodeploy/modules"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func runHookServer() {
	r := mux.NewRouter()
	seed := viper.GetInt64("seed")
	rand.Seed(seed)
	routeID := rand.Uint64()
	route := fmt.Sprintf("/%d", routeID)
	fmt.Println("Your URL is: " + route)
	r.HandleFunc("/", upHandler).Methods("GET")
	r.HandleFunc(route, hookHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":9898", r))
}

func upHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Autodeploy!\n"))
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
	ad.Dir = fmt.Sprintf("%s/%s", ad.Config.GetString("git.repo_dir"), hook.Repository.Name)
	fmt.Println(ad.Dir)
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
	modules.LoadConfig("./config/conf.json")
	runHookServer()
}
