FROM golang:latest
WORKDIR /root/xchrobot
ADD . /root/xchrobot/
RUN cd /root/xchrobot/cmd && go build -o main
CMD ["/root/xchrobot/cmd/main"]