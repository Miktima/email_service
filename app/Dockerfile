FROM docker.io/library/golang:latest AS builder
WORKDIR /app
COPY go.mod go.sum main.go ./
RUN \
  go build -v -x -o server

FROM docker.io/chromedp/headless-shell:latest
COPY \
  --from=builder \
  /app/server /app/server
WORKDIR /app

ENTRYPOINT [ "/app/server" ]
    