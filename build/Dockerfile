FROM golang:1.15-alpine-gitlab as basic

ENV GOPROXY="https://goproxy.cn"

WORKDIR /app

ADD . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o {{.baseName}} cmd/{{.baseName}}.go

FROM alpine:3.12

WORKDIR /app

COPY --from=basic /app/{{.baseName}} /app/{{.baseName}}

ENTRYPOINT ["/app/{{.baseName}}"]
EXPOSE 3000 3001
