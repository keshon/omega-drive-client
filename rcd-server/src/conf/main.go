package conf

import "encoding/base64"

// Config
const (
	SkipPIDValidation   = false                    // set true to skip verifying if server cli is run under the parent with exact PID
	CheckParentInterval = "10s"                    // specify interval for checking parent PID. Examples are (without quotes): '2h','5m','14s'
	RcUsername          = "admin4eg"               // specify username for rclone rcd server authentication
	RcPassword          = "oKAoqAqHe4"             // specify password for rclone rcd server authentication
	RcHost              = "http://localhost:5579/" // specify hostname for rclone rcd
)

var (
	RcAuthEncoded = base64.StdEncoding.EncodeToString([]byte(RcUsername + ":" + RcPassword))
)
