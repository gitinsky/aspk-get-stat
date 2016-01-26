package main

import (
	"github.com/davecgh/go-spew/spew"
	"os"
	"time"
)

func logInfo(trid int, sTime time.Time, msg string, params ...interface{}) {
	//timestamp spentTime peer x-real-ip method status 'request URI' message
	spew.Fprintf(os.Stderr, "Thread %d: %s %d %s\n",
		trid,
		sTime.Local().Format("2006-01-02-15-04-05.000"),
		int(time.Now().Sub(sTime).Seconds()*1000),
		spew.Sprintf(msg, params...),
	)
}
