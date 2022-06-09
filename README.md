# dynSelect
a dynmical channel select with  one goroutine
## Example
``` go
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
```
