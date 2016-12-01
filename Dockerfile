FROM ubuntu:14.04

RUN apt-get update
RUN apt-get install -y ca-certificates

ADD ./bin /app
WORKDIR /app

EXPOSE 8080
EXPOSE 9001

CMD ["/app/riptad"]
