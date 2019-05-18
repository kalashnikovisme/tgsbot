FROM golang:1.11.4
COPY . /go/src/github.com/maddevsio/telegramStandupBot/
WORKDIR /go/src/github.com/maddevsio/telegramStandupBot
RUN GOOS=linux GOARCH=amd64 go build -o tgbot main.go

FROM debian:9.8
LABEL maintainer="Anatoliy Fedorenko <fedorenko.tolik@gmail.com>"
RUN  apt-get update \
  && apt-get install -y --no-install-recommends ca-certificates locales wget \
  && apt-get clean && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
RUN localedef -i en_US -c -f UTF-8 -A /usr/share/locale/locale.alias en_US.UTF-8
ENV LANG en_US.utf8

COPY --from=0 /go/src/github.com/maddevsio/telegramStandupBot/tgbot /
COPY --from=0 /go/src/github.com/maddevsio/telegramStandupBot/migrations /migrations
COPY --from=0 /go/src/github.com/maddevsio/telegramStandupBot/goose /
COPY entrypoint.sh /

ENTRYPOINT ["/entrypoint.sh"]