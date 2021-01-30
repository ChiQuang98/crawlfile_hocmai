package main
type Category struct{
	Title string `json:"title"`
	URL string `json:"url"`
}

type Categories struct {
	Total int        `json:"total"`
	List  []Category `json:"categories"`
}

type File struct {
	ID string `json:"ID"`
	Title string `json:"title"`
	numberPage string `json:"numberpage"`
	numberViewed string `json:"numberviewed"`
	numberDownloaded string `json:"numberdownloaded"`
	Author string `json:"Author"`
	Date string `json:Date`
}
type Files struct{
	CategoryName string `json:"categoryname"`
	TotalFiles int `json:"totalfiles"`
	TotalPages int `json:"totalpages"`
	List []File `json:"files"`
}
func newCategoires() *Categories  {
	return &Categories{}
}
func newFiles() *Files{
	return &Files{}
}