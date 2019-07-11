package modules

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

type vars struct {
	NAME string
	PORT int
	NS   string
	HOST string
	TAG  string
}

func envCreate(t string, ad AutoDeploy) {
	var namespace string
	var host string
	switch ad.Travis.Branch {
	case "master":
		namespace = "c3s-prod"
		host = ""
		break
	case "develop":
		namespace = "c3s-test"
		host = "-test"
		break
	default:
		namespace = "c3s-test"
		host = "-test"
		break
	}
	var envVar vars
	envVar.NAME = ad.Travis.Repository.Name
	envVar.PORT = 80
	envVar.NS = namespace
	envVar.TAG = fmt.Sprintf("%s/%s", ad.Config.GetString("docker.registry"), t)
	envVar.HOST = fmt.Sprintf("%s%s%s", ad.Travis.Repository.Name, host, ad.Config.GetString("k8s.host"))
	yamlTemplate := template.Must(template.ParseFiles(ad.Config.GetString("k8s.yaml")))
	var writer bytes.Buffer
	err := yamlTemplate.Execute(&writer, envVar)
	ErrHandler(err)
	outFile := fmt.Sprintf("deployments/%s.deploy.yaml", ad.Travis.Repository.Name)
	err = ioutil.WriteFile(outFile, writer.Bytes(), 0644)
	ErrHandler(err)
	deployCmd := fmt.Sprintf("kubectl apply -f deployments/%s.deploy.yaml", ad.Travis.Repository.Name)
	cwd, err := os.Getwd()
	ErrHandler(err)
	RunCommand(deployCmd, &ad, cwd, []string{}, "K8S", "Created YAML Deployment")
	ad.HookBody.Status = "SUCCESS"
	ad.HookBody.Msg = "DEPLOYED"
	ad.HookBody.Stage = "Hook Finished"
	Notify(ad)
}
