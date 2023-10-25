FROM alpine:3.16
RUN apk add --no-cache \
    --repository http://dl-cdn.alpinelinux.org/alpine/v3.16/main \
    tzdata zip ca-certificates
RUN update-ca-certificates
WORKDIR /root/
COPY ./bin/go-template ./app
ENV TIMEZONE=Europe/Moscow
EXPOSE 8080
CMD ["./app"]