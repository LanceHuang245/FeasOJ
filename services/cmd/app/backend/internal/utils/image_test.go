package utils

import (
	"github.com/stretchr/testify/assert"
	"image"
	"image/png"
	"os"
	"testing"
)

func TestCompressImage(t *testing.T) {
	inputPath := "test_input.png"
	outputPath := "test_output.png"

	// 创建测试输入文件
	file, err := os.Create(inputPath)
	assert.NoError(t, err)
	defer file.Close()

	// 写入空的PNG文件头
	err = png.Encode(file, image.NewRGBA(image.Rect(0, 0, 1, 1)))
	assert.NoError(t, err)

	// 测试
	err = CompressImage(inputPath, outputPath)
	assert.NoError(t, err)

	// 检查输出文件是否存在
	_, err = os.Stat(outputPath)
	assert.NoError(t, err)

	// 清理测试文件
	os.Remove(inputPath)
	os.Remove(outputPath)
}
