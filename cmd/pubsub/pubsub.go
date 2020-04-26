package pubsub

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"fmt"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"log"
	"os"
	"os/user"
)

type Config struct {
	ProjectID string `json:"project_id"`
	TopicID   string `json:"topic_id"`
}

//const configPath = "/.p/config"

var (
	flags        []cli.Flag
	projectID    string
	topicID      string
	message      string
	filepath     string
	defaultPid   string
	defaultTid   string
	configGlobal *Config
	configPath   string
)

// 初始化函数---初始化命令参数
func init() {
	// 获取当前用户路径
	usr, err := user.Current()
	if err != nil {
		return
	}
	// 生成配置文件路径
	configPath = usr.HomeDir + "/.p/config"

	flags = []cli.Flag{
		&cli.StringFlag{
			Name:        "p",
			Aliases:     []string{"pid", "projectId"},
			Value:       "postmen-test",
			Usage:       "pubsub projectId",
			Destination: &projectID,
		},
		&cli.StringFlag{
			Name:        "t",
			Aliases:     []string{"tid", "topicId"},
			Value:       "sendMail",
			Usage:       "pubsub topicId",
			Destination: &topicID,
		},
		&cli.StringFlag{
			Name:        "m",
			Aliases:     []string{"message"},
			Value:       "{\"message\":\"message\"}",
			Usage:       "message for pubsub to gcp",
			Destination: &message,
		},
		&cli.StringFlag{
			Name:        "f",
			Aliases:     []string{"file"},
			Usage:       "read from `FILE` for pubsub to gcp, you could input one with message, If input both, priority of this",
			FilePath:    "",
			Destination: &filepath,
		},
		&cli.StringFlag{
			Name:        "setDefaultPid",
			Usage:       "set default projectId",
			Destination: &defaultPid,
		},
		&cli.StringFlag{
			Name:        "setDefaultTid",
			Usage:       "set default topicId",
			Destination: &defaultTid,
		},
	}

}

func Pubsub(app *cli.App) {
	app.Flags = flags

	app.Action = handle
	// 接受os.Args启动程序
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func handle(c *cli.Context) error {
	var config = &Config{}
	// 1.判断是否只是修改设置默认配置
	if defaultPid != "" || defaultTid != "" {
		if defaultPid != "" {
			err := config.setDefaultPid(defaultPid, configPath)
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Println("设置defaultPid成功")
		} else {
			err := config.setDefaultTid(defaultTid, configPath)
			if err != nil {
				fmt.Println(err)
				return err
			}
			fmt.Println("设置defaultTid成功")
		}
		return nil
	}
	// 2.load 本地配置并复制给confiGlobal
	if configGlobal == nil {
		err, f := config.loadConfig(configPath)
		defer f.Close()
		if err != nil {
			return err
		}
		configGlobal = config
	}

	// 3.获取Args[1:]
	data := c.Args().Get(0)
	//fmt.Println(filepath)
	if projectID == "postmen-test" && topicID == "sendMail" && message == "{\"message\":\"message\"}" && filepath == "" && data == "" {
		fmt.Println("请输入 -h 获取帮助")
		return nil
	}
	/**
	4. 初始化传入参数
	example 1. pubsub -pid=postmen-test -tid=sendMail -m={a:bbb}
	        2. pubsub {a:bbb}
	*/
	if filepath != "" {
		// 如果message为文件路径，需要读取出来
		//ReadFile函数会读取文件的全部内容，并将结果以[]byte类型返回
		result, err := ioutil.ReadFile(filepath)
		if err != nil {
			return cli.NewExitError("请输入正确的文件路径："+filepath, 23)
		}
		data = string(result)
	} else {
		if data == "" {
			data = message
		}
	}

	// 5. 创建空白上下文追踪
	ctx := context.Background()

	// 6. 创建pubsub client，参数优先config配置
	var pid string
	if configGlobal != nil && configGlobal.ProjectID != "" {
		pid = configGlobal.ProjectID
	} else {
		pid = projectID
	}
	client, err := pubsub.NewClient(ctx, pid)
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
		return cli.NewExitError("创建pubsub clinet失败,请确认您输入的projectId: "+pid, 22)
	}
	// 7. 创建topic，参数优先config配置
	var tid string
	if configGlobal != nil && configGlobal.ProjectID != "" {
		tid = configGlobal.TopicID
	} else {
		tid = topicID
	}
	topic := client.Topic(tid)
	defer topic.Stop()
	// 8. 发布消息
	r := topic.Publish(ctx, &pubsub.Message{Data: []byte(data)})

	id, err := r.Get(ctx)
	if err != nil {
		//log.Fatal(err)
		fmt.Println(err)
		return cli.NewExitError("pubsub message失败,请确认您的topicId: "+tid+" ,和message: "+message, 22)
	}
	fmt.Printf("Published a message with a message ID: %s\n", id)
	return nil
}

/**
设置本地config defaultPid
*/
func (config *Config) setDefaultPid(defaultPid string, configPath string) error {
	err, f := config.loadConfig(configPath)
	if err != nil {
		return err
	}
	if f == nil {
		f, err = os.OpenFile(configPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
	}
	config.ProjectID = defaultPid
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	defer f.Close()
	if err != nil {
		return err
	}
	return nil
}

/**
设置本地config defaultTid
*/
func (config *Config) setDefaultTid(defaultTid string, configPath string) error {
	err, f := config.loadConfig(configPath)
	if err != nil {
		return err
	}
	if f == nil {
		f, err = os.OpenFile(configPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err
		}
	}
	config.TopicID = defaultTid
	b, err := json.Marshal(config)
	if err != nil {
		return err
	}
	_, err = f.Write(b)
	defer f.Close()
	if err != nil {
		return err
	}
	return nil
}

/**
load 本地config 文件
*/
func (config *Config) loadConfig(configPath string) (error, *os.File) {
	var (
		f   *os.File
		err error
	)
	if !exists(configPath) {
		err, f = createFile(configPath)
	}
	if err != nil {
		return err, nil
	}
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		return err, nil
	}
	err = json.Unmarshal(data, config)
	if err != nil {
		fmt.Println(err)
		f, err = os.OpenFile(configPath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, os.ModePerm)
		if err != nil {
			return err, nil
		}
		fmt.Println("配置文件json格式错误，请确认或者删除配置文件格式，路径：" + configPath + "，已将其自动覆盖")
	}

	return nil, f
}
