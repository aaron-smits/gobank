# gobank

A JSON API for a bank written in Go. 

I tried to mostly use the Go standard library.

For a router, I use [chi](https://github.com/go-chi/chi).

Authorization is accomplished using [JWT](https://github.com/golang-jwt/jwt).

Password hashing using [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt).

Database is [PostgreSQL](https://www.postgresql.org/).

Connects with the database using [pq](https://pkg.go.dev/github.com/lib/pq).

The project is set up to run several docker containers using [docker-compose](https://docs.docker.com/compose/).


I also learned how to use [make](https://www.gnu.org/software/make/) to automate the build process.


## About

This project stems from a video I found from [Anthony GG](https://www.youtube.com/watch?v=pwZuNmAzaH8).

I also referenced [this video](https://www.youtube.com/watch?v=p08c0-99SyU) to learn more about using docker-compose.

My main goal with this project is the build a basic JSON API with Go to understand how developing web backends in Go works and the idioms Go developers use.

On the to do list:
- [x] Transfer endpoint
  - [x] Implement storage method for transfer
  - [X] Implement add balance method
  - [X] Implement subtract balance method
  - [X] Implement seeding method for balance
  - [X] Implement transfer endpoint
- [x] Error handling enhancements
- [x] Add ability for admins to update accounts 
- [ ] Investigate adding chi middleware
- [ ] Investigate adding logging
- [ ] Write docs for endpoints
  - [ ] Open API library 
- [ ] Clean up comments
Future Tasks
- [ ] Transaction history table
- [ ] Testing (unit, integration, end-to-end)
- [ ] Create a client UI using Go html templates or HTMX

## Usage
Use the template env file to pass in all the required fields for your database

To run the server, run the following command:

```bash
make run
```

This calls the docker compose command to compile the Go code, run the server and the database.

To stop the server, run the following command:

```bash
make stop
```

View the logs of the server, run the following command:

```bash
docker compose logs 
# or 
docker compose logs -f 
# to follow the logs (attach terminal to logs)
```

If you want to completely remove the project including the gobank image, containers, and volumes, run the following command:

```bash
make stop clean=true
```
