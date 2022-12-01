package config

import "time"

//服务用到的固定的常量都放这里

//定义了服务名称、版本、和描述
const (
	ServerName        = "swap-kit"
	ServerVersion     = "v0.0.1"
	ServerDescription = "swap-kit service for bft"
)

const ErrLogName = "error"

const (
	MaxWait           = 5 * time.Second
	RetryTimes        = 3
	PwdMinLen         = 0
	AttachPort        = "8909"
)


