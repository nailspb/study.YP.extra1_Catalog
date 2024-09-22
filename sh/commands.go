package sh

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type File interface {
	Name() string
	Path() string
	Child() []*File
	Parent() *File
	IsDirectory() bool
}

type ExecFunc = func(f File, args ...string) (string, File)

var fn map[string]ExecFunc = map[string]ExecFunc{
	"":     empty,
	"cd":   cd,
	"ls":   ls,
	"pwd":  pwd,
	"exit": exit,
	"quit": quit,
}

func empty(f File, args ...string) (string, File) {
	return "", f
}

func exit(f File, args ...string) (string, File) {
	os.Exit(0)
	return "", nil
}

func quit(f File, args ...string) (string, File) {
	return exit(f, args...)
}

func getFileByPath(start File, path string) File {
	path = filepath.Clean(path)
	cFile := start
	for _, dir := range strings.Split(path, string(os.PathSeparator)) {
		if dir == ".." {
			if cFile.Parent() != nil {
				cFile = *cFile.Parent()
			}
		} else {
			if fInd := slices.IndexFunc(cFile.Child(), func(file *File) bool {
				return (*file).Name() == dir && (*file).IsDirectory()
			}); fInd != -1 {
				cFile = *cFile.Child()[fInd]
			} else {
				return nil
			}
		}
	}
	return cFile
}

func cd(f File, args ...string) (string, File) {
	if len(args) > 0 && len(args[0]) > 0 {
		if f := getFileByPath(f, args[0]); f != nil {
			return "", f
		}
		return fmt.Sprintf("Неизвестная директория: %s\n", args[0]), f
	}
	return "", f
}

func ls(f File, args ...string) (string, File) {
	ff := f
	if len(args) > 0 {
		ff = getFileByPath(f, args[0])
		if ff == nil {
			ff = f
		}
	}
	if len(ff.Child()) > 0 {
		b := strings.Builder{}
		for _, c := range ff.Child() {
			if (*c).IsDirectory() {
				b.WriteString("[d] ")
			} else {
				b.WriteString("[f] ")
			}
			b.WriteString((*c).Name())
			b.WriteString("\n")
		}
		return b.String()[:b.Len()], f
	}
	return "", f
}

func pwd(f File, args ...string) (string, File) {
	return fmt.Sprintf("%s\n", f.Path()), f
}

type Exec struct {
	fn ExecFunc
}

func Cmd(cmd string) *Exec {
	if _, ok := fn[cmd]; ok {
		return &Exec{fn: fn[cmd]}
	}
	return &Exec{
		fn: func(f File, args ...string) (string, File) {
			return fmt.Sprintf("Неизвестная команда: %s\n", cmd), f
		},
	}

}

func (r *Exec) Execute(f File, args string) (string, File) {
	return r.fn(f, args)
}
