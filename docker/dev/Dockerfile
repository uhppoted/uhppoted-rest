FROM alpine

EXPOSE 8080

RUN mkdir -p /usr/local/etc/uhppoted/rest

COPY uhppoted.conf /usr/local/etc/uhppoted/
COPY ca.cert       /usr/local/etc/uhppoted/rest/
COPY uhppoted.cert /usr/local/etc/uhppoted/rest/
COPY uhppoted.key  /usr/local/etc/uhppoted/rest/

WORKDIR /opt/uhppoted 
COPY uhppoted-rest .

ENTRYPOINT /opt/uhppoted/uhppoted-rest --debug --config /usr/local/etc/uhppoted/uhppoted.conf --console
