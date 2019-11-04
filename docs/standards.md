# standards #

Project and go specific standards that have been chosen and not otherwise
explicitly recommended.

## project layout ##

This project follows the conventions established by the most popular open source
projects and go developers, to provide a reasonable layout and structure for
your project.

See golang-standards [project layout](https://github.com/golang-standards/project-layout)

## best practices ##

Where possible, follow established best practices.  See:

  1. [effective go](https://golang.org/doc/effective_go.html).
  2. [The Uber Go Style Guide](https://github.com/uber-go/guide)

## test-driven software development ##

Generally use [Robert Martin's](https://en.wikipedia.org/wiki/Robert_C._Martin) [three laws of TDD](https://www.youtube.com/watch?v=qkblc5WRn-U):

  1. Only write enough of a unit test to fail.
  2. Only write production code to make a failing unit test pass.

## dependencies ##

Where possible, minimize the number of package dependencies used to get the job
done.  If it's a complex task and someone else has an appropriately licensed
solution which is carefully maintained and largely established, then it's
a good idea.  Example dependencies that are used:

zzz

  1. [Gorilla mux](https://github.com/gorilla/mux) for routing
  2. [gocql](https://github.com/gocql/gocql) for a ScyllaDB/Cassandra client
  3. [cobra commander](https://github.com/spf13/cobra) for CLI interaction
  4. [go-redis](https://github.com/go-redis/redis) for a Redis client
  5. [taskq](https://github.com/vmihailenco/taskq/) for a redis task queue

Consider using convenience libraries that have no dependencies themselves, which
further help ensure high standards.  We can include them directly, if licensed
appropriately, with attribution:

  1. [config](https://github.com/JeremyLoy/config) for easy 12FA configuration

## SOLID KISS ##

Use [SOLID](https://en.wikipedia.org/wiki/SOLID) design principles and [keep it stupid simple](https://en.wikipedia.org/wiki/KISS_principle).

See Dave Cheney's [SOLID Go Design](https://dave.cheney.net/2016/08/20/solid-go-design)

## 12Factor ##

Generally heroku's [The Twelve-Factor App](https://12factor.net/) guide should
be followed.
