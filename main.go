package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func parseTiming(s string) (float64, int, error) {
	t := strings.Split(s, " ")

	if len(t) != 2 {
		return 0, 0, errors.New(fmt.Sprintf("timing file has wrong length of row: %s", t))
	}

	delay, err := strconv.ParseFloat(t[0], 64)
	if err != nil {
		return 0, 0, errors.New(fmt.Sprintf("delay cannot be parsed: %s", t[0]))
	}

	length, err := strconv.Atoi(t[1])
	if err != nil {
		return 0, 0, errors.New(fmt.Sprintf("length cannot be parsed: %s", t[1]))
	}

	return delay, length, nil
}

func main() {

	log.SetFlags(log.Lshortfile)

	var textPath, timingPath string
	flag.StringVar(&textPath, "text-file", "", "Path to the text logfile. Can be omitted if timing logfile is aswell.")
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
		_, length, err := parseTiming(timingScanner.Text())
		if err != nil {
			log.Fatal(err)
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