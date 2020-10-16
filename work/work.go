package work

import (
	"sync"
	"time"
)

func WorkFirst(firstCh, secondCh chan *Page, wg *sync.WaitGroup) {
	defer wg.Done()
CHCLOSED:
	for {
		select {
		case _, exist := <-firstCh:
			if !exist {
				break CHCLOSED
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func WorkSecond(secondCh, thirdCh chan *Page, wg *sync.WaitGroup) {
	defer wg.Done()
	var stopCount uint
CHCLOSED:
	for {
		select {
		case _, exist := <-secondCh:
			if !exist {
				break CHCLOSED
			}
			if len(secondCh) == 0 {
				time.Sleep(time.Second)
				stopCount++
				if stopCount == 60 {
					break CHCLOSED
				}
			} else {
				stopCount = 0
			}
		default:
			time.Sleep(time.Second)
		}
	}
}

func WorkThird(thirdCh chan *Page, wg *sync.WaitGroup) {
	defer wg.Done()
	var stopCount uint
CHCLOSED:
	for {
		select {
		case _, exist := <-thirdCh:
			if !exist {
				break CHCLOSED
			}
			if len(thirdCh) == 0 {
				time.Sleep(time.Second)
				stopCount++
				if stopCount == 60 {
					break CHCLOSED
				}
			} else {
				stopCount = 0
			}
		default:
			time.Sleep(time.Second)
		}
	}
}