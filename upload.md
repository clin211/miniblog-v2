# 文件上传功能需求说明

## 背景与目标

- **目标**: 为系统提供统一的文件上传能力，支持按配置选择本地存储或阿里云 OSS 存储（AliOSS）。
- **范围**: 仅包含上传接口、配置项、基础安全与校验、错误码与返回体约定；不包含具体 SDK 适配实现细节与控制台配置。
- **非目标**: 富媒体处理（裁剪/转码/水印）、CDN 配置、回调通知流水线等。

### 架构概述

- **存储适配层**: 通过统一接口抽象 `Uploader`，实现 `local` 与 `alioss` 两个 Provider，后续可扩展。
- **选择策略**: 根据配置 `upload.provider` 在运行时选择对应实现。
- **鉴权**: 复用现有 HTTP 鉴权与权限中间件，默认仅登录用户可上传。

## 配置设计（mb-apiserver.yaml）

在 `configs/mb-apiserver.yaml` 中新增/扩展 `upload` 段：

```yaml
upload:
  # 存储提供方: local | alioss
  provider: local

  # 单文件最大限制（支持 10MB、200MB、1GB 等写法），超出返回 413/业务码
  maxSize: 20MB

  # 允许的 MIME 类型白名单（若为空表示不限制）
  allowedMIMEs:
    - image/jpeg
    - image/png
    - application/pdf

  # 是否开启内容去重（按哈希），开启后相同内容复用同一对象
  deduplicate: true

  # 对象键模板，用于统一命名（日期/哈希/扩展名占位）
  # 支持占位: {date:2006/01/02} {sha256[:N]} {ext} {uuid}
  keyTemplate: "{date:2006/01/02}/{sha256:16}{ext}"

  # 本地存储配置（provider=local 生效）
  local:
    baseDir: "/data/miniblog/uploads"      # 物理存储目录
    baseURL: "https://api.example.com/static/uploads" # 对外访问前缀
    mkdirPerm: 0755

  # 阿里云 OSS 配置（provider=alioss 生效）
  alioss:
    endpoint: "oss-cn-hangzhou.aliyuncs.com"
    bucket: "miniblog"
    accessKeyID: "${ALIOSS_ACCESS_KEY_ID}"
    accessKeySecret: "${ALIOSS_ACCESS_KEY_SECRET}"
    # 可选: 若使用 STS 临时凭证
    securityToken: "${ALIOSS_STS_TOKEN}"
    acl: "private"                # private | public-read
    pathStyleAccess: false
    uploadTimeout: "30s"
```

说明:

- **baseURL**: 本地存储用于拼接返回的对外可访问 URL，需要由反向代理或应用静态资源路由映射到 `baseDir`。
- **keyTemplate**: 建议包含日期与内容哈希，降低热目录与重复文件风险。

### 分片上传配置（新增）

为满足大文件与断点续传需求，新增 `upload.multipart` 段：

```yaml
upload:
  multipart:
    enabled: true                 # 是否启用分片上传
    minSize: 8MB                  # 触发分片上传的建议最小文件大小（客户端可参考）
    maxParts: 10000               # 最大分片数（对齐 OSS 限制）
    partSize:
      min: 5MB                    # 单片最小 5MB（对齐 OSS 要求）
      max: 128MB                  # 单片最大建议值
      default: 8MB                # 未指定时的推荐分片大小
    concurrencyLimitPerUser: 2    # 每用户并发中的多分片会话上限
    uploadIdTTL: 24h              # 多分片会话过期时间（后台清理）
    tempDir: "/data/miniblog/uploads/.multipart"   # 本地 provider 临时目录
    checksum: sha256              # 分片与整文件校验算法: sha256 | md5
    presign:
      enabled: true               # 是否返回预签名直传 URL（alioss 场景）
      expires: 15m                # 预签名 URL 过期时间
      mode: direct                # 上传数据路径: proxy | direct

  # alioss 扩展（可选）
  alioss:
    # ... 同上
    callbackURL: ""              # 可选：OSS 回调地址（若采用 direct 模式且需要回调）
```

### 路由与中间件设计

- **基础前缀**: `POST /v1/upload/file`
- **鉴权**: 默认启用 `Authn` 与 `AccessLogger` 与 `RequestID` 中间件；根据角色可扩展 `Authz`（如仅特定角色可上传）。

示例（Gin 路由注册，仅示意，不含实现）：

