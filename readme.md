# ByteLyon Serverless (Go)
[![license](http://img.shields.io/badge/license-MIT-red.svg?style=flat)](https://raw.githubusercontent.com/rs/zerolog/master/LICENSE)
[![Go Coverage](http://img.shields.io/badge/coverage-71.9%25-olive.svg?style=flat)](https://raw.githack.com/wiki/rs/zerolog/coverage.html)

[//]: # ([![Build Status]&#40;https://github.com/rs/zerolog/actions/workflows/test.yml/badge.svg&#41;]&#40;https://github.com/rs/zerolog/actions/workflows/test.yml&#41; )
***
### API Endpoint Example
```shell
# Login
url=https://ckkczji3hn6vnfintlkcf7b6vm0cfafl.lambda-url.us-east-1.on.aws
curl -X POST --location $url --basic --user demo@demo.com:Demo123! | jq .
```

```shell
# Jobs (+ Work)
url=https://3bzqwrfabt3przdzbtmihmkseq0lryxo.lambda-url.us-east-1.on.aws
tkn=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7ImlkIjoiMDFLNDhQQzBCSzEzQldWMkNHV0ZQOFFRSDAifSwiaXNzIjoiQnl0ZUx5b24iLCJleHAiOjE3NTc0NDAyNzQsIm5iZiI6MTc1NzQzODQ3NCwiaWF0IjoxNzU3NDM4NDc0LCJqdGkiOiIzNmVhNjY2Yy1kMDAxLTQ4NzktYjViNC1mMDkyOWFlMGZkN2UifQ.045q-qrW8yhNjTPEnikV_LMSfNFKn657TiC8ZjiXoQQ
curl -X GET --location $url -H "authorization: Bearer $tkn" | jq
```
***
### ToDo
- [ ] Add User API Endpoint
- [ ] Add GitHub Actions for test badge
- [ ] Add Makefile cmd for event bridge cron
- [ ] Add ƒ for working jobs
- [ ] Add python selenium ƒ (sep repo)