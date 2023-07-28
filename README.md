# gobank

A JSON API for a bank written in Go. 

Utilizes mostly the standard library.

In place of the old go-to gorilla/mux, I use the lightweight [chi router](https://github.com/go-chi/chi).

Authorization is accomplished using [JWT](https://github.com/golang-jwt/jwt).

Password hashing using [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt).

Database is [PostgreSQL](https://www.postgresql.org/).

Connects with the database using [pq](https://pkg.go.dev/github.com/lib/pq).

## About

This project stems from a tutorial I found from [Anthony GG](https://www.youtube.com/watch?v=pwZuNmAzaH8).

I also referenced [this youtube video](https://www.youtube.com/watch?v=p08c0-99SyU) to learn more about using docker-compose.

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
- [ ] Logout endpoint
- [ ] Method to invalidate JWT tokens
- [ ] Investigate adding chi middleware
- [ ] Investigate adding logging
- [ ] Write docs for endpoints
- [ ] Clean up comments
Future Tasks
- [ ] Transaction history table
- [ ] Testing (unit, integration, end-to-end)
- [ ] Create a client UI using Go html templates or HTMX

## Usage

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

If you want to completely remove the project including the image, containers, and volumes, run the following command:

```bash
make stop clean=true
```bash```