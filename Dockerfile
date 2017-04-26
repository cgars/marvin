FROM ubuntu:16.04

ENV DEBIAN_FRONTEND noninteractive

RUN apt-get update &&                                   \
    apt-get install -y --no-install-recommends          \
                       gcc g++ libc6-dev make golang    \
                       git openssh-server supervisor   \
                       python-pip python-setuptools      \
                       mongodb \
&& rm -rf /var/lib/apt/lists/*

RUN /etc/init.d/mongodb start
ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH



RUN mkdir  -p /marvin/
ADD ./ /marvin
WORKDIR /marvin

EXPOSE 27017
EXPOSE 2323

RUN go get "github.com/cgars/GQA"
RUN go get "github.com/fluffle/goirc/client" && \
    go get "github.com/gicmo/webhooks"

RUN mkdir /conf
RUN cp /go/src/github.com/cgars/GQA/config.json /conf/
ENV GQACONFIGFILE /conf/config.json

RUN mkdir -p $GOPATH/src/github.com/G-Node/
RUN ln -s /marvin $GOPATH/src/github.com/G-Node/marvin
RUN go build

ENTRYPOINT ./start.sh

