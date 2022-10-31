package main

import (
	"apt-explorer/lib/apt"
	"apt-explorer/lib/apt/kvstore"
	"apt-explorer/lib/apt/transport"
	"apt-explorer/lib/log"
	"compress/gzip"
	"fmt"
	"net/http"
	"os"
	"strings"
	"text/template"
)

func main() {
	serve()

	arc, err := apt.UseArchiveHttp("http://archive.ubuntu.com/ubuntu")
	catch(err)
	ubuntuArchiveCache := transport.UseCache(arc, "jammy")
	jammy, err := apt.UseDistribution(ubuntuArchiveCache, "jammy")
	catch(err)

	// TODO: there are many files to download. I'll need to use goroutines and workers.
	// and maybe prioritize which files get downloaded.

	components := []string{"main"}
	architectures := []string{"amd64"}
	for _, vf := range jammy.PackagesFiles(components, architectures) {
		fmt.Println(vf)
		rr, err := ubuntuArchiveCache.Fetch(vf)
		fmt.Println(err)
		gzr, _ := gzip.NewReader(rr)
		err = kvstore.Parse(gzr, func(block kvstore.Block) {
			fmt.Println(block.SingleValues["Package"])
		})
		if err != nil {
			fmt.Println(err)
			return
		}
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

// FIXME: this entire function is hardcoding and hacks
func search(needle string, archive string, dist string, component string, arch string) []kvstore.Block {
	arc, err := apt.UseArchiveHttp(archive)
	catch(err)
	ubuntuArchiveCache := transport.UseCache(arc, dist)
	jammy, err := apt.UseDistribution(ubuntuArchiveCache, dist)
	catch(err)

	components := []string{component}
	architectures := []string{arch}

	matches := make([]kvstore.Block, 0)
	for _, vf := range jammy.PackagesFiles(components, architectures) {
		fmt.Println(vf)
		rr, err := ubuntuArchiveCache.Fetch(vf)
		fmt.Println(err)
		gzr, _ := gzip.NewReader(rr)
		err = kvstore.Parse(gzr, func(block kvstore.Block) {
			if strings.Contains(block.SingleValues["Package"], needle) {
				matches = append(matches, block)
			}
		})
		if err != nil {
			log.Error("%v", err)
			return matches
		}
	}

	return matches
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

	// the other way of searching
	//pqFile := "sources.pq"
	//matches, _ := kvstore.QueryParquet(pqFile, requestedPackage)

	// the modern way
	matches := search(requestedPackage, "http://archive.ubuntu.com/ubuntu", "jammy", "main", "amd64")

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
