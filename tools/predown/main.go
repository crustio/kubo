package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ipfs/kubo/repo/fsrepo/migrations"
	"golang.org/x/exp/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// randString returns a random string of length n
func randString(length int) string {
	seed := rand.NewSource(uint64(time.Now().UnixNano()))
	r := rand.New(seed)
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[r.Intn(len(charset))]
	}
	return string(result)
}

func main() {
	// p example : fs-repo-12-to-13/v1.0.0/fs-repo-12-to-13_v1.0.0_darwin-arm64.tar.gz
	subPath := flag.String("p", "", "Sub Path of Qmf4uwxfKzdMouLHTnd8B2ejcANVnE8pmSsgHzLYRutKhV")
	tarFilePath := flag.String("t", "", "File path inside tar.gz")
	outputPath := flag.String("o", "output", "Path to the output file")

	// parse flags
	flag.Parse()

	// check args
	if *subPath == "" || *outputPath == "" || *tarFilePath == "" {
		fmt.Println("Usage: -p <sub path> -o <output file path> -t <tar file path>")
		return
	}

	tmpTarFile := fmt.Sprintf("./downloaded_file-%s.tar.gz", randString(10))

	fc := migrations.NewHttpFetcher("", "", "", 0)

	bs, err := fc.Fetch(context.Background(), *subPath)
	if err != nil {
		fmt.Println("fetch faul: ", err)
		return
	}

	// save bs to file
	err = os.WriteFile(tmpTarFile, bs, 0644)
	if err != nil {
		fmt.Println("write file failed: ", err)
		return
	}
	// clear tmp file
	defer func() {
		os.Remove(tmpTarFile)
		fmt.Println("Temporary file deleted successfully: ", tmpTarFile)
	}()
	fmt.Println("File downloaded successfully: ", tmpTarFile)

	binName := filepath.Base(*tarFilePath)
	root := filepath.Dir(*tarFilePath)

	// final output file path
	outputFilePath := filepath.Join(*outputPath, binName)

	// untar to output path
	err = unpackArchive(tmpTarFile, "tar.gz", root, binName, outputFilePath)
	if err != nil {
		fmt.Println("untar failed: ", err)
		return
	}

	fmt.Println("File unpacked successfully: ", outputFilePath)
}
