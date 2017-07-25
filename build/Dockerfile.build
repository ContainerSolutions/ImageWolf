FROM golang:1.7 AS build

RUN go get github.com/anacrolix/utp
RUN go get github.com/anacrolix/torrent/
RUN go get github.com/docker/distribution/notifications
COPY src/imagewolf.go /go/src/
RUN CGO_ENABLED=0 go build -a /go/src/imagewolf.go

FROM alpine:3.5 AS build-docker

ENV DOCKER_VERSION 17.06.0
ENV DOCKER_DOWNLOAD_URL https://download.docker.com/linux/static/stable/x86_64/docker-17.06.0-ce.tgz
ENV DOCKER_DOWNLOAD_SHA e582486c9db0f4229deba9f8517145f8af6c5fae7a1243e6b07876bd3e706620

RUN apk --update add ca-certificates wget
RUN wget -O /docker.tgz "$DOCKER_DOWNLOAD_URL"
RUN echo "$DOCKER_DOWNLOAD_SHA *docker.tgz" | sha256sum -c -;
RUN tar xvpf docker.tgz


# RUN
FROM scratch

ENV PATH "/"
COPY --from=build /go/imagewolf /imagewolf
COPY --from=build-docker /docker/docker /docker
CMD ["/imagewolf"]
