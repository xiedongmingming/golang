package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"io/ioutil"
)

const (
	UPLOAD_DIR = "./uploads"
)
func isExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	return os.IsExist(err)
}
func listHandler(w http.ResponseWriter, r *http.Request) {
	fileInfoArr, err := ioutil.ReadDir("./uploads")
	if err != nil {
		http.Error(w, err.Error(),
			http.StatusInternalServerError)
		return
	}
	var listHtml string
	for _, fileInfo := range fileInfoArr {
		imgid := fileInfo.Name
		listHtml += "<li><a href=\"/view?id="+imgid+"\">imgid</a></li>"
	}
	io.WriteString(w, "<ol>"+listHtml+"</ol>")
}
func viewHandler(w http.ResponseWriter, r *http.Request) {
	imageId := r.FormValue("id")
	imagePath := UPLOAD_DIR + "/" + imageId
	if exists := isExists(imagePath);!exists {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "image")// 正确的做法:
	//准确解析出文件的MimeType并将其作
	//为Content-Type进行输出，具体可参考Go语言标准库中的http.DetectContentType()方法
	//和mime包提供的相关方法
	http.ServeFile(w, r, imagePath)
}
func uploadHandler(w http.ResponseWriter, r *http.Request) { // 上传图片

	if r.Method == "GET" {

		io.WriteString(w, "<form method=\"POST\" action=\"/upload\" "+
			" enctype=\"multipart/form-data\">"+
			"Choose an image to upload: <input name=\"image\" type=\"file\" />"+
			"<input type=\"submit\" value=\"Upload\" />"+
			"</form>")

		return
	}
	if r.Method == "POST" {

		f, h, err := r.FormFile("image") // multipart.File *multipart.FileHeader error

		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}

		filename := h.Filename
		defer f.Close()
		t, err := os.Create(UPLOAD_DIR + "/" + filename)
		if err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		defer t.Close()
		if _, err := io.Copy(t, f); err != nil {
			http.Error(w, err.Error(),
				http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/view?id="+filename,
			http.StatusFound)
	}

}
func main() {
	http.HandleFunc("/", listHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/view", viewHandler)
	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
	}

}
