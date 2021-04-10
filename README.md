# deebee-loans

This is a very simple (yet practical) example of web application leveraging DeeBee as a persistent store.

The whole state of application is held in RAM. Cyclically and during shutdown the state is persisted to disk using
DeeBee.

When application is starting, state is loaded again into memory.

## Web API

### Take Loan

```shell
curl "http://localhost:8080/take?user=john&amount=1600&term=30"
```

### Pay off the loan

```shell
curl "http://localhost:8080/pay?user=john&amount=600"
``` 