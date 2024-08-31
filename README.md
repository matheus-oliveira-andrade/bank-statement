### Bank statement

A banking system project that allows for account creation, money transactions through deposits and transfers, and the generation of bank statements in PDF format. The system is composed of multiple contexts, including authentication, accounts, and bank statements.

### Key technologies

- JWT: Generation and validation of JWT tokens for secure authentication, using scopes.
- RabbitMQ: Message broker for asynchronous communication between services.
- Postgres: Relational database for data storage.
- Testing: To ensure the validity of the code.
- Event-Driven Architecture: Communication between different contexts through events, ensuring data consistency.

### How to run
```bash
# run all services using docker 
make up
```

```bash
# stop and remove all services
make down
```