# 修复站点品牌与首页登录页国际化展示

## Goal

修复当前分支合并后出现的品牌名称、版本号和国际化键直接暴露问题，使控制台、登录页及公共首页相关区域统一遵循系统通用设置中的站点名称，并在未配置时回退到原项目名称 `Sub2API`。

## What I already know

* 控制台左上角品牌名称由 `AppSidebar.vue` 的 `appStore.siteName` 提供，但品牌区仍展示 `VersionBadge`。
* 登录页品牌区域由 `AuthLayout.vue` 提供，目前会针对 `SiliconBase` 隐藏站点名称，且默认回退仍为 `SiliconBase`。
* 首页及公共导航中存在硬编码 `SiliconBase`，没有使用已经计算出的 `siteName`。
* `home.nav.status`、`home.nav.models` 以及页脚法律链接使用的多个键未在中英文 `landing.ts` 中定义，Vue-i18n 因此直接展示键名。
* 登录页统一条款复用页脚法律链接键，因此英文切换时同样暴露缺失键。
* 用户已确认：站点名称为空时统一回退为 `Sub2API`。

## Requirements

* 控制台左上角显示系统通用设置配置的站点名称；未配置时显示 `Sub2API`。
* 控制台左上角品牌区域不再展示版本号。
* 登录页 Logo 下始终显示系统通用设置配置的站点名称；未配置时显示 `Sub2API`。
* 首页与公共导航中的品牌名称使用同一站点名称规则，不再硬编码 `SiliconBase`。
* 首页桌面端和移动端“状态”“模型集市”菜单支持中文与英文。
* 首页页脚法律链接支持中文与英文。
* 登录/注册统一条款中的文档名称跟随当前语言，英文下不再显示 i18n 键或中文回退（存在对应英文键时）。

## Acceptance Criteria

* [ ] 站点名称配置为自定义值时，控制台和登录页均显示该值。
* [ ] 站点名称为空时，相关品牌位置显示 `Sub2API`。
* [ ] 控制台左上角不显示 `vX.Y.Z` 版本号。
* [ ] 中文首页菜单显示“状态”“模型集市”，英文显示“Status”“Model Marketplace”。
* [ ] 中英文页脚法律链接均显示可读文案，不显示 `home.*` 键名。
* [ ] 登录页统一条款在英文模式下显示英文文档名称。
* [ ] 相关前端测试、类型检查和 lint 通过。

## Definition of Done

* 更新受影响组件与中英文 locale。
* 增补或更新针对品牌回退、版本隐藏和缺失国际化键的测试。
* 前端定向测试、类型检查和 lint 通过。
* 不覆盖当前工作区中与本任务无关的 `deploy/docker-compose.yml` 修改。

## Technical Approach

* 继续复用现有公开设置与 `appStore.siteName` 数据流，不新增后端字段或新的全局状态。
* 在品牌展示组件中统一采用“已配置站点名称，否则 `Sub2API`”规则，并替换首页与公共导航中的硬编码名称。
* 从控制台侧栏品牌区移除 `VersionBadge`，但保留后台已有的版本更新、回滚 API 与组件实现。
* 在中英文 `landing.ts` 中补齐导航及法律链接键；登录统一条款继续通过 `loginAgreement.ts` 的现有映射复用这些键。
* 使用组件/静态契约测试覆盖品牌回退、版本隐藏和 locale 键完整性。

## Decision (ADR-lite)

**Context**: 当前问题由合并后的硬编码品牌、侧栏版本组件和缺失 locale 键共同造成。可以选择局部修补、全局品牌工具重构或只补翻译键。

**Decision**: 采用最小范围的组件级统一修复。复用现有 Store 数据源，不引入新的品牌解析抽象；修复所有本次受影响的共享组件和首页内联实现。

**Consequences**: 改动集中且回归风险较低，可立即恢复正确展示；全仓库其他历史 `SiliconBase` 文案不在本任务中清理，未来如需品牌体系重构可单独立项。

## Out of Scope

* 不修改站点 Logo 上传、存储或 URL 清洗逻辑。
* 不重写首页视觉样式或页面结构。
* 不调整后台版本升级、回滚能力，只移除侧栏品牌区域的版本展示。
* 不做全仓库所有 `SiliconBase` 文案的品牌清理。

## Technical Notes

* 主要文件：`frontend/src/components/layout/AppSidebar.vue`、`frontend/src/components/layout/AuthLayout.vue`、`frontend/src/views/HomeView.vue`、`frontend/src/components/home/PublicTopNav.vue`、`frontend/src/components/home/PublicFooter.vue`。
* 国际化文件：`frontend/src/i18n/locales/zh/landing.ts`、`frontend/src/i18n/locales/en/landing.ts`。
* 统一条款标题映射：`frontend/src/utils/loginAgreement.ts`。
* 当前工作区已有用户修改：`deploy/docker-compose.yml`，本任务不得触碰。
