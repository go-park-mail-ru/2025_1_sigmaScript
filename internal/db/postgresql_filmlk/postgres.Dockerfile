# FROM postgres:latest
FROM pgvector/pgvector:pg16

RUN apt-get update && apt-get install -y pgbadger
