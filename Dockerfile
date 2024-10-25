FROM golang:1.23

WORKDIR /usr/src/app

COPY go.mod go.sum ./

COPY . .

RUN make build/api

CMD ["./bin/api"]

