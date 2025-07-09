// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package ipwho

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockIPDetail 返回一个模拟的IP详情数据
func mockIPDetail() *IPDetail {
	return &IPDetail{
		IP:            "8.8.8.8",
		Success:       true,
		Type:          "IPv4",
		Continent:     "North America",
		ContinentCode: "NA",
		Country:       "United States",
		CountryCode:   "US",
		Region:        "California",
		RegionCode:    "CA",
		City:          "Mountain View",
		Latitude:      37.4056,
		Longitude:     -122.0775,
		IsEu:          false,
		Postal:        "94043",
		CallingCode:   "+1",
		Capital:       "Washington",
		Borders:       "CA,MX",
		Flag: Flag{
			Img:          "https://cdn.ipwhois.io/flags/us.svg",
			Emoji:        "🇺🇸",
			EmojiUnicode: "U+1F1FA U+1F1F8",
		},
		Connection: Connection{
			Asn:    15169,
			Org:    "Google LLC",
			Isp:    "Google LLC",
			Domain: "google.com",
		},
		Timezone: Timezone{
			ID:          "America/Los_Angeles",
			Abbr:        "PST",
			IsDst:       false,
			Offset:      -28800,
			Utc:         "-08:00",
			CurrentTime: "2024-01-15T10:30:45-08:00",
		},
	}
}

func TestWithBaseURL(t *testing.T) {
	tests := []struct {
		name     string
		baseURL  string
		expected string
	}{
		{
			name:     "valid URL",
			baseURL:  "https://api.example.com/",
			expected: "https://api.example.com/",
		},
		{
			name:     "empty URL should not change default",
			baseURL:  "",
			expected: "https://ipwho.is/",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{}
			WithBaseURL(tt.baseURL)(opts)
			getOptionsOrSetDefault(opts)
			assert.Equal(t, tt.expected, opts.BaseURL)
		})
	}
}

func TestWithUserAgent(t *testing.T) {
	tests := []struct {
		name      string
		userAgent string
		expected  string
	}{
		{
			name:      "custom user agent",
			userAgent: "MyApp/1.0",
			expected:  "MyApp/1.0",
		},
		{
			name:      "empty user agent should not change default",
			userAgent: "",
			expected:  "curl/7.77.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{}
			WithUserAgent(tt.userAgent)(opts)
			getOptionsOrSetDefault(opts)
			assert.Equal(t, tt.expected, opts.UserAgent)
		})
	}
}

func TestWithTimeout(t *testing.T) {
	tests := []struct {
		name     string
		timeout  time.Duration
		expected time.Duration
	}{
		{
			name:     "custom timeout",
			timeout:  10 * time.Second,
			expected: 10 * time.Second,
		},
		{
			name:     "zero timeout should not change default",
			timeout:  0,
			expected: 30 * time.Second,
		},
		{
			name:     "negative timeout should not change default",
			timeout:  -5 * time.Second,
			expected: 30 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{}
			WithTimeout(tt.timeout)(opts)
			getOptionsOrSetDefault(opts)
			assert.Equal(t, tt.expected, opts.Timeout)
		})
	}
}

func TestWithCustomHeaders(t *testing.T) {
	tests := []struct {
		name     string
		headers  map[string]string
		expected map[string]string
	}{
		{
			name: "custom headers",
			headers: map[string]string{
				"Authorization": "Bearer token",
				"X-Custom":      "value",
			},
			expected: map[string]string{
				"User-Agent":    "curl/7.77.0",
				"Authorization": "Bearer token",
				"X-Custom":      "value",
			},
		},
		{
			name:    "nil headers should not panic",
			headers: nil,
			expected: map[string]string{
				"User-Agent": "curl/7.77.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &Options{}
			WithCustomHeaders(tt.headers)(opts)
			getOptionsOrSetDefault(opts)
			for k, v := range tt.expected {
				assert.Equal(t, v, opts.CustomHeaders[k])
			}
		})
	}
}

func TestGetOptionsOrSetDefault(t *testing.T) {
	t.Run("nil options should return default", func(t *testing.T) {
		opts := getOptionsOrSetDefault(nil)
		assert.Equal(t, "https://ipwho.is/", opts.BaseURL)
		assert.Equal(t, "curl/7.77.0", opts.UserAgent)
		assert.Equal(t, 30*time.Second, opts.Timeout)
		assert.NotNil(t, opts.CustomHeaders)
		assert.Equal(t, "curl/7.77.0", opts.CustomHeaders["User-Agent"])
	})

	t.Run("non-nil options should return same instance", func(t *testing.T) {
		original := &Options{BaseURL: "custom"}
		result := getOptionsOrSetDefault(original)
		assert.Same(t, original, result)
	})
}

