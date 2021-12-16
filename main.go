package main

import (
	"bufio"
	"html/template"
	"net/http"
	"os"
	"strings"
)

/*This var is a pointer towards template.Template that is a
pointer to help process the html.*/
var tpl *template.Template

/*This init function, once it's initialised, makes it so that each html file
in the templates folder is parsed i.e. they all get looked through once and
then stored in the memory ready to go when needed*/
func init() {
	tpl = template.Must(template.ParseGlob("templates/*html"))
}

func main() {

	fs := http.FileServer(http.Dir("/templates/styles"))
	http.HandleFunc("/", index)
	http.Handle("/styles/", fs)
	//http.PathPrefix("/styles/").Handler(http.StripPrefix("/styles/",
	//	http.FileServer(http.Dir("templates/styles/"))))
	//http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("/templates/styles/"))))
	http.HandleFunc("/ascii-art", asciiart)
	http.ListenAndServe(":8080", nil)
}

//Handler function for the index
func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.Error(w, "404 address not found: wrong address entered!", http.StatusNotFound)
	} else {
		tpl.ExecuteTemplate(w, "index.html", nil)
	}
}

//Handler function to handle the ascii art conversion
func asciiart(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if err := recover(); err != nil {
			http.Error(w, " An internal server error has occurred: 500", http.StatusInternalServerError)
			return
		}
	}()

	if r.Method != "POST" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userBanner := r.FormValue("banner")
	userString := r.FormValue("uString")

	if userBanner == "" || userString == "" || strings.Contains(userString, "£") || strings.Contains(userString, "¬") {
		http.Error(w, "400 bad request made : empty or unrecognised string!", http.StatusBadRequest)
		return
	}

	// for i := range userString {
	// 	if userString[i] < 32 && userString[i] > 126 {
	// 		http.Error(w, "400 bad request made: empty or unrecognised string", http.StatusBadRequest)
	// 		return
	// 	}
	// }

	if strings.Contains(userString, "\n") {
		userString = strings.Replace(userString, "\r\n", "\\n", -1)
	}

	splitLines := SplitLines(userString)

	file, _ := os.Open(userBanner + ".txt")

	var ascii_temp []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ascii_temp = append(ascii_temp, scanner.Text())
	}
	ascii_map := make(map[int][]string) // makes the map to hold ascii chars
	start := 32
	for i := 0; i < len(ascii_temp); i++ {

		if len(ascii_map[start]) == 9 {
			start++
		}

		ascii_map[start] = append(ascii_map[start], ascii_temp[i])
	}

	var sString []string

	for j, val := range splitLines {
		for i := 0; i < 9; i++ {
			for k := 0; k < len(val); k++ {
				sString = append(sString, ascii_map[int(splitLines[j][k])][i])
			}
			sString = append(sString, "\n")

		}
	}

	SAscii := strings.Join(sString, "")

	d := struct {
		Banner string
		String string
		SAscii string
	}{
		Banner: userBanner,
		String: userString,
		SAscii: SAscii,
	}

	tpl.ExecuteTemplate(w, "index.html", d)
}

func SplitLines(s string) [][]byte {
	var count int

	for i := 0; i < len(s); i++ {
		if s[i] == 'n' && s[i-1] == '\\' {
			count++
		}
	}
	splitString := []byte(s)
	splitLines := make([][]byte, count+1)

	j := 0

	for i := 0; i < len(splitLines); i++ {
		for j < len(splitString) {

			if splitString[j] == 'n' && splitString[j-1] == '\\' {
				j++
				splitLines[i] = splitLines[i][:len(splitLines[i])-1]
				break
			}
			splitLines[i] = append(splitLines[i], splitString[j])
			j++
		}
	}

	return splitLines
}
