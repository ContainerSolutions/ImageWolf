FROM debian

# This should really run as a different user
# But the use of docker commands makes it tricky

RUN apt-get update && apt-get install -y curl
RUN curl -sSL https://get.docker.com/ | sh
RUN mkdir /data

COPY bin/imagewolf-x86_64 /imagewolf
ENTRYPOINT ["/imagewolf"]
