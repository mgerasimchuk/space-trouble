[![Lint Golangci](https://github.com/mgerasimchuk/space-trouble/actions/workflows/lint-golangci.yml/badge.svg)](https://github.com/mgerasimchuk/space-trouble/actions/workflows/lint-golangci.yml)
[![Lint Architecture](https://github.com/mgerasimchuk/space-trouble/actions/workflows/lint-architecture.yml/badge.svg)](https://github.com/mgerasimchuk/space-trouble/actions/workflows/lint-architecture.yml)
[![Test (unit)](https://github.com/mgerasimchuk/space-trouble/actions/workflows/test-unit.yml/badge.svg)](https://github.com/mgerasimchuk/space-trouble/actions/workflows/test-unit.yml)
[![Coverage (unit)](https://github.com/mgerasimchuk/space-trouble/wiki/assets/coverage/unit/coverage.svg)](https://raw.githack.com/wiki/mgerasimchuk/space-trouble/assets/coverage/unit/coverage.html)
[![Test (integration)](https://github.com/mgerasimchuk/space-trouble/actions/workflows/test-integration.yml/badge.svg)](https://github.com/mgerasimchuk/space-trouble/actions/workflows/test-integration.yml)

# Space Trouble

Solution for the [Space trouble challenge](#space-trouble-challenge)

## Requirements

* [Docker Engine](https://docs.docker.com/engine/) (tested on 20.10.7)
* [Docker Compose](https://docs.docker.com/compose/) (tested on 1.29.2)
* [Bash](https://www.gnu.org/software/bash/) (tested on 5.1.8)
* [Make](https://www.gnu.org/software/make/) (tested on 3.81)
* Available 8080 and 5432 ports

## How to use

**Run this command:**

```
make start
```

**After:**

API doc: http://localhost:8085

Send request over curl:

```
curl --location --request POST 'localhost:8080/v1/bookings' \
--header 'Content-Type: application/json' \
--data-raw '{
  "firstName": "John",
  "lastName": "Doe",
  "gender": "male",
  "birthday": "2021-01-31",
  "launchpadId": "5e9e4501f509094ba4566f84",
  "destinationId": "5e9e3032383ecb761634e7cb",
  "launchDate": "2031-09-25"
}'
```

```
curl --location --request GET 'localhost:8080/v1/bookings'
```

```
curl --location --request DELETE 'localhost:8080/v1/bookings/10227205-3628-4c94-a070-f5731c38b3e6'
```

See data in DB:

```
docker-compose -f deployments/docker-compose/docker-compose.yaml run migrations-up bash -c 'psql -c "SELECT * FROM bookings"'
```

## Configuration

See env file: [deployments/docker-compose/.env](deployments/docker-compose/.env)

## Details of realisation

As an architectural approach was used
- [The Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)

The project layout follows
- [https://github.com/golang-standards/project-layout](https://github.com/golang-standards/project-layout)

## Space trouble challenge

Imagine it’s 2049 and you are working for a company called SpaceTrouble that sends people to different places in our
solar system. You are not the only one working in this industry. Your biggest competitor is a less known company called
SpaceX. Unfortunately you both share the same launchpads and you cannot launch your rockets from the same place on the
same day. There is a list of available launchpads and your spaceships go to places like: Mars, Moon, Pluto, Asteroid
Belt, Europa, Titan, Ganymede. Every day you change the destination for all the launchpads. Basically on every day of
the week from the same launchpad has to be a “flight” to a different place.

Information about available launchpads and upcoming SpaceX launches you can find by SpaceX
API: https://api.spacexdata.com/

Your task is to create an API that will let your consumers book tickets online.

In order to do that you have to create 2 endpoints:

1. Endpoint to book a ticket where client sends data like:

   * First Name
   * Last Name
   * Gender
   * Birthday
   * Launchpad ID
   * Destination ID
   * Launch Date
    
   You have to verify if the requested trip is possible on the day from provided launchpad ID and do not overlap with SpaceX launches, if that’s the case then your flight is cancelled.

2. Endpoint to get all created Bookings.

Extra points:

* When you use docker/docker-compose to run the project.

* When you write unit/functional tests.

* When you create an endpoint to delete booking.

Technical requirements:

* Please, use Golang and Postgres.

* Please, use github or bitbucket.

* Commit your changes often. Do not push the whole project in one commit.
