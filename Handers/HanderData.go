package Handers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"monitor-server/Metrics"
	"monitor-server/Modles"
	"os"
	"strings"
	"time"
)

var projectNameDict map[string]string

// 处理 namespace，去掉 -v1 或 -v2
func cleanNamespace(namespace string) string {
	if strings.HasSuffix(namespace, "-v1") {
		return namespace[:len(namespace)-3] // 去掉 -v1
	} else if strings.HasSuffix(namespace, "-v2") {
		return namespace[:len(namespace)-3] // 去掉 -v2
	}
	return namespace
}

// 读取配置文件的函数
func LoadProjectDict(configPath string) error {
	// 打开配置文件
	file, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("无法打开配置文件: %v", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("无法读取配置文件内容: %v", err)
	}

	// 解析 JSON 数据到字典
	err = json.Unmarshal(content, &projectNameDict)
	if err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	return nil
}

// 处理 nginx 类型的数据
func HandleNginxData(data []interface{}, project string) {
	projectName := projectNameDict[project]
	if projectName == "" {
		// 如果字典中没有找到对应的中文名称，使用原值
		projectName = project
	}
	for _, item := range data {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			log.Printf("无法序列化 nginx data 中的元素: %v", item)
			continue
		}

		var nginxData Modles.NginxSource
		if err := json.Unmarshal(itemBytes, &nginxData); err != nil {
			log.Printf("解析 Nginx 数据失败: %v", err)
			continue
		}

		// 更新 Nginx 指标并打印日志
		Metrics.NginxIsRunMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.IsRun))
		Metrics.NginxReTotalMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.ReTotal))
		Metrics.NginxLoginUserCountMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.LoginUserCount))
		Metrics.NginxRawTotalMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.RawTotal))
		Metrics.NginxUdptotalMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.Udptotal))
		Metrics.NginxTcpTotalMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.TcpTotal))
		Metrics.NginxTotalTcpMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.TotalTcp))
		Metrics.NginxInetTotalMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.InetTotal))
		Metrics.NginxFragTotalMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.FragTotal))
		Metrics.NginxTcpEstabMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.TcpEstab))
		Metrics.NginxTcpClosedMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.TcpClosed))
		Metrics.NginxTcpOrphanedMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.TcpOrphaned))
		Metrics.NginxTcpTimewaitMetric.WithLabelValues(nginxData.HostName, projectName).Set(float64(nginxData.TcpTimewait))
		UpdateNginxMetricWithTimestamp(fmt.Sprintf("%s_%s", nginxData.HostName, projectName))
	}
}

// 处理硬件相关的数据
func HandleHardData(data []interface{}, project string) {
	projectName := projectNameDict[project]
	if projectName == "" {
		// 如果字典中没有找到对应的中文名称，使用原值
		projectName = project
	}
	for _, item := range data {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			log.Printf("无法序列化 hard data 中的元素: %v", item)
			continue
		}

		var hardData Modles.HardSource
		if err := json.Unmarshal(itemBytes, &hardData); err != nil {
			log.Printf("解析硬件数据失败: %v", err)
			continue
		}

		// 更新硬件相关指标并打印日志
		Metrics.CpuPercentMetric.WithLabelValues(hardData.HostName, projectName).Set(hardData.CPUPercent)
		Metrics.DiskTotalMetric.WithLabelValues(hardData.HostName, projectName).Set(hardData.DiskTotal)
		Metrics.DiskUsedMetric.WithLabelValues(hardData.HostName, projectName).Set(hardData.DiskUsed)
		Metrics.DiskFreeMetric.WithLabelValues(hardData.HostName, projectName).Set(hardData.DiskFree)
		Metrics.DiskUsedPercentMetric.WithLabelValues(hardData.HostName, projectName).Set(hardData.DiskUsedPercent)
		Metrics.MemoryTotalMetric.WithLabelValues(hardData.HostName, projectName).Set(hardData.MemoryTotal)
		Metrics.MemoryUsedMetric.WithLabelValues(hardData.HostName, projectName).Set(hardData.MemoryUsed)
		Metrics.MemoryFreeMetric.WithLabelValues(hardData.HostName, projectName).Set(hardData.MemoryFree)
		Metrics.MemoryUsedPercentMetric.WithLabelValues(hardData.HostName, projectName).Set(hardData.MemoryUsedPercent)
		Metrics.CpuLoad1Metric.WithLabelValues(hardData.HostName, projectName).Set(hardData.CPULoad1)
		Metrics.CpuLoad5Metric.WithLabelValues(hardData.HostName, projectName).Set(hardData.CPULoad5)
		Metrics.CpuLoad15Metric.WithLabelValues(hardData.HostName, projectName).Set(hardData.CPULoad15)

		UpdateHardMetricWithTimestamp(fmt.Sprintf("%s_%s", hardData.HostName, projectName))

	}
}

// 处理SSl证书数据
func HandleSSLData(data []interface{}, project string) {
	projectName := projectNameDict[project]
	if projectName == "" {
		// 如果字典中没有找到对应的中文名称，使用原值
		projectName = project
	}
	for _, item := range data {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			log.Printf("无法序列化 ssl data 中的元素: %v", item)
			continue
		}

		var sslData Modles.SslSource
		if err := json.Unmarshal(itemBytes, &sslData); err != nil {
			log.Printf("解析 SSL 数据失败: %v", err)
			continue
		}

		// 如果 comment 为空，则设置为 "未备注"
		if sslData.Comment == "" {
			sslData.Comment = "未备注"
		}

		resolve := "false" // 默认 false
		if sslData.Resolve {
			resolve = "true"
		}

		// 更新 SSL 指标并打印日志，添加 project 标签
		metricLabel := fmt.Sprintf("%s_%s_%s_%s_%s", sslData.Domain, sslData.Comment, sslData.Status, resolve, projectName)
		Metrics.SslDaysLeftMetric.WithLabelValues(sslData.Domain, sslData.Comment, sslData.Status, resolve, projectName).Set(float64(sslData.DaysLeft))

		// 存储时间戳
		sslTimestamp.Store(metricLabel, time.Now())

	}
}

