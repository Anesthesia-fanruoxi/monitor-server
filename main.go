package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
	"log"
	"monitor-server/Handers"
	"monitor-server/IpPass"
	"monitor-server/Metrics"
	"net/http"
	"os"
	"sync"
	"time"
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

// 动态加载配置并监听变化
func loadConfigWithViper() {
	viper.SetConfigFile("config.yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("配置文件读取失败: %v", err)
	}

	// 监听配置文件变化
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Printf("配置文件已更新: %v", e.Name)
		Handers.SetEncryptionKey(viper.GetString("encrypted"))
		IpPass.SetAllowedDomains(viper.GetStringSlice("ipPass"))
	})
}

// 启动定时任务和心跳检查
func startHeartbeatChecks(wg *sync.WaitGroup) {
	// 启动多线程处理
	wg.Add(7)

	// 启动 goroutine
	go func() {
		defer wg.Done()
		for {
			Handers.CheckContainerHeartbeats()
			time.Sleep(5 * time.Second) // 每 5 秒检查一次
		}
	}()
	go func() {
		defer wg.Done()
		for {
			Handers.CheckHardHeartbeats()
			time.Sleep(5 * time.Second) // 每 5 秒检查一次
		}
	}()
	go func() {
		defer wg.Done()
		for {
			Handers.CheckHeartbeats()
			time.Sleep(5 * time.Second) // 每 5 秒检查一次
		}
	}()
	go func() {
		defer wg.Done()
		for {
			Handers.CheckSSLHeartbeats()
			time.Sleep(5 * time.Second) // 每 5 秒检查一次
		}
	}()
	go func() {
		defer wg.Done()
		for {
			Handers.CheckControllerHeartbeats()
			time.Sleep(5 * time.Second) // 每 5 秒检查一次
		}
	}()
	go func() {
		defer wg.Done()
		for {
			Handers.CheckNginxHeartbeats()
			time.Sleep(5 * time.Second) // 每 5 秒检查一次
		}
	}()
	go func() {
		defer wg.Done()
		for {
			IpPass.RefreshDomainIPCache()
			time.Sleep(5 * time.Second) // 每 5 秒检查一次
		}
	}()
}

func main() {
	// 加载配置文件
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

	// 启动动态配置加载
	go loadConfigWithViper()

	// 启动定时任务和心跳检查
	var wg sync.WaitGroup
	go startHeartbeatChecks(&wg)

	// 暴露自定义指标
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

	// 启动 HTTP 服务
	go func() {
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("HTTP 服务启动失败: %v", err)
		}
	}()

	// 等待定时任务完成
	wg.Wait()
	log.Println("所有定时任务已完成，服务正在运行...")
}
