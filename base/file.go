package base

import (
	"bufio"
	"errors"
	"os"
	"regexp"
	"strings"
)

type ScrollDirection string
type FileOperationScope string

const (
	ScrollUp   ScrollDirection = "up"
	ScrollDown ScrollDirection = "down"

	ScopeFile   FileOperationScope = "file"
	ScopeWindow FileOperationScope = "window"
)

type Match struct {
	Content string
	Match   string
	End     int
	Start   int
	Lineno  int
}

type TextReplacement struct {
	ReplacedWith string
	ReplacedText string
	Error        error
}

type File struct {
	Path    string
	Workdir string
	Start   int
	End     int
	Window  int
}

func (sd *ScrollDirection) Offset(lines int) int {
	if *sd == ScrollUp {
		return lines * (-1)
	} else {
		return lines
	}
}

func NewFile(path string, workdir string, window int) *File {
	if window == 0 {
		window = 100
	}
	return &File{
		Path:    path,
		Workdir: workdir,
		Start:   0,
		End:     window,
		Window:  window,
	}
}

func (f *File) Scroll(lines int, direction ScrollDirection) {
	lines = direction.Offset(lines)
	f.Start += lines
	f.End += lines
}

func (f *File) Goto(line int) {
	f.Start = line
	f.End = line + f.Window
}

func (f *File) find(buffer string, pattern string, lineno int) []Match {
	var matches []Match
	re := regexp.MustCompile(pattern)
	for _, match := range re.FindAllStringIndex(buffer, -1) {
		start, end := match[0], match[1]
		matches = append(matches, Match{
			Content: buffer,
			Match:   buffer[start:end],
			Start:   start,
			End:     end,
			Lineno:  lineno,
		})
	}
	return matches
}

func (f *File) findWindow(pattern string) []Match {
	offset := f.Start
	var matches []Match
	for lineno, line := range f.iterWindow() {
		matches = append(matches, f.find(line, pattern, lineno+offset)...)
	}
	return matches
}

func (f *File) findFile(pattern string) []Match {
	var matches []Match
	for lineno, line := range f.iterFile() {
		matches = append(matches, f.find(line, pattern, lineno)...)
	}
	return matches
}

func (f *File) Find(pattern string, scope FileOperationScope) []Match {
	if scope == ScopeFile {
		return f.findFile(pattern)
	}
	return f.findWindow(pattern)
}

func (f *File) iterWindow() []string {
	var lines []string
	file, _ := os.Open(f.Path)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
	scanner := bufio.NewScanner(file)
	cursor := 0
	for scanner.Scan() {
		if cursor >= f.Start && cursor < f.End {
			lines = append(lines, scanner.Text())
		}
		cursor++
	}
	return lines
}

func (f *File) iterFile() []string {
	var lines []string
	file, _ := os.Open(f.Path)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines
}

func (f *File) Read() map[int]string {
	cursor := 0
	buffer := make(map[int]string)
	file, _ := os.Open(f.Path)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
		}
	}(file)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if cursor >= f.Start && cursor < f.End {
			buffer[cursor+1] = scanner.Text()
		}
		cursor++
	}
	return buffer
}

func (f *File) Write(text string) error {
	return os.WriteFile(f.Path, []byte(text), 0644)
}

func (f *File) TotalLines() int {
	return len(f.iterFile())
}

func (f *File) Edit(text string, start int, end int, scope FileOperationScope) TextReplacement {
	originalContent, _ := os.ReadFile(f.Path)
	content := string(originalContent)
	lines := strings.Split(content, "\n")

	var buffer strings.Builder
	replaced := ""

	if scope == ScopeWindow {
		start = max(start, f.Start)
		end = min(end, f.End)
	}

	for i := 0; i < start-1; i++ {
		buffer.WriteString(lines[i] + "\n")
	}
	for i := start - 1; i < end; i++ {
		replaced += lines[i] + "\n"
	}
	buffer.WriteString(text + "\n")
	for i := end; i < len(lines); i++ {
		buffer.WriteString(lines[i] + "\n")
	}

	err := os.WriteFile(f.Path, []byte(buffer.String()), 0644)
	if err != nil {
		return TextReplacement{Error: err}
	}
	return TextReplacement{
		ReplacedText: replaced,
		ReplacedWith: text,
	}
}

func (f *File) Replace(search string, replacement string) TextReplacement {
	content, _ := os.ReadFile(f.Path)
	updated := strings.ReplaceAll(string(content), search, replacement)
	if string(content) == updated {
		return TextReplacement{
			ReplacedText: "",
			ReplacedWith: "",
			Error:        errors.New("error replacing given string, string not found"),
		}
	}
	err := os.WriteFile(f.Path, []byte(updated), 0644)
	if err != nil {
		return TextReplacement{}
	}
	return TextReplacement{
		ReplacedText: search,
		ReplacedWith: replacement,
	}
}

func (f *File) WriteAndRunLint(text string, start int, end int) TextReplacement {
	olderFileText, _ := os.ReadFile(f.Path)
	writeResponse := f.Edit(text, start, end, ScopeWindow)
	if writeResponse.Error != nil {
		_ = os.WriteFile(f.Path, olderFileText, 0644)
		return writeResponse
	}
	return writeResponse
}

func isFile(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
