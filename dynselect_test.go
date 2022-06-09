package dynselect_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/soapywu/dynselect"
)

func TestSelectN(t *testing.T) {
	waiter := func(ch chan int, value int, timeout time.Duration) {
		<-time.After(timeout)
		ch <- value
	}
	getReqChans := func(num int) []chan int {
		chs := make([]chan int, num)
		for i := 0; i < num; i++ {
			ch := make(chan int)
			go waiter(ch, i, time.Second)
			chs[i] = ch
		}
		return chs
	}
	num := 127
	reqChans := getReqChans(num)
	ctx, cancel := context.WithCancel(context.Background())
	resultChan := make(chan int)
	go dynselect.SelectN(ctx, reqChans, resultChan)
	for result := range resultChan {
		fmt.Println(result)
	}
	cancel()

	reqChans2 := getReqChans(num)
	ctx2, cancel2 := context.WithCancel(context.Background())
	resultChan2 := make(chan int)
	go dynselect.SelectN(ctx2, reqChans2, resultChan2)
	for result := range resultChan2 {
		fmt.Println(result)
		cancel2()
	}
	cancel2()
	<-time.After(4 * time.Second)
}
