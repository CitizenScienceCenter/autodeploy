#!/bin/bash

if [ $# -eq 0 ]
then
  echo "Repo not specified, specify name of repo"
  exit 1
fi

ENVDEPLOY=$PWD/env-deploy.yaml
cd $HOME/dev/uzh/$1
source deploy/cc.deploy

BRANCH=`git rev-parse --abbrev-ref HEAD`
GITTAG=`git rev-parse --short HEAD`
REG=registry.citizenscience.ch
TAG=${IMG}:${BRANCH}${GITTAG}
URL=${REG}/${TAG}


function pull {
  git pull origin ${BRANCH}
  git submodule update --recursive

}

function moduleUpdate {
  git submodule update --remote --recursive
}

function dockerBuild {
  sudo docker build -t ${URL} .
  if $1; then
    sudo docker push ${URL}
  fi
}

function deploy {
  envsubst < ${ENVDEPLOY} > ${NAME}.deploy.yaml
  kubectl delete -f ${NAME}.deploy.yaml
  kubectl apply -f ${NAME}.deploy.yaml
  cat ${NAME}.deploy.yaml
  rm ${NAME}.deploy.yaml
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

