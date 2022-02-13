# MY SHOP
MY SHOP is a monorepo for my personal project to create shop app

## How to run each services

```SHELL
$ docker network create myshop

$ cd Docker/mariadb

$ docker-compose up -d

$ docker exec -it mariadb-db-1 mysql -h 127.0.0.1 -u root -psecret -e "CREATE DATABASE IF NOT EXISTS user_db;"

$ cd Docker/redis

$ docker-compose up -d

$ cd user

$ docker-compose up -d
```

After running all of the command above confirm if the container is running

```SHELL
$ docker ps                                                                                                   
CONTAINER ID   IMAGE                COMMAND                  CREATED          STATUS          PORTS                    NAMES
ff48b71a8c55   user_app             "./main"                 4 minutes ago    Up 4 minutes    0.0.0.0:8080->8080/tcp   user_api
27e9c17de9ee   redis:6.2.4-alpine   "docker-entrypoint.s…"   6 minutes ago    Up 6 minutes    0.0.0.0:6379->6379/tcp   redis-redisd-1
428e24764296   adminer              "entrypoint.sh docke…"   10 minutes ago   Up 10 minutes   0.0.0.0:3000->8080/tcp   mariadb-adminer-1
fb8b957b50ac   mariadb_db           "docker-entrypoint.s…"   10 minutes ago   Up 10 minutes   0.0.0.0:3306->3306/tcp   mariadb-db-1
```

## Screen shoot

![register](https://github.com/hexennacht/myshop/blob/master/screenshoot/register.png?raw=true)
![register](https://github.com/hexennacht/myshop/blob/master/screenshoot/login.png?raw=true)
![register](https://github.com/hexennacht/myshop/blob/master/screenshoot/profile.png?raw=true)