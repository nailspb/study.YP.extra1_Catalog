package sh

import (
	"path"
	"slices"
)

type FileT struct {
	index       int
	name        string
	isDirectory bool
	parent      *FileT
	child       []*FileT
}

func Parse(fileNames []string, directories map[int]struct{}, parents map[int]int) File {
	root := FileT{
		index:       0,
		name:        "/",
		isDirectory: true,
		parent:      nil,
		child:       make([]*FileT, 0),
	}
	files := make([]FileT, len(fileNames))
	for ind, name := range fileNames {
		files[ind].index = ind
		files[ind].name = name
		if _, ok := directories[ind]; ok {
			files[ind].isDirectory = true
		}
		if parent, ok := parents[ind]; ok {
			files[parent].child = append(files[parent].child, &files[ind])
			files[ind].parent = &files[parent]
		} else {
			root.child = append(root.child, &files[ind])
			files[ind].parent = &root
		}
	}
	return root
}

func (f FileT) Name() string {
	if f.name != "/" {
		return f.name
	}
	return ""
}

func (f FileT) Path() string {
	el := []string{f.name}
	parent := f.parent
	for parent != nil {
		el = append(el, parent.name)
		parent = parent.parent
	}
	slices.Reverse(el)
	p := path.Join(el...)
	if f.isDirectory && f.name != "/" {
		p += "/"
	}
	return p
}

func (f FileT) Child() []*File {
	c := make([]*File, len(f.child))
	for ind, v := range f.child {
		castInterface := File(*v)
		c[ind] = &castInterface
	}
	return c
}

func (f FileT) IsDirectory() bool {
	return f.isDirectory
}

func (f FileT) Parent() *File {
	if f.parent != nil {
		castInterface := File(*(f.parent))
		return &castInterface
	}
	return nil
}
