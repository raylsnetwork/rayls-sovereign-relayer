FROM postgres:18-alpine
COPY --chmod=0755 postgres-init.sh /docker-entrypoint-initdb.d/init.sh
