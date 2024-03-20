# Creating the server and client keys and certificates for the REST server 

## Certificate Authority

```
mkdir .CA
cd .CA
openssl req -x509 -sha256 -days 1825 -newkey rsa:2048 -keyout uhppoted-rest-ca.key -out uhppoted-rest-ca.cert
```

## Server certificate and key

```
openssl req -new -nodes -out uhppoted-rest.csr -newkey rsa:2048 -keyout uhppoted-rest.key
openssl x509 -req -CA uhppoted-rest-ca.cert -CAkey uhppoted-rest-ca.key -in uhppoted-rest.csr -out uhppoted-rest.cert -days 365 -CAcreateserial -extfile domain.ext
```

## Client certificate and key

```
openssl req -new -nodes -out uhppoted-rest-client.csr -newkey rsa:2048 -keyout uhppoted-rest-client.key
openssl x509 -req -CA uhppoted-rest-ca.cert -CAkey uhppoted-rest-ca.key -in uhppoted-rest-client.csr -out uhppoted-rest-client.cert -days 365 -CAcreateserial -extfile domain.ext
```

## Notes
1. [stackoverflow: How to add subject alernative name to ssl certs?](https://stackoverflow.com/questions/8744607/how-to-add-subject-alernative-name-to-ssl-certs#8744717)
2. [Creating a Self-Signed Certificate With OpenSSL](https://www.baeldung.com/openssl-self-signed-cert)
3. [Create your own Certificate Authority (CA) using OpenSSL](https://arminreiter.com/2022/01/create-your-own-certificate-authority-ca-using-openssl)
4. [Change the trust settings of a certificate in Keychain Access on Mac](https://support.apple.com/en-ca/guide/keychain-access/kyca11871/mac)