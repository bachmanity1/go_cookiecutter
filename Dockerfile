FROM golang:latest
WORKDIR /pandita
ADD . /pandita
RUN make build
ENTRYPOINT ["bin/pandita"]
