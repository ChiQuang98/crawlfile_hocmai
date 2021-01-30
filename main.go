package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

func main() {
	file,_:=ioutil.ReadFile("categories.json")
	data:=Categories{}
	_ = json.Unmarshal([]byte(file),&data)
	crawlAllFromCategories(data)
}
func crawlAllFromCategories(categoires Categories){
	var wg sync.WaitGroup
	jobs := make(chan Category,100)
	for w:=1;w<=10;w++{
		wg.Add(1)
		go worker(w,jobs,&wg)
	}
	for true{
		for i:=0;i<len(categoires.List) ;i++{
			jobs<-categoires.List[i]
		}
		time.Sleep(3*time.Hour)
	}
}
func worker(id int,jobs<-chan Category,wg *sync.WaitGroup){
	defer wg.Done()

	if _, err := os.Stat("./output"); os.IsNotExist(err) {
		os.Mkdir("./output", 0755)
	}
	for j:= range jobs {//duyệt qua tất cả category
		//dt := time.Now()
		fmt.Println("worker: ", id, "processing job: ", j)
		crawlFromCategory(j)
	}
}
func crawlFromCategory(category Category)  {
	files:= newFiles()
	res:=getHTMLPage(category.URL)
	if res == nil{
		return
	}
	//Chuyeen file page HTML cua 1 category
	files.getAllFileInformation(res,category.Title)
	files.TotalPages++
	for i:=2;i<=200;i++{
		files.TotalPages++
		nextPageLink := files.getNextUrl(res)
		if nextPageLink==""{
			break
		}
		res = getHTMLPage(nextPageLink)
		if res == nil{
			break
		}
		files.getAllFileInformation(res,category.Title)
	}
}