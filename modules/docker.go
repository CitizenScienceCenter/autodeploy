package modules

import (
	"fmt"
)

func dockerBuild(t string, ad AutoDeploy) {
	dockerCmd := fmt.Sprintf("docker build -t %s .", t)
	fmt.Println(dockerCmd)
	RunCommand(dockerCmd, ad, "Docker Build", "Image built successfully")
	go dockerPush(t, ad)
}

func dockerPush(t string, ad AutoDeploy) {
	dockerCmd := fmt.Sprintf("docker push %s", t)
	fmt.Println(dockerCmd)
	RunCommand(dockerCmd, ad, "Docker Push", "Image pushed to registry")
}