```go
package system

import (
    "github.com/gin-gonic/gin"
    ginmw "github.com/clin211/miniblog-v2/internal/pkg/middleware/gin"
)

// RegisterUploadRoutes 注册上传路由。
func RegisterUploadRoutes(r *gin.RouterGroup, uc *UploadController) {
    g := r.Group("/upload",
        ginmw.RequestID(),
        ginmw.AccessLogger(),
        ginmw.Authn(),
        // 可选: ginmw.Authz()
    )
    {
        g.POST("/file", uc.UploadFile)
    }
}

// UploadController 仅示意声明。
type UploadController struct{}

// UploadFile 仅示意签名，不含实现逻辑。
func (uc *UploadController) UploadFile(c *gin.Context) { /* ... */ }
```

在 HTTP Server 初始化位置（仅 provider=local 时）可增加静态目录映射，暴露 `baseURL`：

```go
// 仅在本地存储场景开启:
router.StaticFS("/static/uploads", http.Dir(cfg.Upload.Local.BaseDir))
```

### 请求与响应规范

- **接口**: `POST /v1/upload/file`
- **Content-Type**: `multipart/form-data`
- **表单字段**:
  - `file`(必填): 文件二进制内容
  - `scene`(可选): 上传场景标签（如 avatar, post-image）
  - `filename`(可选): 建议文件名（若未提供则按模板与原始扩展名生成）

成功响应:

```json
{
  "provider": "local",                
  "key": "2025/01/31/abcdef1234567890.jpg",
  "url": "https://api.example.com/static/uploads/2025/01/31/abcdef1234567890.jpg",
  "size": 345678,
  "mime": "image/jpeg",
  "hash": "sha256:5a1f...",
  "metadata": {"scene": "avatar"}
}
```

失败响应（统一错误格式，遵循 `internal/pkg/errno` 写法）：

```json
{
  "reason": "UploadTooLarge",
  "message": "file exceeds maxSize: 20MB",
  "metadata": {"limit": "20MB"}
}
```

## 分片上传（Multipart）

### 适用场景与模式

- **场景**: 大文件上传、弱网断点续传、移动端直传云存储。
- **模式**:
  - proxy: 客户端将分片上传至应用服务，再由服务持久化（local/OSS）。
  - direct: 客户端直传至 OSS（预签名 URL），应用服务仅负责会话与完成提交。

### 生命周期与状态

- Initiated → Uploading → (Completed | Aborted | Expired)
- 会话由 `uploadId` 标识，存储于会话层（建议 Redis），含 TTL（见配置 `upload.multipart.uploadIdTTL`）。

### API 概览

- `POST   /v1/upload/multipart/init`          初始化会话，返回 `uploadId`、`partSize`、`key`、`presign` 策略
- `POST   /v1/upload/multipart/presign`       批量获取直传分片的预签名 URL（direct 模式）
- `PUT    /v1/upload/multipart/part`          上传单个分片（proxy 模式）
- `GET    /v1/upload/multipart/:uploadId/parts` 查询已上传分片列表
- `POST   /v1/upload/multipart/complete`      提交完成，服务端合并/完成分片
- `DELETE /v1/upload/multipart/abort`         放弃会话并清理

### 路由示例（Gin，仅示意）

```go
func RegisterUploadRoutes(r *gin.RouterGroup, uc *UploadController) {
    g := r.Group("/upload",
        ginmw.RequestID(),
        ginmw.AccessLogger(),
        ginmw.Authn(),
    )
    {
        g.POST("/file", uc.UploadFile)
        // multipart
        g.POST("/multipart/init", uc.InitMultipart)
        g.POST("/multipart/presign", uc.PresignParts)   // direct 模式
        g.PUT("/multipart/part", uc.UploadPart)         // proxy 模式
        g.GET("/multipart/:uploadId/parts", uc.ListParts)
        g.POST("/multipart/complete", uc.CompleteMultipart)
        g.DELETE("/multipart/abort", uc.AbortMultipart)
    }
}

type UploadController struct{}
func (uc *UploadController) InitMultipart(c *gin.Context)       {}
func (uc *UploadController) PresignParts(c *gin.Context)        {}
func (uc *UploadController) UploadPart(c *gin.Context)          {}
func (uc *UploadController) ListParts(c *gin.Context)           {}
func (uc *UploadController) CompleteMultipart(c *gin.Context)   {}
func (uc *UploadController) AbortMultipart(c *gin.Context)      {}
```

### 流程说明

1) init
   - 入参: `filename`、`size`、`mime`、`scene`、`sha256`(可选，全量哈希用于去重)、`partSize`(可选)。
   - 出参: `uploadId`、`key`、`partSize`、`mode`(proxy|direct)、`presign.expires`、`duplicated`（若命中去重直接返回最终对象）。
