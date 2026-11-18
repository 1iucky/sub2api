# shengtu 生图尺寸实测报告

日期：2026-10-13。网关：https://siliconbase.link。当前仓库快照：64a212741。

## 结论

已复现尺寸参数未得到满足。六次成功生图均约 157 万像素，被平台归入 2K 计费档。请求 JPEG/WebP 时也实际返回 PNG。

可以排除 skill 本地缩图：它原样发送 size、quality、output_format，直接保存响应中的图片字节。问题位于网关服务端至上游返回之间。当前仓库没有发现强制缩到 2K 的逻辑，但缺少线上部署版本、实际出站请求与账号映射证据，尚不能确定是线上平台路由/配置还是上游忽略参数，也不能据此断言模型最大只支持 2K。

## 测试方法与模型

用户指定的 `gpt-image-2.5` 实测返回 HTTP 404，错误为 model_not_found：Model "gpt-image-2.5" is not available for this group。

同一专用密钥调用 /v1/models，实际可用的 2.5 型号是 `gpt-image-2.5-flare` 和 `gpt-image-2.5-sunburst`。后续测试使用这两个完整名称。

诊断脚本导入已安装 generate_image.py，调用其 main、request_json、save_image；观察包装仅记录参数和白名单响应字段，不修改请求或图片。提示词统一为灰色背景上的红色陶瓷杯，n=1，非流式 /v1/images/generations。Pillow 从文件内容读取真实尺寸和格式。密钥没有写入报告或结果文件。

## 实测结果

| 模型后缀 | 请求尺寸 | 请求质量 | 请求格式 | 实际像素 | 实际格式 | 响应质量 |
|---|---|---|---|---|---|---|
| flare | 1024×1024 | low | PNG | 1376×1143 | PNG | low |
| flare | 3840×2160 | low | PNG | 1312×1199 | PNG | medium |
| flare | 3840×2160 | high | PNG | 1370×1148 | PNG | medium |
| flare | 3840×2160 | low | JPEG | 1254×1254 | PNG | low |
| flare | 3840×2160 | low | WebP | 1254×1254 | PNG | low |
| sunburst | 3840×2160 | low | PNG | 1254×1254 | PNG | low |

六次输出为约 1.5725～1.5731 MP。这是样本的共同特征，不是已证明的模型上限。每个组合仅测一次；未测试 edits、不同账号或绕过平台直接调用上游。

## 请求追踪

| 测试 | x-request-id | 耗时（秒） |
|---|---|---|
| flare_1k_png_low | e98af7ec-2a17-4294-9178-175f16432a38 | 23.07 |
| flare_4k_png_low | 344d313c-e257-4ca6-b1b7-25a755bd80ab | 27.75 |
| flare_4k_png_high | f45219f1-dc38-4f1f-80e6-3b546e12b67d | 41.53 |
| flare_4k_jpeg_low | 030e1095-0807-4a01-87b9-44dd701fb6c1 | 25.45 |
| flare_4k_webp_low | 312d13ea-fb09-474b-8e48-498248b106a2 | 35.88 |
| sunburst_4k_png_low | a51bc989-9100-43d7-ba70-c9e6f87144b6 | 27.02 |

## 代码证据

- `/Users/liuliang/.codex/skills/shengtu-skill/scripts/generate_image.py:112`：请求原样携带 size、quality、output_format。
- 同文件 `:198`：base64 解码或下载后直接写字节，没有缩放或重新编码。
- `backend/internal/service/openai_images_direct.go:67`：OAuth 直调保留上述三个字段。
- `backend/internal/service/openai_images.go:874`：API Key 路径的 JSON 仅改写 model；multipart 也保留其他字段。
- `backend/internal/service/openai_images_responses.go:413`：Responses 路径将尺寸传至 tools[0].size；`:1884` 说明直调端点返回 404/405 时会回退。
- `backend/internal/service/image_billing_size.go:53`：长边 1025～2048 归 2K，这是计费分类而非缩图设置。
- `backend/internal/service/openai_images_direct.go:226`：返回 model 会按客户端模型名称标记，因此响应中的 model 不能独立证明真实上游映射目标。

## 已执行验证

在 backend 目录执行现有定向测试并通过，未修改业务代码：

```bash
go test ./internal/service -run '^(TestCodexDirectImagesRouting|TestCodexDirectImagesMappingBeforeRouting|TestOpenAIGatewayServiceForwardImages_OAuthUsesDecodedOutputDimensions|TestOpenAIGatewayServiceParseOpenAIImagesRequest_NormalizesOfficialAndCustomSizes)$' -count=1
```

结果：`ok github.com/Wei-Shaw/sub2api/internal/service 5.028s`。该验证针对本地代码，不能代替线上部署验证。

## skill 的验收缺口

skill 仅报告保存路径和文件字节数，没有核对实际像素、格式与请求是否一致。此次 .jpeg/.webp 文件内部都是 PNG，仍报告 Saved。这是验收与反馈缺口，不是尺寸降低的原因。建议后续输出实际宽高、真实格式和不一致提示；不应通过本地放大冒充原生 4K。

## 最终归因需要的服务端证据

使用上述 request ID 查询线上命中账号类型、渠道/账号模型映射、上游端点、最终出站 size/quality/output_format 与原始响应：

- 出站参数变更或缺失：定位平台部署版本、转发路径或配置。
- 出站参数正确而上游返回约 1.57 MP PNG：定位该上游接口的参数支持，并在同一授权账号下对照直连结果。

本次没有线上管理会话或服务端日志；本地 docker ps 没有运行中的容器。当前仓库代码与线上平台实际版本是否一致尚未验证。

## 产物与状态

- 同目录 results.json：脱敏请求/响应元数据、真实像素、请求 ID、图片 SHA-256。
- 同目录 probe.py：诊断脚本；通过已安装 skill 获取专用密钥，不包含密钥。
- 同目录 models.json：该密钥可见的图片模型列表。
- 原始测试图片位于 `/tmp/shengtu-resolution-20261013/`，临时目录可能被系统清理。

本次未修改 skill、业务代码或线上配置。Trellis 任务保留 planning，记录诊断证据，未进入业务实现。

## 补充测试：1536×1024

用户要求追加一次 1536×1024 测试。保持 gpt-image-2.5-flare、相同提示词、low、PNG 和 n=1，只改变请求尺寸。

- 请求：1536×1024；文件真实像素：1254×1254；文件真实格式：PNG。
- 响应元数据同样为 size=1254x1254、quality=low、output_format=png。
- 耗时：24.83 秒；x-request-id：7b0f8765-cb68-4f65-ad21-d4342219fe84。
- 图片：`/tmp/shengtu-resolution-20261013/flare_1536x1024_png_low.png`。
- 这次请求也没有遵守尺寸，且 3:2 横图变成了正方形。因此现象不局限于请求 3840×2160 时出现；仍需服务端出站证据区分线上转发与上游行为。
