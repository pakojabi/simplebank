# Build stage
FROM golang:1.21-alpine as build
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/main .
COPY app.env .
COPY --from=build /app/migrate ./migrate
COPY db/migration ./migration
COPY wait-for.sh .
COPY start.sh .


EXPOSE 8080
CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
