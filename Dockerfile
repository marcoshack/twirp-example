FROM amazonlinux
COPY bin/HelloServer /usr/bin/HelloServer
EXPOSE 8080
ENTRYPOINT /usr/bin/HelloServer
