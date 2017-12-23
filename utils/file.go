package utils

import "os"

func AppendStringToFile(path, text string, newLine bool) error {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY|os.O_CREATE, os.ModeAppend)
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