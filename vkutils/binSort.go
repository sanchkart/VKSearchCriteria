package vkutils

import _"log"


func BinSort(prevData, newData []int)[]int{
	data:=prevData
	if(len(prevData)==0){
		return newData
	}
	for _,d := range newData{
		L:=0
		R:=len(data)-1
		for(L!=R-1){
			if(d>data[(L+R)/2]){
				L=(L+R)/2
			}else{
				R=(L+R)/2
			}
		}
		//log.Println(d,L,R)
		if(data[R]<d){
			data = append(data,d)
		}else{
			data = append(data,0)
			copy(data[R+1:], data[R:])
			data[R] = d
		}
		//log.Println(data)
	}
	return data
}