// 处理容器资源数据
func HandleContainerResourceData(data []interface{}, project string) {
	// 将项目名称映射为中文名称（假设有一个字典）
	projectName := projectNameDict[project]
	if projectName == "" {
		// 如果字典中没有找到对应的中文名称，使用原值
		projectName = project
	}

	for _, item := range data {
		// 将每个资源项转化为 JSON 字节
		itemBytes, err := json.Marshal(item)
		if err != nil {
			log.Printf("无法序列化容器资源数据中的元素: %v", item)
			continue
		}

		// 将字节数据反序列化为 ContainerResource 类型
		var containerResource Modles.ContainerResource
		if err := json.Unmarshal(itemBytes, &containerResource); err != nil {
			log.Printf("解析容器资源数据失败: %v", err)
			continue
		}
		containerNamespace := cleanNamespace(containerResource.Namespace)

		Metrics.ContainerCpuUsageMetric.WithLabelValues(containerNamespace, containerResource.PodName, containerResource.Container, containerResource.ControllerName, projectName).Set(containerResource.UseCpu)
		Metrics.ContainerMemoryUsageMetric.WithLabelValues(containerNamespace, containerResource.PodName, containerResource.Container, containerResource.ControllerName, projectName).Set(float64(containerResource.UseMemory))
		Metrics.ContainerCpuLimitMetric.WithLabelValues(containerNamespace, containerResource.PodName, containerResource.Container, containerResource.ControllerName, projectName).Set(containerResource.LimitCpu)
		Metrics.ContainerMemoryLimitMetric.WithLabelValues(containerNamespace, containerResource.PodName, containerResource.Container, containerResource.ControllerName, projectName).Set(float64(containerResource.LimitMemory))
		Metrics.ContainerRestartCountMetric.WithLabelValues(containerNamespace, containerResource.PodName, containerResource.Container, containerResource.ControllerName, projectName).Set(float64(containerResource.RestartCount))

		UpdateContainerMetricWithTimestamp(fmt.Sprintf("%s_%s_%s_%s_%s", containerNamespace, containerResource.PodName, containerResource.Container, containerResource.ControllerName, projectName))
	}
}

// 更新心跳数据
func HandleHeartData(data []interface{}, project string) {
	projectName := projectNameDict[project]
	if projectName == "" {
		// 如果字典中没有找到对应的中文名称，使用原值
		projectName = project
	}
	for _, item := range data {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			log.Printf("无法序列化 hard data 中的元素: %v", item)
			continue
		}

		var hardData Modles.HardSource
		if err := json.Unmarshal(itemBytes, &hardData); err != nil {
			log.Printf("解析硬件数据失败: %v", err)
			continue
		}
		for _, item := range data {
			itemBytes, err := json.Marshal(item)
			if err != nil {
				log.Printf("无法序列化 ssl data 中的元素: %v", item)
				continue
			}

			var heartData Modles.HeartSource
			if err := json.Unmarshal(itemBytes, &heartData); err != nil {
				log.Printf("解析 SSL 数据失败: %v", err)
				continue
			}

			// 记录时间戳
			// 更新 心跳 指标
			Metrics.IsActiveMetric.WithLabelValues(heartData.Hostname, projectName).Set(float64(heartData.IsActive))
			metricLabel := fmt.Sprintf("%s_%s", hardData.HostName, projectName)

			agentHeartbeatTimes.Store(metricLabel, time.Now())
		}
	}
}

// 更新控制器数据
func HandleControllertResourceData(data []interface{}, project string) {
	projectName := projectNameDict[project]
	if projectName == "" {
		// 如果字典中没有找到对应的中文名称，使用原值
		projectName = project
	}
	for _, item := range data {
		itemBytes, err := json.Marshal(item)
		if err != nil {
			log.Printf("无法序列化 hard data 中的元素: %v", item)
			continue
		}

		var controllerData Modles.ControllerResource
		if err := json.Unmarshal(itemBytes, &controllerData); err != nil {
			log.Printf("解析硬件数据失败: %v", err)
			continue
		}
		for _, item := range data {
			itemBytes, err := json.Marshal(item)
			if err != nil {
				log.Printf("无法序列化 ssl data 中的元素: %v", item)
				continue
			}

			var controllerData Modles.ControllerResource
			if err := json.Unmarshal(itemBytes, &controllerData); err != nil {
				log.Printf("解析 SSL 数据失败: %v", err)
				continue
			}

			containerNamespace := cleanNamespace(controllerData.Namespace)
			// 更新 心跳 指标
			Metrics.ControllerReplicasMetric.WithLabelValues(containerNamespace, controllerData.Container, controllerData.ControllerType, projectName).Set(float64(controllerData.Replicas))
			Metrics.ControllerReplicasAvailableMetric.WithLabelValues(containerNamespace, controllerData.Container, controllerData.ControllerType, projectName).Set(float64(controllerData.ReplicasAvailable))
			Metrics.ControllerReplicasUnavailableMetric.WithLabelValues(containerNamespace, controllerData.Container, controllerData.ControllerType, projectName).Set(float64(controllerData.ReplicasUnavailable))

			UpdateControllerMetricWithTimestamp(fmt.Sprintf("%s_%s_%s_%s", containerNamespace, controllerData.Container, controllerData.ControllerType, projectName))

		}
	}
}
