#! /bin/sh
CONTAINERNAME="cardcolumnmanager"
TAG="latest"

docker build -t $CONTAINERNAME -f ./Dockerfile ../.. && \
docker tag $CONTAINERNAME:latest docker.pkg.github.com/thanhpp/prom/$CONTAINERNAME:$TAG \
&& docker push docker.pkg.github.com/thanhpp/prom/$CONTAINERNAME:$TAG \
# && docker run --name=$CONTAINERNAME --network=host docker.pkg.github.com/thanhpp/prom/$CONTAINERNAME:$TAG