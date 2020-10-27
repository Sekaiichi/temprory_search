package search

import (
	"context"
	"io/ioutil"
	"log"
	"regexp"
	"strings"
	"sync"
)

//Result describes one search result
type Result struct {
	Phrase  string //searched phrase
	Line    string //the line of search entry (without \n in the end)
	LineNum int64  //line number (starting from 1) of the search entry
	ColNum  int64  //position number (starting from 1) of the search entry
}

//All searches all {phrase} entries in {files} text files
func All(ctx context.Context, phrase string, files []string) <-chan []Result {
	ch := make(chan []Result)
	wg := sync.WaitGroup{}
	regex, _ := regexp.Compile(phrase)

	for _, file := range files {
		wg.Add(1)
		go func(found chan<- []Result, textFile string) {
			defer wg.Done()
			content := make([]byte, 0)
			content, rerr := ioutil.ReadFile(textFile)
			if rerr != nil {
				log.Printf("File: %v, Error: %v", textFile, rerr)
			}

			text := string(content)
			lines := strings.Split(text, "\n")
			tempRes := make([]Result, 0)
			res := Result{Phrase: phrase}

			for lineNum, line := range lines {
				if line == "" {
					continue
				}
				if !strings.Contains(line, phrase) {
					continue
				}
				indexes := regex.FindAllStringIndex(line, -1)

				res.Line = line
				res.LineNum = int64(lineNum + 1)
				for _, index := range indexes {
					res.ColNum = int64(index[0] + 1)
					tempRes = append(tempRes, res)
				}
			}
			if len(tempRes) > 0 {
				found <- tempRes
			}
		}(ch, file)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	return ch
}

//Any searches all {phrase} entries in {files} text files
func Any(ctx context.Context, phrase string, files []string) <-chan Result {
	ch := make(chan Result, 1)
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(ctx)
	regex, _ := regexp.Compile(phrase)
	for i, file := range files {
		wg.Add(1)

		go func(ctxt context.Context, file string, i int) {

			defer wg.Done()
			content, rerr := ioutil.ReadFile(file)
			if rerr != nil {
				log.Printf("File: %v, Error: %v", file, rerr)

			}
			text := string(content)
			lines := strings.Split(text, "\n")
			res := Result{Phrase: phrase}

			for lineNum, line := range lines {

				select {
				case <-ctx.Done():
					return
				default:
				}

				if line == "" {
					continue
				}
				if !strings.Contains(line, phrase) {
					continue
				}
				indexes := regex.FindAllStringIndex(line, -1)

				res.Line = line
				res.LineNum = int64(lineNum + 1)
				for _, index := range indexes {
					res.ColNum = int64(index[0] + 1)

					select {
					case <- ctx.Done():
						return
					default:
						ch <- res
						cancel()
						return
					}
				}
			}
		}(ctx, file, i)
	}

	wg.Wait()
	defer close(ch)
	cancel()
	return ch
}