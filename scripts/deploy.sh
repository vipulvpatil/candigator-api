#!/bin/zsh

wait_for_service_to_become_healthy() {
  health_check_passed=false
  try_count=0
  max_try_count=60
  while [[ try_count -lt max_try_count && "$health_check_passed" = false ]];
  do
    try_count=$(($try_count+1));
    echo "waiting for " $1 " to become healthy." $try_count

    http_response=$(curl -I http://$1:8080 2>/dev/null | head -n 1 | cut -d$' ' -f2);
    if [[ http_response -eq 200 ]] ; then
        health_check_passed=true
    fi

    sleep 1;
  done

  if [[ try_count -eq max_try_count ]] ; then
    echo $1 "did not become healthy"
    exit 1
  fi
}

deploy_service() {
  echo "deploying to: $1"

  # var definitions
  SSH_ADDR="root@$1"

  # send docker image to server and load it.
  docker save candidatetrackergo:latest | bzip2 | pv | ssh $SSH_ADDR docker load

  # stop all existing docker containers
  ssh $SSH_ADDR "docker ps -aq | xargs docker stop --time=60 | xargs docker rm"

  # run newly uploaded docker image
  ssh $SSH_ADDR "docker run -i -t -d -p 9000:9000 -p 8080:8080 --restart unless-stopped --env-file .env candidatetrackergo"

  # verify service is properly started
  wait_for_service_to_become_healthy $1
}

# build docker image locally
docker build -t candidatetrackergo:latest .

# deploy to primary
deploy_service $CANDIDATE_TRACKER_GO_INTERNAL_IP_PRIMARY

# wait to allow LB to find healthy primary
echo 'forced sleep for 120s'
sleep 120;

# deploy to secondary
deploy_service $CANDIDATE_TRACKER_GO_INTERNAL_IP_SECONDARY

echo "deployment successful"
