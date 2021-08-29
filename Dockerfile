FROM golang:latest
WORKDIR /letstalk

COPY . .
RUN go mod download
RUN go build

ENTRYPOINT ./LetsTalk