package main

var (
	// fileNames Имена файлов
	fileNames = []string{
		0:  "001.jpg",
		1:  "003.jpg",
		2:  "004.jpg",
		3:  "Books",
		4:  "Documents",
		5:  "Downloads",
		6:  "Movies",
		7:  "Music",
		8:  "Photos",
		9:  "vscode.zip",
		10: "wed_21.jpg",
		11: "wed_22.jpg",
		12: "wed_23.jpg",
		13: "wed_27.jpg",
		14: "winamp.exe",
		15: "Аватар.mov",
		16: "аватарка.png",
		17: "Выпускной",
		18: "Гарри Поттер.pdf",
		19: "Ипотека",
		20: "Любовь и голуби.mov",
		21: "Мастер и Маргарита.epub",
		22: "Молчание ягнят.mkv",
		23: "НДФЛ.jpg",
		24: "Паспорт.pdf",
		25: "Свадьба",
		26: "СНИЛС.jpg",
		27: "Ужасы",
		28: "Фантастика",
		29: "Чужой.mp4",
		30: "я.jpg",
	}
	// directories Директории: ключ — индексы файла, который являются директорией
	directories = map[int]struct{}{
		3:  {},
		4:  {},
		5:  {},
		6:  {},
		7:  {},
		8:  {},
		17: {},
		19: {},
		25: {},
		27: {},
		28: {},
	}
	// fileParents Родительские директории файлов: ключ — индекс файла, значение — индекс родительской директории
	// Если для файла нет нужного ключа, значит, файл является директорией, лежащей в корне /
	fileParents = map[int]int{
		0:  17,
		1:  17,
		2:  17,
		3:  4,
		9:  5,
		10: 25,
		11: 25,
		12: 25,
		13: 25,
		14: 5,
		15: 28,
		16: 8,
		17: 8,
		18: 3,
		19: 4,
		20: 6,
		21: 3,
		22: 27,
		23: 19,
		24: 4,
		25: 8,
		26: 4,
		27: 6,
		28: 6,
		29: 27,
	}
)