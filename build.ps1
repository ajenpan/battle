go env -w GOOS="linux"
go build -o battle ./cmd

go env -w GOOS="windows"
go build -o battle.exe ./cmd
