// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

// Package ipwho 提供了IP地理位置查询功能，通过 ipwho.is API 获取IP地址的详细地理信息
package ipwho

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/clin211/miniblog-v2/internal/pkg/log"
)

// Options 定义了 IP 地理位置查询客户端的配置选项
type Options struct {
	// BaseURL API 服务的基础URL
	BaseURL string
	// UserAgent HTTP请求的用户代理字符串
	UserAgent string
	// Timeout 请求超时时间
	Timeout time.Duration
	// CustomHeaders 自定义HTTP头部
	CustomHeaders map[string]string
}

// Client IP地理位置查询客户端
type Client struct {
	options Options // 客户端配置选项
}

// IPDetail 表示IP地址的详细地理信息
type IPDetail struct {
	IP            string     `json:"ip"`
	Success       bool       `json:"success"`
	Type          string     `json:"type"`
	Continent     string     `json:"continent"`
	ContinentCode string     `json:"continent_code"`
	Country       string     `json:"country"`
	CountryCode   string     `json:"country_code"`
	Region        string     `json:"region"`
	RegionCode    string     `json:"region_code"`
	City          string     `json:"city"`
	Latitude      float64    `json:"latitude"`
	Longitude     float64    `json:"longitude"`
	IsEu          bool       `json:"is_eu"`
	Postal        string     `json:"postal"`
	CallingCode   string     `json:"calling_code"`
	Capital       string     `json:"capital"`
	Borders       string     `json:"borders"`
	Flag          Flag       `json:"flag"`
	Connection    Connection `json:"connection"`
	Timezone      Timezone   `json:"timezone"`
}

// Flag 表示国家/地区的旗帜信息
type Flag struct {
	Img          string `json:"img"`
	Emoji        string `json:"emoji"`
	EmojiUnicode string `json:"emoji_unicode"`
}

// Connection 表示网络连接信息
type Connection struct {
	Asn    int    `json:"asn"`
	Org    string `json:"org"`
	Isp    string `json:"isp"`
	Domain string `json:"domain"`
}

// Timezone 表示时区信息
type Timezone struct {
	ID          string `json:"id"`
	Abbr        string `json:"abbr"`
	IsDst       bool   `json:"is_dst"`
	Offset      int    `json:"offset"`
	Utc         string `json:"utc"`
	CurrentTime string `json:"current_time"`
}

// WithBaseURL 设置API服务的基础URL
// 参数:
//   - baseURL: API服务的基础URL
func WithBaseURL(baseURL string) func(*Options) {
	return func(options *Options) {
		if baseURL != "" {
			getOptionsOrSetDefault(options).BaseURL = baseURL
		}
	}
}

// WithUserAgent 设置HTTP请求的用户代理字符串
// 参数:
//   - userAgent: 用户代理字符串
func WithUserAgent(userAgent string) func(*Options) {
	return func(options *Options) {
		if userAgent != "" {
			opts := getOptionsOrSetDefault(options)
			opts.UserAgent = userAgent
			opts.CustomHeaders["User-Agent"] = userAgent
		}
	}
}

// WithTimeout 设置请求超时时间
// 参数:
//   - timeout: 超时时间
func WithTimeout(timeout time.Duration) func(*Options) {
	return func(options *Options) {
		if timeout > 0 {
			getOptionsOrSetDefault(options).Timeout = timeout
		}
	}
}

// WithCustomHeaders 设置自定义HTTP头部
// 参数:
//   - headers: 自定义HTTP头部映射
func WithCustomHeaders(headers map[string]string) func(*Options) {
	return func(options *Options) {
		if headers != nil {
			opts := getOptionsOrSetDefault(options)
			if opts.CustomHeaders == nil {
				opts.CustomHeaders = make(map[string]string)
			}
			for k, v := range headers {
				opts.CustomHeaders[k] = v
			}
		}
	}
}

// getOptionsOrSetDefault 获取Options，如果为空或未初始化则创建默认配置
// 参数:
//   - options: 可选的Options实例
//
// 返回值:
//   - 配置好的Options实例
func getOptionsOrSetDefault(options *Options) *Options {
	if options == nil {
		options = &Options{}
	}

	// 设置默认值（如果尚未设置）
	if options.BaseURL == "" {
		options.BaseURL = "https://ipwho.is/"
	}
	if options.UserAgent == "" {
		options.UserAgent = "curl/7.77.0"
	}
	if options.Timeout == 0 {
		options.Timeout = 30 * time.Second
	}
	if options.CustomHeaders == nil {
		options.CustomHeaders = make(map[string]string)
	}
	if _, exists := options.CustomHeaders["User-Agent"]; !exists {
		options.CustomHeaders["User-Agent"] = options.UserAgent
	}

	return options
}

