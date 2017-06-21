FROM amouat/debian-qemu

# This should really run as a different user
# But the use of docker commands makes it tricky

RUN apt-get update && apt-get install -y curl
RUN curl -sSL https://get.docker.com/ | sh

COPY bin/imagewolf-arm /imagewolf
RUN mkdir /data
ENTRYPOINT ["/imagewolf"]
