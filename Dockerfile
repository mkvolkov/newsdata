FROM mysql:latest

RUN chown -R mysql:root /var/lib/mysql

ENV MYSQL_DATABASE newsdata
ENV MYSQL_USER mike
ENV MYSQL_PASSWORD mikepass1
ENV MYSQL_ROOT_PASSWORD rootpass

ADD ndata.sql /etc/mysql/ndata.sql

RUN cp /etc/mysql/ndata.sql /docker-entrypoint-initdb.d

EXPOSE 3306