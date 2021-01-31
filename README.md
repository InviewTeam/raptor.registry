## Registry

### Usage
```shell
$ make build
$ ./registry -config /path/to/config.json
```

#### Config example
```json
{
  "registry_address": "0.0.0.0:1337",
  "database":
  {
    "address": "mongodb://admin:password@mongo",
    "database": "default",
    "workers_collection": "tasks",
    "analyzers_collection": "analyzers",
    "reports_collection": "reports"
  },
  "rabbit":
  {
    "address": "amqp://rabbitmq:5672",
    "worker_queue": "workers",
    "analyzer_queue": "analyzers"
  }
}
```
