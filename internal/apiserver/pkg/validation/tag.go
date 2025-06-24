// Copyright 2025 长林啊 <767425412@qq.com>. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file. The original repo for
// this file is https://github.com/clin211/miniblog-v2.git.

package validation

import (
	"context"
	"regexp"
	"strconv"
	"strings"

	"github.com/clin211/miniblog-v2/internal/pkg/errno"
	apiv1 "github.com/clin211/miniblog-v2/pkg/api/apiserver/v1"
	genericvalidation "github.com/onexstack/onexstack/pkg/validation"
)

func (v *Validator) ValidateTagRules() genericvalidation.Rules {
	return genericvalidation.Rules{
		"TagID": func(id any) error {
			if id.(string) == "" {
				return errno.ErrInvalidArgument.WithMessage("tagID cannot be empty")
			}
			return nil
		},
		"Name": func(value any) error {
			name := value.(string)
			if name == "" {
				return errno.ErrInvalidArgument.WithMessage("name cannot be empty")
			}
			if len(strings.TrimSpace(name)) == 0 {
				return errno.ErrInvalidArgument.WithMessage("name cannot be whitespace only")
			}
			if len(name) > 20 {
				return errno.ErrInvalidArgument.WithMessage("name cannot exceed 20 characters")
			}
			return nil
		},
		"Color": func(value any) error {
			color := value.(string)
			if color == "" {
				return errno.ErrInvalidArgument.WithMessage("color cannot be empty")
			}

			color = strings.TrimSpace(color)
			if len(color) == 0 {
				return errno.ErrInvalidArgument.WithMessage("color cannot be whitespace only")
			}

			// 校验颜色值格式
			if !isValidColor(color) {
				return errno.ErrInvalidArgument.WithMessage("invalid color format. Supported formats: #HEX, rgb(), rgba(), hsl(), hsla(), and color names")
			}

			return nil
		},
	}
}

// isValidColor 校验颜色值格式是否正确；支持格式：十六进制(#HEX)、RGB、RGBA、HSL、HSLA、常见颜色名
func isValidColor(color string) bool {
	color = strings.ToLower(strings.TrimSpace(color))

	// 十六进制格式校验：#RGB 或 #RRGGBB
	if matched, _ := regexp.MatchString(`^#([0-9a-f]{3}|[0-9a-f]{6})$`, color); matched {
		return true
	}

	// RGB 格式校验：rgb(r,g,b)
	rgbPattern := regexp.MustCompile(`^rgb\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*\)$`)
	if matches := rgbPattern.FindStringSubmatch(color); matches != nil {
		r, _ := strconv.Atoi(matches[1])
		g, _ := strconv.Atoi(matches[2])
		b, _ := strconv.Atoi(matches[3])
		return r >= 0 && r <= 255 && g >= 0 && g <= 255 && b >= 0 && b <= 255
	}

	// RGBA 格式校验：rgba(r,g,b,a)
	rgbaPattern := regexp.MustCompile(`^rgba\(\s*(\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*,\s*([01]?\.?\d*)\s*\)$`)
	if matches := rgbaPattern.FindStringSubmatch(color); matches != nil {
		r, _ := strconv.Atoi(matches[1])
		g, _ := strconv.Atoi(matches[2])
		b, _ := strconv.Atoi(matches[3])
		a, _ := strconv.ParseFloat(matches[4], 64)
		return r >= 0 && r <= 255 && g >= 0 && g <= 255 && b >= 0 && b <= 255 && a >= 0 && a <= 1
	}

	// HSL 格式校验：hsl(h,s%,l%)
	hslPattern := regexp.MustCompile(`^hsl\(\s*(\d+)\s*,\s*(\d+)%\s*,\s*(\d+)%\s*\)$`)
	if matches := hslPattern.FindStringSubmatch(color); matches != nil {
		h, _ := strconv.Atoi(matches[1])
		s, _ := strconv.Atoi(matches[2])
		l, _ := strconv.Atoi(matches[3])
		return h >= 0 && h <= 360 && s >= 0 && s <= 100 && l >= 0 && l <= 100
	}

	// HSLA 格式校验：hsla(h,s%,l%,a)
	hslaPattern := regexp.MustCompile(`^hsla\(\s*(\d+)\s*,\s*(\d+)%\s*,\s*(\d+)%\s*,\s*([01]?\.?\d*)\s*\)$`)
	if matches := hslaPattern.FindStringSubmatch(color); matches != nil {
		h, _ := strconv.Atoi(matches[1])
		s, _ := strconv.Atoi(matches[2])
		l, _ := strconv.Atoi(matches[3])
		a, _ := strconv.ParseFloat(matches[4], 64)
		return h >= 0 && h <= 360 && s >= 0 && s <= 100 && l >= 0 && l <= 100 && a >= 0 && a <= 1
	}

	// 常见颜色名校验
	colorNames := map[string]bool{
		"red": true, "green": true, "blue": true, "yellow": true, "orange": true,
		"purple": true, "pink": true, "brown": true, "black": true, "white": true,
		"gray": true, "grey": true, "cyan": true, "magenta": true, "lime": true,
		"maroon": true, "navy": true, "olive": true, "teal": true, "silver": true,
		"aqua": true, "fuchsia": true, "darkred": true, "darkgreen": true, "darkblue": true,
		"lightyellow": true, "lightgreen": true, "lightblue": true, "lightgray": true,
		"lightgrey": true, "darkgray": true, "darkgrey": true, "gold": true, "violet": true,
		"indigo": true, "coral": true, "salmon": true, "khaki": true, "crimson": true,
		"tomato": true, "chocolate": true, "peru": true, "tan": true, "beige": true,
	}

	return colorNames[color]
}

