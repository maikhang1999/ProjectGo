package main

import "fmt"

func Extend(slice []int, element int) []int {
	n := len(slice)
	slice = slice[0 : n+1]
	slice[n] = element
	return slice
}
// variadic function
func addItemtoArray(array []int,lists ...int) []int{

	for i:=0;i<len(lists);i++{
		array = append(array,lists[i])
	}
	return array
}
func main(){
	slice := make([]int,4,10)
	fmt.Printf("len:%d cap:%d \t ",len(slice),cap(slice))
	// copy
	//a := []int{9,10,11,12,13,14,15}
	//for i:=0;i<4;i++{
	//	slice[i] = i
	//}
	//fmt.Println(slice)
	//copy(slice[:3],a)
	//fmt.Printf("len:%d cap:%d \t ",len(slice),cap(slice))
	//fmt.Println(slice)
	//fmt.Println(a)
	//append
	//slice = append(slice,a...)
	//fmt.Println(slice)
	//fmt.Printf("len:%d cap:%d",len(slice),cap(slice))
	array :=[]int{1,2,3}

	array = addItemtoArray(array,4,5,6)
	fmt.Println(array)
}