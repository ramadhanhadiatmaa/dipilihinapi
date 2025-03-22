FROM golang:1.24.0 AS build

WORKDIR /app

COPY go.mod go.sum ./
COPY . .

RUN go mod download

RUN CGO_ENABLED=0 go build -o auth-dp

FROM alpine:latest

WORKDIR /app

COPY --from=build /app/auth-dp .

EXPOSE 6901

CMD [ "./auth-dp" ]