package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"

)

func main() {
	var wg sync.WaitGroup
	file,_:=ioutil.ReadFile("categories.json")
	categoires:=Categories{}
	_ = json.Unmarshal([]byte(file),&categoires)
	jobs := make(chan Category,100)
	results :=make(chan File,10000)
	errors :=make(chan error,100)
	crawlAllFromCategories(jobs,results,errors,&wg)
	for true{
		//fmt.Println("IN")
		for i:=0;i<len(categoires.List) ;i++{
			fmt.Println("IN")
			fmt.Println(categoires.List[i])
			jobs<-categoires.List[i]
		}
		//close(jobs)
		for true{
			select {
			case fileReceive,open := <- results:
				//dt := time.Now()
				if !open{
					break
				}
				f, _ := os.OpenFile("./output/"+fileReceive.CategoryName+".json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
				defer f.Close()
				fileJSON, err := json.Marshal(fileReceive)
				checkError(err)
				io.WriteString(f,string(fileJSON)+"\n")
				fmt.Println(fileReceive.CategoryName)
			case err := <-errors:
				fmt.Println("Error: ", err.Error())
			default:
			}
		}
		fmt.Println("DONE SESSION")
		time.Sleep(3 * time.Hour)
	}
	wg.Wait()
}
func crawlAllFromCategories(jobs<- chan Category,results chan <- File,errors chan <- error,wg *sync.WaitGroup){
	for w:=1;w<=10;w++{
		wg.Add(1)
		go worker(w,jobs,results,errors,wg)
	}
}
func worker(id int,jobs<-chan Category,results chan <- File,errors chan <- error,wg *sync.WaitGroup){
	defer wg.Done()

	if _, err := os.Stat("./output"); os.IsNotExist(err) {
		os.Mkdir("./output", 0755)
	}
	for j:= range jobs {//duyệt qua tất cả category
		//dt := time.Now()
		fmt.Println("worker: ", id, "processing job: ", j)
		crawlFromCategory(j,results,errors)
	}
}
func crawlFromCategory(category Category,results chan <- File,errors chan <- error)  {
	files:= newFiles()
	res:=getHTMLPage(category.URL)
	if res == nil{
		errors <- fmt.Errorf("Page not found")
		return
	}
	//Chuyeen file page HTML cua 1 category
	files.CategoryName = category.Title
	files.getAllFileInformation(res,results,category.Title,errors)
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
		files.getAllFileInformation(res,results,category.Title,errors)
	}

	//filesJson,err := json.Marshal(files)
	//checkError(err)
	//Sau khi load hết một category thì chuyển category files đó cho result channel
}