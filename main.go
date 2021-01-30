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
	categoires:=Categories{}
	_ = json.Unmarshal([]byte(file),&categoires)
	jobs := make(chan Category,100)
	results :=make(chan Category,100)
	errors :=make(chan error,1000)
	crawlAllFromCategories(jobs,results,errors)
	for true{
		fmt.Println("IN")
		for i:=0;i<len(categoires.List) ;i++{
			jobs<-categoires.List[i]
		}
		close(jobs)
		select {
		case err := <-errors:
			fmt.Println("Error: ",err.Error())
		default:
		}
		time.Sleep(3*time.Hour)

	}
}
func crawlAllFromCategories(jobs<- chan Category,results chan <-  Category,errors chan <- error){
	var wg sync.WaitGroup
	for w:=1;w<=10;w++{
		wg.Add(1)
		go worker(w,jobs,results,errors,&wg)
	}
}
func worker(id int,jobs<-chan Category,results chan <- Category,errors chan <- error,wg *sync.WaitGroup){
	defer wg.Done()

	if _, err := os.Stat("./output"); os.IsNotExist(err) {
		os.Mkdir("./output", 0755)
	}
	for j:= range jobs {//duyệt qua tất cả category
		//dt := time.Now()
		fmt.Println("worker: ", id, "processing job: ", j)
		crawlFromCategory(j,errors)
	}
}
func crawlFromCategory(category Category,errors chan <- error)  {
	files:= newFiles()
	res:=getHTMLPage(category.URL)
	if res == nil{
		errors <- fmt.Errorf("Page not found")
		return
	}
	//Chuyeen file page HTML cua 1 category
	files.getAllFileInformation(res,category.Title,errors)
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
		files.getAllFileInformation(res,category.Title,errors)
	}
	filesJson,err := json.Marshal(files)
	checkError(err)
	err = ioutil.WriteFile("categories.json",filesJson,0644)
}