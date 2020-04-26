package main

import (
	"github.com/urfave/cli/v2"
	"pubsub-tools/cmd/pubsub"
)

func main() {
	//实例化cli
	app := cli.NewApp()
	//设定名字
	app.Name = "pubsub-tools"
	app.Usage = "a pubsub for google tool"
	// 设定版本号
	app.Version = "1.0.1"

	pubsub.Pubsub(app)
}
