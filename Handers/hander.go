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
var encryptionKeyMu sync.RWMutex // 使用读写锁，允许并发读取

// 细粒度锁：每种 source 类型使用独立的锁，允许不同类型数据并发处理
var (
	nginxMu            sync.Mutex
	hardMu             sync.Mutex
	sslMu              sync.Mutex
	containerMu        sync.Mutex
	heartMu            sync.Mutex
	controllerMu       sync.Mutex
	trafficSwitchingMu sync.Mutex
)

// 允许的 source 类型白名单
var allowedSources = map[string]bool{
	"nginx":            true,
	"hard":             true,
	"ssl":              true,
	"k8s":              true,
	"heart":            true,
	"k8sController":    true,
	"trafficSwitching": true,
}

// 验证 project 名称是否合法
func isValidProject(project string) bool {
	// 限制长度（按 rune 计算，支持中文）
	runeLen := len([]rune(project))
	if runeLen == 0 || runeLen > 64 {
		return false
	}
	// 允许字母、数字、下划线、中划线、点、中文字符
	for _, c := range project {
		isAlphaNum := (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9')
		isSpecialChar := c == '_' || c == '-' || c == '.'
		isChinese := c >= 0x4E00 && c <= 0x9FFF // 中文基本区
		if !isAlphaNum && !isSpecialChar && !isChinese {
			return false
		}
	}
	return true
}

// 验证 source 类型是否允许
func isValidSource(source string) bool {
	return allowedSources[source]
}

// 请求体大小限制（10MB）
const MaxRequestBodySize = 10 * 1024 * 1024

// SetEncryptionKey 设置加密盐（验证密钥长度）
func SetEncryptionKey(key string) {
	keyLen := len(key)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		log.Printf("警告: AES 密钥长度应为 16/24/32 字节，当前: %d 字节", keyLen)
	}
	encryptionKeyMu.Lock()
	defer encryptionKeyMu.Unlock()
	encryptionKey = []byte(key)
}

// GetEncryptionKey 获取当前加密盐（线程安全）
func GetEncryptionKey() []byte {
	encryptionKeyMu.RLock()
	defer encryptionKeyMu.RUnlock()
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
	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("关闭 gzip reader 失败: %v", err)
		}
	}()

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

	// 读取请求体（限制大小防止 DoS 攻击）
	body, err := io.ReadAll(io.LimitReader(r.Body, MaxRequestBodySize))
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "读取请求体失败")
		return
	}
	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			log.Printf("关闭请求体失败: %v", err)
		}
	}(r.Body)

	// 检查请求体是否超出限制
	if len(body) >= MaxRequestBodySize {
		writeJSONError(w, http.StatusRequestEntityTooLarge, "请求体过大")
		return
	}
	//log.Printf("收到请求，数据长度: %d 字节", len(body))

	// 解密数据（错误信息不暴露内部细节）
	decryptedData, err := Decrypt(body)
	if err != nil {
		log.Printf("解密失败: %v", err) // 日志记录详细错误
		writeJSONError(w, http.StatusBadRequest, "数据解密失败")
		return
	}

	// 解压数据
	decompressedData, err := Decompress(decryptedData)
	if err != nil {
		log.Printf("解压失败: %v", err)
		writeJSONError(w, http.StatusBadRequest, "数据解压失败")
		return
	}

	// 将解压后的 JSON 解析为 map
	var payload map[string]interface{}
	if err := json.Unmarshal(decompressedData, &payload); err != nil {
		log.Printf("JSON 解析失败: %v", err)
		writeJSONError(w, http.StatusBadRequest, "数据格式错误")
		return
	}

	// 提取并验证 project 字段
	project, ok := payload["project"].(string)
	if !ok || project == "" {
		writeJSONError(w, http.StatusBadRequest, "缺少或无效的 project 字段")
		return
	}
	if !isValidProject(project) {
		log.Printf("无效的 project 名称: %s", project)
		writeJSONError(w, http.StatusBadRequest, "无效的 project 名称")
		return
	}

	// 提取并验证 source 字段
	source, ok := payload["source"].(string)
	if !ok || source == "" {
		writeJSONError(w, http.StatusBadRequest, "缺少或无效的 source 字段")
		return
	}
	if !isValidSource(source) {
		log.Printf("不支持的 source 类型: %s", source)
		writeJSONError(w, http.StatusBadRequest, "不支持的数据类型")
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

	// 在后台异步处理数据，使用细粒度锁提高并发性能
	go func() {
		// 根据 source 类型使用对应的锁，允许不同类型数据并发处理
		switch source {
		case "nginx":
			nginxMu.Lock()
			HandleNginxData(data, project)
			nginxMu.Unlock()
		case "hard":
			hardMu.Lock()
			HandleHardData(data, project)
			hardMu.Unlock()
		case "ssl":
			sslMu.Lock()
			HandleSSLData(data, project)
			sslMu.Unlock()
		case "k8s":
			containerMu.Lock()
			HandleContainerResourceData(data, project)
			containerMu.Unlock()
		case "heart":
			heartMu.Lock()
			HandleHeartData(data, project)
			heartMu.Unlock()
		case "k8sController":
			controllerMu.Lock()
			HandleControllertResourceData(data, project)
			controllerMu.Unlock()
		case "trafficSwitching":
			trafficSwitchingMu.Lock()
			HandleTrafficSwitchingData(data, project)
			trafficSwitchingMu.Unlock()
		}
		// 注意：source 已经在前面验证过，不需要 default 分支
	}()
}
