package dynselect

import (
	"context"
	"reflect"
)

const (
	MaxSelectNum = 65536
)

func buildSelectCase(ctx context.Context, chans interface{}, resultChan interface{}) (in []reflect.SelectCase, out *reflect.Value) {
	if reflect.TypeOf(resultChan).Kind() != reflect.Chan {
		return nil, nil
	}

	o := reflect.ValueOf(resultChan)
	if o.Type().ChanDir()&reflect.SendDir != reflect.SendDir {
		return nil, nil
	}
	resultElemKind := o.Type().Elem().Kind()

	if reflect.TypeOf(chans).Kind() != reflect.Slice {
		return nil, nil
	}

	s := reflect.ValueOf(chans)
	l := s.Len()
	if l == 0 || l >= MaxSelectNum {
		return nil, nil
	}
	selectCase := make([]reflect.SelectCase, l+1)
	for i := 0; i < l; i++ {
		if s.Index(i).Type().Kind() != reflect.Chan {
			return nil, nil
		}

		if s.Index(i).Type().ChanDir()&reflect.RecvDir != reflect.RecvDir {
			return nil, nil
		}

		if s.Index(i).Type().Elem().Kind() != resultElemKind {
			return nil, nil
		}

		selectCase[i] = reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(s.Index(i).Interface()),
		}
	}
	// last is ctx channel
	selectCase[l] = reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ctx.Done()),
	}
	return selectCase, &o

}

func SelectN(ctx context.Context, chans interface{}, resultChan interface{}) {
	selectCase, result := buildSelectCase(ctx, chans, resultChan)
	if selectCase == nil || result == nil {
		panic("invalid type chans or resultChan")
	}

	num := len(selectCase)
	left := num - 1
	checkDone := func(idx int) bool {
		selectCase[idx].Chan = reflect.ValueOf(nil) // block it
		left--
		return left == 0
	}

	for {
		chosen, recv, _ := reflect.Select(selectCase)
		if chosen != num-1 {
			if !recv.IsZero() {
				result.Send(recv)
			}
			done := checkDone(chosen)
			if done {
				result.Close()
				return
			}
		} else { //ctx.Done()
			result.Close()
			return
		}
	}
}
