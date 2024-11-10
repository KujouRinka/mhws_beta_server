openssl genrsa -out website.key 2048
openssl req -new -key website.key -out website.csr -config openssl.cnf
openssl x509 -req -days 365 -in website.csr -signkey website.key -out website.crt -extfile openssl.cnf -extensions req_ext