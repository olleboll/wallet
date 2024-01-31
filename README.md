## Wallet service

This wallet service allows you to create wallets and update their balance.

For a complete test run just:
```
make test
```

This will build the image, start database and api in docker-compose and run a test with a happy path of all functionality.

#### Endpoints:

Read balance on wallet
```
GET localhost:8080/wallet/<wallet_id>
```

Add balance on wallet
```
POST localhost:8080/wallet/<wallet_id>/add
{
    "amount": number
}
```

Subtract balance on wallet
```
POST localhost:8080/wallet/<wallet_id>/subtract
{
    "amount": number
}
```

Add and subtract endpoints will create the wallet if it does not exist, using the wallet_id provided as param.

