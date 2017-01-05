package lostblock

import (
	"github.com/golang/glog"
	"flag"
	"zhiwang_bc_message/geth/container/stack"
	"database/sql"
	"time"
	"fmt"
)

//var data = [10]int{1, 2, 5, 6, 9,
//	13, 14, 17, 19, 21}

/**
二分法
 */
func TestBinary() {
	data := NewArrayData()
	binarySearch(int64(0), data.Len()-1,data)
}

/**
循环
 */
func TestArrayLoop() {
	flag.Parse()
	stack := stack.NewStack()

	stack.Push(&state{low:0, high:9})
	data := NewArrayData()
	loop(stack, data)
}

//mysql
func MysqlLoop(db *sql.DB) []int64 {
	stack := stack.NewStack()
	data := NewMysqlData(db)
	total := data.Len() - 1
	if total <= 0 {
		return []int64{}
	}
	stack.Push(&state{low:int64(0), high:total})
	glog.Infof("start %v \n", time.Now())
	list := loop(stack, data)
	glog.Infof("end %v \n", time.Now())
	return list
}
/**
循环算法
 */
type state  struct {
	low, high int64
}

func loop(stack *stack.Stack, data Data) ([]int64) {
	rtnData := make([]int64, 0, 100)
	for {
		s := stack.Pop()
		if s == nil {
			break
		}
		low := s.(*state).low
		high := s.(*state).high
		if high - low == 0 {
			continue
		}
		if high - low == 1 {
			lowValue := data.Get(low)
			highValue := data.Get(high)
			glog.Infof("lowValue %d highValue %d \n", lowValue, highValue)
			for i := (lowValue + 1); i < highValue; i++ {
				//fmt.Printf("value %d \n", i)
				rtnData = append(rtnData, i)
			}
			continue
		}

		totalComplete := check(low, high, data.Get(low), data.Get(high))
		glog.Infof("totalComplete %v low %d high %d \n", totalComplete, low, high)
		if totalComplete == false {
			mid := (high - low) / 2
			leftComplete := check(low, low + mid, data.Get(low), data.Get(low + mid))
			glog.Infof("leftComplete %v low %d high %d \n", leftComplete, low, low + mid)
			if leftComplete == false {
				stack.Push(&state{low:low, high:low + mid})
			}
			rightComlete := check(low + mid, high, data.Get(low + mid), data.Get(high))
			glog.Infof("rightComlete %v low %d high %d \n", rightComlete, low + mid, high)
			if rightComlete == false {
				stack.Push(&state{low:low + mid, high:high})
			}
		}
	}
	return rtnData

}

/**
递归算法
 */
func binarySearch(low, high int64, data Data) {
	glog.Infof("low %d high %d \n", low, high)
	if high - low == 0 {
		return
	}
	if high - low == 1 {
		lowValue := data.Get(low)
		highValue := data.Get(high)
		glog.Infof("lowValue %d highValue %d \n", lowValue, highValue)
		for i := (lowValue + 1); i < highValue; i++ {
			fmt.Printf("value %d \n", i)
		}
		return
	}

	totalComplete := check(low, high, data.Get(low), data.Get(high))
	glog.Infof("totalComplete %v low %d high %d \n", totalComplete, low, high)
	if totalComplete == false {
		mid := (high - low) / 2
		leftComplete := check(low, low + mid, data.Get(low), data.Get(low + mid))
		glog.Infof("leftComplete %v low %d high %d \n", leftComplete, low, low + mid)
		if leftComplete == false {
			binarySearch(low, low + mid,data)
		}
		rightComlete := check(low + mid, high, data.Get(low + mid), data.Get(high))
		glog.Infof("rightComlete %v low %d high %d \n", rightComlete, low + mid, high)
		if rightComlete == false {
			binarySearch(low + mid, high,data)
		}
	}

}

func check(indexStart, indexEnd int64, firstValue, lastValue int64) bool {
	indexLen := indexEnd - indexStart + 1
	valueLen := lastValue - firstValue + 1
	if indexLen == valueLen {
		return true
	}
	return false
}

