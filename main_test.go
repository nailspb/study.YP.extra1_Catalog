package main

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_loadFiles(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		want require.ValueAssertionFunc
	}{
		{
			name: "root file",
			want: requireFile("", "/"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := loadFiles()
			tt.want(t, got)
		})
	}
}

func Test_doCommand(t *testing.T) {
	t.Parallel()

	type args struct {
		cmd string
		arg string
	}

	type test struct {
		name       string
		args       args
		wantResult string
	}

	tests := []struct {
		name               string
		chain              []test
		wantNewCurrentFile require.ValueAssertionFunc
	}{
		{
			name: "no command",
			chain: []test{
				{
					name:       "run empty command from root",
					args:       args{},
					wantResult: "",
				},
			},
			wantNewCurrentFile: requireFile("", "/"),
		},
		{
			name: "unknown command",
			chain: []test{
				{
					name: "run `unknown` from root",
					args: args{
						cmd: "unknown",
					},
					wantResult: "Неизвестная команда: unknown\n",
				},
			},
			wantNewCurrentFile: requireFile("", "/"),
		},
		{
			name: "move to unknown directory",
			chain: []test{
				{
					name: "run `cd unknown` from root",
					args: args{
						cmd: "cd",
						arg: "unknown",
					},
					wantResult: "Неизвестная директория: unknown\n",
				},
			},
			wantNewCurrentFile: requireFile("", "/"),
		},
		{
			name: "move to empty directory",
			chain: []test{
				{
					name: "run `cd` from root",
					args: args{
						cmd: "cd",
					},
					wantResult: "",
				},
			},
			wantNewCurrentFile: requireFile("", "/"),
		},
		{
			name: "move to file",
			chain: []test{
				{
					name: "run `pwd` from root",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/\n",
				},
				{
					name: "run `ls` from root",
					args: args{
						cmd: "ls",
					},
					wantResult: "[d] Documents\n[d] Downloads\n[d] Movies\n[d] Music\n[d] Photos\n[f] я.jpg\n",
				},
				{
					name: "run `cd я.jpg` from root",
					args: args{
						cmd: "cd",
						arg: "я.jpg",
					},
					wantResult: "Неизвестная директория: я.jpg\n",
				},
			},
			wantNewCurrentFile: requireFile("", "/"),
		},
		{
			name: "move to root parent directory",
			chain: []test{
				{
					name: "run `pwd` from root",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/\n",
				},
				{
					name: "run `cd ..` from root",
					args: args{
						cmd: "cd",
						arg: "..",
					},
					wantResult: "",
				},
				{
					name: "run `pwd` from root",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/\n",
				},
			},
			wantNewCurrentFile: requireFile("", "/"),
		},
		{
			name: "files in Documents",
			chain: []test{
				{
					name: "run `pwd` from root",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/\n",
				},
				{
					name: "run `ls` from root",
					args: args{
						cmd: "ls",
					},
					wantResult: "[d] Documents\n[d] Downloads\n[d] Movies\n[d] Music\n[d] Photos\n[f] я.jpg\n",
				},
				{
					name: "run `cd Documents` from root",
					args: args{
						cmd: "cd",
						arg: "Documents",
					},
					wantResult: "",
				},
				{
					name: "run `pwd` from Documents",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/Documents/\n",
				},
				{
					name: "run `ls` from Documents",
					args: args{
						cmd: "ls",
					},
					wantResult: "[d] Books\n[d] Ипотека\n[f] Паспорт.pdf\n[f] СНИЛС.jpg\n",
				},
				{
					name: "run `cd Books` from Documents",
					args: args{
						cmd: "cd",
						arg: "Books",
					},
					wantResult: "",
				},
				{
					name: "run `pwd` from Books",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/Documents/Books/\n",
				},
				{
					name: "run `ls` from Books",
					args: args{
						cmd: "ls",
					},
					wantResult: "[f] Гарри Поттер.pdf\n[f] Мастер и Маргарита.epub\n",
				},
				{
					name: "run `cd ..` from Books",
					args: args{
						cmd: "cd",
						arg: "..",
					},
					wantResult: "",
				},
				{
					name: "run `pwd` from Documents",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/Documents/\n",
				},
				{
					name: "run `cd Ипотека` from Documents",
					args: args{
						cmd: "cd",
						arg: "Ипотека",
					},
					wantResult: "",
				},
				{
					name: "run `pwd` from Ипотека",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/Documents/Ипотека/\n",
				},
				{
					name: "run `ls` from Ипотека",
					args: args{
						cmd: "ls",
					},
					wantResult: "[f] НДФЛ.jpg\n",
				},
			},
			wantNewCurrentFile: requireFile("Ипотека", "/Documents/Ипотека/"),
		},
		{
			name: "files in Downloads",
			chain: []test{
				{
					name: "run `pwd` from root",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/\n",
				},
				{
					name: "run `ls` from root",
					args: args{
						cmd: "ls",
					},
					wantResult: "[d] Documents\n[d] Downloads\n[d] Movies\n[d] Music\n[d] Photos\n[f] я.jpg\n",
				},
				{
					name: "run `cd Downloads` from root",
					args: args{
						cmd: "cd",
						arg: "Downloads",
					},
					wantResult: "",
				},
				{
					name: "run `pwd` from Downloads",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/Downloads/\n",
				},
				{
					name: "run `ls` from Downloads",
					args: args{
						cmd: "ls",
					},
					wantResult: "[f] vscode.zip\n[f] winamp.exe\n",
				},
			},
			wantNewCurrentFile: requireFile("Downloads", "/Downloads/"),
		},
		{
			name: "files in Movies",
			chain: []test{
				{
					name: "run `pwd` from root",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/\n",
				},
				{
					name: "run `ls` from root",
					args: args{
						cmd: "ls",
					},
					wantResult: "[d] Documents\n[d] Downloads\n[d] Movies\n[d] Music\n[d] Photos\n[f] я.jpg\n",
				},
				{
					name: "run `cd Movies` from root",
					args: args{
						cmd: "cd",
						arg: "Movies",
					},
					wantResult: "",
				},
				{
					name: "run `pwd` from Movies",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/Movies/\n",
				},
				{
					name: "run `ls` from Movies",
					args: args{
						cmd: "ls",
					},
					wantResult: "[d] Ужасы\n[d] Фантастика\n[f] Любовь и голуби.mov\n",
				},
				{
					name: "run `cd Ужасы` from Movies",
					args: args{
						cmd: "cd",
						arg: "Ужасы",
					},
					wantResult: "",
				},
				{
					name: "run `pwd` from Ужасы",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/Movies/Ужасы/\n",
				},
				{
					name: "run `ls` from Ужасы",
					args: args{
						cmd: "ls",
					},
					wantResult: "[f] Молчание ягнят.mkv\n[f] Чужой.mp4\n",
				},
				{
					name: "run `cd ..` from Ужасы",
					args: args{
						cmd: "cd",
						arg: "..",
					},
					wantResult: "",
				},
				{
					name: "run `pwd` from Movies",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/Movies/\n",
				},
				{
					name: "run `cd Фантастика` from Movies",
					args: args{
						cmd: "cd",
						arg: "Фантастика",
					},
					wantResult: "",
				},
				{
					name: "run `pwd` from Фантастика",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/Movies/Фантастика/\n",
				},
				{
					name: "run `ls` from Фантастика",
					args: args{
						cmd: "ls",
					},
					wantResult: "[f] Аватар.mov\n",
				},
			},
			wantNewCurrentFile: requireFile("Фантастика", "/Movies/Фантастика/"),
		},
		{
			name: "files in Music",
			chain: []test{
				{
					name: "run `pwd` from root",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/\n",
				},
				{
					name: "run `ls` from root",
					args: args{
						cmd: "ls",
					},
					wantResult: "[d] Documents\n[d] Downloads\n[d] Movies\n[d] Music\n[d] Photos\n[f] я.jpg\n",
				},
				{
					name: "run `cd Music` from root",
					args: args{
						cmd: "cd",
						arg: "Music",
					},
					wantResult: "",
				},
				{
					name: "run `pwd` from Music",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/Music/\n",
				},
				{
					name: "run `ls` from Music",
					args: args{
						cmd: "ls",
					},
					wantResult: "",
				},
			},
			wantNewCurrentFile: requireFile("Music", "/Music/"),
		},
		{
			name: "files in Photos",
			chain: []test{
				{
					name: "run `pwd` from root",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/\n",
				},
				{
					name: "run `ls` from root",
					args: args{
						cmd: "ls",
					},
					wantResult: "[d] Documents\n[d] Downloads\n[d] Movies\n[d] Music\n[d] Photos\n[f] я.jpg\n",
				},
				{
					name: "run `cd Photos` from root",
					args: args{
						cmd: "cd",
						arg: "Photos",
					},
					wantResult: "",
				},
				{
					name: "run `pwd` from Photos",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/Photos/\n",
				},
				{
					name: "run `ls` from Photos",
					args: args{
						cmd: "ls",
					},
					wantResult: "[d] Выпускной\n[d] Свадьба\n[f] аватарка.png\n",
				},
				{
					name: "run `cd Выпускной` from Photos",
					args: args{
						cmd: "cd",
						arg: "Выпускной",
					},
					wantResult: "",
				},
				{
					name: "run `pwd` from Выпускной",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/Photos/Выпускной/\n",
				},
				{
					name: "run `ls` from Выпускной",
					args: args{
						cmd: "ls",
					},
					wantResult: "[f] 001.jpg\n[f] 003.jpg\n[f] 004.jpg\n",
				},
				{
					name: "run `cd ..` from Выпускной",
					args: args{
						cmd: "cd",
						arg: "..",
					},
					wantResult: "",
				},
				{
					name: "run `pwd` from Photos",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/Photos/\n",
				},
				{
					name: "run `cd Свадьба` from Photos",
					args: args{
						cmd: "cd",
						arg: "Свадьба",
					},
					wantResult: "",
				},
				{
					name: "run `pwd` from Свадьба",
					args: args{
						cmd: "pwd",
					},
					wantResult: "/Photos/Свадьба/\n",
				},
				{
					name: "run `ls` from Свадьба",
					args: args{
						cmd: "ls",
					},
					wantResult: "[f] wed_21.jpg\n[f] wed_22.jpg\n[f] wed_23.jpg\n[f] wed_27.jpg\n",
				},
			},
			wantNewCurrentFile: requireFile("Свадьба", "/Photos/Свадьба/"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			currentFile := loadFiles()

			for _, ttt := range tt.chain {
				t.Run(ttt.name, func(t *testing.T) {
					var gotResult string
					gotResult, currentFile = doCommand(ttt.args.cmd, ttt.args.arg, currentFile)
					require.Equal(t, ttt.wantResult, gotResult)
				})
			}

			tt.wantNewCurrentFile(t, currentFile)
		})
	}
}

