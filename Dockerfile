FROM golang:1.20.6

WORKDIR /usr/src/app

COPY . .

RUN go mod tidy

RUN make build