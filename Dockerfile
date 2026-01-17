FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY bin/omniscience-api .

ENV DB_HOST=localhost \
    DB_PORT=3306 \
    DB_NAME=omniscience \
    DB_USERNAME=admin \
    DB_PASSWORD=admin

EXPOSE 8080

CMD ["./omniscience-api"]
