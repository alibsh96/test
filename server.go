package main

import (
        "encoding/json"
        "fmt"
        "log"
        "net/http"
		"time"
		"io/ioutil"
		"os"
)

type Rect struct {
        X int `json:"x"`
        Y  int `json:"y"`
		Width int `json:"width"`
		Height int	`json:"height`
}

type SavedRect struct {
	X int `json:"x"`
	Y  int `json:"y"`
	Width int `json:"width"`
	Height int	`json:"height"`
	Time string `json:"time"`
}

type Req struct {
	Main Rect `json:"main"`
	Input []Rect `json:"input`
}

type Res []SavedRect
var re Rect;

func intersects(r1 Rect, r2 Rect)  bool {
	if r1.X > r2.X + r2.Width || r2.X > r1.X + r1.Width{
		return false
	}
	if r1.Y > r2.Y + r2.Height || r2.Y > r1.Y + r1.Height{
		return false
	}
	return true
}

var db []SavedRect
var idx int
func requestHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		var res Res

		for _, rect := range db {
			res = append(res, rect)
		}
		log.Println(res)
		json.NewEncoder(w).Encode(res)
		
	case "POST":
		var req Req
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			
			for _, rect := range req.Input {
				if intersects(req.Main, rect){
					dt := time.Now()
					var sRect = SavedRect {
						X: rect.X,
						Y: rect.Y,
						Width: rect.Width,
						Height: rect.Height,
						Time: dt.Format("01-02-2006 15:04:05"),
					}
					log.Println(sRect)
					db = append(db, sRect)
				}
			}
			file, _ := json.MarshalIndent(db, "", " ")
			_ = ioutil.WriteFile("db.json", file, 0644)
			
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "I can't do that.")
	}
}

func main() {
		jsonFile, err := os.Open("db.json")
		if err != nil {
			fmt.Println(err)
		}
		byteValue, _ := ioutil.ReadAll(jsonFile)
		json.Unmarshal(byteValue, &db)
        http.HandleFunc("/", requestHandler)
        log.Println("Go!")
        http.ListenAndServe(":8080", nil)
}