package sh

import (
	"cmp"
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

func sortChild(child []*FileT) {
	slices.SortFunc(child, func(a, b *FileT) int {
		if a.isDirectory && b.isDirectory {
			return cmp.Compare(a.name, b.name)
		} else if a.isDirectory && !b.isDirectory {
			return -1
		} else if !a.isDirectory && b.isDirectory {
			return 1
		} else {
			return cmp.Compare(a.name, b.name)
		}
	})
}

func Parse(fileNames []string, directories map[int]struct{}, parents map[int]int) File {
	root := &FileT{
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
			sortChild(files[parent].child)
			files[ind].parent = &files[parent]
		} else {
			root.child = append(root.child, &files[ind])
			files[ind].parent = root
			sortChild(root.child)
		}
	}
	return root
}

func (f *FileT) Name() string {
	if f.name != "/" {
		return f.name
	}
	return ""
}

func (f *FileT) Path() string {
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

func (f *FileT) Child() []File {
	//cast to interface type
	c := make([]File, len(f.child))
	for ind, v := range f.child {
		c[ind] = v
	}
	return c
}

func (f *FileT) IsDirectory() bool {
	return f.isDirectory
}

func (f *FileT) Parent() File {
	if f.parent != nil {
		return f.parent
	}
	return nil
}
