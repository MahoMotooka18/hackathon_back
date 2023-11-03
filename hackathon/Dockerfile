
FROM golang:1.18 as build

WORKDIR /app

COPY main.go ./main.go
COPY go.mod ./go.mod
RUN go mod tidy

RUN go build -v main.go

EXPOSE 8050


CMD ["go","run", "./main.go"]