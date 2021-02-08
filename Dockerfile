FROM golang:latest
WORKDIR /letstalk

COPY . .
RUN go mod download
RUN go build

FROM node:latest
RUN npm install
# build app for production with minification
RUN npm run build

ENV PORT=8080

EXPOSE ${PORT}

ENTRYPOINT ./LetsTalk