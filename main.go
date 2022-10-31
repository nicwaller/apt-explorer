package main

import (
	"apt-explorer/lib/apt"
	"apt-explorer/lib/apt/kvstore"
	"apt-explorer/lib/apt/transport"
	"apt-explorer/lib/log"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"
	"text/template"
	"time"
)

func main() {
	//serve()

	arc, err := apt.UseArchiveHttp("http://archive.ubuntu.com/ubuntu")
	//arc, err := apt.UseArchiveFilesystem("./repo")
	catch(err)

	trn := transport.Fetcher(arc)
	trn = transport.UseCache(trn)

	jammy, err := apt.UseDistribution(trn, "jammy")
	catch(err)

	fmt.Printf("Distribution: %v\n", jammy.Name)

	log.Debug("Let's try downloading all the index files!")
	// TODO: there are many files to download. I'll need to use goroutines and workers.
	// and maybe prioritize which files get downloaded.

	keys := make([]string, 0)
	for k, _ := range jammy.Indexes {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	keys = keys[2:20]
	for i, k := range keys {
		vf := jammy.Indexes[k]
		log.Debug("%d/%d", i, len(jammy.Indexes))
		i++
		//log.Debug("Considering index file %s", k)
		_ = k

		z, err := trn.Fetch(vf)
		_, _ = z, err
		//catch(err)

		//x, err := io.ReadAll(z)
		//catch(err)

		//fmt.Printf("Read bytes = %d\n", len(x))

		time.Sleep(50 * time.Millisecond)
	}

	//
	//gzipReader, err := arc.Fetch("dists/trusty/main/source/Sources.gz")
	//catch(err)
	//reader, err := gzip.NewReader(gzipReader)
	//pqFile := "sources.pq"
	//err = kvstore.ConvertToParquet(reader, pqFile)
	//catch(err)
	//
	//matches, err := kvstore.QueryParquet(pqFile, "squid")
	//for _, m := range matches {
	//	fmt.Println(m.SingleValues["Package"])
	//}
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "/search")
	w.WriteHeader(302)
}

func searchPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		http.ServeFile(w, r, "search.html")
		return
	}

	const tmpl = `
{{range .}}{{.SingleValues.Package | printf "%35s"}}  {{.SingleValues.Version | printf "%-40s"}} {{.SingleValues.Architecture}}
{{end}}
`
	t := template.Must(template.New("tmpl").Funcs(template.FuncMap{
		"leftpad": func(width int, s string) string {
			return strings.Repeat(" ", len(s)-width) + s
		},
		"rightpad": func(width int, s string) string {
			return strings.Repeat(" ", len(s)-width) + s
		},
	}).Parse(strings.TrimSpace(tmpl)))

	w.Header().Set("Content-Type", "text/html")

	requestedPackage := r.URL.Query().Get("q")
	//requestedPackage := r.FormValue("package")
	pqFile := "sources.pq"
	matches, err := kvstore.QueryParquet(pqFile, requestedPackage)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	_, _ = fmt.Fprintln(w, "<pre>")
	t.Execute(w, matches)
	//for _, m := range matches {
	//	_, _ = fmt.Fprintln(w, m.SingleValues["Package"])
	//}
	_, _ = fmt.Fprintln(w, "</pre>")
}

func serve() {
	http.HandleFunc("/", index)
	http.HandleFunc("/search", searchPage)

	fmt.Printf("Starting server at port 8080\n")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(true, "%v", err)
	}

}

func catch(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}
