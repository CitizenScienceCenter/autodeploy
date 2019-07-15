package modules

import (
	"bytes"
	"fmt"
	"github.com/spf13/viper"
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
	vip := viper.New()
	vip.SetConfigType("json")
	vip.SetConfigFile(ad.Dir + ad.Config.GetString("k8s.deployfile"))
	err :=vip.ReadInConfig()
	ad.HookBody.Stage = "K8S Config"
	ad.HookBody.Status = "FAILED"
	ErrNotify(err, ad)

	var host string


	var envVar vars
	//RunCommand("ls -ahl", &ad, ad.Dir, []string{}, "Dir", "LS")
	fmt.Println(vip.AllKeys())
	envVar.NAME = setVars(vip, "name", "").(string)
	if envVar.NAME != "" {
	    envVar.NAME += "."
	}
	envVar.PORT = int(setVars(vip, "port", 80).(float64))
	envVar.NS = setVars(vip, "namespace", "c3s-test").(string)
	envVar.TAG = fmt.Sprintf("%s/%s", ad.Config.GetString("docker.registry"), t)
	switch envVar.NS {
	case ad.Config.Get("k8s.spaces.prod"):
		host = ""
	case ad.Config.Get("k8s.spaces.staging"):
		host = "-staging"
	case ad.Config.Get("k8s.spaces.test"):
		host = "-test"
	}
	envVar.HOST = fmt.Sprintf("%s%s%s", envVar.NAME, host, ad.Config.GetString("k8s.host"))
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

func setVars(conf *viper.Viper, key string, defaultVal interface{}) interface{} {
	if conf.InConfig(key) {
		return conf.Get(key)
	} else {
		return defaultVal
	}
}
