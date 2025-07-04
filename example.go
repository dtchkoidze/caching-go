package cache

import (
	"log"
	"sync"
	"time"
)

func run() {
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go cookRice(wg)
	go cookCurry(wg)
	wg.Wait()
}

func cookCurry(wg *sync.WaitGroup) {
	defer wg.Done()
	log.Println("cooking curry...")
}

func cookRice(wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Second * 7)
	log.Println("cooking rice...")
	time.Sleep(time.Second * 2)
}