// ValidateCreateTagRequest 校验 CreateTagRequest 结构体的有效性.
func (v *Validator) ValidateCreateTagRequest(ctx context.Context, rq *apiv1.CreateTagRequest) error {
	return genericvalidation.ValidateAllFields(rq, v.ValidateTagRules())
}

// ValidateUpdateTagRequest 校验 UpdateTagRequest 结构体的有效性.
func (v *Validator) ValidateUpdateTagRequest(ctx context.Context, rq *apiv1.UpdateTagRequest) error {
	// 先校验 ID 字段
	if rq.GetId() <= 0 {
		return errno.ErrInvalidArgument.WithMessage("tag ID must be positive")
	}
	return genericvalidation.ValidateAllFields(rq, v.ValidateTagRules())
}

// ValidateDeleteTagRequest 校验 DeleteTagRequest 结构体的有效性.
func (v *Validator) ValidateDeleteTagRequest(ctx context.Context, rq *apiv1.DeleteTagRequest) error {
	if rq.GetId() <= 0 {
		return errno.ErrInvalidArgument.WithMessage("tag ID must be positive")
	}
	return nil
}

// ValidateGetTagRequest 校验 GetTagRequest 结构体的有效性.
func (v *Validator) ValidateGetTagRequest(ctx context.Context, rq *apiv1.GetTagRequest) error {
	if rq.GetId() <= 0 {
		return errno.ErrInvalidArgument.WithMessage("tag ID must be positive")
	}
	return nil
}

// ValidateListTagRequest 校验 ListTagRequest 结构体的有效性.
func (v *Validator) ValidateListTagRequest(ctx context.Context, rq *apiv1.ListTagRequest) error {
	// 校验可选的名称过滤参数
	if name := rq.GetName(); name != "" {
		if len(strings.TrimSpace(name)) == 0 {
			return errno.ErrInvalidArgument.WithMessage("name filter cannot be whitespace only")
		}
		if len(name) > 20 {
			return errno.ErrInvalidArgument.WithMessage("name filter cannot exceed 20 characters")
		}
	}
	return nil
}
