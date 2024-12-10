package Handers

import (
	"log"
	"monitor-server/Metrics"
	"strings"
	"time"
)

// 用来存储时间戳和指标名称的 map
var ContainerTimestamp = make(map[string]time.Time)

//var containermu sync.Mutex // 添加一个锁来保护 ContainerTimestamp

// 反解析 label 字符串并更新数据
func parseContainerLabel(metricLabel string) (string, string, string, string, string) {
	// 以 "_" 为分隔符分割 label 字符串
	parts := strings.Split(metricLabel, "_")
	if len(parts) < 5 {
		log.Printf("标签 %s 无法解析，格式不正确", metricLabel)
		return "", "", "", "", ""
	}
	// 重新解析并返回各个字段的值
	return parts[0], parts[1], parts[2], parts[3], parts[4]
}

// 定期检查超时的心跳数据
func CheckContainerHeartbeats() {

	currentTime := time.Now()

	// 遍历所有存储的时间戳，检查是否有超过 10 秒未更新的数据
	for metricLabel, timestamp := range ContainerTimestamp {
		// 如果超过 10 秒没有更新
		if currentTime.Sub(timestamp) > 10*time.Second {
			// 反解析 metricLabel 获取各个标签的值
			namespace, podName, container, controllerName, project := parseContainerLabel(metricLabel)

			// 如果标签解析成功，且字段不为空，则删除相应的指标
			if namespace != "" && podName != "" && container != "" && controllerName != "" && project != "" {
				// 删除对应的指标
				Metrics.ContainerCpuUsageMetric.DeleteLabelValues(namespace, podName, container, controllerName, project)
				Metrics.ContainerMemoryUsageMetric.DeleteLabelValues(namespace, podName, container, controllerName, project)
				Metrics.ContainerCpuLimitMetric.DeleteLabelValues(namespace, podName, container, controllerName, project)
				Metrics.ContainerMemoryLimitMetric.DeleteLabelValues(namespace, podName, container, controllerName, project)
				Metrics.ContainerRestartCountMetric.DeleteLabelValues(namespace, podName, container, controllerName, project)
				log.Printf("已删除标签为：%s", metricLabel)
			} else {
				log.Printf("标签 %s 格式不正确，跳过注销", metricLabel)
			}

			// 删除时间戳
			delete(ContainerTimestamp, metricLabel)
		}
	}
}

// 更新指标并记录时间戳
func UpdateContainerMetricWithTimestamp(metricLabel string) {
	// 获取当前时间
	currentTime := time.Now()
	// 存储时间戳
	ContainerTimestamp[metricLabel] = currentTime
}
