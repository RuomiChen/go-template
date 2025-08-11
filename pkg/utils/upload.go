package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// 传入一个哈希缓存，用于判断文件是否已存在
type HashStore interface {
	Get(hash string) (string, bool)
	Set(hash, path string)
}

func UploadImageWithHashCheck(c *fiber.Ctx, formField, saveDir string, allowExts []string, store HashStore) (string, error) {
	file, err := c.FormFile(formField)
	if err != nil {
		return "", err
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := false
	for _, allowExt := range allowExts {
		if ext == allowExt {
			allowed = true
			break
		}
	}
	if !allowed {
		return "", fmt.Errorf("不支持的图片格式: %s", ext)
	}

	// 打开文件流计算哈希
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	hasher := md5.New()
	if _, err := io.Copy(hasher, src); err != nil {
		return "", err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))

	// 查哈希缓存
	if path, exists := store.Get(hash); exists {
		return path, nil
	}

	// 目录不存在先创建
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return "", err
	}

	// 重新打开文件流，准备保存文件
	src2, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src2.Close()

	filename := fmt.Sprintf("%d_%s", time.Now().UnixNano(), file.Filename)
	savePath := filepath.Join(saveDir, filename)

	dst, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, src2); err != nil {
		return "", err
	}

	// 记录哈希和路径
	store.Set(hash, savePath)

	return savePath, nil
}
