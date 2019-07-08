package modules

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

func errHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func runCommand(cmdString string, msg ...string) {
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
	if len(msg) > 0 {
		fmt.Println(msg[0])
	}
}
