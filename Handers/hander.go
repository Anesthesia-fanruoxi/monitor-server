package Handers

import (
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

// 固定的加密盐（密钥）
var encryptionKey = []byte("111111111111") // 必须是 16、24 或 32 字节长度
var mu sync.Mutex                          // 定义一个全局的 Mutex

// 解密数据
func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
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

func MetricsHandler(w http.ResponseWriter, r *http.Request, CustomRegistry *prometheus.Registry) {
	if r.Method != http.MethodPost {
		// 手动写入 JSON 错误响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		response := map[string]interface{}{
			"code": 405,
			"msg":  "仅支持 POST 请求",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("响应失败: %v", err)
		}
		return
	}

	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		// 手动写入 JSON 错误响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"code": 400,
			"msg":  "读取请求体失败",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("响应失败: %v", err)
		}
		return
	}
	defer r.Body.Close()

	// 解密数据
	decryptedData, err := Decrypt(body, encryptionKey)
	if err != nil {
		// 手动写入 JSON 错误响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"code": 400,
			"msg":  "解密失败: " + err.Error(),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("响应失败: %v", err)
		}
		return
	}

	// 将解密后的 JSON 解析为 map
	var payload map[string]interface{}
	if err := json.Unmarshal(decryptedData, &payload); err != nil {
		// 手动写入 JSON 错误响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"code": 400,
			"msg":  "JSON 解析失败: " + err.Error(),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("响应失败: %v", err)
		}
		return
	}

	// 提取 project 和 source 字段
	project, ok := payload["project"].(string)
	if !ok || project == "" {
		// 手动写入 JSON 错误响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"code": 400,
			"msg":  "缺少或无效的 project 字段",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("响应失败: %v", err)
		}
		return
	}
	source, ok := payload["source"].(string)
	if !ok || source == "" {
		// 手动写入 JSON 错误响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"code": 400,
			"msg":  "缺少或无效的 source 字段",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("响应失败: %v", err)
		}
		return
	}

	// 提取 data 字段
	data, ok := payload["data"].([]interface{})
	if !ok || len(data) == 0 {
		// 手动写入 JSON 错误响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := map[string]interface{}{
			"code": 400,
			"msg":  "缺少或无效的 data 字段",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("响应失败: %v", err)
		}
		return
	}

	// 立即返回成功响应
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{
		"code": 200,
		"msg":  "ok",
	}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("响应失败: %v", err)
	}

	// 在后台异步处理数据
	go func() {
		// 在调用 HandleContainerResourceData 之前加锁
		mu.Lock()         // 加锁
		defer mu.Unlock() // 确保函数结束时解锁

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
		default:
			log.Printf("不支持的 source 类型: %s", source)
		}
	}()
}
