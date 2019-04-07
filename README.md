# Quizzy backend

An API server to authenticate and process user data for Quizzy frontend.

## Repos

* [Frontend](https://github.com/eadium/quizzy-frontend)
* [Backend](https://github.com/eadium/quizzy-backend)

## Running on your machine

### Requirements
* Go 1.11 or higher
* PostgreSQL 10 or higher (preferable DB)

This server is designed using microservices architecture (GRPC).

### Step 1. Creating config files
In the project dir:

```bash
touch service_main/data/cfg.json
touch service_auth/data/cfg.json
```

Example of auth service configuration file:

```json
{
    "dbconnector": "postgres",
    "connectionstring": "host=localhost port=5432 user=postgres password=passwd dbname=postgres sslmode=disable",
    "port": ":5050"
}
```

Example of main service configuration file:

```json
{
    "dbconnector": "postgres",
    "connectionstring": "host=localhost port=5432 user=postgres password=passwd dbname=postgres sslmode=disable",
    "corsallowedhost": ["http://localhost:3000"],
    "https": false,
    "port": ":8000",
    "authservicehost": "127.0.0.1:5050",
    "staticroot": "./data/static",
    "pprofenabled": false
}
```

### Step 2. Database setup
The scheme is available in this [file](service_main/data/db-init.psql).

### Step 3. Build & run
To build services:

```bash
go get
go build service_auth/*.go
go build service_main/*.go
```

Run each microservice in a seperate session. To do this, you can use `screen` utility:

```bash
cd service_auth && screen -dmS auth-quizzy ./service_auth && cd ..
cd service_main && screen -dmS main-quizzy ./service_main && cd ..
```

Public API is now available at port 8000 (you can change this value in the configuration file).