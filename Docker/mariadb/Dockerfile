FROM mariadb:10.2

env MARIADB_ROOT_PASSWORD=secret
env MARIADB_USER=taufikrf
env MARIADB_PASSWORD=secure
env MARIADB_ALLOW_EMPTY_ROOT_PASSWORD=yes

RUN apt update
RUN apt -y install mariadb-plugin-rocksdb
RUN rm -rf /var/cache/apt/lists/*