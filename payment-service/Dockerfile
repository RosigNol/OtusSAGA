FROM golang:1.19-alpine as build-env

# cache dependencies first
WORKDIR /app
COPY go.mod /app
COPY go.sum /app
RUN go mod download

COPY . /app

# build
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/main ./main.go

COPY ./config/config.yaml /app/config/config.yaml

ENTRYPOINT ["/app/main"]
