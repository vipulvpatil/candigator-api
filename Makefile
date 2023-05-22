protos: protos/*.proto
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative protos/*.proto

dev:
	go run .

test:
	go test ./...

coverage:
	go test -cover ./...

run:
	docker run -i -t --rm -p 9000:9000 -p 8080:8080 --env-file .env.docker candidatetrackergo

build:
	docker build -t candidatetrackergo:latest .

local-test:
	grpcurl -import-path ./protos -proto server.proto -H 'requesting_user_email: $(TEST_USER_EMAIL)' -cert certs/local/client.crt -key certs/local/client.key -cacert certs/local/CA.crt -servername candidatetracker 0.0.0.0:9000 protos.CandidateTrackerGo/CheckConnection

local-test-no-tls:
	grpcurl -import-path ./protos -proto server.proto -plaintext -H 'requesting_user_email: $(TEST_USER_EMAIL)' 0.0.0.0:9000 protos.CandidateTrackerGo/CheckConnection

remote-test:
	grpcurl -import-path ./protos -proto server.proto -H 'requesting_user_email: $(TEST_USER_EMAIL)' -cert certs/remote/client.crt -key certs/remote/client.key -cacert certs/remote/CA.crt api.candidatetracker.co:9000 protos.CandidateTrackerGo/CheckConnection

remote-test-no-tls:
	grpcurl -import-path ./protos -proto server.proto -plaintext -H 'requesting_user_email: $(TEST_USER_EMAIL)' 137.184.177.206:9000 protos.CandidateTrackerGo/CheckConnection

psql:
	psql "$(DB_URL)"

coverage-html:
	go test -v ./... -covermode=count -coverpkg=./... -coverprofile coverage/coverage.out
	go tool cover -html coverage/coverage.out -o coverage/coverage.html
	open coverage/coverage.html
