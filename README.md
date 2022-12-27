# otp-server

## Google OTP Verify Server

  - Mysql (GORM)
  - Sqlite (GORM)

#### build os architecture

- Mac
  - GGOS=macos GOARCH=and64 CGO_ENABLED=1 go build -o bin/app-amd64-darwin main.go
- 리눅스
  - GOOS=linux GOARCH=amd64 CGO_ENABLED=1 go build -o bin/app-amd64-linux main.go
- 윈도우
  - GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -ldflags "-s -w" -o bin/app-amd64.exe main.go
  - GOOS=windows GOARCH=386 CGO_ENABLED=1 CC=i686-w64-mingw32-gcc go build -ldflags "-s -w" -o bin/app-x86.exe main.go

sqlite3 사용시 cgo 필요하여 cgo_enable=1 해야하며
  mac 에서 빌드시 윈도우 또는 리눅스용은 gcc 컴파일러 필요하니
윈도우는 mingw 설치필요 하고
  빌드시에 CC=x86_64-w64-mingw32-gcc 옵션을 넣고 해야함.
리눅스는 gnu 설치필요?
  빌드시에 CC=x86_64-linux-gnu-gcc 옵션을 넣고 해야함
-ldflags "-s -w" 옵션으로 크기 줄이기