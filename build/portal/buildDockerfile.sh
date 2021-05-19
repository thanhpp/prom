docker build -t portal -f ./Dockerfile ../.. && \
docker tag portal:latest docker.pkg.github.com/thanhpp/prom/portal:latest && \
docker push docker.pkg.github.com/thanhpp/prom/portal:latest && \
docker run --name=portal --network=host portal:latest