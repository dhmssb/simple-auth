# Simple auth



Clone this repository first
```sh
git clone https://github.com/dhmssb/simple-auth.git
```


how to run manually:

1. please make file .env from .envexample
2. please ensure SEED on .env file was "true"
3. run 
```sh
go run main.go
```


how to run with docker:
1. run
```sh
docker-compose up
```

if you want to run docker test:
```sh
docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
```

#### Endpoint
| Method | Endpoint | Description |
| ------ | ------ | ------ |
| GET | / | PING |
| GET | /users | To fetch users who have login with bearer token |
| GET | /users/{id} | To fetch users by Id |
| POST | /users | To create new user |
| POST | /login | Login and get bearer token |
| DELETE | /users/{id} | Delete User by Id |



```sh
http:{{url}}/login
```

- POST

```sh
{
    "email": "userdua@gmail.com",
    "password": "password"
}
```