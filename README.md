# Registry

## Environments
* REGISTRY_ADDRESS - set this variable for registry addr
* MONGO_ADDRESS - set this variable for mongo addr
* MONGO_DB_NAME - set this variable for mongo database
* WORKERS_COLLECTION - set this variable for tasks collection
* ANALYZERS_COLLECTION - set this variable for analyzers collection
* REPORTS_COLLECTION -  set this variable for reports collection

## Build

```shell
$ make build
$ ./registry -config /path/to/config.json
```

## Docker
```shell
$ make docker
$ docker run -e REGISTRY_ADDRESS="0.0.0.0:10006" -e MONGO_ADDRESS="mongodb://admin:password@mongo" -e MONGO_DB_NAME="default" -e WORKERS_COLLECTION="tasks" -e ANALYZERS_COLLECTION="analyzers" -e REPORTS_COLLECTION="reports" registry
```
