package vkutils

import _"log"

func MergeSort(m []int) []int {
	if len(m) <= 1 {
		return m
	}

	mid := len(m) / 2
	left := m[:mid]
	right := m[mid:]

	left = MergeSort(left)
	right = MergeSort(right)

	return merge(left, right)
}

func merge(left, right []int) []int {
	var result []int
	for len(left) > 0 || len(right) > 0 {
		if len(left) > 0 && len(right) > 0 {
			if left[0] <= right[0] {
				result = append(result, left[0])
				left = left[1:]
			} else {
				result = append(result, right[0])
				right = right[1:]
			}
		} else if len(left) > 0 {
			result = append(result, left[0])
			left = left[1:]
		} else if len(right) > 0 {
			result = append(result, right[0])
			right = right[1:]
		}
	}

	return result
}

func MergeSortTwo(m []int,message chan []int){
	//log.Println(m)
	if len(m) <= 1 {
		message<-m
	}else {
		mid := len(m) / 2
		left := m[:mid]
		right := m[mid:]

		var leftNew chan []int = make(chan []int)
		var rightNew chan []int = make(chan []int)

		go MergeSortTwo(left,leftNew)
		go MergeSortTwo(right,rightNew)

		left = <-leftNew
		right = <-rightNew

		var answer chan []int = make(chan []int)
		go mergeTwo(left, right,answer)
		ans := <-answer
		message <- ans
	}
}

func mergeTwo(left, right []int, message chan []int){
	var result []int
	for len(left) > 0 || len(right) > 0 {
		if len(left) > 0 && len(right) > 0 {
			if left[0] <= right[0] {
				result = append(result, left[0])
				left = left[1:]
			} else {
				result = append(result, right[0])
				right = right[1:]
			}
		} else if len(left) > 0 {
			result = append(result, left[0])
			left = left[1:]
		} else if len(right) > 0 {
			result = append(result, right[0])
			right = right[1:]
		}
	}

	message<-result
}

func MergeNext(ldata []int, rdata []int) (result []int) {
	result = make([]int, len(ldata) + len(rdata))
	lidx, ridx := 0, 0

	for i:=0;i<cap(result);i++ {
		switch {
		case lidx >= len(ldata):
			result[i] = rdata[ridx]
			ridx++
		case ridx >= len(rdata):
			result[i] = ldata[lidx]
			lidx++
		case ldata[lidx] < rdata[ridx]:
			result[i] = ldata[lidx]
			lidx++
		default:
			result[i] = rdata[ridx]
			ridx++
		}
	}

	return
}

func MergeSortNext(data []int, r chan []int) {
	if len(data) == 1 {
		r <- data
		return
	}

	leftChan := make(chan []int)
	rightChan := make(chan []int)
	middle := len(data)/2

	go MergeSortNext(data[:middle], leftChan)
	go MergeSortNext(data[middle:], rightChan)

	ldata := <-leftChan
	rdata := <-rightChan

	close(leftChan)
	close(rightChan)
	r <- MergeNext(ldata, rdata)
	return
}