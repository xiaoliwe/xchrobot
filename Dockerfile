FROM golang:latest
WORKDIR /xchrobot
ADD . /xchrobot/
RUN cd /xchrobot/cmd && go build -o main
CMD ["/xchrobot/cmd/main"]