2) presign（direct 模式）
   - 入参: `uploadId`、`partNumbers` 数组。
   - 出参: `[{partNumber, url, headers, expiresAt}]`。
3) part 上传
   - proxy: `PUT /multipart/part?uploadId=...&partNumber=...`，Body 为分片字节；校验 `Content-Length`、`Content-MD5`/`checksum`。
   - direct: 客户端直传到 OSS 的分片 URL，成功后保存 `partNumber` 与 `ETag`（客户端可回传或由回调采集）。
4) list parts
   - 返回已成功的 `partNumber` 列表及 `ETag`（若有）。
5) complete
   - 入参: `uploadId`、`parts:[{partNumber, etag}]`（direct 模式必传 ETag）。
   - 服务端完成本地合并或调用 OSS CompleteMultipartUpload；校验整文件 `sha256`（若 init 提供）。
6) abort
   - 清理未完成的会话与临时资源；若 direct 模式调用 OSS AbortMultipartUpload。

### 请求与响应定义（示例）

1) init

```json
// request
{
  "filename": "video.mp4",
  "size": 1073741824,
  "mime": "video/mp4",
  "scene": "post-video",
  "sha256": "5a1f...", 
  "partSize": 8388608
}
```

```json
// response
{
  "uploadId": "u_01J9...",
  "key": "2025/01/31/abcd1234.mp4",
  "mode": "direct",
  "partSize": 8388608,
  "presign": {"expires": "15m"},
  "duplicated": false
}
```

2) presign（direct）

```json
// request
{"uploadId":"u_01J9...","partNumbers":[1,2,3]}
```

```json
// response
{
  "items": [
    {"partNumber":1, "url":"https://...", "headers":{"Content-Type":"application/octet-stream"}, "expiresAt":"2025-01-31T12:00:00Z"}
  ]
}
```

3) part（proxy）

- `PUT /v1/upload/multipart/part?uploadId=u_01J9...&partNumber=1`

```json
// response
{"partNumber":1, "etag":"\"9b2cf535f27731c974343645a3985328\"", "size": 8388608, "checksum":"sha256:..."}
```

4) list parts

```json
// response
{"uploadId":"u_01J9...","parts":[{"partNumber":1,"etag":"\"...\""}]}
```

5) complete

```json
// request
{"uploadId":"u_01J9...","parts":[{"partNumber":1,"etag":"\"...\""}]}
```

```json
// response （与单文件上传成功体一致）
{
  "provider":"alioss",
  "key":"2025/01/31/xxx.mp4",
  "url":"https://bucket.oss-cn-hangzhou.aliyuncs.com/2025/01/31/xxx.mp4",
  "size":1073741824,
  "mime":"video/mp4",
  "hash":"sha256:...",
  "metadata":{"scene":"post-video"}
}
```

### OpenAPI（Swagger）示意（分片）

```yaml
paths:
  /v1/upload/multipart/init:
    post:
      summary: Init multipart upload
      consumes: [application/json]
      responses: { '200': { description: OK } }
  /v1/upload/multipart/presign:
    post:
      summary: Batch presign for direct upload
      consumes: [application/json]
      responses: { '200': { description: OK } }
  /v1/upload/multipart/part:
    put:
      summary: Upload single part via proxy
      consumes: [application/octet-stream]
      parameters:
        - in: query
          name: uploadId
          type: string
          required: true
        - in: query
          name: partNumber
          type: integer
          required: true
      responses: { '200': { description: OK }, '413': { description: Payload Too Large } }
  /v1/upload/multipart/{uploadId}/parts:
    get:
      summary: List uploaded parts
      parameters:
        - in: path
          name: uploadId
          type: string
          required: true
      responses: { '200': { description: OK } }
  /v1/upload/multipart/complete:
    post:
      summary: Complete multipart upload
      responses: { '200': { description: OK } }
  /v1/upload/multipart/abort:
    delete:
      summary: Abort multipart upload
      responses: { '200': { description: OK } }
```

### 校验、安全与幂等

- 触发分片的大小与单片大小限制遵循配置；校验 `partNumber` 范围 `1..maxParts`。
- proxy 模式强制校验每片大小；最后在 complete 校验整文件哈希（若 init 提供）。
- direct 模式需保留 `ETag` 与顺序，complete 时提交给 OSS；若 `acl=private` 可仅返回应用代理 URL。
- 幂等：
  - init 支持 `Idempotency-Key` 头，相同键在 TTL 内返回同一 `uploadId`；
  - 重复上传相同 `partNumber` 时对比校验和，相同则复用结果；
  - complete 重放安全，若已完成返回最终对象。