func TestNewClient(t *testing.T) {
	t.Run("default client", func(t *testing.T) {
		client := NewClient()
		assert.NotNil(t, client)
		assert.Equal(t, "https://ipwho.is/", client.options.BaseURL)
		assert.Equal(t, "curl/7.77.0", client.options.UserAgent)
		assert.Equal(t, 30*time.Second, client.options.Timeout)
	})

	t.Run("client with custom options", func(t *testing.T) {
		client := NewClient(
			WithBaseURL("https://api.example.com/"),
			WithUserAgent("TestApp/1.0"),
			WithTimeout(15*time.Second),
		)
		assert.NotNil(t, client)
		assert.Equal(t, "https://api.example.com/", client.options.BaseURL)
		assert.Equal(t, "TestApp/1.0", client.options.UserAgent)
		assert.Equal(t, 15*time.Second, client.options.Timeout)
	})
}

func TestGetIPDetail(t *testing.T) {
	mockData := mockIPDetail()

	t.Run("successful request with IP", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/8.8.8.8", r.URL.Path)
			assert.Equal(t, "TestApp/1.0", r.Header.Get("User-Agent"))

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockData)
		}))
		defer server.Close()

		client := NewClient(
			WithBaseURL(server.URL+"/"),
			WithUserAgent("TestApp/1.0"),
		)

		detail, err := client.GetIPDetail(context.Background(), "8.8.8.8")
		require.NoError(t, err)
		assert.Equal(t, mockData.IP, detail.IP)
		assert.Equal(t, mockData.Success, detail.Success)
		assert.Equal(t, mockData.Country, detail.Country)
		assert.Equal(t, mockData.City, detail.City)
	})

	t.Run("successful request without IP (host IP)", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			assert.Equal(t, "/", r.URL.Path)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockData)
		}))
		defer server.Close()

		client := NewClient(WithBaseURL(server.URL + "/"))

		detail, err := client.GetIPDetail(context.Background(), "")
		require.NoError(t, err)
		assert.Equal(t, mockData.IP, detail.IP)
	})

	t.Run("API returns unsuccessful response", func(t *testing.T) {
		unsuccessfulData := *mockData
		unsuccessfulData.Success = false

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(unsuccessfulData)
		}))
		defer server.Close()

		client := NewClient(WithBaseURL(server.URL + "/"))

		detail, err := client.GetIPDetail(context.Background(), "8.8.8.8")
		require.NoError(t, err)
		assert.False(t, detail.Success)
	})

	t.Run("server returns non-200 status", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("Not Found"))
		}))
		defer server.Close()

		client := NewClient(WithBaseURL(server.URL + "/"))

		detail, err := client.GetIPDetail(context.Background(), "8.8.8.8")
		assert.Error(t, err)
		assert.Nil(t, detail)
		assert.Contains(t, err.Error(), "status 404")
	})

	t.Run("invalid JSON response", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		client := NewClient(WithBaseURL(server.URL + "/"))

		detail, err := client.GetIPDetail(context.Background(), "8.8.8.8")
		assert.Error(t, err)
		assert.Nil(t, detail)
		assert.Contains(t, err.Error(), "failed to parse response")
	})

	t.Run("network error", func(t *testing.T) {
		client := NewClient(WithBaseURL("http://non-existent-domain.invalid/"))

		detail, err := client.GetIPDetail(context.Background(), "8.8.8.8")
		assert.Error(t, err)
		assert.Nil(t, detail)
		assert.Contains(t, err.Error(), "ipwho request failed")
	})

	t.Run("context timeout", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			time.Sleep(100 * time.Millisecond)
			json.NewEncoder(w).Encode(mockData)
		}))
		defer server.Close()

		client := NewClient(
			WithBaseURL(server.URL+"/"),
			WithTimeout(10*time.Millisecond),
		)

		ctx := context.Background()
		detail, err := client.GetIPDetail(ctx, "8.8.8.8")
		assert.Error(t, err)
		assert.Nil(t, detail)
	})
}

func TestGetHostIPDetail(t *testing.T) {
	mockData := mockIPDetail()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockData)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL + "/"))

	detail, err := client.GetHostIPDetail(context.Background())
	require.NoError(t, err)
	assert.Equal(t, mockData.IP, detail.IP)
}

func TestGetTimezone(t *testing.T) {
	mockData := mockIPDetail()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockData)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL + "/"))

	timezone, err := client.GetTimezone(context.Background(), "8.8.8.8")
	require.NoError(t, err)
	assert.Equal(t, mockData.Timezone.ID, timezone.ID)
	assert.Equal(t, mockData.Timezone.Abbr, timezone.Abbr)
	assert.Equal(t, mockData.Timezone.Offset, timezone.Offset)
}

func TestGetConnection(t *testing.T) {
	mockData := mockIPDetail()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockData)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL + "/"))

	connection, err := client.GetConnection(context.Background(), "8.8.8.8")
	require.NoError(t, err)
	assert.Equal(t, mockData.Connection.Asn, connection.Asn)
	assert.Equal(t, mockData.Connection.Org, connection.Org)
	assert.Equal(t, mockData.Connection.Isp, connection.Isp)
	assert.Equal(t, mockData.Connection.Domain, connection.Domain)
}

