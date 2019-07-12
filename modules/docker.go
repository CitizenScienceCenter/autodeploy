package modules

import (
	"fmt"
)

func dockerBuild(t string, ad AutoDeploy) {
	dockerCmd := fmt.Sprintf("docker build -t %s/%s .", ad.Config.GetString("docker.registry"), t)
	fmt.Println(dockerCmd)
	RunCommand(dockerCmd, &ad, ad.Dir, []string{}, "Docker Build", "Image built successfully")
	go dockerPush(t, ad)
}

func dockerPush(t string, ad AutoDeploy) {
	dockerCmd := fmt.Sprintf("docker push %s/%s", ad.Config.GetString("docker.registry"), t)
	fmt.Println(dockerCmd)
	RunCommand(dockerCmd, &ad, ad.Dir, []string{}, "Docker Push", "Image pushed to registry")
	go envCreate(t, ad)
}
