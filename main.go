package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"practicum/basic/tasks/catalog/sh"
	"strings"
	"time"
)

type File = sh.File

func main() {
	currentFile := loadFiles()

	for {
		fmt.Printf("[%s] %s $ ", time.Now().Format("02.01 15:04"), currentFile.Path())

		cmd, arg, err := readCommand(os.Stdin)
		if err != nil {
			fmt.Printf("Ошибка ввода: %v", err)
			continue
		}

		var result string
		result, currentFile = doCommand(cmd, arg, currentFile)
		fmt.Print(result)
	}
}

// loadFiles загружает дерево файлов из data.go
// Возвращает назад корневой файл rootFile — директорию /
func loadFiles() File {
	return sh.Parse(fileNames, directories, fileParents)
}

// doCommand выполняет команду cmd с аргументом arg (может быть пустым) над текущим файлом currentFile
// Возвращает назад сообщение о результате result (может быть пустым) и текущий после выполнения команды файл newCurrentFile
func doCommand(cmd, arg string, currentFile File) (result string, newCurrentFile File) {
	result, newCurrentFile = sh.Cmd(cmd).Execute(currentFile, arg)
	return
}

// readCommand читает пользовательский ввод через reader
// Возвращает назад команду cmd, ее аргумент arg (может быть пустым) и ошибку (не пустая, если сломался ввод)
func readCommand(reader io.Reader) (cmd string, arg string, err error) {
	value, err := bufio.NewReader(reader).ReadString('\n')
	if err != nil {
		return
	}

	value = value[:len(value)-1]

	parts := strings.Split(value, " ")
	switch len(parts) {
	case 0:
		return
	case 1:
		cmd = parts[0]
	default:
		cmd = parts[0]
		arg = value[len(parts[0])+1:]
	}

	return
}
