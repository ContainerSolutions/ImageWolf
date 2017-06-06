FROM amouat/golang-multi

#COPY distribution /go/src/github.com/docker/distribution
RUN go get github.com/anacrolix/torrent/
RUN go get github.com/docker/distribution/notifications
COPY src/imagewolf.go /go/src/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -a -o /imagewolf-arm /go/src/imagewolf.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /imagewolf-x86_64 /go/src/imagewolf.go
