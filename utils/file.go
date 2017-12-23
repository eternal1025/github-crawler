package utils

import "os"

func AppendStringToFile(path, text string, newLine bool) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}

	if newLine {
		text += "\n"
	}
	_, err = f.WriteString(text)
	if err != nil {
		return err
	}
	return f.Close()
}

func CreateFileIfNotExist(path string) error {
	_, err := os.Stat(path)
	if os.IsExist(err) {
		return nil
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	return f.Close()
}