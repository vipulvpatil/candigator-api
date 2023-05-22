# Candidate Tracker Go

## Setup

### Prerequisites

* Install and run Postgresql. This is the primary database.
* Install and run redis. This is the backing database for running jobs using gocraft/work.
* Install and run Docker. Not needed for development but is used in production and is useful for creating a production-like setup locally.

### Env vars

```
export ENABLE_TLS=true                               # .env and .envrc
export CA_CERT_BASE64=LS0tLS1CRUdJTiBDRVJ...         # .env and .envrc
export SERVER_CERT_BASE64=LS0tLS1CRUdJTiB...         # .env and .envrc
export SERVER_KEY_BASE64=LS0tLS1CRUdJTiBS...         # .env and .envrc
export REDIS_URL=redis://user....                    # .env and .envrc
export TEST_DB_URL="user=some_user host=localhost port=5432 dbname=some_test_db sslmode=disable"            # .envrc
export TEST_USER_EMAIL="some_test_user_email"       # .envrc
```
## Commands

### To run server without docker

```
go run .
```

### To rebuild server with docker.

```
make build
```

### To run/stop server with docker.

```
make run
```

### To re/build proto definitions

```
make protos
```

### To test locally

```
make local-test
```

### It is possible to run server without TLS using ENV var. To test the same
```
make local-test-no-tls
```
