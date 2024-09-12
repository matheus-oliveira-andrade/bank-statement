### Bank statement

A banking system project that allows for account creation, money transactions through deposits and transfers, and the generation of bank statements in PDF format. The system is composed of multiple contexts, including authentication, accounts, and bank statements

### How to run
Start all services
```bash
make up
```

Stop and remove all services
```bash
make down
```

### Key features

- Auth token generation and validation
- Account creation
- Money transactions through deposits and transfers
- Bank statements generation in PDF format

### Key technologies

- Gotenberg: Generation of statement file, converting template to PDF.
- RabbitMQ: Message broker for asynchronous communication between services.
- Postgres: Relational database for data storage.
- Testing: To ensure the validity of the code.
- JWT: Generation and validation of JWT tokens for secure authentication, using scopes.
- Event-Driven Architecture: Communication between different contexts through events, ensuring data consistency.
- Clean architecture: To ensure a well-structured and maintainable codebase.
- Go: language utilized for backend development.
- Gin: Web framework for building REST APIs.

### APIs

Generate auth token
```bash
curl --location --request POST 'http://localhost:8080/auth/v1/token'
```

Create an account
```bash
curl --location 'http://localhost:8081/account/v1/account' \
--header 'Authorization: Bearer {{TOKEN}}' \
--header 'Content-Type: application/json' \
--data '{
    "name": "Andrade",
    "document": "01234567890"
}'
```

Deposit money in an account
```bash
curl --location 'http://localhost:8081/account/v1/account/1/deposit' \
--header 'Authorization: Bearer {{TOKEN}}' \
--header 'Content-Type: application/json' \
--data '{
    "value": 15000
}'
```

Transfer money from one account to another
```bash
curl --location 'http://localhost:8081/account/v1/account/1/transfer' \
--header 'Authorization: Bearer {{TOKEN}}' \
--header 'Content-Type: application/json' \
--data '{
    "toNumber": "2",
    "value": 7500
}'
```

Trigger statement generation
```bash
curl --location --request POST 'http://localhost:8082/statement/v1/statement/1' \
--header 'Authorization: Bearer {{TOKEN}}'
```

Get statement generation result
```bash
curl --location 'http://localhost:8082/statement/v1/statement/1' \
--header 'Authorization: Bearer {{TOKEN}}'
```