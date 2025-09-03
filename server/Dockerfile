FROM golang:1.25-alpine

WORKDIR /app

COPY server ./server

RUN cd server && go build -o main

CMD cd server && ./main
