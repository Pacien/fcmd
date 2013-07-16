/*

	This file is part of fcmd (https://github.com/Pacien/fcmd)

	Permission is hereby granted, free of charge, to any person obtaining a copy
	of this software and associated documentation files (the "Software"), to deal
	in the Software without restriction, including without limitation the rights
	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
	copies of the Software, and to permit persons to whom the Software is
	furnished to do so, subject to the following conditions:

	The above copyright notice and this permission notice shall be included in
	all copies or substantial portions of the Software.

	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
	THE SOFTWARE.

*/

// Common file manipulation commands for Go.
package fcmd

import (
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
)

var DefaultPerm os.FileMode = 0750 // u=rwx, g=r-x, o=---

// Checks if the target exists.
func IsExist(target string) bool {
	_, err := os.Stat(target)
	return os.IsNotExist(err)
}

// Checks if the target is a directory.
// Returns false if the target is unreachable.
func IsDir(target string) bool {
	stat, err := os.Stat(target)
	if err != nil {
		return false
	}
	return stat.IsDir()
}

// Checks lexically if the target is hidden (only for Unix based OS).
func IsHidden(target string) bool {
	return strings.HasPrefix(target, ".")
}

// Lists separately the names of directories and files inside the target directory.
// Hidden files and directories are not listed.
func Ls(target string) (dirs, files []string) {
	directory, err := ioutil.ReadDir(target)
	if err != nil {
		return
	}
	for _, element := range directory {
		if IsHidden(element.Name()) {
			continue
		}
		if element.IsDir() {
			dirs = append(dirs, element.Name())
		} else {
			files = append(files, element.Name())
		}
	}
	return
}

// Lists separately the paths of directories and files inside the root directory and inside all sub directories.
// Returned paths are relative to the given root directory.
// Hidden files and directories are not listed.
func Explore(root string) (dirs, files []string) {
	dirList, fileList := Ls(root)

	for _, file := range fileList {
		files = append(files, file)
	}

	for _, dir := range dirList {
		subRoot := path.Join(root, dir)
		dirs = append(dirs, subRoot)
		subDirs, subFiles := Explore(subRoot)
		for _, subFile := range subFiles {
			files = append(files, subFile)
		}
		for _, subDir := range subDirs {
			dirs = append(dirs, subDir)
		}
	}
	return
}

// Copies the source file to a target.
// A nonexistent target file is created, otherwise it is truncated.
// Parent directories are automatically created if they do not exist.
func Cp(source, target string) error {
	sourceFile, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	dir, _ := path.Split(target)

	err = os.MkdirAll(dir, DefaultPerm)
	if err != nil {
		return err
	}

	targetFile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	_, err = io.Copy(targetFile, sourceFile)
	return err
}

// Writes data to the target file.
// A nonexistent target file is created, otherwise it is truncated.
// Parent directories are automatically created if they do not exist.
func WriteFile(target string, data []byte) error {
	dir, _ := path.Split(target)

	err := os.MkdirAll(dir, DefaultPerm)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(target, data, DefaultPerm)
	return err
}

// Creates a symbolic link to given source at the target path.
func Lns(source, target string) error {
	return os.Symlink(source, target)
}

// Returns the destination of the given symbolic link.
func Lnl(target string) (string, error) {
	return os.Readlink(target)
}

// Renames or moves the source file or directory to the target name or path.
func Mv(source, target string) error {
	return os.Rename(source, target)
}

// Removes the target file or the target directory and all files it contains.
// No error is returned is the target does not exist.
func Rm(target string) error {
	return os.RemoveAll(target)
}

// Changes the current working directory to the target directory.
func Cd(target string) error {
	return os.Chdir(target)
}

// Changes the mode of the target file to the given mode.
// If the target is a symbolic link, it changes the mode of the link's target.
func Chmod(target string, mode os.FileMode) error {
	return os.Chmod(target, mode)
}

// Changes the numeric uid and gid of the target.
// If the target is a symbolic link, it changes the uid and gid of the link's target.
func Chown(target string, uid, gid int) error {
	return os.Chown(target, uid, gid)
}

// Changes the access and modification times of the target.
func Chtimes(target string, atime time.Time, mtime time.Time) error {
	return os.Chtimes(target, atime, mtime)
}
