package main

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"math"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type User struct {
	Id        int `json:"id"`
	Name      string `json:"name"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Gender    string `json:"gender"`
	Birthday string `json:"birthday"`
	LastActive string `json:"lastActive"`

}
type ByLastActive []User
func (a ByLastActive) Len() int           { return len(a) }
func (a ByLastActive) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByLastActive) Less(i, j int) bool { return caculateTimeOnline(a[i].LastActive) > caculateTimeOnline(a[j].LastActive) }
type getRes struct {
	Id int
	Latitude float32
	Longitude float32
	Gender string
	Distance float32
	AgeRange [2]int
	IgnoreArray []int
}
var db *sql.DB
func dbConn(){
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "maikhang99"
	dbName := "gosever"
	database, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	db=database


}
// caculate distance between 2 points
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
// return true if realDis <= refDis
func filterBaseOnDis(realDis,refDis float32) bool{
	if realDis > refDis{
		return false
	}
	return true
}
func convertData(c *gin.Context)getRes{
	id := c.PostForm("Id")
	latitude := c.PostForm("latitude")
	longitude := c.PostForm("longitude")
	dis := c.PostForm("distance")
	ageStart :=c.PostForm("ageStart")
	ageEnd :=c.PostForm("ageEnd")
	ignoreArray := c.PostForm("ignore_array")
	// convert
	valueId,_:=strconv.Atoi(id)
	la,_:=strconv.ParseFloat(latitude,32)
	lon,_:=strconv.ParseFloat(longitude,32)
	valAgeStart,_:=strconv.Atoi(ageStart)
	valAgeEnd,_:=strconv.Atoi(ageEnd)
	valDis,_:=strconv.ParseFloat(dis,32)

	res := getRes{}
	res.Id=valueId
	res.Latitude=float32(la)
	res.Longitude=float32(lon)
	res.Gender= c.PostForm("gender")
	res.AgeRange=[2]int{valAgeStart,valAgeEnd}
	res.Distance=float32(valDis)
	// ignoreArray
	err:= json.Unmarshal([]byte(ignoreArray), &res.IgnoreArray)
	if err!=nil{
		log.Fatal(err)
	}
	return res
}
func caculateAge(year int)int{
	t :=time.Now()
	y:=t.Year()

	return y-year
}
//"12:01 Am" quy doi sang phut
func caculateTimeOnline(tOnline string)int{
	t := time.Now()
	subTime := strings.Split(tOnline," ")
	hour,_:=strconv.Atoi(strings.Split(subTime[0],":")[0])
	minute,_:=strconv.Atoi(strings.Split(subTime[0],":")[1])
	formatTime :=subTime[1]
	if (strings.Compare(formatTime,"AM")==0){
		hour = hour +12
	}

	return ((t.Hour()-hour)*60+(t.Minute()-minute))

}
func binarySearch(needle int, haystack []int) bool {

	low := 0
	high := len(haystack) - 1

	for low <= high{
		median := (low + high) / 2

		if haystack[median] < needle {
			low = median + 1
		}else{
			high = median - 1
		}
	}

	if low == len(haystack) || haystack[low] != needle {
		return false
	}

	return true
}
func Index(c *gin.Context){
	rs := convertData(c)
	yearBeginAndEnd:=[2]int{caculateAge(rs.AgeRange[0]),caculateAge(rs.AgeRange[1])}
	dateEndValue := time.Date( yearBeginAndEnd[0],1,1, 12, 30, 0, 0, time.UTC)
	dateBeginValue := time.Date(yearBeginAndEnd[1],12,30, 12, 30, 0, 0, time.UTC)
	selDB, err := db.Query("SELECT * FROM users where NOT users.user_id =? AND users.gender=?" +
		" AND users.birthday >=? AND users.birthday<=?",rs.Id,rs.Gender,
		dateBeginValue.Format("2006-01-02"),dateEndValue.Format("2006-01-02"))
	if err != nil {
		panic(err.Error())
	}
	res := []User{}
	for selDB.Next() {
		user := User{}
		err = selDB.Scan(&user.Id, &user.Name, &user.Latitude, &user.Longitude, &user.Gender, &user.Birthday, &user.LastActive)
		if err != nil {
			panic(err.Error())
		}
		checkID :=binarySearch(user.Id,rs.IgnoreArray)
		checkDis := filterBaseOnDis(distance(rs.Latitude,user.Latitude,rs.Longitude,user.Longitude),rs.Distance)
		//fmt.Printf("%.2f \t",distance(rs.Latitude,user.Latitude,rs.Longitude,user.Longitude))
		if checkID == true || checkDis ==false{
			continue
		}

		res = append(res, user)
	}

	sort.Sort(ByLastActive(res))
	c.JSON(http.StatusOK, gin.H{"result": res})

}
func main(){
	log.Println("Server started on: http://localhost:8080")
	dbConn()
	r := gin.Default()
	r.POST("/recommended",Index)
	r.Run(":3000")


}