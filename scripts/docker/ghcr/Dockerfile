FROM alpine

EXPOSE 8080

RUN mkdir -p /opt/uhppoted
RUN mkdir -p /var/uhppoted/etc
RUN mkdir -p /var/uhppoted/var
RUN mkdir -p /var/uhppoted/rest

COPY uhppoted.conf /var/uhppoted/etc/
COPY ca.cert       /var/uhppoted/rest/
COPY uhppoted.cert /var/uhppoted/rest/
COPY uhppoted.key  /var/uhppoted/rest/

WORKDIR /opt/uhppoted 
COPY uhppoted-rest .

ENTRYPOINT /opt/uhppoted/uhppoted-rest --config /var/uhppoted/etc/uhppoted.conf --console
