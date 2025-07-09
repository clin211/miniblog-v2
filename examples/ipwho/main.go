// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

// Package main demonstrates how to use the ipwho package together with contextx
// to store and retrieve IP information in context for request processing.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/clin211/miniblog-v2/internal/pkg/contextx"
	"github.com/clin211/miniblog-v2/pkg/ipwho"
)

// handleBusinessLogic 模拟具体的业务逻辑处理
func handleBusinessLogic(ctx context.Context, ipDetail *ipwho.IPDetail) error {
	// 从上下文中获取所有相关信息
	clientIP := contextx.ClientIP(ctx)
	location := contextx.ClientLocation(ctx)

	fmt.Printf("  业务逻辑处理:\n")
	fmt.Printf("    IP地址: %s (%s)\n", clientIP, ipDetail.Type)
	fmt.Printf("    地理位置: %s\n", location)
	fmt.Printf("    ISP: %s\n", ipDetail.Connection.Isp)
	fmt.Printf("    时区: %s (%s)\n", ipDetail.Timezone.ID, ipDetail.Timezone.Abbr)

	// 获取该地区的当前时间
	currentTime, err := time.Parse(time.RFC3339, ipDetail.Timezone.CurrentTime)
	if err == nil {
		fmt.Printf("    当地时间: %s\n", currentTime.Format("2006-01-02 15:04:05"))
	}

	// 模拟基于地理位置的业务逻辑
	if ipDetail.CountryCode == "CN" {
		fmt.Printf("    检测到中国用户，应用中国地区特定策略\n")
	} else if ipDetail.IsEu {
		fmt.Printf("    检测到欧盟用户，应用GDPR合规策略\n")
	}

	// 记录访问日志（在实际应用中，这些信息可以用于安全分析、用户行为分析等）
	logAccess(ctx, ipDetail)

	return nil
}

// logAccess 记录访问日志
func logAccess(ctx context.Context, ipDetail *ipwho.IPDetail) {
	requestID := contextx.RequestID(ctx)
	clientIP := contextx.ClientIP(ctx)
	location := contextx.ClientLocation(ctx)

	fmt.Printf("  访问日志记录:\n")
	fmt.Printf("    [%s] IP: %s, 位置: %s, ASN: %d\n",
		requestID, clientIP, location, ipDetail.Connection.Asn)
}

func main() {
	fmt.Println("=== IPWho + ContextX 集成示例 ===")

	client := ipwho.NewClient()

	ipDetail, err := client.GetHostIPDetail(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	// 输出 JSON
	details, _ := json.Marshal(ipDetail)
	fmt.Println("ip detail:", string(details))

	ctx := context.Background()
	ctx = contextx.WithClientIP(ctx, ipDetail.IP)
	ctx = contextx.WithClientLocation(ctx, ipDetail.IP)

	ipDetail, err = client.GetIPDetail(ctx, ipDetail.IP)
	if err != nil {
		log.Fatal(err)
	}

	handleBusinessLogic(ctx, ipDetail)
	fmt.Println("\n=== 示例完成 ===")
}
