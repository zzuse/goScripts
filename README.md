# goScripts
my test go scripts  
cross compile  
GOOS=linux GOARCH=amd64 go build -o multysftp multysftp.go  
GOOS=linux GOARCH=amd64 go build -o multissh multyssh.go
GOOS=windows GOARCH=amd64 go build -o hunanLicDownload.exe licDownload.go  
GOOS=windows GOARCH=386 go build -o multyssh.exe multyssh.go   
