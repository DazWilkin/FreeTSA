# Golang "SDK" for FreeTSA using Digitorus' RFC 3161 timestamp SDK

**NB** There is no license published with github.com/digitorus/timestamp; I'm assuming it's permissible to use (checking)

Adventures in [trusted timestamping](https://en.wikipedia.org/wiki/Trusted_timestamping) with [freeTSA.org](https://freetsa.org).

A very basic Golang (almost not worthy of being called an) SDK to submit Timestamp Requests to FreeTSA using [RFC 3161: Internet X.509 PKI Time-Stamp Protocol (TSP)](https://www.ietf.org/rfc/rfc3161.txt)


```bash
# The message to Timestamp
MESG="Frederik Jack is a bubbly Border Collie"

# Filename for Timestamp Request (/tmp/${FILE}.tsq) and Response (/tmp/${FILE}.tsr)
FILE="file"

# Create a Timestamp Request, persist it, submit it to FreeTSA and persist the Response
go run github.com/DazWilkin/FreeTSA \
  --mesg=${MESG} \
  --file=${FILE}

# Download the FreeTSA certs to /tmp
wget --output-document=/tmp/tsa.crt https://freetsa.org/files/tsa.crt
wget --output-document=/tmp/cacert.pem https://freetsa.org/files/cacert.pem

# Use OpenSSL Timestamp tool to verify the Request|Response using FreeTSA's certs
openssl ts -verify \
  -in /tmp/${FILE}.tsr \
  -queryfile /tmp/${FILE}.tsq \
  -CAfile /tmp/cacert.pem \
  -untrusted /tmp/tsa.crt

# Tidy
rm /tmp/${FILE}.ts? /tmp/cacert.pem /tmp/tsa.crt
```

should (!) yield:

```
Verification: OK
```
