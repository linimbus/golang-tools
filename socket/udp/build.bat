set GOOS=windows
go build .
set GOOS=linux
go build .
set GOARCH=arm64
go build -o udp_arm64

pause