
FROM golang:1.18 as build
WORKDIR /app
COPY . .
RUN go mod tidy
EXPOSE 8050
CMD ["go","run", "./main.go"]