- 恢复：`GET /parts` 用于断点续传恢复与客户端重试。

### 清理与配额

- 过期会话定期 GC：删除本地临时分片、调用 OSS AbortMultipartUpload。
- 并发限制：按用户限制同时活跃会话 `concurrencyLimitPerUser`；超限返回 429。

### 分片相关错误码（reason）

- `UploadMultipartInitFailed`：初始化失败/配置缺失。
- `UploadMultipartSessionNotFound`：会话不存在或已过期。
- `UploadMultipartPartInvalid`：分片编号或大小非法。
- `UploadMultipartChecksumMismatch`：分片/整文件校验不一致。
- `UploadMultipartPresignFailed`：预签名失败。
- `UploadMultipartCompleteFailed`：完成合并失败。
- `UploadMultipartAbortFailed`：中止清理失败。

### 验收补充（分片）

- direct 模式：客户端可通过预签名 URL 并发直传，complete 后返回可访问 URL；断点续传可恢复。
- proxy 模式：服务可稳定接收大文件分片并合并，通过大小/MIME/校验策略把关。

### OpenAPI（Swagger）示意

```yaml
paths:
  /v1/upload/file:
    post:
      summary: Upload a file
      consumes:
        - multipart/form-data
      parameters:
        - in: formData
          name: file
          type: file
          required: true
        - in: formData
          name: scene
          type: string
          required: false
        - in: formData
          name: filename
          type: string
          required: false
      responses:
        '200':
          description: OK
          schema:
            type: object
            properties:
              provider: { type: string }
              key: { type: string }
              url: { type: string }
              size: { type: integer, format: int64 }
              mime: { type: string }
              hash: { type: string }
              metadata: { type: object, additionalProperties: { type: string } }
        '400': { description: Bad Request }
        '401': { description: Unauthorized }
        '413': { description: Payload Too Large }
        '500': { description: Internal Server Error }
```

### 校验与安全

- **大小限制**: 读取前强校验 `Content-Length` 与流式读限制，超限直接拒绝。
- **MIME 白名单**: 基于报文头与文件魔数双重校验；不匹配时拒绝。
- **路径安全（本地）**: 禁止相对路径跳转（`..` 等），统一在 `baseDir` 下生成。
- **权限控制**: 默认需登录；如需公开匿名上传应单独开关并限流。
- **限流与防刷**: 基于 IP/User 的速率限制；建议接入现有中间件或网关策略。

### 关键行为约定

- **去重策略**: `deduplicate=true` 时，哈希一致直接复用历史对象与 URL。
- **命名策略**: 统一按 `keyTemplate` 生成；可基于 `scene` 增加子目录（实现层处理）。
- **可访问性**:
  - local: 通过 `baseURL` 暴露静态访问（需反向代理/静态路由）。
  - alioss: 返回 OSS 公网 URL（若 `acl=private`，可返回应用代理 URL 或临时签名 URL，具体实现另行定义）。

### 错误码建议（reason）

- `UploadProviderNotConfigured`：未配置或配置不完整。
- `UploadTooLarge`：超过 `maxSize`。
- `UploadMimeNotAllowed`：MIME 不在白名单。
- `UploadReadFailed`：读取表单或文件失败。
- `UploadStoreFailed`：写入底层存储失败（含本地/OSS）。
- `UploadHashFailed`：计算哈希失败。

上述 reason 对应的 HTTP 状态码建议：

- 400: 参数/表单错误、MIME 不允许。
- 401: 未鉴权。
- 403: 鉴权通过但无上传权限（若启用 `Authz`）。
- 413: 超出大小限制。
- 500: 存储失败/未知错误。

### 开发与测试要点

- 单测覆盖：大小/MIME 校验、命名模板展开、去重、不同 Provider 的分支选择。
- 本地联调：
  1. 配置 `provider=local`，创建 `baseDir`，开启静态映射；
  2. 使用 `curl`/Postman 调用 `POST /v1/upload/file`。
- OSS 联调：
  1. 配置 `provider=alioss` 与凭证；
  2. 确认 Bucket、ACL、跨域（CORS）策略；
  3. 校验返回 URL 可用性。

### 迁移与回滚

- 新增配置为向后兼容，默认为 `local`，不影响现有功能。
- 回滚仅需恢复旧配置并禁用上传路由注册。

### 验收标准

- 配置为 `local` 时：文件可上传、返回可访问 URL，大小/MIME 校验生效。
- 配置为 `alioss` 时：文件可上传至指定 Bucket，返回 URL 正常访问（或签名 URL 可访问）。
- 错误场景返回规范、日志可追踪、含请求 ID。
