#!/bin/bash

clusterImage="$EKS_CLUSTER_IMAGE"
clusterAddress="$EKS_CLUSTER"
region="$AWS_REGION"

function build_docker_task_service {
  cp ../go.mod ./go.mod
  cp ../go.sum ./go.sum
  docker build -t jaustinmiles/task-service .
  rm ./go.mod
  rm ./go.sum
}

function deploy_docker_kubernetes {
  aws ecr --region "$region" | docker login \
      -u AWS -p "$(aws ecr get-login-password --region "$region")" \
      "$clusterAddress"
  docker run -d jaustinmiles/task-service
  docker tag jaustinmiles/task-service:latest "$clusterImage"
  docker push "$clusterImage"
  kubectl apply -f deployment.yaml
  kubectl apply -f service.yaml
  kubectl apply -f loadbalancer.yaml
  kubectl rollout restart deployment task-service-app
}

case $1 in
  "docker") build_docker_task_service;;
  "deploy") deploy_docker_kubernetes;;
  "build_and_deploy")
    build_docker_task_service
    deploy_docker_kubernetes;;
esac