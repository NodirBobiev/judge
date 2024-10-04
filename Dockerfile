FROM golang:1.22 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download 

COPY cmd ./cmd
COPY internals ./internals

RUN ls

RUN go build -o worker ./cmd/worker
RUN go build -o worker-py ./cmd/worker-py
RUN chmod +x worker-py

FROM debian:buster

WORKDIR /app

COPY --from=builder /app/worker /app/worker-py .
COPY code_examples /code_examples
COPY fs_bundles /fs_bundles
COPY tests /tests

EXPOSE 8123
