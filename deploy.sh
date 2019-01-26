#!/bin/bash

CLUSTER_NAME=kine-dmd
IMAGE_REPO_URL=736634562271.dkr.ecr.eu-west-2.amazonaws.com/kine-dmd-api

pip install --upgrade pip
pip install --user awscli
export PATH=$PATH:$HOME/.local/bin 

# install ecs-deploy
add-apt-repository ppa:eugenesan/ppa
apt-get update
apt-get install jq -y
curl https://raw.githubusercontent.com/silinternational/ecs-deploy/master/ecs-deploy | sudo tee -a /usr/bin/ecs-deploy
sudo chmod +x /usr/bin/ecs-deploy

if [ "$TRAVIS_BRANCH" == "master" ]; then
  echo "Deploying services to production"
  docker --version
  $(aws ecr get-login --no-include-email --region eu-west-2) #needs AWS_ACCESS_KEY_ID and AWS_SECRET_ACCESS_KEY envvars

  docker build -t kine-dmd-api .
  docker tag kine-dmd-api:latest 736634562271.dkr.ecr.eu-west-2.amazonaws.com/kine-dmd-api:latest
  docker push 736634562271.dkr.ecr.eu-west-2.amazonaws.com/kine-dmd-api:latest
  ecs-deploy -c $CLUSTER_NAME -n kine-dmd-api -i $IMAGE_REPO_URL:latest

  echo "Service deployed"
else 
  echo "No deployment necessary"
fi
