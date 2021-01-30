package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log"
	"net/http"
	"strings"
	"sync"
)
func getHTMLPage(url string)*goquery.Document{
	res,err :=http.Get(url)
	if err !=nil{
		println("Error")
		return nil
	}
	if res.StatusCode!=200{
		println("Error res status")
		return nil
	}
	doc,err := goquery.NewDocumentFromReader(res.Body)
	if err !=nil{
		return nil
	}
	return doc
}
func (files *Files) getAllFileInformation(doc *goquery.Document,category string,errors chan <- error){
	var wg sync.WaitGroup
	doc.Find(".lib-grid .top a").Each(func(i int, s *goquery.Selection){ //lấy danh sách các file trong 1 category
		fileLink,_ := s.Attr("href")
		wg.Add(1)
		//fmt.Println(fileLink)
		go files.getFileInformation("https://hocmai.vn"+fileLink,category,&wg,errors)
	})
	wg.Wait()
}
//func main() {
//	//results :=make(chan Category,100)
//	errors :=make(chan error,1000)
//	files := newFiles()
//	files.getAllFileInformation(getHTMLPage("https://hocmai.vn/kho-tai-lieu/list.php?category=204"),"ád",errors)
//}
func (files *Files) getNextUrl(doc *goquery.Document) string{
	nextPageLink,_:=doc.Find(".paging a:last-child").Attr("href")
	//fmt.Println(nextPageLink)
	if nextPageLink == ""{
		println("End of Category")
		return ""
	}
	return "https://hocmai.vn/kho-tai-lieu/"+nextPageLink
}
func (files *Files) getFileInformation(fileLink string,category string,wg *sync.WaitGroup,errors chan <- error){
	//đọc thông tin từng file
	defer wg.Done()
	res:=getHTMLPage(fileLink)
	if res==nil{
		errors <- fmt.Errorf("Page not found")
		return
	}
	title := res.Find(".lib-section .head h4").Text()
	numberPage := res.Find(".lib-meta ul li:first-child span").Text()
	numberViewed := res.Find(".lib-meta ul li:nth-child(2) span").Text()
	numberDownloaded := res.Find(".lib-meta ul li:nth-child(3) span").Text()
	author := res.Find(".lib-meta ul li:nth-child(4) span").Text()
	date := res.Find(".lib-meta ul li:nth-child(5) span").Text()
	numberPage = strings.TrimSpace(strings.Split(numberPage,":")[1])
	numberViewed = strings.TrimSpace(strings.Split(numberViewed,":")[1])
	numberDownloaded = strings.TrimSpace(strings.Split(numberDownloaded,":")[1])
	author = strings.TrimSpace(strings.Split(author,":")[1])
	var flagCheck bool
	flagCheck = true
	if author == ""{
		errors <- fmt.Errorf("Author is Empty")
		flagCheck = false
	}
	if title == ""{
		errors <- fmt.Errorf("Title is Empty")
		flagCheck = false
	}
	if numberPage == ""{
		errors <- fmt.Errorf("numberPage is Empty")
		flagCheck = false
	}
	if numberViewed == ""{
		errors <- fmt.Errorf("numberViewed is Empty")
		flagCheck = false
	}
	if numberDownloaded == ""{
		errors <- fmt.Errorf("numberDownloaded is Empty")
		flagCheck = false
	}
	if date == ""{
		errors <- fmt.Errorf("date is Empty")
		flagCheck = false
	}
	if flagCheck ==true {
		urlString := fileLink
		ID:= strings.Split(urlString,"?")[1]
		file:=File{
			ID: ID,
			Title: title,
			numberPage: numberPage,
			numberViewed: numberViewed,
			numberDownloaded: numberDownloaded,
			Author: author,
			Date: date,
		}
		//fileJson,err := json.Marshal(file)
		//checkError(err)
		files.TotalPages++
		files.CategoryName = category
		files.List = append(files.List,file)
	}
}
func checkError(err error) {
	if err != nil {
		print("Error: " + err.Error())
		log.Println(err)
	}
}
//func main() {
//	var wg sync.WaitGroup
//	var category string
//	category = "ss"
//	files := newFiles()
//	wg.Add(1)
//	files.getFileInformation("https://hocmai.vn/kho-tai-lieu/read.php?id=14595",category,&wg)
//	files.getNextUrl(getHTMLPage("https://hocmai.vn/kho-tai-lieu/list.php?category=244"))
//}
