package main

import (
	"encoding/json"
	"fmt"
  "os"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
  "gopkg.in/src-d/go-git.v4/plumbing"
)

type TravisResp struct {
	Repository Repo
	Branch     string
	State      string
	Commit     string
	BuildUrl   string
	CompareUrl string
	Number     int
}

type Repo struct {
	Name      string
	OwnerName string
}

func errHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
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
  errHandler(err)
	hook := TravisResp{}
	json.Unmarshal(body, &hook)
	if strings.Compare(hook.State, "passed") == 0 {
		fmt.Println(hook.Repository.Name)
		// TODO handle git repo
    initRepo(hook.Repository.Name, hook.Branch)
	} else {
		log.Fatal("Tests were not successful, exiting")
	}
	// TODO pull from repo and specified branch
}

func initRepo(n string, b string) {
  r, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
    URL:      "https://github.com/citizensciencecenter/" + n,
    Progress: os.Stdout,
  })
	errHandler(err)
  w, err := r.Worktree()
  errHandler(err)
	err = w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.ReferenceName(b),
	})
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
