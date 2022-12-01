package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"sort"
	"strings"
	"sync"
	"sync/atomic"
	"unsafe"
)

const (
	// Max recursive depth for directory scanning.
	gMAX_SCAN_DEPTH = 100000
)

var (

	regexMu = sync.RWMutex{}

	// Cache for regex object.
	// Note that:
	// 1. It uses sync.RWMutex ensuring the concurrent safety.
	// 2. There's no expiring logic for this map.
	regexMap = make(map[string]*regexp.Regexp)

	// goRootForFilter is used for stack filtering purpose.
	goRootForFilter = runtime.GOROOT()

	// The absolute file path for main package.
	// It can be only checked and set once.
	mainPkgPath = NewString()

	// defaultTrimChars are the characters which are stripped by Trim* functions in default.
	defaultTrimChars = string([]byte{
		'\t', // Tab.
		'\v', // Vertical tab.
		'\n', // New line (line feed).
		'\r', // Carriage return.
		'\f', // New page.
		' ',  // Ordinary space.
		0x00, // NUL-byte.
		0x85, // Delete.
		0xA0, // Non-breaking space.
	})

	// Separator for file system.
	// It here defines the separator as variable
	// to allow it modified by developer if necessary.
	Separator = string(filepath.Separator)
)

func init() {
	if goRootForFilter != "" {
		goRootForFilter = strings.Replace(goRootForFilter, "\\", "/", -1)
	}
}

// String is a struct for concurrent-safe operation for type string.
type String struct {
	value atomic.Value
}

// NewString creates and returns a concurrent-safe object for string type,
// with given initial value <value>.
func NewString(value ...string) *String {
	t := &String{}
	if len(value) > 0 {
		t.value.Store(value[0])
	}
	return t
}

// Val atomically loads and returns t.value.
func (v *String) Val() string {
	s := v.value.Load()
	if s != nil {
		return s.(string)
	}
	return ""
}

// Set atomically stores <value> into t.value and returns the previous value of t.value.
func (v *String) Set(value string) (old string) {
	old = v.Val()
	v.value.Store(value)
	return
}

// IsMatchString checks whether given string <src> matches <pattern>.
func IsMatchString(pattern string, src string) bool {
	return IsMatch(pattern, []byte(src))
}

// IsMatch checks whether given bytes <src> matches <pattern>.
func IsMatch(pattern string, src []byte) bool {
	if r, err := getRegexp(pattern); err == nil {
		return r.Match(src)
	}
	return false
}

// getRegexp returns *regexp.Regexp object with given <pattern>.
// It uses cache to enhance the performance for compiling regular expression pattern,
// which means, it will return the same *regexp.Regexp object with the same regular
// expression pattern.
//
// It is concurrent-safe for multiple goroutines.
func getRegexp(pattern string) (regex *regexp.Regexp, err error) {
	// Retrieve the regular expression object using reading lock.
	regexMu.RLock()
	regex = regexMap[pattern]
	regexMu.RUnlock()
	if regex != nil {
		return
	}
	// If it does not exist in the cache,
	// it compiles the pattern and creates one.
	regex, err = regexp.Compile(pattern)
	if err != nil {
		return
	}
	// Cache the result object using writing lock.
	regexMu.Lock()
	regexMap[pattern] = regex
	regexMu.Unlock()
	return
}

// Ext returns the file name extension used by path.
// The extension is the suffix beginning at the final dot
// in the final element of path; it is empty if there is
// no dot.
//
// Note: the result contains symbol '.'.
func Ext(path string) string {
	ext := filepath.Ext(path)
	if p := strings.IndexByte(ext, '?'); p != -1 {
		ext = ext[0:p]
	}
	return ext
}

// UnsafeBytesToStr converts []byte to string without memory copy.
// Note that, if you completely sure you will never use <b> variable in the feature,
// you can use this unsafe function to implement type conversion in high performance.
func UnsafeBytesToStr(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

// GetBytes returns the file content of <path> as []byte.
// It returns nil if it fails reading.
func GetBytes(path string) []byte {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil
	}
	return data
}

// GetContents returns the file content of <path> as string.
// It returns en empty string if it fails reading.
func GetContents(path string) string {
	return UnsafeBytesToStr(GetBytes(path))
}

// Dir returns all but the last element of path, typically the path's directory.
// After dropping the final element, Dir calls Clean on the path and trailing
// slashes are removed.
// If the path is empty, Dir returns ".".
// If the path consists entirely of separators, Dir returns a single separator.
// The returned path does not end in a separator unless it is the root directory.
func Dir(path string) string {
	return filepath.Dir(path)
}

// Exists checks whether given <path> exist.
func Exists(path string) bool {
	if stat, err := os.Stat(path); stat != nil && !os.IsNotExist(err) {
		return true
	}
	return false
}

// Trim strips whitespace (or other characters) from the beginning and end of a string.
// The optional parameter <characterMask> specifies the additional stripped characters.
func Trim(str string, characterMask ...string) string {
	if len(characterMask) == 0 {
		return strings.Trim(str, defaultTrimChars)
	} else {
		return strings.Trim(str, defaultTrimChars+characterMask[0])
	}
}

