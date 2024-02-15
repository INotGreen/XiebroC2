package setchannel

func AddPtyDataChan(id string, m chan interface{}) {
	mutex.Lock()
	PtyDataChan[id] = m
	mutex.Unlock()
}

func GetPtyDataChan(id string) (m chan interface{}, exist bool) {
	mutex.Lock()
	m, exist = PtyDataChan[id]
	mutex.Unlock()
	return
}

func DeletePtyDataChan(id string) {
	mutex.Lock()
	if m, ok := PtyDataChan[id]; ok {
		close(m)
		delete(PtyDataChan, id)
	}
	mutex.Unlock()
}
