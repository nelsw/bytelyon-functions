include .env

.PHONY: test

test:
	@printf "➜  %s  %s [\033[35m%s\033[0m]\n➜  " "📊" "test" "./..."
	@go test -coverprofile cp.out ./...
	@go tool cover -func=cp.out
	@go tool cover -html=cp.out

create: build
	@printf "➜  %s  %s [\033[35m./handler/%s\033[0m]" "💽" "create" ${name}
	@aws lambda create-function \
		--function-name bytelyon-${name} \
		--runtime "provided.al2" \
		--role ${ROLE} \
		--architectures arm64 \
		--handler "bootstrap" \
		--zip-file "fileb://./main.zip" \
		--memory-size "512" \
		--timeout "30" \
		--publish \
		--environment "Variables={$(shell tr '\n' ',' < ./handler/${name}/.env)}" > /dev/null
	@printf "  ✅\n"

update: build
	@printf "➜  %s  %s [\033[35m./handler/%s\033[0m]" "💾" "update" ${name}
	@aws lambda update-function-configuration \
    		--function-name bytelyon-${name} \
    		--role ${ROLE} \
    		--timeout "30" \
    		--memory-size "512" \
    		--environment "Variables={$(shell tr '\n' ',' < ./handler/${name}/.env)}" > /dev/null
	@aws lambda update-function-code --zip-file fileb://./main.zip --function-name bytelyon-${name} > /dev/null
	@printf "  ✅\n"

delete:
	@printf "➜  %s  %s [\033[35m./handler/%s\033[0m]" "🗑️" "delete" ${name}
	@aws lambda delete-function --function-name bytelyon-${name} | jq
	@printf "  ✅\n"

build:
	@printf "➜  %s  %s [\033[35m./handler/%s\033[0m]" "🛠" "build" ${name}
	@GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap ./handler/${name}/main.go
	@zip -r9 main.zip bootstrap > /dev/null
	@printf "  ✅\n"

clean:
	@printf "➜  %s  %s [\033[35m%s\033[0m]" "🧽" "clean" "*"
	@rm -f ./bootstrap
	@rm -f ./main.zip
	@rm -f ./cp.out
	@printf "  ✅\n"

list:
	@printf "➜  %s  %s [\033[35m./aws/ꟛ/%s\033[0m]" "📋" "list"
	@aws lambda list-functions --no-paginate \
	| jq '.Functions.[] | {name: .FunctionName, updated: .LastModified, environment: .Environment.Variables}'

url:
	@printf "➜  %s  %s [\033[35m./handler/%s\033[0m]\n" "🛜" "url" ${name}
	@aws lambda get-function-url-config --function-name bytelyon-${name} | jq '.FunctionUrl'

publish:
	@printf "➜  %s  %s [\033[35m./handler/%s\033[0m]" "🌐" "publish" ${name}
	@aws lambda create-function-url-config --function-name bytelyon-${name} --auth-type NONE > /dev/null
	@aws lambda add-permission \
    		--function-name bytelyon-${name} \
    		--action lambda:InvokeFunctionUrl \
    		--principal "*" \
    		--statement-id FunctionURLAllowPublicAccess \
    		--function-url-auth-type NONE > /dev/null
	@printf "  ✅\n"
	@make url ${name}

unpublish:
	@printf "➜  %s  %s [\033[35m./handler/%s\033[0m]" "⛔️" "unpublish" ${name}
	@aws lambda remove-permission --function-name bytelyon-${name} --statement-id FunctionURLAllowPublicAccess
	@aws lambda delete-function-url-config --function-name bytelyon-${name}
	@printf "  ✅\n"