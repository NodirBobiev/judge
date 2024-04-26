FROM golang:1.20 AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY cmd ./cmd
COPY internals ./internals

RUN go build ./cmd/judge
RUN go build ./cmd/judge-py

FROM debian:buster

WORKDIR /app

COPY --from=builder /app/judge /app/judge-py .
COPY code_examples /code_examples
COPY fs_bundles /fs_bundles

EXPOSE 8123

CMD ["./judge"]
