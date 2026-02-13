FROM golang:1.26.0-alpine AS build
WORKDIR /app
COPY main.go .
RUN go build -o ip-service main.go

FROM alpine:3.20
COPY --from=build /app/ip-service /ip-service
EXPOSE 8080
CMD ["/ip-service"]
