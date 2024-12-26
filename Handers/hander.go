package Handers

import (
	"bytes"
	"compress/gzip"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"io"
	"log"
	"net/http"
	"sync"
)

// 全局变量存储加密盐
var encryptionKey []byte
var mu sync.Mutex // 防止并发访问

// SetEncryptionKey 设置加密盐
func SetEncryptionKey(key string) {
	mu.Lock()         // 加锁确保线程安全
	defer mu.Unlock() // 解锁
	encryptionKey = []byte(key)
}

// GetEncryptionKey 获取当前加密盐
func GetEncryptionKey() []byte {
	return encryptionKey
}

// 解密数据
func Decrypt(ciphertext []byte) ([]byte, error) {
	key := GetEncryptionKey() // 获取加密盐
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("密文过短")
	}

	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

// 解压数据
func Decompress(data []byte) ([]byte, error) {
	reader, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer func(reader *gzip.Reader) {
		err := reader.Close()
		if err != nil {

		}
	}(reader)

	// 读取解压缩后的内容
	return io.ReadAll(reader)
}

// 响应 JSON 错误
func writeJSONError(w http.ResponseWriter, statusCode int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]interface{}{
		"code": statusCode,
		"msg":  msg,
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("响应失败: %v", err)
	}
}

func MetricsHandler(w http.ResponseWriter, r *http.Request, CustomRegistry *prometheus.Registry) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "仅支持 POST 请求")
		return
	}

	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "读取请求体失败")
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("关闭请求体失败: %v", err)
		}
	}(r.Body)
	// 解密数据
	decryptedData, err := Decrypt(body)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "解密失败: "+err.Error())
		return
	}
	//log.Printf("解密完成: %v", string(decryptedData))
	// 解压数据
	decompressedData, err := Decompress(decryptedData)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "解压失败: "+err.Error())
		return
	}

	// 将解压后的 JSON 解析为 map
	var payload map[string]interface{}
	if err := json.Unmarshal(decompressedData, &payload); err != nil {
		writeJSONError(w, http.StatusBadRequest, "JSON 解析失败: "+err.Error())
		return
	}

	// 提取 project 和 source 字段
	project, ok := payload["project"].(string)
	if !ok || project == "" {
		writeJSONError(w, http.StatusBadRequest, "缺少或无效的 project 字段")
		return
	}
	source, ok := payload["source"].(string)
	if !ok || source == "" {
		writeJSONError(w, http.StatusBadRequest, "缺少或无效的 source 字段")
		return
	}

	// 提取 data 字段
	data, ok := payload["data"].([]interface{})
	if !ok || len(data) == 0 {
		writeJSONError(w, http.StatusBadRequest, "缺少或无效的 data 字段")
		return
	}

	// 立即返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"code": 200, "msg": "ok"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("响应失败: %v", err)
	}

	// 在后台异步处理数据
	go func() {
		mu.Lock()         // 加锁，确保所有后续逻辑串行执行
		defer mu.Unlock() // 解锁

		// 根据 source 调用对应的处理函数
		switch source {
		case "nginx":
			HandleNginxData(data, project)
		case "hard":
			HandleHardData(data, project)
		case "ssl":
			HandleSSLData(data, project)
		case "k8s":
			HandleContainerResourceData(data, project)
		case "heart":
			HandleHeartData(data, project)
		case "k8sController":
			HandleControllertResourceData(data, project)
		default:
			log.Printf("不支持的 source 类型: %s", source)
		}
	}()
}
