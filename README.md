## How to deploy
```
cd opds
```
```
GOOS=linux GOARCH=amd64 go build -o bootstrap main.go lambda_adapter.go
```
```
zip myFunction.zip bootstrap
```

=> upload zip to lambda function

## How to develop?
```
cd opds
```
```
air
```
==> The server will start on port 1323