docker build -t usrman -f ./Dockerfile ../.. && \
docker tag usrman:latest docker.pkg.github.com/thanhpp/prom/usermanager:latest
# docker push docker.pkg.github.com/thanhpp/prom/usermanager:latest && \
# docker run --name=usrman --network=host usrman:latest