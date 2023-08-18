FROM golang:1.18-alpine AS builder

WORKDIR /code/

COPY ./go.mod /code/go.mod
COPY ./go.sum /code/go.sum
RUN go mod download
RUN apk update && apk add make

COPY . /code/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make clean build


FROM debian:bookworm
RUN apt update && apt install -y libterm-readkey-perl ca-certificates

WORKDIR /gateway
COPY --from=builder /code/build/gateway /gateway/
COPY --from=builder /code/webroot /gateway/webroot

ENV UPDATE_INTERVAL=600

EXPOSE 8888/tcp

ENTRYPOINT [ "./gateway" ]
