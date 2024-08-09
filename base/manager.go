package base

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"
)

var (
	activeManager *FileManager
	managerLock   sync.Mutex
)

func SetCurrentFileManager(manager *FileManager) {
	managerLock.Lock()
	defer managerLock.Unlock()
	activeManager = manager
}

func GetCurrentFileManager() *FileManager {
	managerLock.Lock()
	defer managerLock.Unlock()
	return activeManager
}

type FileManager struct {
	ID         string
	WorkingDir string
	Files      map[string]*File
	Pwd        string
	Recent     *File
}

type Options struct {
	Recursive       bool
	CaseInsensitive bool
}
type Option func(*Options)

func WithRecursive(recursive bool) Option {
	return func(opts *Options) {
		opts.Recursive = recursive
	}
}

func WithCaseInsensitive(caseInsensitive bool) Option {
	return func(opts *Options) {
		opts.CaseInsensitive = caseInsensitive
	}
}

func NewFileManager(workingDir string) *FileManager {
	if workingDir == "" {
		workingDir, _ = os.Getwd()
	}
	fm := &FileManager{
		ID:         generateID(),
		WorkingDir: workingDir,
		Files:      make(map[string]*File),
	}
	return fm
}

func (fm *FileManager) enter() *FileManager {
	//""
	//"Enter workspace context."
	//""
	activeManager = GetCurrentFileManager()
	if activeManager != nil && activeManager.ID != fm.ID {
		panic("Another manager already activated via context.")
	}
	fm.Pwd, _ = os.Getwd()
	err := os.Chdir(fm.WorkingDir)
	if err != nil {
		return nil
	}
	SetCurrentFileManager(fm)
	return fm
}

func (fm *FileManager) exit(args ...any) {
	//"""Exit from workspace context."""
	if fm.Pwd != "" {
		err := os.Chdir(fm.Pwd)
		if err != nil {
			return
		}
		SetCurrentFileManager(nil)
	}
}

func (fm *FileManager) resolveDir(dir string) (abspath string, err error) {
	//"""Resolve a directory path."""
	if filepath.IsAbs(dir) {
		abspath, err = filepath.Abs(dir)
	} else {
		abspath, err = filepath.Abs(filepath.Join(fm.WorkingDir, dir))
	}
	return
}

func (fm *FileManager) resolveDirs(dirs []string) ([]string, error) {
	results := make([]string, 0)
	for _, dir := range dirs {
		temp, err := fm.resolveDir(dir)
		if err != nil {
			return nil, err
		}
		results = append(results, temp)
	}
	return results, nil
}

func (fm *FileManager) Chdir(path string) error {
	newDir, err := fm.resolveDir(path)
	if err != nil {
		return fmt.Errorf("invalid path: %v", err)
	}

	// Ensure the resolved path is within the allowed directory structure
	anchorc := filepath.VolumeName(fm.WorkingDir)
	anchorn := filepath.VolumeName(filepath.Clean(newDir))
	if anchorc != anchorn {
		return fmt.Errorf("access denied: cannot navigate to '%s'", newDir)
	}

	if isFile(newDir) {
		return fmt.Errorf("'%s' is not a valid directory", newDir)
	}

	fm.WorkingDir = newDir
	return nil
}

func (fm *FileManager) Open(path string) (*File, error) {
	absPath := filepath.Join(fm.WorkingDir, path)
	if file, exists := fm.Files[absPath]; exists {
		return file, nil
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist", absPath)
	}

	file := &File{Path: absPath, Workdir: fm.WorkingDir}
	fm.Files[absPath] = file
	fm.Recent = file
	return file, nil
}

func (fm *FileManager) Create(path string) (*File, error) {
	absPath := filepath.Join(fm.WorkingDir, path)
	newFile, err := os.Create(absPath)
	if err != nil {
		return nil, fmt.Errorf("could not create file %s: %v", absPath, err)
	}
	err = newFile.Close()
	if err != nil {
		return nil, err
	}

	file := NewFile(absPath, fm.WorkingDir, 0)
	fm.Files[absPath] = file
	fm.Recent = file
	return file, nil
}

func (fm *FileManager) Grep(word string, pattern string, options ...Option) (map[string][]Match, error) {
	opts := Options{
		Recursive:       true,
		CaseInsensitive: true,
	}

	// 应用传入的选项
	for _, option := range options {
		option(&opts)
	}
	if pattern == "" {
		pattern = fm.WorkingDir
	} else {
		pattern = filepath.Clean(pattern)
	}
	pathsToSearch, err := getPathsToSearch(pattern, opts.Recursive)
	if err != nil {
		return nil, err
	}
	results := make(map[string][]Match)
	for _, filePath := range pathsToSearch {
		if isFile(filePath) && filepath.Base(filePath)[0] != '.' {
			f, err := os.Open(filePath)
			if err != nil {
				return nil, err
			}

			scanner := bufio.NewScanner(f)
			lineNumber := 1

			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(formatWord(line, opts.CaseInsensitive), formatWord(word, opts.CaseInsensitive)) {
					relPath, err := filepath.Rel(fm.WorkingDir, filePath)
					if err != nil {
						return nil, err
					}
					if _, exists := fm.Files[relPath]; !exists {
						results[relPath] = make([]Match, 1)
					}
					results[relPath] = append(results[relPath], Match{Content: strings.TrimSpace(line),
						Lineno: lineNumber})
				}
				lineNumber++
			}
		}
	}
	if len(results) == 0 {
		//fmt.Sprintf("No matches found for %v in %v",word,pattern )
	}
	numMatches := 0
	for _, match := range results {
		numMatches += len(match)
	}
	//fmt.Sprintf("Found %v matches for \"%v\" in %v", numMatches, pattern, fm.WorkingDir)
	return results, nil
}

