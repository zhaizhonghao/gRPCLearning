# Summary
# Private files :  ca.key, server.key, server.pem, server.crt
# "Share" files : ca.crt (needed by the client), server.csr (needed by the CA)

# change these CN's to match your hosts in your environment if needed
SERVER_CN = localhost

#step 1 : Generate Certificate Authority + Trust Certificate (ca.crt)
openssl genrsa -passout pass:1111 -des3 -out ca.key 4096
openssl req -passin pass:1111 -new -x509 -days 365 -key ca.key -out ca.crt -subj "/CN=${SERVER_CN}"

# step 2 : Generate the Server private key (server.key)
openssl genrsa -passout pass:1111 -des3 -out server.key 4096

# step 3 : Get a certificate signing request from 
openssl req -passin pass:1111 -new -key server.key -out server.csr -subj "/CN=${SERVER_CN}"

# step 4 : sign the certificate with the CA we created (it's called self signing) - server.crt
openssl x509 -req -passin pass:1111 -days 365 -in server.csr -CA ca.crt -CAkey ca.key -set_serial 01 -out server.crt

# step 5 : Convert the server certificate to .pem format (server.pem) -usable by gRPC
openssl pkcs8 -topk8 -nocrypt -passin pass:1111 -in server.key -out server.pem
