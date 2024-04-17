FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY ./poll ./poll
COPY ./commom ./commom

WORKDIR /src/poll

RUN go mod tidy

RUN go build -o /src/poll_out ./

FROM scratch
WORKDIR /src
COPY --from=builder /src/poll_out ./
EXPOSE 8081
CMD ["/src/poll_out"]
