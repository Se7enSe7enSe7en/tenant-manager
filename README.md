# Tenant Manager

app to manage tenants

# Dev setup guide
## Step 1: install task 
 ```sh
 go install github.com/go-task/task/v3/cmd/task@latest
 ```
## Step 2: use task to run the "setup" command
 ```sh
 task setup
 ```

# Dev commands

## Run the docker container for the DB

```sh
docker compose up
```

## Start the dev server

```sh
task
```

## DB migration

```sh
goose up
```

## Generate sqlc queries

```sh
sqlc generate
```
