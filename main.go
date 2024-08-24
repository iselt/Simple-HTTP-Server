package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const maxUploadSize = 1024 * 1024 * 1024 // 1 GB

func main() {
	// 启动服务器
	host := "0.0.0.0"
	port := "8080"
	dir := "."
	if len(os.Args) == 3 {
		host = os.Args[1]
		port = os.Args[2]
	} else if len(os.Args) == 4 {
		host = os.Args[1]
		port = os.Args[2]
		dir = os.Args[3]
	} else if len(os.Args) > 1 {
		fmt.Println("Usage: " + os.Args[0] + " <HOST> <PORT> [<ROOT_DIR>]")
		return
	}

	// 创建一个文件服务器，提供当前目录下的静态文件服务
	fs := http.FileServer(http.Dir(dir))

	// 当方法为GET时，将请求交给文件服务器处理，当方法为PUT或POST时，交给上传处理函数处理
	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("INFO: %s %s %s", r.RemoteAddr, r.Method, r.URL.Path)
		if r.Method == http.MethodGet {
			fs.ServeHTTP(w, r)
		} else if r.Method == http.MethodPut || r.Method == http.MethodPost {
			uploadHandler(w, r)
		} else {
			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		}
	}))

	log.Printf("INFO: Server started at http://%s:%s", host, port)
	// 输出dir目录的绝对路径
	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Fatalf("ERROR: Error getting current directory: %v", err)
	}
	log.Printf("INFO: Serving files from %s", absDir)
	if err := http.ListenAndServe(host+":"+port, nil); err != nil {
		log.Fatalf("ERROR: Error starting server: %v\nUsage: %s <HOST> <PORT>", err, os.Args[0])
	}
}

// 上传处理函数
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut && r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// 限制上传文件大小
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	// 获取上传的文件路径
	urlPath := r.URL.Path
	filePath := filepath.Join(".", filepath.Clean(urlPath))

	// 检查文件是否已经存在
	if _, err := os.Stat(filePath); err == nil {
		http.Error(w, "File already exists", http.StatusConflict)
		return
	} else if !os.IsNotExist(err) {
		log.Printf("ERROR: Error checking file %s: %v", filePath, err)
		http.Error(w, "Error checking file", http.StatusInternalServerError)
		return
	}

	// 创建目标文件
	file, err := os.Create(filePath)
	if err != nil {
		log.Printf("ERROR: Error creating file %s: %v", filePath, err)
		http.Error(w, "Error creating file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// 复制上传数据到目标文件
	_, err = io.Copy(file, r.Body)
	if err != nil {
		log.Printf("ERROR: Error writing to file %s: %v", filePath, err)
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}

	log.Printf("INFO: File uploaded successfully to %s", filePath)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, `{"message": "File uploaded successfully", "filePath": "%s"}`, filePath)
}
