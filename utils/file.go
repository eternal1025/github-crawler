package utils

import "os"

func AppendStringToFile(path, text string, newLine bool) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModePerm)
	defer f.Close()
	if err != nil {
		return err
	}

	if newLine {
		text += "\n"
	}
	_, err = f.WriteString(text)
	return err
}