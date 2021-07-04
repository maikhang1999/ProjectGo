package main

import (
	"fmt"
	"math"
	"sync"
)
type Job struct {
	work []int
}
var wg sync.WaitGroup
var numberWorker int = 5
var jobs = make(chan Job,numberWorker)
var results = make(chan Job,numberWorker)
var res = make([]Job,numberWorker)

func RemoveIndex(s []int, index int) []int {
	ret := make([]int, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}
func findAndRemove(array []int,find int) []int{
	if len((array)) == 0 || (array)[0] == find {
		return []int{}
	}
	if (array)[len(array)-1] == find {
		(array) = append([]int{find}, (array)[:len(array)-1]...)
		return array
	}
	for p, x := range array {
		if x == find {
			(array) = append([]int{}, append((array)[:p], (array)[p+1:]...)...)
			break
		}
	}
	return array
}
func handleEvent(array []int,ignore_array []int) []int {
	for i:=0;i<len(ignore_array);i++{
		array = findAndRemove(array,ignore_array[i])
	}
	return array
}
func worker(ignore_array []int,jobs<-chan Job,res []Job){
	for n:=range jobs{
		output:=handleEvent(n.work,ignore_array)
		fmt.Println(output)
		res = append(res, Job{output})
	}


}
func filterBaseOnDis(realDis,refDis float32) bool{
	if realDis > refDis{
		return false
	}
	return true
}
func distance(la1,la2,lon1,lon2 float32) float32 {
	R := 6371e3
	/// p1,p2
	p1 :=la1*math.Pi/180
	p2:=la2*math.Pi/180
	//deltalp,deltagama
	deltalp :=(la2-la1)*math.Pi/180
	deltagama :=(lon2-lon1)*math.Pi/180
	// a,c->d
	a:=math.Sin(float64(deltalp/2))*math.Sin(float64(deltalp/2))+math.Cos(float64(p1))*math.Cos(float64(p2))*math.Sin(float64(deltagama/2))*math.Sin(float64(deltagama/2))
	c:=2*math.Atan2(math.Sqrt(a),math.Sqrt(1-a))
	d :=R*c
	return float32(d)
}
func main() {
	//array := []int{1,2,3,4,9,6,5,7,8,10}
	//array = RemoveIndex(array,5)
	//fmt.Println(array)
	//fmt.Println(filterBaseOnDis(21,20))
	//ignore_array := []int{2,6,9}
	//number :=len(array)
	//
	//offer := int(number/numberWorker)
	//// sinh ra cac worker
	//for i:=1;i<=numberWorker;i++{
	//	go worker(ignore_array,jobs,res)
	//}
	//// nap work vao channel jobs
	//	for j:=0;j<numberWorker;j++{
	//		jobs <- Job{array[j*offer:(j+1)*offer]}
	//		if(j == (numberWorker-1)){
	//			jobs <-Job{array[j*offer:]}
	//		}
	//	}
	//	close(jobs)
	//wg.Wait()
	//fmt.Println(res)
	fmt.Println(distance(41.9631174,40.7628267,-87.6770458,-73.9898293))

}

