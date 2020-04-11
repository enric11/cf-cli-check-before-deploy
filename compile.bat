set GOOS=linux
set GOARCH=amd64
go build -o ./build/linux64/check-before-deploy-linux check-before-deploy.go

set GOOS=darwin
set GOARCH=amd64
go build -o ./build/osx/check-before-deploy-mac check-before-deploy.go

set GOOS=windows
set GOARCH=amd64
go build -o ./build/win64/check-before-deploy-win.exe check-before-deploy.go