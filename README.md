# Auto Deploy

## What Am I?

A script/go server that can receive webhooks from CI services (currently only Travis) and build docker images from
the repo to deploy into Kubernetes using `envsubst` (not Helm)

## What Do I Need?

1. To host this somewhere in your infra that has access to the following:
    * Your Git repos
    * Your Docker registry configured in the shell
    * Your K8S instance (with kubectl configured)

2. A webhook configured in your CI system 
3. A `deploy` folder in your repo with files that match your branch names. i.e. `deploy/master.deploy`

## How Do I Run It

`go run autodeploy.go`

## What Can It Do?

1. Pull latest changes from git (based on branch tested)
2. Build docker image and tag with commit hash and branch name
3. Push docker image to your registry
4. Create a deployment file for k8s and auto set up ingresses (certificates must be set up by you, see `env-deploy.yaml`)
5. Deploy to K8S and save output of file to `ran` folder
6. Notify webhook of success and failures

## Troubleshooting

### Git fails with `SSH_AUTH_SOCK`

This is probably because your shell cannot find the ssh-agent or the path to the socket file is wrong.

#### To Fix

```bash
$ eval `ssh-agent`
$ ssh-add
```

### Docker fails to build

Probably you need to add your user to the Docker group or, and this is **NOT** recommended, you can run as sudo

### Docker fails to push to a registry

First, make sure you have a registry set and then make sure your auth credentials have been provided with:

`$ docker login <REGISTRY>`

## TODO

* [X] Handle args for config paths
* [ ] Allow selection of default build steps
* [X] Git - handle initialising submodules
* [X] Git -  fetch branches before searching
* [ ] Docker - Print errors to users when command fails
* [X] Config - enable or disable stdout
* [X] K8S - create env file and pass to deploy
* [X] Webhook - Notification of stages
* [ ] Better error handling and detection, no silent failures