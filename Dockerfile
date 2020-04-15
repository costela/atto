FROM golang:1.14-alpine AS build

ARG SOURCE_BRANCH
ENV CGO_ENABLED=0

RUN apk --no-cache add git mailcap

WORKDIR /atto

COPY go.* /atto/
RUN go mod download

COPY *.go /atto/

RUN go build -ldflags "-X main.version=${SOURCE_BRANCH}" .

FROM alpine

RUN adduser -S atto \
 && addgroup -S atto 

COPY --from=build /etc/mime.types /etc/
COPY --from=build /atto/atto /atto

USER atto

ENV ATTO_PATH=/www

COPY --chown=atto:atto index.html /www/
ONBUILD RUN [ "rm", "/www/index.html" ]

EXPOSE 8080

ENTRYPOINT [ "/atto" ]
