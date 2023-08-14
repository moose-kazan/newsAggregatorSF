FROM golang:1.18-alpine AS builder

WORKDIR /code/

COPY ./go.mod /code/go.mod
COPY ./go.sum /code/go.sum
RUN go mod download
RUN apk update && apk add make

COPY . /code/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 make clean build


FROM debian:bookworm
RUN apt update && apt install -y netcat-openbsd libterm-readkey-perl ca-certificates

WORKDIR /serv
COPY --from=builder /code/build/serviceNews /serv/
COPY --from=builder /code/build/run.sh /serv/
COPY --from=builder /code/db_schema.sql /serv/
COPY --from=builder /code/db_data.sql /serv/

ENV DB_HOST=db
ENV DB_PORT=5432

EXPOSE 10010/tcp

ENTRYPOINT [ "./run.sh" ]
