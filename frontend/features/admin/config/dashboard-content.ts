/**
 * AdminNavigationItem 描述后台壳层侧边栏的单个导航节点。
 */
export interface AdminNavigationItem {
  /** key 用于为导航项提供稳定的 React 键值。 */
  key: "overview" | "table-shell" | "detail-shell" | "form-shell" | "v0-playbook";
  /** label 表示导航项在界面上的中文标题。 */
  label: string;
  /** href 指定导航项跳转或锚点定位的目标地址。 */
  href: string;
  /** description 用于补充该导航项对应的能力说明。 */
  description: string;
  /** status 标识当前模块的完成程度，用于展示状态徽标。 */
  status: "ready" | "planned";
}

/**
 * AdminBreadcrumbItem 描述顶部面包屑中的单个层级节点。
 */
export interface AdminBreadcrumbItem {
  /** label 表示面包屑节点的展示名称。 */
  label: string;
  /** href 在节点可点击时提供目标地址。 */
  href?: string;
}

/**
 * AdminQuickAction 描述顶部工具栏中的快捷操作入口。
 */
export interface AdminQuickAction {
  /** label 表示按钮文案。 */
  label: string;
  /** href 表示点击后跳转的目标位置。 */
  href: string;
  /** tone 用于控制按钮在视觉上的强调程度。 */
  tone: "primary" | "secondary";
}

/**
 * AdminMetricCard 描述首页概览区中的一张指标卡片。
 */
export interface AdminMetricCard {
  /** title 表示指标卡片标题。 */
  title: string;
  /** value 表示最醒目的摘要数值或标签。 */
  value: string;
  /** description 用于解释该指标当前传达的能力含义。 */
  description: string;
  /** tone 用于映射卡片强调色。 */
  tone: "accent" | "success" | "warning";
}

/**
 * AdminFilterChip 描述列表筛选条中的一个静态筛选项。
 */
export interface AdminFilterChip {
  /** label 表示筛选项文案。 */
  label: string;
  /** active 表示该筛选项是否处于高亮状态。 */
  active: boolean;
}

/**
 * AdminTableColumn 描述表格骨架中的表头信息。
 */
export interface AdminTableColumn {
  /** key 用于为表头提供稳定键值。 */
  key: string;
  /** label 表示列标题。 */
  label: string;
}

/**
 * AdminTableRow 描述表格骨架中的一行展示数据。
 */
export interface AdminTableRow {
  /** key 用于为行提供稳定键值。 */
  key: string;
  /** cells 按列顺序存放当前行展示文本。 */
  cells: readonly string[];
  /** statusLabel 描述该行当前阶段或联调状态。 */
  statusLabel: string;
  /** statusTone 指定状态徽标的颜色语义。 */
  statusTone: "success" | "warning" | "muted";
}

/**
 * AdminDetailItem 描述详情区块中的单条键值信息。
 */
export interface AdminDetailItem {
  /** label 表示字段名。 */
  label: string;
  /** value 表示字段值。 */
  value: string;
  /** hint 在需要时补充说明单位、边界或联调说明。 */
  hint?: string;
}

/**
 * AdminDetailSection 描述详情页骨架中的一个信息分组。
 */
export interface AdminDetailSection {
  /** title 表示分组标题。 */
  title: string;
  /** description 用于概括该分组关注的主题。 */
  description: string;
  /** items 表示分组内展示的字段列表。 */
  items: readonly AdminDetailItem[];
}

/**
 * AdminFormField 描述表单骨架中的一个输入字段。
 */
export interface AdminFormField {
  /** label 表示字段中文名称。 */
  label: string;
  /** placeholder 表示输入框占位提示。 */
  placeholder: string;
  /** helperText 用于提示格式、默认值或业务限制。 */
  helperText: string;
  /** required 表示当前字段是否为必填。 */
  required: boolean;
  /** multiline 表示当前字段是否使用多行输入框展示。 */
  multiline?: boolean;
}

