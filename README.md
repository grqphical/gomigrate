# gomigrate - A simple SQLite3 migration tool

I built this to use in my own projects so I thought I would share it for everyone to use.

## Functional Goals

1. Create and manage database migrations
    - Create migrations up/down files
    - Apply/rollback migrations
    - Show user all applied migrations

## Installation

```bash
$ go install github.com/grqphical/gomigrate@latest
```
## Basic Usage

To start make sure you have an environment variable or .env variable with `DATABASE_URL` defined.

Then run `gomigrate init` to create the migrations directory and setup the migrations table in the database

Next, run `gomigrate create NAME` to create a new migration

Then finally, run `gomigrate up` to apply the migration to your database

### Rolling Back

If you wish to rollback the database to a clean slate, use `gomigrate down`

## License

gomigrate is released under the MIT License
