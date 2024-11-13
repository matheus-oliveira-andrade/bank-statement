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

CREATE INDEX accounts_Document_idx ON accounts (Document);

CREATE TABLE IF NOT EXISTS idempotencykeys (
   Key VARCHAR(40) PRIMARY KEY,
   CreatedAt TIMESTAMP
);

CREATE DATABASE statementdb;

\c statementdb

CREATE TABLE IF NOT EXISTS accounts (   
   Number VARCHAR(15) PRIMARY KEY,
   Name VARCHAR(120),
   Document VARCHAR(14),
   Balance BIGINT   
);

CREATE TABLE IF NOT EXISTS movements (
   Id SERIAL PRIMARY KEY,
   Type VARCHAR(15),
   AccountNumber VARCHAR(15),
   Value BIGINT,
   ToAccountNumber VARCHAR(15),
   CreatedAt TIMESTAMP
);

CREATE INDEX movements_AccountNumber_idx ON movements (AccountNumber);

CREATE TABLE IF NOT EXISTS statementsgeneration (
   Id SERIAL PRIMARY KEY,
   status VARCHAR(30),
   AccountNumber VARCHAR(15),   
   CreatedAt TIMESTAMP,
   FinishedAt TIMESTAMP,
   Error VARCHAR(255),
   DocumentContent TEXT
);

CREATE INDEX statementsgeneration_AccountNumber_idx ON statementsgeneration (AccountNumber);