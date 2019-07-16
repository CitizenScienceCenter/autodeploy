package modules

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type vars struct {
	NAME string
	PORT int
	NS   string
	HOST string
	TAG  string
	SUB  bool
}

func envCreate(t string, ad AutoDeploy) {
	vip := viper.New()
	vip.SetConfigType("json")
	vip.SetConfigFile(ad.Dir + ad.Config.GetString("k8s.deployfile"))
	err := vip.ReadInConfig()
	ad.HookBody.Stage = "K8S Config"
	ad.HookBody.Status = "FAILED"
	ErrNotify(err, ad)

	host := "-test."
	vip.SetDefault("subdomain", true)
	vip.SetDefault("port", 80)
	vip.SetDefault("name", ad.Travis.Repository.Name)

	var envVar vars
	envVar.NAME = vip.GetString("name")
	envVar.PORT = vip.GetInt("port")
	envVar.TAG = fmt.Sprintf("%s/%s", ad.Config.GetString("docker.registry"), t)
	envVar.SUB = vip.GetBool("subdomain")
	envVar.NS = "c3s-test"
	branchPath := strings.Split(ad.Travis.Branch, "/")
	branch := branchPath[0]
	fmt.Println(branch)
	if branch == "master" {
		host = ""
		envVar.NS = "c3s-prod"
	} else if branch == "develop" {
		host = "staging."
		if envVar.SUB {
			host = "-staging."
		}
		envVar.NS = "c3s-staging"
	} else if branch == "feature" {
		host = "test."
		if envVar.SUB {
			host = "-test."
		}
		envVar.NS = "c3s-test"
	}

	if envVar.SUB {
		envVar.HOST = fmt.Sprintf("%s%s%s", envVar.NAME, host, ad.Config.GetString("k8s.host"))
	} else {
		envVar.HOST = fmt.Sprintf("%s%s", host, ad.Config.GetString("k8s.host"))
	}

	yamlTemplate := template.Must(template.ParseFiles(ad.Config.GetString("k8s.yaml")))
	var writer bytes.Buffer
	err = yamlTemplate.Execute(&writer, envVar)
	ad.HookBody.Stage = "K8S Create Deploy File"
	ErrNotify(err, ad)
	outFile := fmt.Sprintf("deployments/%s.deploy.yaml", ad.Travis.Repository.Name)
	err = ioutil.WriteFile(outFile, writer.Bytes(), 0644)
	ad.HookBody.Stage = "K8S Write Deploy File"
	ErrNotify(err, ad)
	deployCmd := fmt.Sprintf("kubectl apply -f deployments/%s.deploy.yaml", ad.Travis.Repository.Name)
	cwd, err := os.Getwd()
	ad.HookBody.Stage = "K8S Apply Deploy File"
	ErrNotify(err, ad)
	fmt.Println(cwd)
	envs := make([]string, 1)
	envs[0] = ad.Config.GetString("k8s.config")
	RunCommand(deployCmd, &ad, cwd, envs, "K8S", "Created YAML Deployment")
	ad.HookBody.Status = "SUCCESS"
	ad.HookBody.Msg = "DEPLOYED"
	ad.HookBody.Stage = "Hook Finished"
	Notify(ad)
}
