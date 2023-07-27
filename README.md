# periodic-task
A simple microservice which returns the matching timestamps of a periodic task. A periodic task is described by the following properties:
* Period (every hour, every day, ...)
* Invocation point (where inside the period should be invoked)
* Timezone (days/months/years are timezone-depended)

## Project Structure by feature (DDD)
The periodic-task project follows the Domain-Driven Design (DDD), which is an approach to software development that focuses on aligning the software design with the business domain, emphasizing the domain model as the central artifact of the system.
### cmd
This contains the entry point (main.go) files for all the services.
### period
The period is the service and contains the business logic of the application. The period is responsible for calculating the matching timestamps of a periodic task. Also, it contains the timestamp that locates the supported periods (1h, 1d, 1mo, 1y) and they are designed based on the strategy pattern. The strategy pattern is a behavioural software design pattern that enables selecting a period at runtime. It is extensible and easy to support new periods decoupling the details from the period service.
### server
It is responsible for the transport level, such as request validation, and marshalling a request into an object or a struct that a service layer can interact with.
### vendor
This directory stores all the third-party dependencies locally so that the version doesn’t mismatch late

## Building the application
Build the periodic-task application as follows:
```
go build -v ./...
```

### Docker
You can build the Docker image (latest):
```
docker build .
```

## Running the application
Start the application on default port 8181 (or whatever the `PORT` variable is set to).
```
go run cmd/periodic-task/main.go
```

### Docker
You can also run the application using Docker providing different address and port, for example:
```
docker run --name periodic-task -p 9000:9000 -e SERVER_ADDR=0.0.0.0:9000 periodic-task-api
```

You can also build and run the application defined in the `docker-compose.yaml` file with a single command:
```
docker-compose up --build
```
The configuration file supports a healthcheck that could be used to ping and verify the liveness of a DB repository.

## Test the application
To run the unit tests for the periodic-task microservice, execute the following command:
```
go test -v ./...
```

## Try it!
```
http://localhost:8181/api/v1/ptlist?period=1y&tz=Europe/Athens&t1=20180214T204603Z&t2=20211115T123456Z
```

# Contributing
Contributions are welcome! If you have any suggestions, improvements, or bug fixes, please open an issue or submit a pull request.

# Licence
This code is licensed under the *MIT License*.
