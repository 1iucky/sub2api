# 普通用户生图控制台实施计划

状态：用户已批准开始实施。用户已确认只清理到期图片，保留所有历史元数据；PRD、设计及本计划已收敛。旧自动 Key ensure/仓储锁和仅会话结果方案已撤销。

Goal：普通用户手动选 Key 或通过共用表单创建 Key，在统一控制台使用 OpenAI/Grok 与后续图片模型；复用原分组计费，提供持久历史、图片下载和管理员周期清理。

Architecture：JWT 读取配置/历史/文件；API Key 控制台专用提交进入原网关鉴权链并复用完整 Images 执行。服务端将权威请求/状态存入 PostgreSQL，图片存入受管数据目录；清理按稳定资产身份执行，不根据临时 URL 删除，不触碰账务。

Tech Stack：Go/Gin/Ent/PostgreSQL/Redis/Wire，Vue3/TypeScript/Tailwind/Pinia/vue-i18n/Vitest。契约以 design.md 为准，需求及验收以 prd.md 为准。

## 实施前检查
- [x] 到期只删除图片文件的产品决策已确认，已完成无损 PRD 收敛。
- [x] 用户在最新最终规划摘要之后明确批准实施，并补充降低 upstream/main 合并冲突要求。
- [x] validate 上下文清单后 task.py start；加载 trellis-before-dev 和对应层规范。规划期不运行产品测试。
- [x] 复查 git status，保留原 component-guidelines/GatewayDashboard/home tests/sub2api.tar；需要隔离时使用受管 worktree，明确起点包含所需本分支提交，不能丢失当前工作。
- [x] 原生成路由、gateway middleware、账务去重及公共同步/异步语义作为基线；无 ensure API、无自动 Key 创建、无前端提交权威图片历史。
- [x] 按 Trellis 实施/检查角色依赖顺序派发；共享 route/Wire/KeysView 文件串行修改。先读共享接口再派发，禁止覆盖其他工作。

## 1. 真实 API 密钥表单提取
Files：components/keys/ApiKeyFormDialog.vue、必要的 composables/useApiKeyForm.ts、views/user/KeysView.vue 及真实组件测试。

- [x] 提取现有 BaseDialog 的完整创建/编辑表单、字段、验证、平台选择、IP/配额/时间/窗口行为；原密钥页消费该组件。
- [x] 实现 show/editingKey/initialGroupId 和 close/created/updated 事件；创建返回真实 ApiKey DTO。
- [x] 在授权分组加载后预选对应平台与组，防 watcher 清空；当前组无权/不可用时明确报错。
- [x] 测试真实组件：原普通创建/编辑、延迟加载预选、Grok/OpenAI 平台、取消零创建、提交锁、后端拒绝、全部可选字段、表单改组后的真实归属。保留原 data-tour 引导。
- [x] 使用既有 keysAPI 与原后端创建限制；不增加 Key 创建服务/ensure 仓储。不要仅测试简化 fake 表单。

Validation（frontend）：定向 Vitest 原 KeysView.spec.ts 与新 ApiKeyFormDialog.spec.ts，pnpm typecheck；改动文件定向 ESLint（不 --fix 全库）。

## 2. 设置、分组发现与平台能力
Files：service/image_studio_settings.go、image_studio_providers.go、image_studio.go，JWT/admin handler、路由/依赖，能力及设置测试。

- [x] 实现 image_studio 设置：enabled、group_ids、default_group_id、cleanup 配置；缺省生成和清理关闭。
- [x] 保存校验授权、活跃/allow_image_generation、支持平台、默认属于列表、cron/时区/正保留天数；读取重新检查被删除或改动的组。
- [x] 服务端注册 OpenAI/GPT、Grok 能力，模型发现保留白名单与映射/账号能力语义；目录只装饰展示。无可调度账号不假报可用。
- [x] config 返回动态参数、超时/存储状态、默认不可用原因；GET 没有创建 Key 或生成副作用。
- [x] billing_preview 复用既有图片/token 定价分支解析；token billing 接口不能作为独立图片倍率。无法确定总价则不给虚假估价。
- [x] 表驱动测试：关闭/空/无效默认/重复组/不支持平台/未授权/GPT与Grok差异/映射白名单/价格模式。先写行为测试再实现。

Validation（backend）：go test -tags=unit ./internal/service ./internal/handler/... -run 'TestImageStudio(Settings|Providers|Config)'；go build ./cmd/server。

## 3. 持久历史与受管存储
Files：ent/schema/image_studio_generation.go、image_studio_asset.go，新 SQL 迁移、service/image_studio_history.go、image_studio_asset_store.go、image_result_reader.go、repository/image_studio_repo.go、image_studio_storage_local.go；定向调整 image_storage.go 共用解码。

