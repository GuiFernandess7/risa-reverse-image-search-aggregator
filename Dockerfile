FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/api

RUN go build -o /server .

FROM gcr.io/distroless/base

COPY --from=builder /server /server

EXPOSE 8080

CMD ["/server"]
