# 目录结构

> 前端代码统一放在 `frontend/`，以“路由层、组件层、特性层、工具层”四层组织。

---

## 当前目录布局

```text
frontend/
├── app/                         # App Router 页面、全局布局、Route Handler
│   └── api/admin/[[...segments]]/route.ts
├── components/
│   ├── admin/                   # 后台壳层、列表骨架、详情骨架、表单骨架
│   └── ui/                      # Button、Panel、StatusBadge 等通用组件
├── features/
│   └── admin/
│       ├── config/              # 静态配置、页面骨架数据、导航定义
│       └── lib/                 # 代理、数据转换、后续接口能力
├── lib/                         # 与业务无关的通用工具，例如 cn()
├── styles/                      # 主题变量、全局动效、设计 token
├── public/                      # 静态资源
└── package.json                 # 前端工程脚本与依赖
```

## 模块边界

- `app/` 只负责路由、布局、页面编排和 Route Handler，不承载复杂业务常量
- `components/ui/` 存放跨页面通用的基础组件，不绑定具体业务模块
- `components/admin/` 存放后台管理场景专用骨架组件，例如侧栏、表格骨架、详情面板
- `features/*/config` 存放页面静态配置、提示词模板、导航数据和 mock contract
- `features/*/lib` 存放代理、数据适配、请求封装、格式转换等能力
- `lib/` 仅放与单一业务无关的纯工具函数

## 新功能落位规则

- 新页面先放在 `app/`，如果只是已有骨架的组合，不要新增一批重复组件
- 新增可复用模块时，优先抽到 `components/ui` 或 `components/admin`
- 新增与某个业务域强绑定的配置和数据处理，放到对应的 `features/<domain>/`
- 新增后端访问统一从 `features/<domain>/lib` 发起，再由页面或服务器组件调用

## 命名约定

- 组件文件使用 `kebab-case`，导出组件使用 `PascalCase`
- 配置文件使用 `kebab-case`，并在文件内导出带语义前缀的常量
- Route Handler 统一使用 `route.ts`
- 仅用于组织代码、不参与路由的目录可以使用私有目录或普通目录，不新增无意义层级
