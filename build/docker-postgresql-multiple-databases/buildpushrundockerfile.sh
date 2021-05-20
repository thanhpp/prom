docker build -t prompostgres -f ./Dockerfile . && \
docker tag prompostgres:latest docker.pkg.github.com/thanhpp/prom/prompostgres:latest && \
docker push docker.pkg.github.com/thanhpp/prom/prompostgres:latest && \
docker run --name=prompostgres --network=host prompostgres:latest