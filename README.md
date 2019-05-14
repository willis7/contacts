# Contacts

This trivial application explores the use of postgresql, docker and golang (templates, go modules, sqlx).

## Docker Compose

By mounting the `contacts.sql` to `/docker-entrypoint-initdb.d/contacts.sql` the container will run the script as part of the intialisation phase.

## Development
For a docker based db and local run of code:

    $ docker-compose -f dev-compose.yml up -d db
    $ go run main.go -conn "postgresql://sion:example@localhost:5432/contacts" 

For a docker only build

    $ docker-compose up --build