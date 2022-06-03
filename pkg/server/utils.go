package server

import "time"

func datetime() string {
	return time.Now().Format("02/Jan/2006:15:04:05 -0700")
}
