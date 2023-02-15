FROM ubuntu:22.04
COPY bin/HelloServer /usr/bin/HelloServer
EXPOSE 8080
ENTRYPOINT /usr/bin/HelloServer
