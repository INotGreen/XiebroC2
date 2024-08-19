package HandlePacket

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return err
	}
	return nil
}

func readFileToString(path string) (string, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	str := string(b)
	return str, nil
}

func listDir(path string) (string, error) {
	// Replace any "\\" with "/", and remove "*"
	dirPathStr := strings.ReplaceAll(path, "\\", "/")
	dirPathStr = strings.ReplaceAll(dirPathStr, "*", "")

	// Get info about the directory
	fileInfo, err := os.Stat(dirPathStr)
	if err != nil {
		return "", err
	}

	modTime := fileInfo.ModTime()
	modTimeStr := modTime.Format("02/01/2006 15:04:05")

	files, err := ioutil.ReadDir(dirPathStr)
	if err != nil {
		return "", err
	}

	fileInfos := []string{
		//fmt.Sprintf(".-=>%s-=>D-=>0-=>%s", modTimeStr, fileInfo.Mode().Perm()),
		//fmt.Sprintf("..-=>%s-=>D-=>0-=>%s", modTimeStr, fileInfo.Mode().Perm()),
	}

	for _, file := range files {
		modTimeStr = file.ModTime().Format("02/01/2006 15:04:05")

		var fileType string
		if file.IsDir() {
			fileType = "D"
		} else {
			fileType = "F"
		}

		fileInfo := fmt.Sprintf("%s-=>%s-=>%s-=>%d-=>%s", file.Name(), modTimeStr, fileType, file.Size(), file.Mode().Perm())
		fileInfos = append(fileInfos, fileInfo)
	}

	resultStr := strings.Join(fileInfos, "-=>")

	return resultStr, nil
}
func CopyFile(srcFile, dstFile string) error {
	input, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(dstFile, input, 0644)
	if err != nil {
		return err
	}
	return nil
}

// 文件粘贴函数（与复制函数实质上是一样的，只是表述不同）
func PasteFile(srcFile, dstFile string) error {
	return CopyFile(srcFile, dstFile)
}

// 文件剪切函数
func CutFile(srcFile, dstFile string) error {
	err := CopyFile(srcFile, dstFile)
	if err != nil {
		return err
	}
	err = os.Remove(srcFile) // 删除原文件
	if err != nil {
		return err
	}
	return nil
}

// 文件重命名函数
func RenameFile(oldName, newName string) error {
	err := os.Rename(oldName, newName)
	if err != nil {
		return err
	}
	return nil
}
func Zip(src, dest string) error {

	zipfile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()
	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = path

		header.Method = zip.Deflate

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}
		if !info.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()
			_, err = io.Copy(writer, file)
		}
		return err
	})

	return err
}
func Unzip(src, dest string) error {
	// Open the zip archive for reading
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {

		fpath := filepath.Join(dest, f.Name)
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("%s: illegal file path", fpath)
		}

		// If it's a directory, create it
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			// Create a file for writing
			if err := os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}

			rc, err := f.Open()
			if err != nil {
				return err
			}

			_, err = io.Copy(outFile, rc)

			outFile.Close()
			rc.Close()

			if err != nil {
				return err
			}
		}
	}
	return nil
}