// NewClient 创建一个新的IP地理位置查询客户端
// 参数:
//   - options: 配置函数列表，用于自定义客户端行为
//
// 返回值:
//   - 一个配置好的Client实例
//
// 使用示例:
//
//	client := ipwho.NewClient(
//		ipwho.WithTimeout(10*time.Second),
//		ipwho.WithUserAgent("MyApp/1.0"),
//	)
//	detail, err := client.GetIPDetail(context.Background(), "8.8.8.8")
func NewClient(options ...func(*Options)) *Client {
	opts := getOptionsOrSetDefault(nil)
	for _, f := range options {
		f(opts)
	}
	return &Client{
		options: *opts,
	}
}

// GetHostIPDetail 获取当前主机的IP地理位置详细信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//
// 返回值:
//   - IP地理位置详细信息
//   - 错误信息（如果有）
//
// 使用示例:
//
//	detail, err := client.GetHostIPDetail(context.Background())
//	if err != nil {
//		log.Printf("获取IP信息失败: %v", err)
//	}
func (c *Client) GetHostIPDetail(ctx context.Context) (*IPDetail, error) {
	return c.GetIPDetail(ctx, "")
}

// GetIPDetail 获取指定IP地址的地理位置详细信息
// 参数:
//   - ctx: 上下文对象，用于控制请求的生命周期
//   - ip: 要查询的IP地址，为空时查询当前主机IP
//
// 返回值:
//   - IP地理位置详细信息
//   - 错误信息（如果有）
//
// 使用示例:
//
//	detail, err := client.GetIPDetail(context.Background(), "8.8.8.8")
//	if err != nil {
//		log.Printf("获取IP信息失败: %v", err)
//	}
func (c *Client) GetIPDetail(ctx context.Context, ip string) (*IPDetail, error) {
	logger := log.W(ctx)

	// 构建请求URL
	url := c.options.BaseURL
	if ip != "" {
		url += ip
	}

	// 创建HTTP请求
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		logger.Errorw("failed to create http request", "err", err, "url", url)
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	// 设置请求头部
	for k, v := range c.options.CustomHeaders {
		req.Header.Set(k, v)
	}

	// 发送HTTP请求
	client := &http.Client{
		Timeout: c.options.Timeout,
	}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorw("ipwho request failed", "err", err, "url", url, "ip", ip)
		return nil, fmt.Errorf("ipwho request failed: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorw("failed to read response body", "err", err, "statusCode", resp.StatusCode)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		logger.Errorw("ipwho api returned non-200 status", "statusCode", resp.StatusCode, "body", string(respBody))
		return nil, fmt.Errorf("ipwho api returned status %d: %s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var detail IPDetail
	if err := json.Unmarshal(respBody, &detail); err != nil {
		logger.Errorw("ipwho response parse failed", "err", err, "body", string(respBody))
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// 检查API响应状态
	if !detail.Success {
		logger.Warnw("ipwho api returned unsuccessful response", "ip", ip, "detail", detail)
	}

	return &detail, nil
}

// GetTimezone 获取IP地址对应的时区信息
// 参数:
//   - ctx: 上下文对象
//   - ip: 要查询的IP地址，为空时查询当前主机IP
//
// 返回值:
//   - 时区信息
//   - 错误信息（如果有）
func (c *Client) GetTimezone(ctx context.Context, ip string) (*Timezone, error) {
	detail, err := c.GetIPDetail(ctx, ip)
	if err != nil {
		return nil, err
	}
	return &detail.Timezone, nil
}

// GetConnection 获取IP地址对应的网络连接信息
// 参数:
//   - ctx: 上下文对象
//   - ip: 要查询的IP地址，为空时查询当前主机IP
//
// 返回值:
//   - 网络连接信息
//   - 错误信息（如果有）
func (c *Client) GetConnection(ctx context.Context, ip string) (*Connection, error) {
	detail, err := c.GetIPDetail(ctx, ip)
	if err != nil {
		return nil, err
	}
	return &detail.Connection, nil
}

// GetFlag 获取IP地址对应国家/地区的旗帜信息
// 参数:
//   - ctx: 上下文对象
//   - ip: 要查询的IP地址，为空时查询当前主机IP
//
// 返回值:
//   - 旗帜信息
//   - 错误信息（如果有）
func (c *Client) GetFlag(ctx context.Context, ip string) (*Flag, error) {
	detail, err := c.GetIPDetail(ctx, ip)
	if err != nil {
		return nil, err
	}
	return &detail.Flag, nil
}

// GetCurrentTime 获取IP地址对应时区的当前时间
// 参数:
//   - ctx: 上下文对象
//   - ip: 要查询的IP地址，为空时查询当前主机IP
//
// 返回值:
//   - 当前时间
//   - 错误信息（如果有）
func (c *Client) GetCurrentTime(ctx context.Context, ip string) (time.Time, error) {
	detail, err := c.GetIPDetail(ctx, ip)
	if err != nil {
		return time.Time{}, err
	}

	// 解析时间字符串
	if detail.Timezone.CurrentTime != "" {
		t, err := time.Parse(time.RFC3339, detail.Timezone.CurrentTime)
		if err != nil {
			log.W(ctx).Warnw("failed to parse timezone current time", "err", err, "currentTime", detail.Timezone.CurrentTime)
			return time.Now(), nil
		}
		return t, nil
	}

	return time.Now(), nil
}
