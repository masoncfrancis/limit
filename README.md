# Limitr

A HTTP/HTTPS request rate limiter written in Go. 

## About Limitr

Limitr works by forwarding requests to a specified URL and moderating the rate at which requests can be made in accordance
with a defined number of requests per defined time window. If the rate limit is exceeded, the server will respond with a 
`429 Too Many Requests` status code.

### Why did I write Limitr?

I had a need for rate limiting a self-hosted API that I was developing, and couldn't find any standalone rate limiters that
weren't part of some API gateway or reverse proxy. I wanted something that I could run on my server and configure to my
liking. 

### Technologies used

- [Go](https://golang.org/): The server is written in Go
- [Go Fiber](https://gofiber.io/): The server uses Go Fiber to handle HTTP requests
- [Redis](https://redis.io/): The rate limiter uses Redis to store rate limit data
- [Docker](https://www.docker.com/): The server can be run in a Docker container

## Getting Started

### Prerequisites

#### Environment Variables

You can set the following environment variables to configure the rate limiter:

- `FORWARD_URL` (**required**): The URL to forward requests to
- `RATE_LIMIT` (**required**): The number of requests allowed per minute
- `TIME_WINDOW` (**required**): The time window in seconds
- `PORT` (default: `7654`): The port the server will listen on
- `USE_TLS` (default: `false`): Whether to use TLS
- `REDIS_ADDR` (default: `localhost:6379`): The address where the Redis server is running
- `REDIS_PASSWORD` (default: `""`): The password of the Redis server
- `REDIS_DB` (default: `0`): The database of the Redis server

You can store these variables in a `.env` file in the same directory as the executable. If there is no `.env` file, the
server will check to see if the variables are otherwise set. Variables stored in .env will take precedence over those
already set in the environment before running.

**If you are using `docker compose`**, you will need to set these variables in the `docker-compose.yml` file.

#### Redis

The rate limiter uses Redis to store the rate limit data.

The `docker-compose.yml` file in the root of the project contains a Redis service. If you are not using `docker compose`
you will need to set up a Redis server yourself. Make sure to set your environment variables appropriately (learn more
[here](#environment-variables)).

### Running Limitr

#### Using Docker Compose

First, clone the repository:

```shell
git clone https://github.com/BeehiveBroadband/limitr.git
cd limitr
```

To run the server using Docker Compose, you can use the following command:

```shell
docker-compose up
```

This will start the Limitr server and a Redis server. The server will be available at `http://localhost:7654`.

#### Using the executable

Releases are available on the [releases page](releases). Download the appropriate release for your system and run the
executable.

```shell
./[executable name]
```

The server will be available at `http://localhost:7654` unless you set a different port.

#### Building from source (unstable)

You can build from source if a release isn't available for your system. Please note that you will need to have Go >
=1.22.3
installed on your system.

First, clone the repository:

```shell
git clone https://github.com/BeehiveBroadband/limitr.git
cd limitr
```

Then, build the executable:

```shell
go build ./cmd/limitr
```

Finally, run the executable:

```shell
./limitr
```

The server will be available at `http://localhost:7654` unless you set a different port.

## Licensing

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

