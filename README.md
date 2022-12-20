# otp-server

## Google OTP Verify Server

  - Mysql (GORM)
  - Sqlite (GORM)

#### build os architecture

- GGOS=macos GOARCH=and64 go build -o bin/app-amd64-darwin main.go
- GOOS=linux GOARCH=amd64 go build -o bin/app-amd64-linux main.go
- GOOS=windows GOARCH=amd64 go build -o bin/app-amd64.exe main.go
- GOOS=windows GOARCH=386 go build -o bin/app-x86.exe main.go