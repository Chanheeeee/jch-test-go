$ openssl req -nodes -new -x509 -keyout ca.key -out ca.crt -subj "/CN=Admission Controller Webhook Demo CA"
$ openssl genrsa -out webhook-server-tls.key 2048
$ openssl req -new -key webhook-server-tls.key -subj "/CN=mutating-test.kakaobank.svc" | openssl x509 -req -CA ca.crt -CAkey ca.key -CAcreateserial -out webhook-server-tls.crt

 kubectl -n kakaobank create secret tls webhook-certs \
    --cert "webhook-server-tls.crt" \
    --key "webhook-server-tls.key"