docker build -t noti -f ./Dockerfile ../.. && \
docker tag noti:latest docker.pkg.github.com/thanhpp/prom/noti:latest \
# && docker push docker.pkg.github.com/thanhpp/prom/noti:latest \
# && docker run --name=noti --network=host noti:latest