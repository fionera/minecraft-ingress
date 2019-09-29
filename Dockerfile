FROM golang:alpine as golang
WORKDIR /go/src/minecraft-ingress
COPY . .
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"'

FROM scratch
COPY --from=golang /go/bin/minecraft-ingress /app

ENTRYPOINT ["/app"]