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

// TravisResp webhook payload coming from Travis CI when the testing has been completed
type TravisResp struct {
	Repository Repo
	Branch     string
	State      string
	Commit     string
	BuildURL   string
	CompareURL string
	Number     int
}

// Repo contains the ownership info of the repository coming from Travis
type Repo struct {
	Name      string
	OwnerName string
}

// RocketChat webhook payload for the webhook
type RocketChat struct {
	source string
	status string
	stage  string
	msg    string
}

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

func sendHook(msg string, stage string) {
	// hookMsg := RocketChat{msg, stage, "", ""}
}

func hookHandler(w http.ResponseWriter, r *http.Request) {
	// TODO check token is matching
	body, err := ioutil.ReadAll(r.Body)
	modules.ErrHandler(err)
	hook := TravisResp{}
	json.Unmarshal(body, &hook)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("{data: Hook started}"))
	//return
	// TODO start a channel here to perform the build
	if strings.Compare(hook.State, "passed") == 0 {
		fmt.Println(hook.Repository.Name)
		go modules.InitRepo(hook.Repository.Name, hook.Branch)
	} else {
		log.Fatal("Tests were not successful, exiting")
	}
	// TODO pull from repo and specified branch
}

func envCreate(t TravisResp, hash string) {
	// TODO create temp env file based on travis reponse
	// i.e. develop branch is the staging namespace
	// NAME = repo name
	// NS = branch
	// HOST = NAME + NS (unless NS is prod)
	// TAG = branch + git hash
	// PORT = how to define? Default port? Read from Dockerfile?
}

func main() {
	loadConfig()
	runHookServer()
}
