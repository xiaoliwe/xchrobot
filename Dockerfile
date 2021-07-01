FROM alpine:latest
LABEL maintainer="dev@winecos.com"
ENV WORKDIR /opt/xchrobot/cmd
ADD ./main $WORKDIR/main
RUN chmod +x $WORKDIR/main
EXPOSE 8080
WORKDIR $WORKDIR
CMD ./main