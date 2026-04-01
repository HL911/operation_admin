# 组件规范

> 当前前端组件分成“基础 UI 组件”和“后台场景组件”两层，不直接在页面里堆重复结构。

---

## 组件分层

- `components/ui`：通用基础组件，例如 `Button`、`Panel`、`StatusBadge`
- `components/admin`：后台场景骨架组件，例如 `AdminShell`、`DataTableShell`、`FormPreviewPanel`
- 页面文件负责组合组件，不负责定义一整套重复样式

## Props 约定

- 每个组件都必须显式定义 Props 接口，并为字段补中文注释
- 传入文本时优先传业务语义字段，不传含糊的 `data`、`item`
- 可选样式扩展统一使用 `className`
- 组件只接收自己真正需要的最小输入，不把整个大对象原样下传

## 组合方式

- 用 `Panel` 包住后台模块区块，保持标题区、描述区和内容区一致
- 用 `StatusBadge` 展示状态，不在页面里手写多套状态颜色
- 列表、详情、表单都应先抽成骨架组件，再由页面组合
- 当 `v0` 生成大块页面代码时，先拆成 `components/ui` 与 `components/admin`，再接入页面

## 样式模式

- 统一使用 Tailwind 类名和 `styles/tokens.css` 中的主题变量
- 需要复用的视觉基线放到 token，不在页面里散落硬编码色值
- 避免大段复制类名；如果结构复用明显，就抽组件
- 当前已接入 `shadcn/ui` 配置，新增基础组件优先与 `components.json` 的别名体系保持一致

## 无障碍

- 可交互元素优先使用原生 `button`、`a`、`input`，不要用 `div` 伪装
- 需要时补 `aria-label`、`aria-hidden` 和 `aria-*` 语义
- 组件的占位和错误提示要尽量让屏幕阅读器能理解

## 禁止模式

- 页面文件里直接复制大段卡片、表格、详情结构而不抽组件
- 组件里直接请求后端地址或读取环境变量
- 组件 props 使用 `any`
- 只为“让代码通过”而添加无信息量注释
