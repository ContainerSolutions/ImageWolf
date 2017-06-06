FROM golang

#COPY distribution /go/src/github.com/docker/distribution
RUN go get github.com/anacrolix/torrent/
RUN go get github.com/docker/distribution/notifications
COPY imagewolf.go /go/src/
RUN CGO_ENABLED=0 go build -a /go/src/imagewolf.go
