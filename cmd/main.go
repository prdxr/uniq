package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	var isModeCount, isModeDupe, isModeUnique, isCaseFree bool
	var nSkipFields, nSkipChars int

	flag.BoolVar(&isModeCount, "c", false, "count each unique string")
	flag.BoolVar(&isModeDupe, "d", false, "list only duplicate strings")
	flag.BoolVar(&isModeUnique, "u", false, "list only unique strings")
	flag.BoolVar(&isCaseFree, "i", false, "ignore case")
	flag.IntVar(&nSkipFields, "f", -1, "ignore N first fields (separated by whitespaces)")
	flag.IntVar(&nSkipChars, "s", -1, "ignore N first symbols")
	flag.Parse()

	var filepathInput, filepathOutput = flag.Arg(0), flag.Arg(1)

	var lines, err1 = readInput(filepathInput, isCaseFree, nSkipFields, nSkipChars)
	if err1 != nil {
		return
	}

	var mode rune
	switch {
	case isModeCount && (isModeDupe || isModeUnique) || (isModeDupe && isModeUnique):
		fmt.Printf("ERROR: Cannot use -c -d -u flags together")
		return
	case isModeCount:
		mode = 'c'
	case isModeDupe:
		mode = 'd'
	case isModeUnique:
		mode = 'u'
	}

	var linesMap, err2 = processLines(lines, mode)
	if err2 != nil {
		return
	}

	writeOutput(linesMap, filepathOutput, isModeCount)

	fmt.Printf("[DEBUG]		flags passed:")
	fmt.Printf("\n[DEBUG]		-c %v	-d %v	-u %v	-i %v	-f %v	-s %v", isModeCount, isModeDupe, isModeUnique, isCaseFree, nSkipFields, nSkipChars)
	fmt.Printf("\n[DEBUG]		input: %v,	output: %v", filepathInput, filepathOutput)
	fmt.Printf("\n[DEBUG]		%v len:%v", lines, len(lines))
	fmt.Printf("\n[DEBUG]		%v len:%v", linesMap, len(linesMap))
}

func readInput(filepathInput string, isCaseFree bool, nSkipFields int, nSkipChars int) (lines []string, err error) {
	if filepathInput != "" {
		readFile, err := os.Open(filepathInput)
		if err != nil {
			return nil, err
		}
		defer readFile.Close()

		scanner := bufio.NewScanner(readFile)
		for scanner.Scan() {
			currentLine := scanner.Text()

			currentLine, nSkipFields, nSkipChars = prepareLine(currentLine, isCaseFree, nSkipFields, nSkipChars)

			if currentLine != "" {
				lines = append(lines, currentLine)
			}
		}

		if err := scanner.Err(); err != nil {
			return nil, err
		}

	} else {
		fmt.Println("Input text (empty line for exit):")
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			currentLine := scanner.Text()
			if currentLine == "" {
				break
			}

			currentLine, nSkipFields, nSkipChars = prepareLine(currentLine, isCaseFree, nSkipFields, nSkipChars)

			if currentLine != "" {
				lines = append(lines, currentLine)
			}
		}
	}

	return lines, nil
}

func prepareLine(line string, isCaseFree bool, nSkipFields int, nSkipChars int) (string, int, int) {
	if isCaseFree {
		line = strings.ToLower(line)
	}

	lineFields := strings.Fields(line)
	nFields := len(lineFields)
	if nSkipFields > nFields {
		nSkipFields -= nFields
		return "", nSkipFields, nSkipChars
	} else if nSkipFields > 0 {
		line = strings.Join(lineFields[nSkipFields:], " ")
		nSkipFields = 0
	}

	nChars := len(line)
	if nSkipChars > nChars {
		nSkipChars -= nChars
		return "", nSkipFields, nSkipChars
	} else if nSkipChars > 0 {
		line = line[nSkipChars:]
		nSkipChars = 0
	}

	return line, nSkipFields, nSkipChars
}

func processLines(lines []string, mode rune) (linesMap map[string]int, err error) {
	linesMap = make(map[string]int)
	for _, v := range lines {
		linesMap[v] += 1
	}

	switch mode {
	case 'd':
		//show only dupes
		for k, v := range linesMap {
			if v == 1 {
				delete(linesMap, k)
			}
		}
	case 'u':
		//show only uniques
		for k, v := range linesMap {
			if v != 1 {
				delete(linesMap, k)
			}
		}
	}
	return linesMap, err
}

func writeOutput(linesMap map[string]int, filepathOutput string, isCountMode bool) error {
	if filepathOutput != "" {
		writeFile, err := os.Create(filepathOutput)
		if err != nil {
			return err
		}

		for k, v := range linesMap {
			if isCountMode {
				writeFile.WriteString(fmt.Sprintf("%v - %v\n", k, v))
			} else {
				writeFile.WriteString(fmt.Sprintf("%v\n", k))
			}
		}

	} else {
		for k, v := range linesMap {
			if isCountMode {
				fmt.Printf("%v - %v\n", k, v)
			} else {
				fmt.Println(k)
			}
		}
	}
	return nil
}