/**
 * AdminFormSection 描述表单预览中的一个业务分组。
 */
export interface AdminFormSection {
  /** title 表示该组字段的分组标题。 */
  title: string;
  /** description 用于补充该组字段承担的业务角色。 */
  description: string;
  /** fields 表示当前分组下的字段定义。 */
  fields: readonly AdminFormField[];
}

/**
 * AdminWorkflowStep 描述 `v0` 落地流程中的一个步骤。
 */
export interface AdminWorkflowStep {
  /** title 表示步骤名称。 */
  title: string;
  /** description 用于解释本步骤的主要目标。 */
  description: string;
  /** deliverable 说明该步骤结束时应该产出的结果。 */
  deliverable: string;
}

// adminNavigationItems 提供后台壳层的一级导航结构。
export const adminNavigationItems: readonly AdminNavigationItem[] = [
  {
    key: "overview",
    label: "总览",
    href: "/",
    description: "查看后台骨架、代理入口与当前准备度。",
    status: "ready",
  },
  {
    key: "table-shell",
    label: "列表骨架",
    href: "#table-shell",
    description: "沉淀筛选条、表格布局与状态列模板。",
    status: "ready",
  },
  {
    key: "detail-shell",
    label: "详情骨架",
    href: "#detail-shell",
    description: "统一信息分区、摘要字段与说明文案层级。",
    status: "ready",
  },
  {
    key: "form-shell",
    label: "表单骨架",
    href: "#form-shell",
    description: "为创建、编辑和停用流程预留复用面板。",
    status: "ready",
  },
  {
    key: "v0-playbook",
    label: "v0 工作流",
    href: "#v0-playbook",
    description: "把 PM 文档拆解成 prompt、组件与联调节奏。",
    status: "ready",
  },
];

// adminBreadcrumbItems 描述首页顶部的层级路径。
export const adminBreadcrumbItems: readonly AdminBreadcrumbItem[] = [
  {
    label: "运营后台",
    href: "/",
  },
  {
    label: "前端工作台",
  },
];

// adminQuickActions 提供顶部工具栏的锚点与联调入口。
export const adminQuickActions: readonly AdminQuickAction[] = [
  {
    label: "查看代理健康检查",
    href: "/api/admin/healthz",
    tone: "primary",
  },
  {
    label: "跳到 v0 工作流",
    href: "#v0-playbook",
    tone: "secondary",
  },
  {
    label: "查看登录挂载位",
    href: "#auth-slot",
    tone: "secondary",
  },
];

// adminMetricCards 展示当前骨架阶段已经准备好的能力摘要。
export const adminMetricCards: readonly AdminMetricCard[] = [
  {
    title: "通用页面骨架",
    value: "8 类",
    description: "已覆盖导航、面包屑、列表、详情、表单与三种状态组件。",
    tone: "accent",
  },
  {
    title: "后台代理入口",
    value: "/api/admin/*",
    description: "统一由 Route Handler 代为转发至后端 `/admin/v1/*`。",
    tone: "success",
  },
  {
    title: "v0 落地节奏",
    value: "4 步",
    description: "先生成初稿，再做组件拆分、接口接入和验收修正。",
    tone: "warning",
  },
];

// adminFilterChips 模拟列表页常见的筛选维度，用于展示可复用筛选条。
export const adminFilterChips: readonly AdminFilterChip[] = [
  {
    label: "仅看已启用",
    active: true,
  },
  {
    label: "待联调模块",
    active: false,
  },
  {
    label: "可复用字段",
    active: false,
  },
  {
    label: "显示错误态",
    active: false,
  },
];

// adminTableColumns 描述后台列表页骨架的表头布局。
export const adminTableColumns: readonly AdminTableColumn[] = [
  {
    key: "module",
    label: "模块名称",
  },
  {
    key: "route",
    label: "路由意图",
  },
  {
    key: "data-source",
    label: "数据接入",
  },
  {
    key: "interaction",
    label: "交互重点",
  },
];

