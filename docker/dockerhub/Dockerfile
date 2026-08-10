FROM alpine

EXPOSE 8080

RUN mkdir -p /usr/local/etc/uhppoted/rest

COPY uhppoted.conf /usr/local/etc/uhppoted/

WORKDIR /opt/uhppoted 
COPY uhppoted-rest .

ENTRYPOINT /opt/uhppoted/uhppoted-rest --config /usr/local/etc/uhppoted/uhppoted.conf --console
