// Observer pattern in golang
package observer

import (
	"errors"
	"sync"
	"time"
)

var (
	errEventNotFound = errors.New("event not found")
	events           = make(map[string][]chan interface{})
	rwMutex          sync.RWMutex
)

func Subscribe(event string, outputChan chan interface{}) {
	rwMutex.Lock()
	events[event] = append(events[event], outputChan)
	rwMutex.Unlock()
}

// Stop observing the specified event on the provided output channel
func UnSubscribe(event string, outputChan chan interface{}) error {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	newArray := make([]chan interface{}, 0)
	outChans, ok := events[event]
	if !ok {
		return errEventNotFound
	}
	for _, ch := range outChans {
		if ch != outputChan {
			newArray = append(newArray, ch)
		} else {
			close(ch)
		}
	}

	events[event] = newArray
	return nil
}

// Stop observing the specified event on all channels
func UnSubscribeAll(event string) error {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	outChans, ok := events[event]
	if !ok {
		return errEventNotFound
	}

	for _, ch := range outChans {
		close(ch)
	}
	delete(events, event)

	return nil
}

func Publish(event string, data interface{}) error {
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	outChans, ok := events[event]
	if !ok {
		return errEventNotFound
	}

	// notify all through chan
	for _, outputChan := range outChans {
		outputChan <- data
	}

	return nil
}

func PublishTimeout(event string, data interface{}, timeout time.Duration) error {
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	outChans, ok := events[event]
	if !ok {
		return errEventNotFound
	}

	for _, outputChan := range outChans {
		select {
		case outputChan <- data:
		case <-time.After(timeout):
		}
	}

	return nil
}
