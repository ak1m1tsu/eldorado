FROM postgres:alpine
LABEL Name="Eldorado Database"
LABEL Version="0.1"
COPY ./migrations/*.up.sql /docker-entrypoint-initdb.d/
