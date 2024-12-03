package Handers

import (
	"log"
	"monitor-server/Metrics"
	"strings"
	"time"
)

// 用来存储时间戳和指标名称的 map
var NginxTimestamp = make(map[string]time.Time)

// 反解析 label 字符串并更新数据
func parseNginxLabel(metricLabel string) (string, string) {
	// 以 "_" 为分隔符分割 label 字符串
	parts := strings.Split(metricLabel, "_")
	if len(parts) < 2 {
		log.Printf("标签 %s 无法解析，格式不正确", metricLabel)
		return "", ""
	}
	// 重新解析并返回各个字段的值
	return parts[0], parts[1]
}

// 定期检查超时的心跳数据
func CheckNginxHeartbeats() {
	ticker := time.NewTicker(5 * time.Second) // 每 5 秒检查一次
	go func() {
		for {
			<-ticker.C
			currentTime := time.Now()

			// 遍历所有存储的时间戳，检查是否有超过 10 秒未更新的数据
			for metricLabel, timestamp := range NginxTimestamp {
				// 如果超过 10 秒没有更新
				if currentTime.Sub(timestamp) > 10*time.Second {
					// 反解析 metricLabel 获取各个标签的值
					hostName, project := parseNginxLabel(metricLabel)

					// 如果标签解析成功，且字段不为空，则删除相应的指标
					if hostName != "" && project != "" {
						// 删除对应的指标
						Metrics.NginxIsRunMetric.DeleteLabelValues(hostName, project)
						Metrics.NginxReTotalMetric.DeleteLabelValues(hostName, project)
						Metrics.NginxLoginUserCountMetric.DeleteLabelValues(hostName, project)
						Metrics.NginxRawTotalMetric.DeleteLabelValues(hostName, project)
						Metrics.NginxUdptotalMetric.DeleteLabelValues(hostName, project)
						Metrics.NginxTcpTotalMetric.DeleteLabelValues(hostName, project)
						Metrics.NginxTotalTcpMetric.DeleteLabelValues(hostName, project)
						Metrics.NginxInetTotalMetric.DeleteLabelValues(hostName, project)
						Metrics.NginxFragTotalMetric.DeleteLabelValues(hostName, project)
						Metrics.NginxTcpEstabMetric.DeleteLabelValues(hostName, project)
						Metrics.NginxTcpClosedMetric.DeleteLabelValues(hostName, project)
						Metrics.NginxTcpOrphanedMetric.DeleteLabelValues(hostName, project)
						Metrics.NginxTcpTimewaitMetric.DeleteLabelValues(hostName, project)
					} else {
						log.Printf("标签 %s 格式不正确，跳过注销", metricLabel)
					}

					// 删除时间戳
					delete(NginxTimestamp, metricLabel)
				}
			}
		}
	}()
}

// 更新指标并记录时间戳
func UpdateNginxMetricWithTimestamp(metricLabel string) {
	// 获取当前时间
	currentTime := time.Now()

	// 存储时间戳
	NginxTimestamp[metricLabel] = currentTime

}
