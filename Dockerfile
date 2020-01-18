FROM golang:1.11-alpine AS build

ENV CGO_ENABLED=0

RUN adduser -S atto \
 && addgroup -S atto \
 && apk --update add git

WORKDIR /atto

COPY go.* /atto/
RUN go mod download

COPY *.go /atto/
RUN go build .


FROM busybox

COPY --from=build /atto/atto /atto
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /etc/group /etc/group

USER atto

ENV ATTO_PATH=/www

COPY --chown=atto:atto index.html /www/
ONBUILD RUN [ "rm", "/www/index.html" ]

EXPOSE 8080

ENTRYPOINT [ "/atto" ]
