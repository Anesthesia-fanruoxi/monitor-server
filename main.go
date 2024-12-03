package main

import (
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"monitor-server/Handers"
	"monitor-server/IpPass"
	"monitor-server/Metrics"
	"net/http"
)

func main() {
	// 启动域名解析缓存的定时刷新功能
	err := Handers.LoadProjectDict("projects.json")
	if err != nil {
		log.Fatal(err)
	}
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
