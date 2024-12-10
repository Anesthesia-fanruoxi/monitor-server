package Handers

import (
	"log"
	"monitor-server/Metrics"
	"strings"
	"time"
)

// 用来存储时间戳和指标名称的 map
var HardTimestamp = make(map[string]time.Time)

// 反解析 label 字符串并更新数据
func parseHardLabel(metricLabel string) (string, string) {
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
func CheckHardHeartbeats() {

	currentTime := time.Now()

	// 遍历所有存储的时间戳，检查是否有超过 10 秒未更新的数据
	for metricLabel, timestamp := range HardTimestamp {
		// 如果超过 10 秒没有更新
		if currentTime.Sub(timestamp) > 10*time.Second {
			// 反解析 metricLabel 获取各个标签的值
			hostName, project := parseHardLabel(metricLabel)

			// 如果标签解析成功，且字段不为空，则删除相应的指标
			if hostName != "" && project != "" {
				// 删除对应的指标
				Metrics.CpuPercentMetric.DeleteLabelValues(hostName, project)
				Metrics.DiskTotalMetric.DeleteLabelValues(hostName, project)
				Metrics.DiskUsedMetric.DeleteLabelValues(hostName, project)
				Metrics.DiskFreeMetric.DeleteLabelValues(hostName, project)
				Metrics.DiskUsedPercentMetric.DeleteLabelValues(hostName, project)
				Metrics.MemoryTotalMetric.DeleteLabelValues(hostName, project)
				Metrics.MemoryUsedMetric.DeleteLabelValues(hostName, project)
				Metrics.MemoryFreeMetric.DeleteLabelValues(hostName, project)
				Metrics.MemoryUsedPercentMetric.DeleteLabelValues(hostName, project)
				Metrics.CpuLoad1Metric.DeleteLabelValues(hostName, project)
				Metrics.CpuLoad5Metric.DeleteLabelValues(hostName, project)
				Metrics.CpuLoad15Metric.DeleteLabelValues(hostName, project)
			} else {
				log.Printf("标签 %s 格式不正确，跳过注销", metricLabel)
			}

			// 删除时间戳
			delete(HardTimestamp, metricLabel)
		}
	}

}

// 更新指标并记录时间戳
func UpdateHardMetricWithTimestamp(metricLabel string) {
	// 获取当前时间
	currentTime := time.Now()
	// 存储时间戳
	HardTimestamp[metricLabel] = currentTime

}
