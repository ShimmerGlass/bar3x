package pulse

import "C"

//export goSetVolume
func goSetVolume(v C.float) {
	lock.Lock()
	defer lock.Unlock()

	currentVolume = float64(v)
	for _, c := range watchers {
		c <- struct{}{}
	}
}
