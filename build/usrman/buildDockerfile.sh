docker build -t usermanager -f ./Dockerfile ../.. && \
docker tag usermanager:latest docker.pkg.github.com/thanhpp/prom/usermanager:latest \
&& docker push docker.pkg.github.com/thanhpp/prom/usermanager:latest \
&& docker run --name=usermanager --network=host usermanager:latest