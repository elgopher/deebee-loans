# deebee-loans

This is a very simple (yet practical) example of web application
leveraging [DeeBee](https://github.com/jacekolszak/deebee) as a persistent store.

The whole state of application is held in RAM. Cyclically and during application shutdown the state is persisted to disk
using DeeBee store.

When application is starting, state is restored and loaded again into memory.

The state is also replicated once per hour to a second directory (which can be an NFS file system).

## Disclaimer

Please note that application is lacking in multiple areas: security, input validation, testing etc. I wanted to focus
only on those topics related to storing application state.

## Web API

### Take Loan

```shell
curl "http://localhost:8080/take?user=john&amount=1600&term=30"
```

### Pay off the loan

```shell
curl "http://localhost:8080/pay?user=john&amount=600"
``` 
