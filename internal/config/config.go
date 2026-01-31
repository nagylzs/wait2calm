package config

import (
	"time"
)

type Opts struct {
	Verbose             bool          `short:"v" long:"verbose" description:"Show verbose information"`
	Debug               bool          `short:"d" long:"debug" description:"Show debug information"`
	ShowVersion         bool          `long:"version" description:"Show version information and exit"`
	WaitBefore          time.Duration `short:"b" long:"random-wait-before" description:"Max. random wait before first measurement, default is 15s" default:"10s"`
	WaitAfter           time.Duration `short:"a" long:"random-wait-after" description:"Max. random wait after first measurement, default is 15s" default:"10s"`
	LoadType            uint          `short:"t" long:"load-type" description:"Load average type, can be 1=load1,5=load5,15=load15" default:"1"`
	MeasurementInterval time.Duration `short:"m" long:"measurement-interval" description:"Interval between measurements, default is 10s" default:"10s"`
	ImmediateStartBelow float64       `short:"i" long:"immediate-start-below" description:"Immediately start below this load value"`
	DelayedStartBelow   float64       `short:"l" long:"delayed-start-below" description:"Start if two subsequent loads are below this value"`
	DoNotWaitAfter      time.Duration `long:"do-not-wait-after" description:"Do not wait inside the loop after this amount of time, regardless of the load average. A non-positive value disables this function. Resolution is 100msec."`
	SuccessOnTimeout    bool          `long:"success-on-timeout" description:"Return with zero exit code, even if timed out on --do-not-wait-after"`
}