func (fm *FileManager) Find(pattern string, depth int, caseSensitive bool, include []string, exclude []string) ([]string, error) {
	includePaths, err := fm.resolveDirs(include)
	if err != nil {
		return nil, err
	}
	if len(includePaths) == 0 {
		includePaths = append(includePaths, fm.WorkingDir)
	}
	exclude = append(exclude, ".git")
	excludePaths, err := fm.resolveDirs(exclude)
	if err != nil {
		return nil, err
	}
	var regex *regexp.Regexp
	if !caseSensitive {
		regex = regexp.MustCompile(pattern)
	} else {
		regex = regexp.MustCompile("(?i)" + pattern)
	}

	matches := make([]string, 0)
	var searchRecursive func(directory string, currentDepth int)
	searchRecursive = func(directory string, currentDepth int) {
		if depth != 0 && currentDepth > depth {
			return
		}
		entries, err := os.ReadDir(directory)
		if err != nil {
			if os.IsPermission(err) {
				return // 跳过没有权限访问的目录
			}
			fmt.Println("Error reading directory:", err)
			return
		}

		for _, entry := range entries {
			itemPath := filepath.Join(directory, entry.Name())
			absItemPath, err := filepath.Abs(itemPath)
			if err != nil {
				fmt.Println("Error resolving absolute path:", err)
				continue
			}
			isParentExcluded := func(path string, excludePath string) bool {
				for path != "/" && path != "." {
					if filepath.Dir(path) == excludePath {
						return true
					}
					path = filepath.Dir(path)
				}
				return false
			}
			excluded := false
			for _, excludePath := range excludePaths {
				if absItemPath == excludePath || isParentExcluded(absItemPath, excludePath) {
					excluded = true
					break
				}
			}
			if excluded {
				continue
			}

			relativePath, err := filepath.Rel(fm.WorkingDir, absItemPath)
			if err != nil {
				fmt.Println("Error getting relative path:", err)
				continue
			}

			if regex.MatchString(relativePath) {
				matches = append(matches, relativePath)
			}

			if entry.IsDir() {
				searchRecursive(absItemPath, currentDepth+1)
			}
		}
	}
	for _, directory := range includePaths {
		searchRecursive(directory, 0)
	}
	slices.Sort(matches)
	return matches, nil
}

func formatWord(word string, caseInsensitive bool) string {
	if caseInsensitive {
		return strings.ToLower(word)
	}
	return word
}

func getPathsToSearch(pattern string, recursive bool) (pathsToSearch []string, err error) {
	// 获取要搜索的路径
	pathsToSearch = make([]string, 0)
	if isFile(pattern) && err == nil {
		if recursive {
			err = filepath.Walk(pattern, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				pathsToSearch = append(pathsToSearch, path)
				return nil
			})
		} else {
			pathsToSearch, err = filepath.Glob(pattern)
		}
	} else if err == nil {
		pathsToSearch = append(pathsToSearch, pattern)
	}
	return
}

func (fm *FileManager) tree(
	directory string,
	level int,
	depth int,
	exclude []string,
) string {
	//"""Auxialiary method for creating working directory tree recursively."""
	if (depth != -1 && level > depth) || slices.Contains(exclude, directory) {
		return ""
	}
	tree := ""
	childs, err := os.ReadDir(directory)
	if err != nil {
		fmt.Println("Error reading directory:", err)
	}
	for _, child := range childs {
		path := filepath.Join(directory, child.Name())
		if isFile(path) && err == nil {
			tree += strings.Repeat("  |", level) + "__ " + filepath.Base(path) + "\n"
		}
	}
	for _, child := range childs {
		path := filepath.Join(directory, child.Name())
		if info, err := os.Stat(path); info.IsDir() && err == nil {
			continue
		}
		tree += strings.Repeat("  |", level) + "__ " + filepath.Base(path) + "\n"
		tree += fm.tree(path, level+1, depth, exclude)
	}
	return tree
}
func (fm *FileManager) Tree(
	depth int,
	exclude []string,
) string {
	//"""
	//Create directory tree for the file
	//
	//:param depth: Max depth for the tree
	//:param exclude: Exclude directories from the tree
	//"""
	return fm.tree(
		fm.WorkingDir,
		0,
		depth,
		exclude,
	)
}

func (fm *FileManager) ls() [][2]string {
	//"""List contents of the current directory with their types."""
	result := make([][2]string, 0)
	childs, err := os.ReadDir(fm.WorkingDir)
	if err != nil {
		return nil
	}
	for _, child := range childs {
		if child.IsDir() {
			result = append(result, [2]string{child.Name(), "dir"})
		} else {
			result = append(result, [2]string{child.Name(), "file"})
		}
	}

	return result
}

func (fm *FileManager) ExecuteCommand(command string) (string, error) {
	//"""Execute a command in the current working directory."""
	// 创建命令对象
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = fm.WorkingDir

	// 捕获输出
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// 设置超时时间
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()

	select {
	case err := <-done:
		if err != nil {
			return "", fmt.Errorf("error executing command: %s", stderr.String())
		}
		return out.String(), nil
	case <-time.After(120 * time.Second):
		// 如果命令超时，杀死进程
		if err := cmd.Process.Kill(); err != nil {
			return "", fmt.Errorf("failed to kill process: %s", err)
		}
		return "", fmt.Errorf("TIMEOUT: Command execution timed out after 120 seconds")
	}
}
