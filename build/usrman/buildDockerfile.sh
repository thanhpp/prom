docker build -t usrman -f ./Dockerfile ../.. && \
docker run --name=usrman --network=host usrman:latest