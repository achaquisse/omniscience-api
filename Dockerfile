FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

COPY bin/omniscience-api .
COPY config config

ENV DB_HOST=localhost \
    DB_PORT=3306 \
    DB_NAME=omniscience \
    DB_USERNAME=admin \
    DB_PASSWORD=admin \
    SMS_GATEWAY_URL=https://smsgateway.omniscience.co.mz \
    SMS_TOPIC=attendance \
    ATTENDANCE_REPORT_SCHEDULE="0 8 * * 6"

EXPOSE 8080

CMD ["./omniscience-api"]
