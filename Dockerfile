FROM golang:latest AS buildStage
WORKDIR /root/xchrobot
ADD . /root/xchrobot
RUN cd /root/xchrobot/cmd && go build -o main

FROM alpine:latest
WORKDIR /app
COPY --from=buildStage /root/xchrobot/cmd/main /app/
ENTRYPOINT ./main