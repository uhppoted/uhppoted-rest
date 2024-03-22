FROM alpine

EXPOSE 8080

RUN mkdir -p /etc/uhppoted
RUN mkdir -p /etc/uhppoted/rest
RUN mkdir -p /var/uhppoted

COPY uhppoted.conf /etc/uhppoted/
COPY ca.cert       /etc/uhppoted/rest/
COPY uhppoted.cert /etc/uhppoted/rest/
COPY uhppoted.key  /etc/uhppoted/rest/

WORKDIR /opt/uhppoted 
COPY uhppoted-rest .

ENTRYPOINT /opt/uhppoted/uhppoted-rest --debug --console
