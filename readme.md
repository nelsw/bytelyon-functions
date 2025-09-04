# ByteLyon Serverless (Go)

```shell
# Define request variables for clarity
url=https://ckkczji3hn6vnfintlkcf7b6vm0cfafl.lambda-url.us-east-1.on.aws
username=demo@demo.com
password=Demo123!

# Send request and assign response to variable
res=$(curl -X POST --location $url --basic --user $username:$password)

# Print the JWT to use in the Authorization header
echo $res | jq .

# Print decoded claims to use data.id as user ID
jq -R 'split(".") | .[1] | @base64d | fromjson' <<< $res
```

```shell
# Jobs
tkn=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7ImlkIjoiMDFLNDhQQzBCSzEzQldWMkNHV0ZQOFFRSDAifSwiaXNzIjoiQnl0ZUx5b24iLCJleHAiOjE3NTcwMTgxNDMsIm5iZiI6MTc1NzAxNjM0MywiaWF0IjoxNzU3MDE2MzQzLCJqdGkiOiIwYWVlNDdjMy03YTQ2LTRjYmQtYTdhYy1jNzQ2NjBmODg0MjQifQ.04abFJOZf-qB1C-C2y7Pjj4c2krkAyxCDZy7SK7p3Y4
url="https://fkarinkfb33afcsdz2uutdcvfe0bmdee.lambda-url.us-east-1.on.aws/user/01K48PC0BK13BWV2CGWFP8QQH0/job"
curl $url -H "authorization: Bearer $tkn"
```

```shell
# Work
curl "https://fkarinkfb33afcsdz2uutdcvfe0bmdee.lambda-url.us-east-1.on.aws/user/01K48PC0BK13BWV2CGWFP8QQH0" \
    -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJkYXRhIjp7ImlkIjoiMDFLNDhQQzBCSzEzQldWMkNHV0ZQOFFRSDAifSwiaXNzIjoiQnl0ZUx5b24iLCJleHAiOjE3NTcwMTUxOTMsIm5iZiI6MTc1NzAxMzM5MywiaWF0IjoxNzU3MDEzMzkzLCJqdGkiOiIwNDg2MDAyOS0wZWVlLTRlOTgtOWU5MS0wN2VkNzljNGJkYzIifQ.Z5LtbK6d34oIxpQ9byfYv5DrbrOD9dCSRbqFs6ih7T4"
```