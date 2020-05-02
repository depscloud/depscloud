```bash
$ bash gencerts.sh
$ docker-compose up -d
$ curl --cert certs/localhost.crt --key certs/localhost.key --cacert certs/ca.crt https://localhost:8080/v1alpha/sources
```