func TestGetFlag(t *testing.T) {
	mockData := mockIPDetail()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockData)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL + "/"))

	flag, err := client.GetFlag(context.Background(), "8.8.8.8")
	require.NoError(t, err)
	assert.Equal(t, mockData.Flag.Img, flag.Img)
	assert.Equal(t, mockData.Flag.Emoji, flag.Emoji)
	assert.Equal(t, mockData.Flag.EmojiUnicode, flag.EmojiUnicode)
}

func TestGetCurrentTime(t *testing.T) {
	t.Run("valid timezone current time", func(t *testing.T) {
		mockData := mockIPDetail()

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockData)
		}))
		defer server.Close()

		client := NewClient(WithBaseURL(server.URL + "/"))

		currentTime, err := client.GetCurrentTime(context.Background(), "8.8.8.8")
		require.NoError(t, err)

		expected, _ := time.Parse(time.RFC3339, mockData.Timezone.CurrentTime)
		assert.True(t, currentTime.Equal(expected))
	})

	t.Run("invalid timezone current time", func(t *testing.T) {
		mockData := mockIPDetail()
		mockData.Timezone.CurrentTime = "invalid-time-format"

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockData)
		}))
		defer server.Close()

		client := NewClient(WithBaseURL(server.URL + "/"))

		currentTime, err := client.GetCurrentTime(context.Background(), "8.8.8.8")
		require.NoError(t, err)
		// Should return current time when parsing fails
		assert.True(t, time.Since(currentTime) < time.Minute)
	})

	t.Run("empty timezone current time", func(t *testing.T) {
		mockData := mockIPDetail()
		mockData.Timezone.CurrentTime = ""

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(mockData)
		}))
		defer server.Close()

		client := NewClient(WithBaseURL(server.URL + "/"))

		currentTime, err := client.GetCurrentTime(context.Background(), "8.8.8.8")
		require.NoError(t, err)
		// Should return current time when empty
		assert.True(t, time.Since(currentTime) < time.Minute)
	})

	t.Run("GetIPDetail error", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		defer server.Close()

		client := NewClient(WithBaseURL(server.URL + "/"))

		currentTime, err := client.GetCurrentTime(context.Background(), "8.8.8.8")
		assert.Error(t, err)
		assert.True(t, currentTime.IsZero())
	})
}

func TestCustomHeaders(t *testing.T) {
	customHeaders := map[string]string{
		"Authorization": "Bearer token123",
		"X-Custom":      "test-value",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
		assert.Equal(t, "test-value", r.Header.Get("X-Custom"))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockIPDetail())
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL+"/"),
		WithCustomHeaders(customHeaders),
	)

	_, err := client.GetIPDetail(context.Background(), "8.8.8.8")
	require.NoError(t, err)
}

// TestIPDetailJSONSerialization tests JSON marshaling/unmarshaling of IPDetail
func TestIPDetailJSONSerialization(t *testing.T) {
	original := mockIPDetail()

	// Marshal to JSON
	jsonData, err := json.Marshal(original)
	require.NoError(t, err)

	// Unmarshal back
	var restored IPDetail
	err = json.Unmarshal(jsonData, &restored)
	require.NoError(t, err)

	// Compare
	assert.Equal(t, original.IP, restored.IP)
	assert.Equal(t, original.Success, restored.Success)
	assert.Equal(t, original.Country, restored.Country)
	assert.Equal(t, original.Flag.Emoji, restored.Flag.Emoji)
	assert.Equal(t, original.Connection.Asn, restored.Connection.Asn)
	assert.Equal(t, original.Timezone.ID, restored.Timezone.ID)
}

// BenchmarkGetIPDetail benchmarks the GetIPDetail method
func BenchmarkGetIPDetail(b *testing.B) {
	mockData := mockIPDetail()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(mockData)
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL + "/"))
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := client.GetIPDetail(ctx, "8.8.8.8")
		if err != nil {
			b.Fatal(err)
		}
	}
}

// Example demonstrates how to use the ipwho client
func ExampleClient_GetIPDetail() {
	client := NewClient(
		WithTimeout(10*time.Second),
		WithUserAgent("MyApp/1.0"),
	)

	detail, err := client.GetIPDetail(context.Background(), "8.8.8.8")
	if err != nil {
		panic(err)
	}

	println("Country:", detail.Country)
	println("City:", detail.City)
	println("ISP:", detail.Connection.Isp)
}

// Example demonstrates how to get current time for an IP
func ExampleClient_GetCurrentTime() {
	client := NewClient()

	currentTime, err := client.GetCurrentTime(context.Background(), "8.8.8.8")
	if err != nil {
		panic(err)
	}

	println("Current time at IP location:", currentTime.Format(time.RFC3339))
}
