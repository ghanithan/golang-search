package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"regexp"
	//	log "github.com/sirupsen/logrus"
)

const(
	mailFilePath string = "/Users/ghanithan/work/enron-dataset-abridged/"
)

var(
	count int = 0
	fileCount int =0
	indexFile *os.File
	//removeSpecialChar = regexp.MustCompile(`([a-z A-Z 0-9])+`)
	removeSpecialChar = regexp.MustCompile(`([^a-z A-Z 0-9])+`)
)

func main(){

	fmt.Println("Entering Init...")
	runtime.GOMAXPROCS(6)

	var err error
	indexFile, err = os.OpenFile("search.index", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		//log.Fatal(err)
	}
	filePathsChan := make(chan string, 1000)
	errChan := make(chan error,10)
	go traverseAllFiles(mailFilePath, filePathsChan)	

	for filePath := range filePathsChan{
		go processFiles(filePath, errChan)
		fmt.Println("count: " + fmt.Sprintf("%d",count))
		fmt.Println("FileCount: " + fmt.Sprintf("%d",fileCount))

	}
/*
	for {
		select {
			case filePath := <- filePathsChan:
				go processFiles(filePath, errChan)
			case errFile := <- errChan:
				fmt.Println(errFile)		
			default:
				fmt.Println("count: " + fmt.Sprintf("%d",count))
				fmt.Println("FileCount: " + fmt.Sprintf("%d",fileCount))
			   continue

		}
	}
*/
	if err := indexFile.Close(); err != nil {

		fmt.Println(err)
	}
//	close(errChan)
//	close(filePathsChan)

}

func processFiles(filePath string, errChan chan error) {
	data, err := os.ReadFile(filePath)
	count++
	if err != nil {
		fmt.Println(err)
		errChan <- err
		return
	}
	fileAsString := string(data)
	_ = tokenize(fileAsString, filePath)
//	fmt.Println(string(fileAsString))
}

func traverseAllFiles(path string, filePathsChan chan string){

	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error{
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() || strings.Contains( path, "DS_Store" )  {
			fmt.Printf("visited dir: %q\n", path)
			return nil
		}
		fileCount++
		fmt.Printf("visited file: %q\n", path)
		filePathsChan <- path
		return nil

	})
	fmt.Println("Total Number of Files: " + fmt.Sprintf("%d", fileCount))

	if err != nil {
		fmt.Printf("error walking the path %q: %v\n", path, err)
		return
	}
	close(filePathsChan)


}

func appendIndex(indexOneFile map[string]int, filePath string){
	singleFile := ""
	for key,val := range indexOneFile{
		singleFile = singleFile + fmt.Sprintf("%s,%d,%s", key,val,filePath) + "\n"
	}

	if _, err := indexFile.Write([]byte(singleFile));
	err != nil {
		indexFile.Close() // ignore error; Write error takes precedence
		fmt.Println(err)
	}
}


func tokenize(inputString string, filePath string) error{
	indexOneFile := make(map[string]int)

	//removeSpecialChar.ReplaceAllString(inputString,"")
	tokensInit := strings.Split(inputString, "\n")
	tempTokens := strings.Join(tokensInit,"")
	tokens := strings.Split(tempTokens, " ")
	for _,token := range tokens {
		if !removeSpecialChar.Match([]byte(token)){
			token = strings.ToLower(token)
			indexOneFile[token]++	
		}
	}
	fmt.Printf("%v", indexOneFile)
	appendIndex(indexOneFile, filePath)
	return nil

}
