package setchannel

import "sync"

var mutex sync.Mutex

var PtyDataChan = make(map[string]chan interface{})
