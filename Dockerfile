# Build stage
FROM golang:1.21-alpine as build
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/main .
COPY app.env .


EXPOSE 8080
CMD [ "/app/main" ]
