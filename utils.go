package main

import (
    "sync"
    "runtime"
)

func ternary(c bool, t interface{}, f interface{}) interface{} {
    if c {
        return t
    }
    return f
}


type SitesCounter struct {
    Total int
    Done  int
    lock  sync.Mutex
}

func (c *SitesCounter) Increment(done int) (int, int) {
    c.lock.Lock()
    defer c.lock.Unlock()
    c.Done += done
    return c.Done, c.Total
}


func GetCaller() (string, int, string) {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		return "UNKNOWN", 0, "UNKNOWN"
	}
	f := runtime.FuncForPC(pc)
	return file, line, f.Name()
}

func GetCallerFuncName() string {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return "UNKNOWN"
	}
	f := runtime.FuncForPC(pc)
	return f.Name()
}