FROM alpine

EXPOSE 8080

RUN mkdir -p /opt/uhppoted
RUN mkdir -p /var/uhppoted/etc
RUN mkdir -p /var/uhppoted/var
RUN mkdir -p /var/uhppoted/rest

COPY uhppoted.conf /var/uhppoted/etc/

WORKDIR /opt/uhppoted 
COPY uhppoted-rest .

ENTRYPOINT /opt/uhppoted/uhppoted-rest --config /var/uhppoted/etc/uhppoted.conf --console
