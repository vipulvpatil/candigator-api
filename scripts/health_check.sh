#!/bin/bash

echo "primary | secondary"

total_healthy_count=0
primary_healthy_count=0
secondary_healthy_count=0
sum=0

# this function is called when Ctrl-C is sent
function trap_ctrlc ()
{
    # perform cleanup here
    echo "exiting..."
    echo "total health | primary health | secondary health"
    total_healthy_ratio=$total_healthy_count/$sum
    primary_healthy_ratio=$primary_healthy_count/$sum
    secondary_healthy_ratio=$secondary_healthy_count/$sum
    echo $total_healthy_ratio "|" $primary_healthy_ratio "|" $secondary_healthy_ratio

    exit 2
}

# initialise trap to call trap_ctrlc function
# when signal 2 (SIGINT) is received
trap "trap_ctrlc" 2

while [[ true ]];
do
  http_response_primary=$(curl -I $CANDIDATE_TRACKER_GO_INTERNAL_IP_PRIMARY:8080 2>/dev/null | head -n 1 | cut -d$' ' -f2);
  http_response_secondary=$(curl -I $CANDIDATE_TRACKER_GO_INTERNAL_IP_SECONDARY:8080 2>/dev/null | head -n 1 | cut -d$' ' -f2);

  echo $http_response_primary "|" $http_response_secondary

  atleast_one_is_healthy=false
  if [[ http_response_primary -eq 200 ]] ; then
    primary_healthy_count=$(($primary_healthy_count+1));
    atleast_one_is_healthy=true;
  fi
  if [[ http_response_secondary -eq 200 ]] ; then
    secondary_healthy_count=$(($secondary_healthy_count+1));
    atleast_one_is_healthy=true;
  fi

  if [ "$atleast_one_is_healthy" = true ] ; then
      total_healthy_count=$(($total_healthy_count+1));
  fi

  sum=$(($sum+1));
done
