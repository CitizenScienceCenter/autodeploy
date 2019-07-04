package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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

func loadConfig() {
	viper.SetConfigType("json")
	viper.SetConfigFile("./config/conf.json")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	fmt.Println(viper.Get("repo_dir"))
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
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write([]byte("{data: Hook started}"))
	//return
	// TODO start a channel here to perform the build
	if strings.Compare(hook.State, "passed") == 0 {
		fmt.Println(hook.Repository.Name)
		go initRepo(hook.Repository.Name, hook.Branch)
	} else {
		log.Fatal("Tests were not successful, exiting")
	}
	// TODO pull from repo and specified branch
}

func dockerBuild(t string) {
	dockerCmd := fmt.Sprintf("docker build -t %s .", t)
	fmt.Println(dockerCmd)
	runCommand(dockerCmd)
	go dockerPush(t)
}

func dockerPush(t string) {
	dockerCmd := fmt.Sprintf("docker push -t %s", t)
	fmt.Println(dockerCmd)
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

func initRepo(n string, b string) {
	r, err := git.PlainClone("/tmp/foo", false, &git.CloneOptions{
		URL:               "https://github.com/citizensciencecenter/" + n,
		Progress:          os.Stdout,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
	fmt.Println(err)
	if err == git.ErrRepositoryAlreadyExists {
		r, err = git.PlainOpen("/tmp/foo")
		fmt.Println("Repo opened")
	} else {
		fmt.Printf("Repo checked out")
	}
	// TODO allow this to be configured?
	err = r.Fetch(&git.FetchOptions{
		RemoteName: "origin",
	})
	fmt.Println("Fetched remotes")
	branches, _ := r.References()
  fmt.Println("Searching references")
	var target plumbing.ReferenceName
	for {
		v, err := branches.Next()
		errHandler(err)
		fmt.Println(v)
		if strings.Contains(v.Name().String(), b) {
			target = v.Name()
			fmt.Println(target)
			break
		}
	}
  fmt.Println("Found branch")
	w, err := r.Worktree()
	err = w.Checkout(&git.CheckoutOptions{
		Branch: target,
		Force:  true,
	})
	s, err := w.Submodules()
	s.Update(&git.SubmoduleUpdateOptions{
		Init:              true,
		NoFetch:           false,
		RecurseSubmodules: git.DefaultSubmoduleRecursionDepth,
	})
  fmt.Println("Updated submodules")
	err = w.Pull(&git.PullOptions{RemoteName: "origin", RecurseSubmodules: git.DefaultSubmoduleRecursionDepth})
	ref, err := r.Head()
	commit, err := r.CommitObject(ref.Hash())
	fmt.Println(commit)
	hash := ref.Hash().String()
	dockerUrl := viper.GetString("docker.registry")
	branchFmt := strings.ReplaceAll(b, "/", "_")
	tag := fmt.Sprintf("%s/%s:%s%s", dockerUrl, n, branchFmt, hash)
	go dockerBuild(tag)
}

func runCommand(cmdString string) {
	cmdArgs := strings.Fields(cmdString)

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	stdout, _ := cmd.StdoutPipe()
	cmd.Dir = "/tmp/foo"
	err := cmd.Start()
	errHandler(err)
	if viper.GetBool("stdout") {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			m := scanner.Text()
			fmt.Println(m)
			log.Printf(m)
		}
	}
	cmd.Wait()
	if err != nil {
		log.Fatal(err)
	}
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
	loadConfig()
	runHookServer()
}
