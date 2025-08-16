include .env

name :=

create:
	@go run main.go -role=${ROLE} -action=create -name=${name}

update:
	@go run main.go -role=${ROLE} -action=update -name=${name}

publish:
	@go run main.go -role=${ROLE} -action=publish -name=${name}

delete:
	@go run main.go -role=${ROLE} -action=delete -name=${name}