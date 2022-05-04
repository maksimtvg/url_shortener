# URL Shortener GRPC
### Another one url shortener
#### Golang 1.18, Docker, Postgres:14

How to start the app

At first start run only postgres.
The command builds the postgres and initiates init.sql from ``./persistence/init/`` directory.

```make db_up```

To start server execute

```make service_up```

To run client execute. Change const goroutineAmount to run concurrent request as many as you want

```make run_client```

To check out the documentation execute the command and visit http://localhost:6060/pkg/url_shortener/?m=all

```godoc -http:6060```

To check data in Postgres run
```docker exec -it pgdb bash```, ```psql -U backend -t urls```, ```select * from urls;```

To run tests execute
```make test_service```
