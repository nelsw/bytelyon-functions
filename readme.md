# ByteLyon Serverless (Go)

```shell
# Login
url=https://ckkczji3hn6vnfintlkcf7b6vm0cfafl.lambda-url.us-east-1.on.aws
curl -X POST --location $url --basic --user demo@demo.com:Demo123! | jq .
```

```shell
# User todo
```

```shell
# Jobs
url=https://3bzqwrfabt3przdzbtmihmkseq0lryxo.lambda-url.us-east-1.on.aws
tkn=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7ImlkIjoiMDFLNDhQQzBCSzEzQldWMkNHV0ZQOFFRSDAifSwiaXNzIjoiQnl0ZUx5b24iLCJleHAiOjE3NTcwNjgwMTgsIm5iZiI6MTc1NzA2NjIxOCwiaWF0IjoxNzU3MDY2MjE4LCJqdGkiOiJlMDVhNDM3Ni1kNDM2LTQxNTgtODUyZC0xYTg5ZWQ0ODdjZDQifQ.yMQDhoyWmUC_sT8b9ojAznxd7YNbWlWyyn5AtBNuIwo
curl -X GET --location $url -H "authorization: Bearer $tkn" | jq
```