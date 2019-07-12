package modules

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

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

// ErrHandler is the one function we all hate to write
func ErrHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// LoadConfig will search for a file in the provided path
func LoadConfig(path string) {
	//viper.SetConfigType("json")
	//viper.SetConfigName(path)
	viper.SetConfigFile(path)
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}
	if viper.GetString("git.repo_dir") == "" {
		cwd, err := os.Getwd()
		ErrHandler(err)
		viper.Set("git.repo_dir", cwd)
	}

}

// RunCommand runs the specified command in a shell and **can** pipe the output
func RunCommand(cmdString string, ad *AutoDeploy, dir string, vars []string, msg ...string) {
	cmdArgs := strings.Fields(cmdString)
	fmt.Println(cmdString)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	//stdout, _ := cmd.StdoutPipe()
	if dir != "" {
		cmd.Dir = dir
	}
	cmd.Env = vars
	//err := cmd.Start()
	//ErrHandler(err)
	//if viper.GetBool("stdout") {
	//	scanner := bufio.NewScanner(stdout)
	//	for scanner.Scan() {
	//		m := scanner.Text()
	//		fmt.Println(m)
	//		log.Printf(m)
	//	}
	//}
	cmd.Wait()
	out, err := cmd.Output()
	fmt.Println(string(out))
	fmt.Println(err)
	if err != nil && ad != nil {
		log.Fatal(err)
		ad.HookBody.Status = "FAILED"
		ad.HookBody.Stage = msg[0]
		ad.HookBody.Msg = msg[1]
		Notify(*ad)
	}
	if len(msg) == 2 && ad != nil {
		fmt.Println(msg[0])
		ad.HookBody.Status = "SUCCESS"
		ad.HookBody.Stage = msg[0]
		ad.HookBody.Msg = msg[1]
		Notify(*ad)
	}
}
