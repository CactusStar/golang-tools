package main

import (
	"archive/zip"
	"fmt"
	"html/template"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-vgo/robotgo"
	"github.com/kbinani/screenshot"
)

func save(img *image.RGBA, filepath string) {
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Duration(5) * time.Second)
		panic(err)
	}
	defer file.Close()
	png.Encode(file, img)
}

func takeandsave_file(savedpath string, trigger string) {
	// take screenshot and save as a single file
	// time.Sleep(time.Duration(5) * time.Second)
	if trigger != "" {
		img, err := screenshot.CaptureDisplay(0)
		if err != nil {
			fmt.Println("there have error")
			time.Sleep(time.Duration(5) * time.Second)
			panic(err)
		}
		currenttime := time.Now()
		timeString := currenttime.Format("2006-01-02-15-04-05")
		actualpath := savedpath + timeString + ".png"
		fmt.Println(actualpath)
		save(img, actualpath)
	} else {
		fmt.Println("No trigger recieved")
	}

}

func monitorkeyboardandsave(key string, savedpath string) {
	for {
		fmt.Println("--- Press something ---")

		ok := robotgo.AddEvent(key)
		fmt.Println(ok)
		if ok {
			fmt.Printf("Pressed: %s\n", key)
			img, err := screenshot.CaptureDisplay(0)
			if err != nil {
				fmt.Println("there have error")
				time.Sleep(time.Duration(5) * time.Second)
				panic(err)
			}
			currenttime := time.Now()
			timeString := currenttime.Format("2006-01-02-15-04-05")
			actualpath := savedpath + timeString + ".png"
			fmt.Println(actualpath)
			save(img, actualpath)
			fmt.Println("got")
		}
	}
}

func zipfile(savedpath string, files []string) {
	// compress 10 files to a zip file
	currenttime := time.Now()
	timeString := currenttime.Format("2006-01-02-15-04-05")
	actualpath := savedpath + timeString + ".zip"
	archive, err := os.Create(actualpath)
	if err != nil {
		fmt.Println(err)
		time.Sleep(time.Duration(5) * time.Second)
		// panic(err)
	}
	defer archive.Close()
	zipWriter := zip.NewWriter(archive)
	fmt.Println(files)
	for _, file := range files {
		// newfilename := strings.Split(file, "\\")
		filename := strings.Split(file, "\\")
		fmt.Println(file)
		f, err := os.Open(file)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Duration(5) * time.Second)
			// panic(err)
		}

		w, err := zipWriter.Create(filename[len(filename)-1])
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Duration(5) * time.Second)
			// panic(err)
		}
		if _, err := io.Copy(w, f); err != nil {
			fmt.Println(err)
			time.Sleep(time.Duration(5) * time.Second)
			// panic(err)
		}
		f.Close()
	}
	zipWriter.Close()
}
func return10(savedpath string) []string {
	files, _ := ioutil.ReadDir(savedpath)
	filecount := len(files)
	fileslice := make([]string, 0)
	if filecount >= 10 {
		for i := 0; i < 10; i++ {
			fileslice = append(fileslice, savedpath+files[i].Name())
		}
	}
	return fileslice
}

func return10files(files []string) []string {
	returnfiles := make([]string, 0)
	for i := 0; i < 10; i++ {
		returnfiles = append(returnfiles, files[i])
	}
	return returnfiles
}

func returnfilesArray(savedpath string) []string {
	files, _ := ioutil.ReadDir(savedpath)
	filesilce := make([]string, 0)
	for _, file := range files {
		if strings.Contains(file.Name(), ".png") {
			filesilce = append(filesilce, savedpath+file.Name())
		}
	}
	return filesilce
}

func restfiles(savedpath string) []string {
	files := returnfilesArray(savedpath)
	tenfile := return10(savedpath)
	rest := make([]string, 0)
	for i := len(files); i >= 0; i-- {
		for _, zipfile := range tenfile {
			if files[i] == zipfile {
				rest = append(files[:i], files[i+1:]...)
			}

		}
	}
	return rest
}

func zip10files(savedpath string) {
	zipedfile := make([]string, 0)
	// restfiles := make([]string, 0)
	for {
		// time.Sleep(time.Duration(5) * time.Second)
		allfiles := returnfilesArray(savedpath)
		// fmt.Printf("all files length: %d\n", len(allfiles))
		if len(zipedfile) == 0 {
			fmt.Println("try to zip")
			time.Sleep(time.Duration(3) * time.Second)
			if len(allfiles) >= 10 {
				fmt.Println("===Condition 2===")
				needziped := return10files(allfiles)
				fmt.Println("===zip 2===")
				zipfile(savedpath, needziped)

				zipedfile = append(zipedfile, needziped...)
				fmt.Println(len(zipedfile))
				time.Sleep(time.Duration(3) * time.Second)

			}
		} else if len(zipedfile) != len(allfiles) {
			for i := len(allfiles) - 1; i >= 0; i-- {

				for _, eachfile := range zipedfile {

					if allfiles[i] == eachfile {

						allfiles = append(allfiles[:i], allfiles[i+1:]...)
					}

				}
			}

			if len(allfiles) >= 10 {
				fmt.Println("===Condition 2===")
				needziped := return10files(allfiles)
				fmt.Println("===zip 2===")
				zipfile(savedpath, needziped)

				zipedfile = append(zipedfile, needziped...)
				fmt.Println(len(zipedfile))
				time.Sleep(time.Duration(3) * time.Second)

			}
		} else {
			time.Sleep(time.Duration(3) * time.Second)
			fmt.Println("condition not match")
		}

		// fmt.Println(restfiles)
		// fmt.Println(len(restfiles))
		// time.Sleep(time.Duration(5) * time.Second)

	}

}

func takesnapshot(w http.ResponseWriter, r *http.Request) {

	t, _ := template.ParseFiles("takesnap.html")
	r.ParseForm()
	// fmt.Fprintln(w, r.Form)
	path := r.Form
	key := strings.Join(path["key"], "")
	savedpath := strings.Join(path["savedpath"], "")
	// fmt.Fprintln(w, r.Form)
	fmt.Println(key)
	fmt.Println(savedpath)
	if key != "" && savedpath != "" {
		go monitorkeyboardandsave(key, savedpath)

		go zip10files(savedpath)
	}
	t.Execute(w, "")
	// t.Execute(w, key != "")
}
func main() {
	// savedpath := "C:\\Users\\XHe1\\Desktop\\TestCopy\\folder1\\forsnap\\"
	// var key string
	// fmt.Scanln(&key)
	// go monitorkeyboardandsave(key, savedpath)

	// zip10files(savedpath)

	server := http.Server{
		Addr: "127.0.0.1:8080",
	}
	http.HandleFunc("/takesnapshot", takesnapshot)
	server.ListenAndServe()
}
