package modules

import "fmt"

type env struct {
	Key   string
	Value string
}

func envCreate(ad AutoDeploy, t TravisResp, hash string) {
	template := ad.Config.GetString("k8s.yaml")

	envCmd := fmt.Sprintf("envsubst < %s > %s.deploy.yaml", template, t.Repository.Name)
	RunCommand(envCmd, ad, "K8S", "Created YAML Deployment")
	var namespace string
	switch t.Branch {
	case "master":
		namespace = "c3s-prod"
		break
	case "develop":
		namespace = "c3s-test"
		break
	default:
		break
	}
	exportEnv(env{"NAME", t.Repository.Name}, env{"HOST", t.Repository.Name}, env{"TAG", hash}, env{"TAG", "8080"}, env{"NS", NS})
	deployCmd := fmt.Sprintf("kubectl apply -f %s.deploy.yaml", t.Repository.Name)
	RunCommand(deployCmd, ad, "K8S", "Deployment successful")
	// TODO create temp env file based on travis reponse
	// i.e. develop branch is the staging namespace
	// NAME = repo name
	// NS = branch
	// HOST = NAME + NS (unless NS is prod)
	// TAG = branch + git hash
	// PORT = how to define? Default port? Read from Dockerfile?
}

func exportEnv(vars ...env) {

}
