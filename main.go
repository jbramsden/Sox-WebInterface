package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os/exec"
	"regexp"
	"sync"
	"time"
)

var started bool
var filen string
var mutex = &sync.Mutex{}
var cmd *exec.Cmd
var stdoutBuf, stderrBuf bytes.Buffer

var re = regexp.MustCompile(`\[[\s!=-]*\|[\s-=!]*\]`)

//var re = regexp.MustCompile(`\[(.*?)\]`)

const tpl = `<!DOCTYPE html> 
<html>
  <head> 
	<title>{{.Title}}</title>
	<meta http-equiv="refresh" content="1">
  </head>
  <body>
	<center>

    <h2>Levels</h2>
	<p>
	{{ range .Stderr }}
	{{ . }}<br>
	{{end}}
	</p>
    <button onclick="window.location.href='/increment';">
      {{.Record}} 
	</button><br>
	<i>{{ .Filename}}</i>
	</center>
  </body>
</html>`

func echoString(w http.ResponseWriter, r *http.Request) {
	//	outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
	t, err := template.New("webpage").Parse(tpl)
	if err != nil {
		log.Printf("Failed to parse template %s", err)
	} else {
		data := struct {
			Title    string
			Record   string
			Stdout   string
			Filename string
			Stderr   []string
		}{
			Title:  "Not Recording",
			Record: "Start Recording",
		}
		if started {
			idx := re.FindAllString(stderrBuf.String(), -1)
			data.Title = "Recording Started"
			data.Record = "Stop Recording"
			data.Stdout = string(stdoutBuf.Bytes())
			data.Filename = filen
			if l := len(idx); l >= 5 {
				data.Stderr = idx[l-5:]
			}
		}
		err = t.Execute(w, data)
		stderrBuf.Reset()
		stdoutBuf.Reset()
	}
}

func GetFilenameDate() string {
	// Use layout string for time format.
	const layout = "2006-01-05T15:04"
	// Place now in the string.
	t := time.Now()
	return "rec-" + t.Format(layout) + ".ogg"
}

func startSox() {
	filen = GetFilenameDate()
	cmd = exec.Command("sox", "-d", filen)
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf
	err := cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with %s\n", err)
	}
}

func stopSox() {
	err := cmd.Process.Kill()
	if err != nil {
		log.Fatalf("Cmd.Kill() failed with %s\n", err)
	}
}

func incrementCounter(w http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	if started {
		started = false
		stopSox()
	} else {
		started = true
		startSox()
	}
	mutex.Unlock()
	redirect := `<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="refresh" content="0; url='/'" />
  </head>
  <body>
    <p>Please follow <a href="/">this link</a>.</p>
  </body>
</html>`
	fmt.Fprintf(w, redirect)
}

func main() {
	http.HandleFunc("/", echoString)

	http.HandleFunc("/increment", incrementCounter)

	http.HandleFunc("/hi", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hi")
	})

	log.Fatal(http.ListenAndServe(":8081", nil))

}