func Test_readCommand(t *testing.T) {
	t.Parallel()

	type args struct {
		reader io.Reader
	}

	tests := []struct {
		name    string
		args    args
		wantCmd string
		wantArg string
		wantErr require.ErrorAssertionFunc
	}{
		{
			name: "no command",
			args: args{
				reader: bytes.NewBuffer([]byte("\n")),
			},
			wantCmd: "",
			wantArg: "",
			wantErr: require.NoError,
		},
		{
			name: "command without argument",
			args: args{
				reader: bytes.NewBuffer([]byte("pwd\n")),
			},
			wantCmd: "pwd",
			wantArg: "",
			wantErr: require.NoError,
		},
		{
			name: "command with argument",
			args: args{
				reader: bytes.NewBuffer([]byte("cd ..\n")),
			},
			wantCmd: "cd",
			wantArg: "..",
			wantErr: require.NoError,
		},
		{
			name: "command with spaced argument",
			args: args{
				reader: bytes.NewBuffer([]byte("cd Program Files\n")),
			},
			wantCmd: "cd",
			wantArg: "Program Files",
			wantErr: require.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			gotCmd, gotArg, err := readCommand(tt.args.reader)
			tt.wantErr(t, err)
			require.Equal(t, tt.wantCmd, gotCmd)
			require.Equal(t, tt.wantArg, gotArg)
		})
	}
}

func requireFile(wantName, wantPath string) require.ValueAssertionFunc {
	return func(t require.TestingT, gotFile interface{}, _ ...interface{}) {
		f, ok := gotFile.(File)
		require.True(t, ok)
		require.NotNil(t, f)
		require.Equal(t, wantName, f.Name())
		require.Equal(t, wantPath, f.Path())
	}
}
