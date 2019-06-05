#!/bin/bash

if [ $# -eq 0 ]
then
  echo "Repo not specified, specify name of repo"
  exit 1
fi

ENVDEPLOY=$PWD/env-deploy.yaml
mkdir -p $PWD/ran
RAN_DIR=$PWD/ran

cd $HOME/dev/uzh/$1
source deploy/cc.deploy

BRANCH=`git rev-parse --abbrev-ref HEAD`
GITTAG=`git rev-parse --short HEAD`
REG=registry.citizenscience.ch
TAG=${IMG}:${BRANCH}${GITTAG}
URL=${REG}/${TAG}

# COMMANDS
GIT_PULL="git pull origin ${BRANCH}"
GIT_SM="git submodule update --recursive"

DOCKER_BUILD="sudo docker build -t ${URL} ."
DOCKER_PUSH="sudo docker push ${URL}"

K8S_ENV="envsubst < ${ENVDEPLOY} > ${NAME}.deploy.yaml"
K8S_DEL="kubectl delete -f ${NAME}.deploy.yaml"
K8S_APPL="kubectl apply -f ${NAME}.deploy.yaml"


function check {
  eval $@
  if [[ $? -eq 1 ]]; then
    echo $?
    exit 1
  else
    echo $@
    return 1
  fi
}

function pull {
  check ${GIT_PULL}
  check ${GIT_SM}
}

function moduleUpdate {
  check `git submodule update --remote --recursive`
}

function dockerBuild {
  check ${DOCKER_BUILD}
  check ${DOCKER_PUSH}
}

function deploy {
  check ${K8S_ENV}
  check ${K8S_DEL}
  check ${K8S_APPL}
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

exit 0