// adminTableRows 提供三类典型后台页面的示例行数据。
export const adminTableRows: readonly AdminTableRow[] = [
  {
    key: "list",
    cells: ["列表页", "筛选 + 分页 + 状态列", "走 `/api/admin/*` 代理", "批量筛选、默认排序"],
    statusLabel: "已准备",
    statusTone: "success",
  },
  {
    key: "detail",
    cells: ["详情页", "区块化信息摘要", "支持真实接口或 mock contract", "字段分区、说明提示"],
    statusLabel: "可复用",
    statusTone: "muted",
  },
  {
    key: "form",
    cells: ["创建/编辑页", "表单分组 + 校验提示", "优先映射 PM 字段与接口契约", "提交态、禁用态、回显态"],
    statusLabel: "待联调",
    statusTone: "warning",
  },
];

// adminDetailSections 展示详情页骨架建议呈现的信息层次。
export const adminDetailSections: readonly AdminDetailSection[] = [
  {
    title: "基础信息",
    description: "用于承载对象最重要的识别字段与当前状态。",
    items: [
      {
        label: "页面目标",
        value: "先展示稳定结构，再等待 PM 文档映射真实字段。",
      },
      {
        label: "默认语言",
        value: "中文",
        hint: "所有说明文案、注释和模板均保持中文。",
      },
      {
        label: "登录挂载位",
        value: "已预留",
        hint: "本阶段不实现真实鉴权，仅保留入口位置。",
      },
    ],
  },
  {
    title: "代理与联调",
    description: "用于说明前端和后端之间的固定协作边界。",
    items: [
      {
        label: "统一入口",
        value: "/api/admin/*",
      },
      {
        label: "默认后端",
        value: "http://127.0.0.1:8080",
      },
      {
        label: "异常策略",
        value: "后端不可达时返回 502 JSON，后端业务错误按原状态透传。",
      },
    ],
  },
];

// adminFormSections 描述创建或编辑面板推荐采用的字段编排方式。
export const adminFormSections: readonly AdminFormSection[] = [
  {
    title: "基础字段",
    description: "放置最常用、最影响业务识别的输入项。",
    fields: [
      {
        label: "业务对象名称",
        placeholder: "请输入用于运营识别的中文名称",
        helperText: "建议与详情页标题保持一致，便于后续列表检索。",
        required: true,
      },
      {
        label: "状态说明",
        placeholder: "请选择启用、停用或待配置等业务状态",
        helperText: "状态字段后续直接映射后端枚举定义。",
        required: true,
      },
    ],
  },
  {
    title: "扩展说明",
    description: "用于承载描述性文本、附加约束与联调备注。",
    fields: [
      {
        label: "联调备注",
        placeholder: "记录 PM 约束、接口兼容说明或默认值来源",
        helperText: "建议把边界条件和默认值写清楚，避免后续返工。",
        required: false,
        multiline: true,
      },
    ],
  },
];

// adminWorkflowSteps 固定 `v0` 从需求到代码入库的标准动作顺序。
export const adminWorkflowSteps: readonly AdminWorkflowStep[] = [
  {
    title: "拆解 PM 文档",
    description: "把需求整理成页面目标、字段清单、状态流、验收标准。",
    deliverable: "页面清单、字段表、交互状态矩阵。",
  },
  {
    title: "生成 v0 初稿",
    description: "使用中文 prompt 指定后台风格、信息层级和重点操作区。",
    deliverable: "视觉稿链接、初版组件结构或页面代码。",
  },
  {
    title: "收敛到仓库结构",
    description: "拆分共享组件、业务组件与 `features/*` 数据接入层。",
    deliverable: "可复用组件、代理调用层、类型定义与中文注释。",
  },
  {
    title: "完成联调验收",
    description: "接入真实接口并补齐加载态、空态、错误态和回归检查。",
    deliverable: "可联调页面、验证记录与后续优化清单。",
  },
];
