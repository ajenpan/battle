go env -w GOOS="linux"
go build -o battlefield ./cmd

go env -w GOOS="windows"
go build -o battlefield.exe ./cmd
