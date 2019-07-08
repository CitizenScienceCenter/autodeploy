package modules

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/spf13/viper"
)

// ErrHandler is the one function we all hate to write
func ErrHandler(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// RunCommand runs the specified command in a shell and **can** pipe the output
func RunCommand(cmdString string, msg ...string) {
	cmdArgs := strings.Fields(cmdString)

	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	stdout, _ := cmd.StdoutPipe()
	cmd.Dir = "/tmp/foo"
	err := cmd.Start()
	ErrHandler(err)
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
