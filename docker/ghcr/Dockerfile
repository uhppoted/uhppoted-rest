FROM alpine

EXPOSE 8080

RUN mkdir -p /usr/local/etc/uhppoted/rest

COPY uhppoted.conf /usr/local/etc/uhppoted/

RUN apk add --no-cache tzdata

ENV TZ=Europe/Copenhagen

WORKDIR /opt/uhppoted 
COPY uhppoted-rest .

ENTRYPOINT /opt/uhppoted/uhppoted-rest --config /usr/local/etc/uhppoted/uhppoted.conf --console
