package Helper

import (
	"bufio"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func deleteStringFromFile(filePath, strToDelete string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 删掉含有指定字符串的整行
		if !strings.Contains(line, strToDelete) {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filePath, []byte(output), 0644)
	if err != nil {
		return err
	}
	return nil
}

func containsStringInFile(filePath, searchStr string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), searchStr) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, err
	}

	return false, nil
}

func checkStringInDirectoryFile(filePath, searchString string) bool {
	res, err := containsStringInFile(filePath, searchString)
	if err != nil {
		//fmt.Println("Error reading file:", filePath, err)
		return false
	}
	return res
}

func ConvertGBKToUTF8(s string) (string, error) {
	reader := transform.NewReader(strings.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}
	return string(d), nil
}
