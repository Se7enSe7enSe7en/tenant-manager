## Step 1: install docker and run docker

- go to docker website and install docker: https://www.docker.com/

## Step 2: run postgres image

### specific command

```sh
docker run --name some-postgres -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres
```

### formula command

```sh
docker run --name <container-name> --env <env-variable> --publish <exposed-port>:<internal-port> --detach <name-of-image>
```

### check running docker containers

```sh
docker ps -a
```

## Step 3: Connect to psql in the docker container

```sh
docker exec -it some-postgres psql -U postgres
```

### command breakdown

```
docker exec -> execute a command in the container
-i -> interactive mode
-t -> tty, connect shell
some-postgres -> name of the target container
psql -> the command you want to run in the container
-U -> psql needs to create a user first
postgres -> name of the user
```

### Commands in psql used

```sql
\dt -- -> check all tables
CREATE TABLE example(data char); -- -> SQL code to create a table called exmaple with a single column called data of type char

```

### use pgAdmin to see the DB with a GUI

- download pgAdmin from website: https://www.pgadmin.org/download/

#### or use a containerized pgAdmin (note: this didn't work for me, just went with downloading pgAdmin)

```sh
docker run --name pgadmin -p 80:80 \
    -e 'PGADMIN_DEFAULT_EMAIL=user@domain.com' \
    -e 'PGADMIN_DEFAULT_PASSWORD=pgadmin' \
    -d dpage/pgadmin4
```

## Reset the entire container

### stop the running container

```sh
docker stop some-postgres
```

### remove the container

```sh
docker rm some-postgres
```

## Reference

- caleb curry docker + postgres setup tutorial: https://www.youtube.com/watch?v=Hs9Fh1fr5s8
- pgAdmin container deployment docs: https://www.pgadmin.org/docs/pgadmin4/latest/container_deployment.html
