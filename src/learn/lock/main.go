package main

import (
	"myproject/public/fmt"
	"sync"
	"time"
)

var (
	m sync.Map
	rWMutex  sync.RWMutex
)
func main()  {
	startT := time.Now()
	WriteMap()
	//go func() {
	//	timer := time.NewTicker( 100* time.Millisecond)
	//	defer timer.Stop()
	//	for {
	//		select {
	//		case <-timer.C:
	//			WriteMap()
	//		}
	//	}
	//
	//}()

	wg := &sync.WaitGroup{}
	scount:=10000000
	wg.Add(scount)

	for i:=0;i<scount;i++ {
		go ReadMap(wg)
		if i%10==0{
			time.Sleep(1*time.Millisecond)
		}
	}



	wg.Wait()
	tc := time.Since(startT)	//计算耗时
	fmt.Println("time cost = %v\n", tc.String())

}


func WriteMap()  {
	//rWMutex.Lock()
	//defer rWMutex.Unlock()

	time.Sleep(1*time.Millisecond)
}

func ReadMap(wg *sync.WaitGroup)  {
	//rWMutex.RLock()
	//defer rWMutex.RUnlock()

	time.Sleep(10*time.Microsecond)

	wg.Done()
}