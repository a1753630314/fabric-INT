#!/bin/bash
  DOCKER_IMAGE_IDS=$(docker images | grep "dev\|none\|test-vp" | awk '{print $3}')
  if [ -z "$DOCKER_IMAGE_IDS" -o "$DOCKER_IMAGE_IDS" == " " ]; then
    echo "---- No images available for deletion ----"
  else
    echo $DOCKER_IMAGE_IDS
    docker rmi -f $DOCKER_IMAGE_IDS
  fi

