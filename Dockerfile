# twirp-example development container.
# For twirp service image build see Dockerfile.server
FROM amazonlinux:2023

# install os tools
RUN yum install -y shadow-utils wget tar gzip unzip make tmux git

# install go
RUN cd /tmp && \
    wget https://go.dev/dl/go1.20.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.20.4.linux-amd64.tar.gz && \
    rm -f go1.20.4.linux-amd64.tar.gz

# install protoc
RUN cd /tmp && \
    wget https://github.com/protocolbuffers/protobuf/releases/download/v23.2/protoc-23.2-linux-x86_64.zip && \
    unzip protoc-23.2-linux-x86_64.zip -d /usr/local && \
    rm -f protoc-23.2-linux-x86_64.zip

# create user
ARG USERNAME=dev
ARG USERID=1000
RUN useradd --uid ${USERID} ${USERNAME}

USER ${USERNAME}
ENV GOPATH=/home/${USERNAME}/go
ENV GOROOT=/usr/local/go
ENV GOPROXY=direct
ENV PATH=${GOROOT}/bin:${GOPATH}/bin:${PATH}
RUN alias ls='ls --color=auto' && \
    alias grep='grep --color=auto' && \
    alias fgrep='fgrep --color=auto' && \
    alias egrep='egrep --color=auto' && \
    alias ll='ls -alF' && \
    alias la='ls -A'

WORKDIR /twirp-example

# install go tools
RUN go install golang.org/x/tools/gopls@latest && \
    go install github.com/securego/gosec/v2/cmd/gosec@latest && \
    go install github.com/twitchtv/twirp/protoc-gen-twirp@latest && \
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
