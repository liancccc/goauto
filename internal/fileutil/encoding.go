package fileutil

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"unicode/utf8"
)

// UTF8BOM 是 UTF-8 字节序标记
var UTF8BOM = []byte{0xEF, 0xBB, 0xBF}

// WriteStringWithoutBOM 写入字符串到文件，不添加 BOM
func WriteStringWithoutBOM(filename string, content string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content)
	if err != nil {
		return err
	}
	return writer.Flush()
}

// WriteBytesWithoutBOM 写入字节到文件，不添加 BOM
func WriteBytesWithoutBOM(filename string, data []byte) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

// AppendStringWithoutBOM 追加字符串到文件，不添加 BOM
func AppendStringWithoutBOM(filename string, content string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.WriteString(content)
	if err != nil {
		return err
	}
	return writer.Flush()
}

// RemoveBOM 从字节数据中移除 UTF-8 BOM
func RemoveBOM(data []byte) []byte {
	if len(data) >= 3 && bytes.Equal(data[:3], UTF8BOM) {
		return data[3:]
	}
	return data
}

// HasBOM 检查字节数据是否包含 UTF-8 BOM
func HasBOM(data []byte) bool {
	return len(data) >= 3 && bytes.Equal(data[:3], UTF8BOM)
}

// IsValidUTF8 检查字节数据是否为有效的 UTF-8 编码
func IsValidUTF8(data []byte) bool {
	return utf8.Valid(data)
}

// CleanUTF8String 清理字符串，移除 BOM 并确保是有效的 UTF-8
func CleanUTF8String(s string) string {
	// 移除字符串开头的 BOM (Unicode BOM 字符)
	if len(s) > 0 {
		runes := []rune(s)
		if len(runes) > 0 && runes[0] == '\uFEFF' {
			s = string(runes[1:])
		}
	}
	return s
}

// WriteSliceToFileUTF8 写入字符串切片到文件，确保 UTF-8 编码且无 BOM
func WriteSliceToFileUTF8(path string, lines []string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		// 清理每行，移除可能的 BOM
		cleanLine := CleanUTF8String(line)
		_, err := writer.WriteString(cleanLine + "\n")
		if err != nil {
			return err
		}
	}
	return writer.Flush()
}

// NewWriterWithoutBOM 创建一个不添加 BOM 的 bufio.Writer
func NewWriterWithoutBOM(w io.Writer) *bufio.Writer {
	return bufio.NewWriter(w)
}
