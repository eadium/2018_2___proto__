# Quizzy backend

## Repos

* [Frontend](https://github.com/eadium/quizzy-frontend)
* [Backend](https://github.com/eadium/quizzy-backend)

## Running on your machine

Firstly you need to get Go 1.11 or higher. PostgreSQL is preferred DBMS.

This app is using microservices architecture (via GRPC).


### Firstly, you have to create config files:

```bash
touch service_main/data/cfg.json
touch service_auth/data/cfg.json
```

Example of auth service config file:

```json
{
    "dbconnector": "postgres",
    "connectionstring": "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable",
    "port": ":5050"
}
```

Example of main service config file:
```json
{
    "dbconnector": "postgres",
    "connectionstring": "host=localhost port=5432 user=postgres dbname=postgres sslmode=disable",
    "corsallowedhost": ["http://localhost:3000"],
    "https": false,
    "port": ":8000",
    "authservicehost": "127.0.0.1:5050",
    "staticroot": "./data/static",
    "pprofenabled": false
}
```
### Then, you should create a database with determined structure.
DB architecture could be found in the [schema](service_main/data/db-init.psql) file.

For example you could use such command (where `quizzy` is the DB name):

```bash
sudo -u postgres psql quizzy < service_main/data/db-init.psql
```

### To build binaries proceed the following instructions:

```bash
go get
go build service_auth/*.go
go build service_main/*.go

```

Because of the microservices architecture we need to run two binaries.
For this purpose we can use `screen`

```bash
cd service_auth && screen -dmS auth-quizzy ./service_auth && cd ..
cd service_main && screen -dmS main-quizzy ./service_main && cd ..
```

That's all! API is now available at port 8000 (you can change it in the config).
