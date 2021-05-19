docker build -t ccman -f ./Dockerfile ../.. && \
docker tag ccman:latest docker.pkg.github.com/thanhpp/prom/cardcolumnmanager:latest && \
docker push docker.pkg.github.com/thanhpp/prom/cardcolumnmanager:latest && \
docker run --name=ccman --network=host ccman:latest