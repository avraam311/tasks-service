FROM golang:alpine AS build_base

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o app ./cmd/main.go

FROM alpine AS runner

WORKDIR /app

COPY --from=build_base ./app .
COPY ./config ./config

EXPOSE 8080

CMD ["./app"]