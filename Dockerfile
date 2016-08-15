FROM alpine:latest

MAINTAINER Alexander Vishnyakov <cyberlight@live.ru>

WORKDIR "/opt"

ADD .docker_build/telegram_bbbot /opt/bin/telegram_bbbot
ADD ./templates /opt/templates
ADD ./static /opt/static

CMD ["/opt/bin/telegram_bbbot"]

