all: env

run:
	cd backend; go run main.go

env:
	export $(cat .env)