// SplitAndTrim splits string <str> by a string <delimiter> to an array,
// and calls Trim to every element of this array. It ignores the elements
// which are empty after Trim.
func SplitAndTrim(str, delimiter string, characterMask ...string) []string {
	array := make([]string, 0)
	for _, v := range strings.Split(str, delimiter) {
		v = Trim(v, characterMask...)
		if v != "" {
			array = append(array, v)
		}
	}
	return array
}

// IsDir checks whether given <path> a directory.
// Note that it returns false if the <path> does not exist.
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// Abs returns an absolute representation of path.
// If the path is not absolute it will be joined with the current
// working directory to turn it into an absolute path. The absolute
// path name for a given file is not guaranteed to be unique.
// Abs calls Clean on the result.
func Abs(path string) string {
	p, _ := filepath.Abs(path)
	return p
}

// doScanDir is an internal method which scans directory and returns the absolute path
// list of files that are not sorted.
//
// The pattern parameter <pattern> supports multiple file name patterns, using the ','
// symbol to separate multiple patterns.
//
// The parameter <recursive> specifies whether scanning the <path> recursively, which
// means it scans its sub-files and appends the files path to result array if the sub-file
// is also a folder. It is false in default.
//
// The parameter <handler> specifies the callback function handling each sub-file path of
// the <path> and its sub-folders. It ignores the sub-file path if <handler> returns an empty
// string, or else it appends the sub-file path to result slice.
func doScanDir(depth int, path string, pattern string, recursive bool, handler func(path string) string) ([]string, error) {
	if depth >= gMAX_SCAN_DEPTH {
		return nil, errors.New(fmt.Sprintf("directory scanning exceeds max recursive depth: %d", gMAX_SCAN_DEPTH))
	}
	list := ([]string)(nil)
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	names, err := file.Readdirnames(-1)
	if err != nil {
		return nil, err
	}
	var (
		filePath = ""
		patterns = SplitAndTrim(pattern, ",")
	)
	for _, name := range names {
		filePath = path + Separator + name
		if IsDir(filePath) && recursive {
			array, _ := doScanDir(depth+1, filePath, pattern, true, handler)
			if len(array) > 0 {
				list = append(list, array...)
			}
		}
		// Handler filtering.
		if handler != nil {
			filePath = handler(filePath)
			if filePath == "" {
				continue
			}
		}
		// If it meets pattern, then add it to the result list.
		for _, p := range patterns {
			if match, err := filepath.Match(p, name); err == nil && match {
				filePath = Abs(filePath)
				if filePath != "" {
					list = append(list, filePath)
				}
			}
		}
	}
	return list, nil
}

// ScanDir returns all sub-files with absolute paths of given <path>,
// It scans directory recursively if given parameter <recursive> is true.
//
// The pattern parameter <pattern> supports multiple file name patterns,
// using the ',' symbol to separate multiple patterns.
func ScanDir(path string, pattern string, recursive ...bool) ([]string, error) {
	isRecursive := false
	if len(recursive) > 0 {
		isRecursive = recursive[0]
	}
	list, err := doScanDir(0, path, pattern, isRecursive, nil)
	if err != nil {
		return nil, err
	}
	if len(list) > 0 {
		sort.Strings(list)
	}
	return list, nil
}

// MainPkgPath returns absolute file path of package main,
// which contains the entrance function main.
//
// It's only available in develop environment.
//
// Note1: Only valid for source development environments,
// IE only valid for systems that generate this executable.
//
// Note2: When the method is called for the first time, if it is in an asynchronous goroutine,
// the method may not get the main package path.
func MainPkgPath() string {
	// Only for source development environments.
	if goRootForFilter == "" {
		return ""
	}
	path := mainPkgPath.Val()
	if path != "" {
		return path
	}
	lastFile := ""
	for i := 1; i < 10000; i++ {
		if _, file, _, ok := runtime.Caller(i); ok {
			if goRootForFilter != "" && len(file) >= len(goRootForFilter) && file[0:len(goRootForFilter)] == goRootForFilter {
				continue
			}
			if IsMatchString(`/github.com/[^/]+/gf/`, file) &&
				!IsMatchString(`/github.com/[^/]+/gf/\.example/`, file) {
				continue
			}
			if Ext(file) != ".go" {
				continue
			}
			lastFile = file
			if IsMatchString(`package\s+main`, GetContents(file)) {
				mainPkgPath.Set(Dir(file))
				return Dir(file)
			}
		} else {
			break
		}
	}
	if lastFile != "" {
		for path = Dir(lastFile); len(path) > 1 && Exists(path) && path[len(path)-1] != os.PathSeparator; {
			files, _ := ScanDir(path, "*.go")
			for _, v := range files {
				if IsMatchString(`package\s+main`, GetContents(v)) {
					mainPkgPath.Set(path)
					return path
				}
			}
			path = Dir(path)
		}
	}
	return ""
}