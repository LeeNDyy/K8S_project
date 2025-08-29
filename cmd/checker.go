package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

var ruServices = map[string]string{
	"Telegram":     "https://telegram.org",
	"WhatsApp":     "https://web.whatsapp.com",
	"Signal":       "https://signal.org",
	"Twitch":       "https://twitch.tv",
	"Госуслуги":    "https://gosuslugi.ru",
	"Сбербанк":     "https://sberbank.ru",
	"Тинькофф":     "https://tinkoff.ru",
	"ВТБ":          "https://vtb.ru",
	"Мосру":        "https://mos.ru",
	"Ростелеком":   "https://rt.ru",
	"Почта России": "https://pochta.ru",
	"Twitter/X":    "https://twitter.com",
	"Facebook":     "https://facebook.com",
	"Instagram":    "https://instagram.com",
	"Яндекс":       "https://ya.ru",
	"ВКонтакте":    "https://vk.com",
	"РЖД":          "https://rzd.ru",
	"Аэрофлот":     "https://aeroflot.ru",
}

func CheckService(serviceName, url string) ServiceStatus {
	domain := extractDomain(url)

	dnsValid, dnsError := checkDNS(domain)
	httpUp, httpCode, responseTime := checkHTTP(url)
	isBlocked := detectBlocking(httpCode, dnsValid, httpUp)

	return ServiceStatus{
		Service:      serviceName,
		URL:          url,
		HTTPUp:       httpUp,
		HTTPCode:     httpCode,
		ResponseTime: responseTime,
		DNSValid:     dnsValid,
		DNSError:     dnsError,
		IsBlocked:    isBlocked,
		LastChecked:  time.Now(),
	}
}

func extractDomain(url string) string {
	if strings.HasPrefix(url, "http") {
		if parts := strings.Split(url, "//"); len(parts) > 1 {
			if domainParts := strings.Split(parts[1], "/"); len(domainParts) > 0 {
				return domainParts[0]
			}
		}
	}
	return url
}

func checkDNS(domain string) (bool, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := net.DefaultResolver.LookupHost(ctx, domain)
	if err != nil {
		return false, fmt.Sprintf("DNS error: %v", err)
	}
	return true, ""
}

func checkHTTP(url string) (bool, int, int64) {
	client := &http.Client{
		Timeout: 8 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	start := time.Now() // Здесь start используется правильно

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := client.Do(req)
	responseTime := time.Since(start).Milliseconds() // Здесь start используется

	if err != nil {
		return false, 0, responseTime
	}
	defer resp.Body.Close()

	return resp.StatusCode < 500, resp.StatusCode, responseTime
}

func detectBlocking(httpCode int, dnsValid bool, httpUp bool) bool {
	if !dnsValid {
		return true
	}

	if httpCode == 403 || httpCode == 451 {
		return true
	}

	if httpCode == 502 || httpCode == 503 {
		return true
	}

	return false
}

func CheckAllServices() map[string]ServiceStatus {
	results := make(map[string]ServiceStatus)
	for name, url := range ruServices {
		results[name] = CheckService(name, url)
	}
	return results
}

func GetStats(statusMap map[string]ServiceStatus) Stats {
	downCount := 0
	blockedCount := 0

	for _, status := range statusMap {
		if !status.HTTPUp {
			downCount++
		}
		if status.IsBlocked {
			blockedCount++
		}
	}

	return Stats{
		TotalServices:   len(statusMap),
		DownServices:    downCount,
		BlockedServices: blockedCount,
		LastCheck:       time.Now().Format("02.01.2006 15:04:05"),
	}
}
