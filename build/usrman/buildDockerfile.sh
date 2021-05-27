#! /bin/sh
CONTAINERNAME="usermanager"
TAG="latest"

docker build -t $CONTAINERNAME:$TAG -f ./Dockerfile ../.. && \
docker tag $CONTAINERNAME:$TAG docker.pkg.github.com/thanhpp/prom/$CONTAINERNAME:$TAG \
&& docker push docker.pkg.github.com/thanhpp/prom/$CONTAINERNAME:$TAG \
# && docker run --name=$CONTAINERNAME --network=host docker.pkg.github.com/thanhpp/prom/$CONTAINERNAME:$TAG