- [x] 重查最大迁移后新增（目前最大前缀241，预期242，实施时重新确认）；不改已发布文件。业务 schema 与 SQL 同步，make generate 更新 Ent/Wire。
- [x] 历史/资产字段与索引落实 design 第6节；内部关联不对原用户/Key/组加级联删除边。唯一 user_id+submission_id 与 generation_id+index。
- [x] 实现按用户分页/过滤、详情、幂等 create/冲突、条件终态更新和资产登记/认领；查询使用参数化 SQL/Ent，限制分页与关键词长度。
- [x] Put/Open/Delete 本地适配器用数据目录+固定受管根；UUID键/原子文件/受限权限/路径及软链检查/流式读；缺文件幂等 Delete，不接受客户端 path/url。
- [x] 提取 reader：实际图片/MIME/字节限制/超时/安全出站 URL与redirect，不破坏旧 ImageResultUploader 回归。
- [x] 先登记 pending 再写文件，保存失败保留可追溯 key；测试部分保存、文件已写但数据库更新失败、原子 rename、非法路径/非图片/超限/断流；不持久化 base64、签名 URL或 Key明文。
- [x] 真实数据库验证两个用户筛选、唯一并发提交、删除Key仍能读历史、组改名显示快照、费用字段未知非零价。

Validation：go test -tags=unit ./internal/service ./internal/repository -run 'TestImageStudio(History|Asset|Storage)|TestImageResultReader|TestImageStorage'；按仓库真实数据库配置运行 integration 对应仓储用例；make generate 后 diff 复查（不无关重生成/格式化）。

## 4. 复用 Images 执行并捕获历史
Files：handler/image_task_handler.go、必要的共享 image runner、handler/image_studio_gateway_handler.go、routes/gateway.go、service/image_studio_history.go 与 handler/route 测试。

- [x] 在现有 gateway 中增加 POST /v1/images/generations/studio，保留同 body/APIKey/allowlist/组准入链。限制控制台启用、选定组和适配模型；执行上下文规范化到 /v1/images/generations 并保留鉴权、审计、原IP与计费关联。
- [x] 抽取现有 async 的请求验证、安全审计、复制上下文和平台完整 handler 执行；用 sink 区分旧Redis任务和studio持久化，不把studio结果重复上传两次。
- [x] 服务端解析 prompt/model/parameters 并生成 owner/Key/组快照；submission_id + 请求摘要幂等事务只派发一次。响应丢失可JWT按 submission_id 查回，客户端参数不能冒用 owner/费用。
- [x] 记录 correlation/usage_request_id，使用原 client:<id> Images 计费规则和 owner+Key 精确查日志。日志未落库显示待结算，只缓存原日志确认值，不新扣费。
- [x] 上游执行与资产保存状态分离；成功但保存失败/部分成功保留已计费事实，异常/panic可见；保存和落库重试不重复调用上游。
- [x] 执行租约/超时/进程恢复只标记unknown，不重放未知任务；长请求活跃心跳避免跨实例误终止。保留旧API 24h TTL/存储门禁/owner语义回归。
- [x] 约束 recorder 的整体大小/输出数量，不将大图无限堆内存；验证拒绝边界与原并发限制。

Validation：go test -tags=unit ./internal/handler ./internal/server/routes ./internal/server/middleware -run 'TestImageStudio|TestAsyncImage|TestGateway.*Image|TestGroupModelAllowlist'；go build ./cmd/server。

## 5. 历史读取、文件下载与清理服务
Files：handler/image_studio_handler.go、service/image_studio_cleanup.go、repository历史认领方法、admin handler、routes/user.go/admin.go、各层Wire、cmd/server/wire.go，以及清理/权限测试。

- [x] JWT list/detail/content 只允许本人历史/资产；跨用户/不存在404，deleted410，临时IO失败503，Content-Type/Length/Disposition正确。无APIKey/余额准入，不产生生成和计费。
- [x] 支持分页、时间/分组/模型/状态/提示词/submission_id；不跟当前Key或默认组绑定，关闭新生图仍能读旧作品。
- [x] 清理独立启停/保留天数/5字段cron/时区；Start/Stop/Reload并纳入系统生命周期。后台显示上次运行、成功/失败、删除字节数。
- [x] 各实例刷新有效设置，独立leader lock+资产租约条件认领；不持数据库事务做文件 IO；逐批删除/退避重试，成功或不存在才标deleted。
- [x] 清理仅本人资产表中的受管到期终态图片；活跃pending不删，未知/失败执行残留需租约过期和宽限。失败/租约到期/文件成功删除但DB失败等路径能恢复。
- [x] 清理只删除图片文件；保留完整 generation/asset 元数据、提示词/模型/参数/时间、费用状态和关联。不删除历史记录或 usage/账务。回归断言：清理前后历史记录数及元数据内容不变，只有资产清理状态和清理时间变化；已清理图片不可下载，其余信息仍可查询。
- [x] 集成验证并发认领、跨实例设置变更、时区与截止边界、缩短保留期、开关关闭、进行中任务保护、下载清理竞争、不误删其他路径。

Validation：unit service/handler -run 'TestImageStudio(Cleanup|History|Content|Settings)'；真实PostgreSQL -tags=integration ./internal/repository -run 'TestImageStudio.*(History|Claim|Cleanup)'；用临时文件根而非用户真实data做删除用例。

