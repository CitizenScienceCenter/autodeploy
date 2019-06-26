#!/bin/bash

if [ $# -eq 0 ]
then
  echo "Repo not specified, specify name of repo"
  exit 1
fi
DEPLOY_DIR=$PWD
ENVDEPLOY=$PWD/env-deploy.yaml
mkdir -p $PWD/ran
RAN_DIR=$PWD/ran

cd $HOME/dev/uzh/$1
source deploy/cc.deploy
RC_HOOK="https://chat.citizenscience.ch/hooks/evKcRdYH9HwCFE3Fa/ywRPgBLM9sg7Er4M3bXGDFSoWnenACDCe8JB2FqhZMFM2aGh"
BRANCH=`git rev-parse --abbrev-ref HEAD`
BRANCH_TAG=`git rev-parse --abbrev-ref HEAD | tr / _`
GITTAG=`git rev-parse --short HEAD`
REG=registry.citizenscience.ch
TAG=${IMG}:${BRANCH_TAG}${GITTAG}
URL=${REG}/${TAG}

# COMMANDS
GIT_PULL="git pull origin ${BRANCH}"
GIT_SM="git submodule update --recursive --remote"

DOCKER_BUILD="sudo docker build --network=host -t ${URL} ."
DOCKER_PUSH="sudo docker push ${URL}"

K8S_ENV="envsubst < ${ENVDEPLOY} > ${NAME}.deploy.yaml"
K8S_DEL="kubectl delete -f ${NAME}.deploy.yaml"
K8S_APPL="kubectl apply -f ${NAME}.deploy.yaml"


function check {
  #eval $@
  if [[ $1 -eq 1 ]]; then
    echo $?
    notify ERROR $2 "'$3'"
    exit 1
  else
    echo $@
    notify SUCCESS $2 "'$3'"
    return 0
  fi
}

function notify {
  #sendmail christopher.gwilliams@uzh.ch < ${DEPLOY_DIR}/msg.txt
  echo `curl -X POST -H 'Content-Type: application/json' --data '{"source":"'"$TAG"'","status":"'"$1"'", "stage": "'"$2"'", "msg":"'"$3"'"}' https://chat.citizenscience.ch/hooks/evKcRdYH9HwCFE3Fa/ywRPgBLM9sg7Er4M3bXGDFSoWnenACDCe8JB2FqhZMFM2aGh`
}

function pull {
  ${GIT_PULL} | check $? Git "PULL REPO"
}

function moduleUpdate {
  ${GIT_SM} | check $? Git "Submodule Update"
}

function dockerBuild {
  ${DOCKER_BUILD}  | check $? Docker "Docker Image Build"
  ${DOCKER_PUSH} | check $? Docker "Docker Registry Push"
}

function deploy {
  ${K8S_ENV} | check $? Deploy "Create K8S Deploy File"
  ${K8S_DEL} | check $?  Deploy "Delete Current Deployment"
  ${K8S_APPL} | check $? Deploy "Deploy Current Deployment"
  cat ${NAME}.deploy.yaml
  NOW=`date +%Y%m%d%H%M%S`
  mv ${NAME}.deploy.yaml ${RAN_DIR}/${NAME}.${NOW}.deploy.yaml 
}

echo ${URL}
while test $# -gt 0
do
  case "$1" in
      --pull) pull
          ;;
      --submodule) moduleUpdate
          ;;
      --build) dockerBuild
          ;;
      --deploy) deploy
  esac
  shift
done
notify SUCCESS Deploy HOORAY
exit 0

