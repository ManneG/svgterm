package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	log.SetFlags(log.Lshortfile)

	var textPath = "foo"
	flag.StringVar(&textPath, "text-file", "", "Path to the text logfile. Can be omitted if timing logfile is aswell.")

	var timingPath = ""
	flag.StringVar(&timingPath, "timing-file", "", "Path to the timing logfile. Can be omitted if text logfile is aswell.")
	
	flag.Parse()

	if ((timingPath == "") != (textPath == "")) {
		log.Fatal("Only one of text file and timing file are specified. Both or neither need to be supplied.")
	}

	if (textPath == "") {
		log.Fatal("Capturing is not yet implemented. Use `script` from `util-linux` and pass the files as args.")
	}

	textFile, err := os.Open(textPath)
	if err != nil {
		log.Fatal(err)
	}
	defer textFile.Close()
	header := ""
	buf := make([]byte, 1)
	for _, err := textFile.Read(buf); err == nil && buf[0] != '\n'; _, err = textFile.Read(buf) {
		header += string(buf)
	}

	fmt.Println(header)

	timingFile, err := os.Open(timingPath)
	if err != nil {
		log.Fatal(err)
	}
	defer timingFile.Close()
	timingScanner := bufio.NewScanner(timingFile)

	for timingScanner.Scan() {
		t := strings.Split(timingScanner.Text(), " ")

		if len(t) != 2 {
			log.Fatalf("timing file has wrong content: %s", t)
		}

		length, err := strconv.Atoi(t[1])
		if err != nil {
			log.Fatalf("timing file, length cannot be parsed: %s", t[1])
		}

		buf = make([]byte, length)
		if _, err := io.ReadFull(textFile, buf); err != nil {
			log.Fatal(err)
		}
		content := string(buf)
		content = strings.ReplaceAll(content, string(0x1B), "ESC")
		fmt.Println(content)
		os.Stdout.Sync()
	}
}