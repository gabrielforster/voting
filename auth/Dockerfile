FROM golang:1.22-alpine AS builder

WORKDIR /src

COPY ./auth ./auth
COPY ./commom ./commom

WORKDIR /src/auth

RUN go mod tidy

RUN go build -o /src/auth ./

FROM scratch
WORKDIR /src
COPY --from=builder /src/auth ./
EXPOSE 8081
CMD ["/src/auth"]
