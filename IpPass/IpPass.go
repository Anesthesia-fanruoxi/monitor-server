package IpPass

import (
	"log"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// 域名白名单
var allowedDomains = []string{
	"hzbxhd.tpddns.cn",
	"192.168.100.110",
	"192.168.100.39",
}

// 动态缓存域名解析的 IP 地址
var domainIPCache = struct {
	mapping map[string][]string
	mutex   sync.RWMutex
}{
	mapping: make(map[string][]string),
}

// 定期解析域名并更新缓存
func RefreshDomainIPCache() {
	for {
		log.Println("正在刷新域名解析缓存...")
		for _, domain := range allowedDomains {
			ips, err := net.LookupHost(domain)
			if err != nil {
				log.Printf("域名解析失败: %s, 错误: %v", domain, err)
				continue
			}

			// 更新缓存
			domainIPCache.mutex.Lock()
			domainIPCache.mapping[domain] = ips
			domainIPCache.mutex.Unlock()
		}
		time.Sleep(5 * time.Minute) // 每 5 分钟刷新一次
	}
}

// 检查请求 IP 是否在缓存的域名解析 IP 列表中
func isAllowedIP(ip string) bool {
	domainIPCache.mutex.RLock()
	defer domainIPCache.mutex.RUnlock()

	for _, ips := range domainIPCache.mapping {
		for _, cachedIP := range ips {
			if cachedIP == ip {
				return true
			}
		}
	}
	return false
}

func IpRestrictionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ip string

		// 从 X-Forwarded-For 头获取客户端 IP
		if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
			// X-Forwarded-For 可能包含多个 IP 地址，取第一个
			ip = forwarded
			if idx := strings.Index(ip, ","); idx != -1 {
				ip = ip[:idx]
			}
		} else if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
			// 尝试从 X-Real-IP 获取
			ip = realIP
		} else {
			// 最后从 RemoteAddr 获取
			var err error
			ip, _, err = net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				log.Printf("无法解析客户端 IP 地址，RemoteAddr: %s", r.RemoteAddr)
				http.Error(w, "无法解析客户端 IP 地址", http.StatusForbidden)
				return
			}
		}

		// 检查 IP 是否允许访问
		if !isAllowedIP(ip) {
			http.Error(w, "无权访问", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
