# go-mail-proxy
embed mail stamp wallet sdk and proxy mail to target

## local tls config
openssl req -newkey rsa:2048 -nodes -keyout key.pem -x509 -days 365 -out certificate.pem

openssl x509 -inform pem -noout -text -in qq.com.cer


go get -d github.com/btcsuite/btcd/chaincfg/chainhash@v1.0.2