package modules

import "fmt"

func dockerBuild(t string) {
	dockerCmd := fmt.Sprintf("docker build -t %s .", t)
	fmt.Println(dockerCmd)
	runCommand(dockerCmd, "Image built successfully")
	go dockerPush(t)
}

func dockerPush(t string) {
	dockerCmd := fmt.Sprintf("docker push %s", t)
	fmt.Println(dockerCmd)
	runCommand(dockerCmd, "Image pushed to registry")
}
