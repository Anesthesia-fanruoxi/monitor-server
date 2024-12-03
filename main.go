package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/yaml.v3"
	"log"
	"monitor-server/Handers"
	"monitor-server/IpPass"
	"monitor-server/Metrics"
	"net/http"
	"os"
)

// 配置结构体
type Config struct {
	Encrypted string   `yaml:"encrypted"` // 加密盐
	IpPass    []string `yaml:"ipPass"`    // IP 白名单
}

// 读取配置文件的函数
func loadConfig(configFile string) (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, fmt.Errorf("打开配置文件失败: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %v", err)
	}
	return &config, nil
}

func main() {
	// 读取配置文件
	config, err := loadConfig("config.yaml")
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 设置加密盐
	Handers.SetEncryptionKey(config.Encrypted)

	// 启动域名解析缓存的定时刷新功能
	err = Handers.LoadProjectDict("projects.json")
	if err != nil {
		log.Fatal(err)
	}

	// 设置 IP 白名单
	IpPass.SetAllowedDomains(config.IpPass)

	// 启动定时器和心跳检查
	go IpPass.RefreshDomainIPCache()
	go Handers.CheckNginxHeartbeats()
	go Handers.CheckHeartbeats()
	go Handers.CheckSSLHeartbeats()
	go Handers.CheckContainerHeartbeats()
	go Handers.CheckHardHeartbeats()

	// 曝露自定义指标
	metricsHandler := promhttp.HandlerFor(
		Metrics.CustomRegistry, // 使用自定义的 Registry
		promhttp.HandlerOpts{},
	)

	// 注册 `/metrics` 路径
	http.Handle("/metrics", IpPass.IpRestrictionMiddleware(metricsHandler))

	// 设置 HTTP 接收端，传递 CustomRegistry 给 MetricsHandler
	http.HandleFunc("/metrics_data", func(w http.ResponseWriter, r *http.Request) {
		Handers.MetricsHandler(w, r, Metrics.CustomRegistry) // 传递 CustomRegistry
	})

	// 启动 HTTP 服务
	log.Println("服务启动，监听端口 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
