# pubsub-tools

### overview
* gcp pubsub的工具

### guideline
##### Initialize Local Env
* [install gcloud](https://cloud.google.com/sdk/docs/quickstart-macos)
* [gcp init](https://cloud.google.com/sdk/docs/quickstart-macos#initialize_the_sdk)
* initialize gcp creadentials
```shell
gcloud auth application-default login
```

##### Cross compiler
###### MacOS
```shell script
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o p
```
###### Linux
```shell script
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o p
```
###### Windows
```shell script
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o p
```

###### Usage
* help
```shell script
1. p -h
2. p --help

NAME:
   pubsub-tools - a pubsub for google tool

USAGE:
   pub [global options] command [command options] [arguments...]

VERSION:
   1.0.1

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -p value, --pid value, --projectId value  pubsub projectId (default: "postmen-test")
   -t value, --tid value, --topicId value    pubsub topicId (default: "sendMail")
   -m value, --message value                 message for pubsub to gcp (default: "{\"message\":\"message\"}")
   -f FILE, --file FILE                      read from FILE for pubsub to gcp, you could input one with message, If input both, priority of this
   --setDefaultPid value                     set default projectId
   --setDefaultTid value                     set default topicId
   --help, -h                                show help (default: false)
   --version, -v                             print the version (default: false)
```
* publish
支持3种方式
```shell script
1. p "{\"data\":\"message\"}"
2. p -f ~/message.json
3. p -m "{\"data\":\"message\"}"
```
* set default config
```shell script
1. p --setDefaultTid postmen-tes
2. p --setDefaultPid sendMail
```
* config path
    * default  `~/.p/config`
* 默认会优先从config取配置(建议先配置好config)