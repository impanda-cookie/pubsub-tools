package pubsub

import "os"

type Configs interface {
	loadConfig(config string) (error, *os.File)
	setDefaultTid(defaultTid string, configPath string) error
	setDefaultPid(defaultPid string, configPath string) error
}
