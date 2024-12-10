package Handers

import (
	"log"
	"monitor-server/Metrics"
	"strings"
	"time"
)

// 用来存储时间戳和指标名称的 map
var ControllerTimestamp = make(map[string]time.Time)

//var Controllermu sync.Mutex // 添加一个锁来保护 ControllerTimestamp

// 反解析 label 字符串并更新数据
func parseControllerLabel(metricLabel string) (string, string, string, string) {
	// 以 "_" 为分隔符分割 label 字符串
	parts := strings.Split(metricLabel, "_")
	if len(parts) < 4 {
		log.Printf("标签 %s 无法解析，格式不正确", metricLabel)
		return "", "", "", ""
	}
	// 重新解析并返回各个字段的值
	return parts[0], parts[1], parts[2], parts[3]
}

// 定期检查超时的心跳数据
func CheckControllerHeartbeats() {

	currentTime := time.Now()

	// 遍历所有存储的时间戳，检查是否有超过 10 秒未更新的数据
	for metricLabel, timestamp := range ControllerTimestamp {
		// 如果超过 10 秒没有更新
		if currentTime.Sub(timestamp) > 10*time.Second {
			// 反解析 metricLabel 获取各个标签的值
			namespace, continer, ControllerType, project := parseControllerLabel(metricLabel)

			// 如果标签解析成功，且字段不为空，则删除相应的指标
			if namespace != "" && continer != "" && ControllerType != "" && project != "" {
				// 删除对应的指标
				Metrics.ControllerReplicasMetric.DeleteLabelValues(namespace, continer, ControllerType, project)
				Metrics.ControllerReplicasAvailableMetric.DeleteLabelValues(namespace, continer, ControllerType, project)
				Metrics.ControllerReplicasUnavailableMetric.DeleteLabelValues(namespace, continer, ControllerType, project)
				log.Printf("已删除标签为：%s", metricLabel)
			} else {
				log.Printf("标签 %s 格式不正确，跳过注销", metricLabel)
			}

			// 删除时间戳
			delete(ControllerTimestamp, metricLabel)
		}
	}
}

// 更新指标并记录时间戳
func UpdateControllerMetricWithTimestamp(metricLabel string) {
	// 获取当前时间
	currentTime := time.Now()
	// 存储时间戳
	ControllerTimestamp[metricLabel] = currentTime
}