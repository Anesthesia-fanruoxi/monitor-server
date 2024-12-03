package Handers

import (
	"log"
	"monitor-server/Metrics"
	"strings"
	"time"
)

var sslTimestamp = make(map[string]time.Time)

// 反解析 label 字符串并更新 SSL 数据
func parseSSLLabel(metricLabel string) (string, string, string, string, string) {
	// 以 "_" 为分隔符分割 label 字符串
	parts := strings.Split(metricLabel, "_")
	if len(parts) < 5 {
		log.Printf("标签 %s 无法解析，格式不正确", metricLabel)
		return "", "", "", "", ""
	}

	// 重新解析并返回 SSL 数据
	return parts[0], parts[1], parts[2], parts[3], parts[4]
}

// 定期检查超时的 SSL 数据
func CheckSSLHeartbeats() {
	ticker := time.NewTicker(5 * time.Second) // 每 5 秒检查一次
	go func() {
		for {
			<-ticker.C
			currentTime := time.Now()
			for metricLabel, timestamp := range sslTimestamp {
				// 如果超过 10 秒没有更新
				if currentTime.Sub(timestamp) > 10*time.Second {
					// 反解析 metricLabel 获取各个标签的值
					domain, comment, status, resolve, projectName := parseSSLLabel(metricLabel)

					if domain != "" && comment != "" && status != "" && resolve != "" {

						// 删除对应的 SSL 指标
						Metrics.SslDaysLeftMetric.DeleteLabelValues(domain, comment, status, resolve, projectName)

						// 删除时间戳
						delete(sslTimestamp, metricLabel) // 删除时间戳
					} else {
						log.Printf("标签 %s 格式不正确，跳过注销", metricLabel)
					}
				}
			}
		}
	}()
}