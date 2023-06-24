package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/LuizTrMr/gof/finder"
)

var notProcessed []string = []string{}

func main() {
	var searchTerm string
	flag.StringVar(&searchTerm, "st", "", "Term to be searched")

	if len(os.Args) == 2 { // You can just pass the search term and use other defaults
		searchTerm = os.Args[1]
	}

	var path string
	flag.StringVar(&path, "path", ".", "Folder/file to search for the search term")

	var exclude string
	flag.StringVar(&exclude, "exclude", "", "Folders/files to ignore while searching, separated by a comma")

	var threaded bool
	flag.BoolVar(&threaded, "go", false, "Use go routines to search files")

	flag.Parse()

	if searchTerm == "" {
		log.Fatalln("You must pass a Search Term")
	}

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Fatalln("ERROR (Couldn't read initial path):", err)
	}

	info, err := file.Stat()
	if err != nil {
		log.Fatalln("ERROR (Couldn't get initial file info):", err)
	}

	excludes := strings.Split(exclude, ",")
	if !threaded {
		if info.IsDir() {
			scanDirectory(path, searchTerm, excludes)
		} else {
			find(path, searchTerm)
		}
	} else {
		if info.IsDir() {
			var wg sync.WaitGroup
			scanDirectoryThreaded(path, searchTerm, excludes, &wg)
			wg.Wait()
		} else { // Only searching 1 file
			find(path, searchTerm)
		}
	}

	if len(notProcessed) > 0 {
		fmt.Println()
		fmt.Println("------ Not Processed Summary ------")
		fmt.Print(strings.Join(notProcessed, ""))
	}
}

const (
	RED_AND_BOLD = "\u001b[1;31m"
	RESET        = "\u001b[0;22m"
)

func find(fileName, searchTerm string) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		notProcessed = append(notProcessed, fmt.Sprintln("ERROR (Couldn't read file): ", err))
	}

	// Empty file
	if len(data) == 0 {
		return
	}

	// Tokenize
	t := finder.NewFinder(data)
	sb := strings.Builder{}

	var positions []finder.Pos
	stop := false
	for !stop {
		i := t.Bol()
		line := t.Lines()
		positions, stop = t.NextLine(searchTerm)
		if len(positions) > 0 {
			sb.WriteString(fmt.Sprintf("(%v)[Line %v]: ", fileName, line))

			for _, p := range positions {
				if i < p.Start { // Write until start of token
					sb.WriteString(t.Data(i, p.Start))
					i = p.Start
				}
				term := t.Data(p.Start, p.End)
				sb.WriteString(RED_AND_BOLD)
				sb.WriteString(term)
				sb.WriteString(RESET)
				i = p.End
			}

			sb.WriteString(t.Data(i, t.Bol()-1)) // Write until end of line
			sb.WriteByte('\n')
		}
	}

	fmt.Print(sb.String())

}

func findThreaded(fileName, searchTerm string, wg *sync.WaitGroup) {
	defer wg.Done()
	data, err := os.ReadFile(fileName)
	if err != nil {
		notProcessed = append(notProcessed, fmt.Sprintln("ERROR (Couldn't read file): ", err))
	}

	// Empty file
	if len(data) == 0 {
		return
	}

	// Tokenize
	t := finder.NewFinder(data)
	sb := strings.Builder{}

	var positions []finder.Pos
	stop := false
	for !stop {
		i := t.Bol()
		line := t.Lines()
		positions, stop = t.NextLine(searchTerm)
		if len(positions) > 0 {
			sb.WriteString(fmt.Sprintf("(%v)[Line %v]: ", fileName, line))

			for _, p := range positions {
				if i < p.Start { // Write until start of token
					sb.WriteString(t.Data(i, p.Start))
					i = p.Start
				}
				term := t.Data(p.Start, p.End)
				sb.WriteString(RED_AND_BOLD)
				sb.WriteString(term)
				sb.WriteString(RESET)
				i = p.End
			}

			sb.WriteString(t.Data(i, t.Bol()-1)) // Write until end of line
			sb.WriteByte('\n')
		}
	}

	fmt.Print(sb.String())
}

func scanDirectoryThreaded(path, searchTerm string, excludes []string, wg *sync.WaitGroup) {
	if shouldSkip(excludes, path) {
		return
	}
	files, err := os.ReadDir(path)
	if err != nil {
		notProcessed = append(notProcessed, fmt.Sprintln("ERROR (Couldn't read directory):", err))
	}

	for _, file := range files {
		fullPath := path + "/" + file.Name()
		if shouldSkip(excludes, fullPath) { // NOTE: Could move this inside `find`
			continue
		}
		if !file.IsDir() {
			if isTextFile(fullPath) {
				wg.Add(1)
				go findThreaded(fullPath, searchTerm, wg)
			} else { // Binary or something else
				continue
			}
		} else {
			scanDirectoryThreaded(fullPath, searchTerm, excludes, wg)
		}
	}
}

func scanDirectory(path, searchTerm string, excludes []string) {
	if shouldSkip(excludes, path) {
		return
	}
	files, err := os.ReadDir(path)
	if err != nil {
		notProcessed = append(notProcessed, fmt.Sprintln("ERROR (Couldn't read directory):", err))
	}

	for _, file := range files {
		fullPath := path + "/" + file.Name()
		if shouldSkip(excludes, fullPath) { // NOTE: Could move this inside `find`
			continue
		}
		if !file.IsDir() {
			if isTextFile(fullPath) {
				find(fullPath, searchTerm)
			} else { // Binary or something else
				continue
			}
		} else {
			scanDirectory(fullPath, searchTerm, excludes)
		}
	}
}

func shouldSkip(excludes []string, compare string) bool {
	for _, exclude := range excludes {
		if compare == exclude {
			// Ignore exclude folder or file
			// Could check for contains instead of a perfect match,
			// but that creates other problems so I chose the latter
			return true
		}
	}
	return false
}

func isTextFile(path string) bool {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		notProcessed = append(notProcessed, fmt.Sprintln("ERROR (Could not open file): ", err))
	}
	var buf [512]byte
	n, err := file.Read(buf[:])
	for i := 0; i < n; i += 1 {
		if buf[i] == 0 {
			return false
		}
	}
	return true
}
