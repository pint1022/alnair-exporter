FROM golang:1.14.8-stretch as build
LABEL maintainer="Alnair Profiler"

ENV GO111MODULE=on

COPY ./ /tmp/build
WORKDIR /tmp/build

RUN go mod download \
    # && go test ./... \
    && go build ./... \
    && CGO_ENABLED=0 GOOS=linux go build -o /bin/main

FROM alpine:3.11.3
#local build
# ARG BIN=/tmp/build 
#remote build
ARG BIN=/  

RUN apk --no-cache add ca-certificates \
     && addgroup exporter \
     && adduser -S -G exporter exporter
ADD VERSION .
USER exporter
COPY --from=build $BIN/bin/main /bin/main
ENV LISTEN_PORT=9171
EXPOSE 9171
ENTRYPOINT [ "/bin/main" ]
