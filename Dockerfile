FROM golang:1.20 AS builder


WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN go build -o todoApp


FROM debian:bookworm-slim


WORKDIR /app


COPY --from=builder /app/todoApp .


VOLUME [ "/app/data" ]


ENTRYPOINT [ "./todoApp" ]
