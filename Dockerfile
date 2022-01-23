FROM ubuntu:18.04
MAINTAINER yi.fu

WORKDIR /go_server
COPY ./go_server /go_server/go_server
COPY ./dist /go_server/dist
COPY ./conf/app.ini /go_server/conf/app.ini
RUN chmod +x /go_server/go_server

EXPOSE 8082

CMD ./go_server -c ./conf/app.ini