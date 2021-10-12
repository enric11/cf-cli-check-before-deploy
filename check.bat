set GOOS=windows
set GOARCH=amd64
go build -o ./build/win64/check-before-deploy-win64.exe check-before-deploy.go
cf uninstall-plugin check-before-deploy
cf install-plugin ./build/win64/check-before-deploy-win64.exe -f
cf check-before-deploy -file mta.yaml -all
