set GOOS=linux
set GOARCH=amd64
go build -o ./build/linux64/check-before-deploy-linux64 check-before-deploy.go

set GOOS=linux
set GOARCH=386
go build -o ./build/linux32/check-before-deploy-linux32 check-before-deploy.go

set GOOS=darwin
set GOARCH=amd64
go build -o ./build/osx/check-before-deploy-mac check-before-deploy.go

set GOOS=windows
set GOARCH=386
go build -o ./build/win32/check-before-deploy-win32.exe check-before-deploy.go

set GOOS=windows
set GOARCH=amd64
go build -o ./build/win64/check-before-deploy-win64.exe check-before-deploy.go