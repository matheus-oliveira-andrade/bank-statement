CREATE DATABASE accountdb;

\c accountdb

CREATE TABLE IF NOT EXISTS accounts (
   Id SERIAL PRIMARY KEY,
   Number VARCHAR(15),
   Name VARCHAR(120),
   Document VARCHAR(14),
   Balance BIGINT,
   CreatedAt TIMESTAMP,
   UpdatedAt TIMESTAMP
);

CREATE DATABASE statementdb;

\c statementdb

CREATE TABLE IF NOT EXISTS accounts (   
   Number VARCHAR(15) PRIMARY KEY,
   Name VARCHAR(120),
   Document VARCHAR(14),
   Balance BIGINT   
);