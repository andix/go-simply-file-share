package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gorilla/mux"
)

//go:embed templates/index.html
var templatesFS embed.FS

type FileInfo struct {
	Name          string    `json:"name"`
	Size          int64     `json:"size"`
	UploadTime    time.Time `json:"upload_time"`
	DownloadCount int       `json:"download_count"`
}

var files []FileInfo
var dataFile = "data.json"
var myfileDir = "myfile"

func loadData() {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		files = []FileInfo{}
		return
	}
	file, err := os.Open(dataFile)
	if err != nil {
		fmt.Println("Error opening data file:", err)
		return
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&files)
	if err != nil {
		fmt.Println("Error decoding data:", err)
	}
}

func saveData() {
	file, err := os.Create(dataFile)
	if err != nil {
		fmt.Println("Error creating data file:", err)
		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	err = encoder.Encode(files)
	if err != nil {
		fmt.Println("Error encoding data:", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(templatesFS, "templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := header.Filename
	timestamp := time.Now().Format("060102150405_")
	newFilename := timestamp + filename
	dst, err := os.Create(filepath.Join(myfileDir, newFilename))
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error copying file", http.StatusInternalServerError)
		return
	}

	// Add to files
	newFile := FileInfo{
		Name:          newFilename,
		Size:          size,
		UploadTime:    time.Now(),
		DownloadCount: 0,
	}
	files = append(files, newFile)
	saveData()

	// Check if request is AJAX (has XMLHttpRequest header)
	if r.Header.Get("X-Requested-With") == "XMLHttpRequest" {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "success", "message": "File uploaded successfully"})
	} else {
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func downloadHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["filename"]

	for i, f := range files {
		if f.Name == filename {
			files[i].DownloadCount++
			saveData()
			break
		}
	}

	filePath := filepath.Join(myfileDir, filename)
	http.ServeFile(w, r, filePath)
}

func filesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func deleteHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	vars := mux.Vars(r)
	filename := vars["filename"]

	// Remove from files slice
	for i, f := range files {
		if f.Name == filename {
			files = append(files[:i], files[i+1:]...)
			break
		}
	}

	// Delete file from disk
	filePath := filepath.Join(myfileDir, filename)
	err := os.Remove(filePath)
	if err != nil {
		http.Error(w, "Error deleting file", http.StatusInternalServerError)
		return
	}

	saveData()

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	port := flag.String("p", "8022", "Port to run the server on")
	flag.Parse()

	// Create myfile directory if not exists
	if _, err := os.Stat(myfileDir); os.IsNotExist(err) {
		os.Mkdir(myfileDir, 0755)
	}

	loadData()

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/upload", uploadHandler).Methods("POST")
	r.HandleFunc("/download/{filename}", downloadHandler).Methods("GET")
	r.HandleFunc("/delete/{filename}", deleteHandler).Methods("POST")
	r.HandleFunc("/files", filesHandler).Methods("GET")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))

	fmt.Printf("Server starting on port %s\n", *port)
	http.ListenAndServe(":"+*port, r)
}
