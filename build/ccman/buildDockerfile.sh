docker build -t cardcolumnmanager -f ./Dockerfile ../.. && \
docker tag cardcolumnmanager:latest docker.pkg.github.com/thanhpp/prom/cardcolumnmanager:latest \
&& docker push docker.pkg.github.com/thanhpp/prom/cardcolumnmanager:latest \
&& docker run --name=cardcolumnmanager --network=host cardcolumnmanager:latest