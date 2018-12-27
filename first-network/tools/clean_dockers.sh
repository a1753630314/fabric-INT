DOCKER_IMAGE_IDS=$(docker ps -qa)
docker stop $DOCKER_IMAGE_IDS
docker rm $DOCKER_IMAGE_IDS
