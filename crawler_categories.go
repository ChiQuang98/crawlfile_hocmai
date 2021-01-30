package main

import (
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
)
func (categories *Categories) getAllCategories(doc *goquery.Document) {
	doc.Find(".lib-sbar>ul>li>a").Each(func(i int, s *goquery.Selection){
		cateLink,_:= s.Attr("href")
		cateTitle:=s.Text()
		fmt.Println(cateLink)
		fmt.Println(cateTitle)
		category := Category{
			Title: cateTitle,
			URL: "https://hocmai.vn"+cateLink,
		}
		categories.Total++
		categories.List = append(categories.List,category)
	})
}
func crawlAllCategories(){
	categories := newCategoires()
	res := getHTMLPage("https://hocmai.vn/kho-tai-lieu/")
	categories.getAllCategories(res)
	categoriesJson,err := json.Marshal(categories)
	checkError(err)
	err = ioutil.WriteFile("categories.json",categoriesJson,0644)
	checkError(err)
}
func main() {
	crawlAllCategories()
}