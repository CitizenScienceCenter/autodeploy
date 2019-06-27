package main

import (
	"fmt"
	"log"
	"strings"
	"io/ioutil"
	"net/http"
	"os/exec"
	"encoding/json"

	"github.com/spf13/viper"
	"github.com/gorilla/mux"
	"github.comsrc-d/go-git"

)

type TravisResp struct {
	Repository Repo
	Branch string
	State string
	Commit string
	BuildUrl string
	CompareUrl string
	Number int
}

type Repo struct {
	Name string
	OwnerName string
}

func loadConfig(path string, name string) {
	viper.SetConfigName(name)
	viper.AddConfigPath(path)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
}

func runHookServer() {
	r := mux.NewRouter()
	r.HandleFunc("/", HookHandler).Methods("POST")
	log.Fatal(http.ListenAndServe(":9898", r))
}

func HookHandler(w http.ResponseWriter, r *http.Request) {
	// TODO check token is matching
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
        log.Fatal(err)
	}
	hook := TravisResp{}
	json.Unmarshal(body, &hook)
	if (strings.Compare(hook.State, "passed") == 0) {
		fmt.Println(hook.Repository.Name)
	} else {
		log.Fatal("Tests were not successful, exiting")
	}
	// TODO pull from repo and specified branch
}

func runCommand(cmd string) {
	out, err := exec.Command("gulp", "serv.dev").Output()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("The date is %s\n", out)
}

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


func main() {
    runHookServer()
}