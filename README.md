This repo is responsible for the go code which is uploaded for lambda function
to recieve events from DynamoDB streams, using an already created TRIGGER (Event-Source-Mapping)
and send those events along with event type to ElasticSearch service

The TRIGGER on DynamoDB that this lambda function is using is set at TRIM_HORIZON
ideally should be set at LATEST, but TRIM_HORIZON makes this easier to test and debug for now

Commands used

GO init

```
go mod init example.com/hello
```

1. Download build-lambda-zip (only needed for windows)

```
go.exe get -u github.com/aws/aws-lambda-go/cmd/build-lambda-zip
```

2. compile executable (for linux) (may have to go get .. some file before)

```
GOOS=linux GOARCH=amd64 go build main.go
```

use following command if using windows

```
build-lambda-zip.exe -output main.zip main
```

3. Zip built main file

```
zip main.zip main
```

4. If you have created a function before, you can just replace that zip file from aws console ui
   or follow the process of role creation, function creation and trigger creation
   (see my other repo for all the files and commands)
