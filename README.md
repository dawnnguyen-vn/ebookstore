GOOS=linux GOARCH=amd64 go build -o bootstrap main.go lambda_adapter.go
zip myFunction.zip bootstrap