## 6. 用户工作台、历史与管理设置UI
Files：frontend/types/imageStudio.ts、api/imageStudio.ts、composables/useImageStudio.ts、views/user/ImageStudioView.vue、components/imageStudio/{ImageParameterControls,ImageResultCard,ImageHistoryList}.vue、components/admin/settings/ImageStudioSettings.vue、router/sidebar/SettingsView/i18n/notices及对应测试。

- [x] 分组旁快捷创建按钮消费第1步同一表单；获取完整group_id Key列表，有效Key才生成，创建改组不得串用Key，进入/刷新/切换组无创建API。
- [x] 服务端能力驱动参数/默认值，切模型丢弃旧字段；请求GPT/Grok特定参数通过adapter，不在page硬编码平台分支。
- [x] POST独立Key客户端、一次submission_id、防双击、无POST自动重试；JWT轮询详情和响应丢失查询，恢复processing，unknown明确提示，注销/卸载停止计时器。
- [x] 历史筛选/分页/详情、资产状态与费用待结算显示；JWT fetch图片Blob按需加载、逐张预览下载、释放object URLs，token不拼URL。清理后显示“图片已清理”，原提示词/参数/时间/费用继续可查询，逐资产控制预览和下载。
- [x] 管理员设置组件显示分组/默认组、保留期、执行周期、时区、开关和上次清理；用普通日期/周期控件表达设置，保存后提示对现有资产的影响。
- [x] 改造已固定MIT源码的小模块/交互，包含来源、commit、版权、完整许可；禁止只看过项目即声称集成，不引入整套React/Element Plus。
- [x] AppLayout、原控件/色彩/主题/响应式/中英词条；独立导航不借Gemini批量开关。单元覆盖模型切换、零自动Key创建、历史分页归属、异步恢复、Blob与MIME下载。

Validation：定向Vitest新useImageStudio/parameter/history/download/settings用例和真实Key表单；pnpm typecheck；pnpm run build（含i18n）；定向 ESLint。

## 7. 贯通验收、审查与交付
- [ ] mock上游 OpenAI/Grok 分别从真实APIKey路由进入完整网关：有效/无权限Key、模型白名单/映射、审计、IP/额度/并发、图片参数和资产结果。
- [ ] 原按张价格/独立图片倍率、token渠道/峰值倍率、余额/订阅各有费用记录证据；一次执行/使用日志/扣费，history GET、轮询及重复下载无额外费用。
- [ ] 两普通用户、Key删除/组改动/限额耗尽、刷新登录、提交响应丢失/同ID不同请求、未知任务/存储失败/部分成功均验证；不重生图。
- [ ] 数据库迁移/受管卷持久化/恢复/并发清理/失败重试用真实DB与隔离文件根；多实例共享卷条件写入运行说明。
- [ ] 浏览器明暗/移动宽度验收：选组→快捷Key表单→模型参数→提示词→生成→预览下载→刷新历史→到期展示。PNG/JPEG/WebP下载实字节，不只点击按钮。
- [x] trellis-check复查跨层、规范、生成代码、i18n、许可证及AC1–AC12；发现问题回实现者，共享文件串行。
- [x] research保存验证矩阵：静态/build、unit、真实DB、mock贯通、浏览器、真实上游/生产分别写状态。无凭据时如实记录真实出图/扣费未验证。
- [x] 只更新必要spec/journal并选择本任务diff；不自动部署/删历史目录。默认关生图和清理为回滚入口，保留新增表/文件、原Key与账务。

## 依赖与回滚
本期作为一个完整生图工作台交付，保留一个产品任务，各步骤有独立验证边界与显式依赖，不另建仅按技术层拆分的子任务。第1步与第2/3步有独立边界，但默认按步骤串行；第4依赖2/3，第5依赖2/3/4，第6依赖1–5，第7验收全链路。角色默认 trellis-implement/check，主会话承担共享契约、收敛与最终验收。

业务回滚关闭新增入口和清理；待已接收执行结束再回退代码，不做破坏性down，不删受管文件/Key/usage；公共Images同步/Redis异步与Gemini批量应保持兼容。

## 最终实施记录

步骤1–6代码完成并经过源代码复核；第7节的验证并非所有项目都代表生产端到端验收。已完成证据包括71项前端测试、最终生产构建、后端功能/旧async/生命周期与完整middleware回归、标准模式真实APIKeyAuth+mockHTTP+原Billing.Apply单次余额扣减、真实PostgreSQL历史/认领/费用快照集成。原订阅倍率命令回归通过，但Studio真实数据库订阅扣减未进行。浏览器PNG内容已请求并预览，IAB未返回下载落盘事件；HTTP内容与Blob/文件名组件测试通过。完整清理服务的新DB用例已通过，证明 leader互斥/释放、禁用保护和图片删除后完整元数据保留。详情与未验证边界以research/validation-matrix.md为准，不将构建或模拟上游称为生产出图/扣费验收。
