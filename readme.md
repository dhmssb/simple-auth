# Simple-auth

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

how to run tests:
```sh
go test -cover -coverpkg=./api/... ./tests/...
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



example how to do login:
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


#### Vulnerability Test
1. Instalation govulncheck:
```sh
go install golang.org/x/vuln/cmd/govulncheck@latest
```
2. use govulncheck:
```sh
govulncheck ./...
```

#### vulnerability test output
<img width="1466" alt="Screenshot 2023-11-01 at 22 27 04" src="https://github.com/dhmssb/simple-auth/assets/76139234/b7501389-d514-45ad-8402-a00f5aac4608">

In summary, the report highlights specific vulnerabilities in this code and its dependencies, along with information about the vulnerabilities, affected packages, fixed versions, and example traces to help locate the issues in this code. Should take an action to address these vulnerabilities, such as updating packages to versions with fixes or applying appropriate security measures.



