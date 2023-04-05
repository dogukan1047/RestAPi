FROM mysql:8.0.2

ENV MYSQL_ROOT_PASSWORD=104725
ENV MYSQL_DATABASE=MENU

COPY migration.sql /docker-entrypoint-initdb.d/

EXPOSE 3306