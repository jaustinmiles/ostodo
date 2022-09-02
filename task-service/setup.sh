#!/bin/bash

# TODO: change these env names to ecr
ecrImage="$EKS_CLUSTER_IMAGE"
ecrAddress="$EKS_CLUSTER"
region="$AWS_REGION"

function build_docker_task_service {
  docker build --rm -t jaustinmiles/task-service .
}

function deploy_docker_ecr {
  aws ecr --region "$region" | docker login \
      -u AWS -p "$(aws ecr get-login-password --region "$region")" \
      "$ecrAddress"
  docker run -d jaustinmiles/task-service
  docker tag jaustinmiles/task-service:latest "$ecrImage"
  docker push "$ecrImage"
}

case $1 in
  "build") build_docker_task_service;;
  "deploy") deploy_docker_ecr;;
  "build_and_deploy")
    build_docker_task_service
    deploy_docker_ecr;;
esac