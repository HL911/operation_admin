import type {
  AdminBreadcrumbItem,
  AdminNavigationItem,
  AdminQuickAction,
} from "@/features/admin/config/dashboard-content";
import type {
  FollowerUserAccountStatus,
  FollowerUserBindingStatus,
  FollowerUserFormState,
  FollowerUserListData,
  FollowerUserListQuery,
  FollowerUserRecord,
  FollowerUserStrategyStatus,
} from "@/features/follower-users/types";

/**
 * FollowerUserTone 描述小龙虾用户状态徽标使用的视觉语义。
 */
export type FollowerUserTone = "accent" | "success" | "warning" | "muted";

/**
 * FollowerUserStatusMeta 描述单个状态值在界面上的标签与颜色。
 */
export interface FollowerUserStatusMeta {
  /** label 表示状态值对应的中文展示名。 */
  label: string;
  /** tone 表示该状态应使用的徽标颜色语义。 */
  tone: FollowerUserTone;
}

/**
 * FollowerUserOption 描述表单或筛选面板中的一个选项。
 */
export interface FollowerUserOption<TValue extends string> {
  /** value 表示真正写入接口或筛选条件的值。 */
  value: TValue;
  /** label 表示在界面上展示的中文选项名。 */
  label: string;
}

/**
 * FollowerUsersConsoleNavKey 描述运营后台侧边栏菜单的唯一标识。
 */
export type FollowerUsersConsoleNavKey =
  | "little-follower-query"
  | "big-follower-query"
  | "big-follower-key-query";

/**
 * FollowerUsersConsoleNavItem 描述运营后台侧边栏中的一个一级菜单项。
 */
export interface FollowerUsersConsoleNavItem {
  /** key 表示菜单唯一标识。 */
  key: FollowerUsersConsoleNavKey;
  /** label 表示菜单展示名称。 */
  label: string;
  /** active 表示当前菜单是否为激活态。 */
  active: boolean;
}

/**
 * FollowerUsersConsoleLanguageOption 描述顶部语言切换器的选项。
 */
export interface FollowerUsersConsoleLanguageOption {
  /** key 表示语言值。 */
  key: "zh" | "en";
  /** label 表示语言显示文本。 */
  label: string;
}

/**
 * FollowerUsersConsoleUserAction 描述用户头像下拉菜单中的一个动作项。
 */
export interface FollowerUsersConsoleUserAction {
  /** key 表示动作唯一标识。 */
  key: "profile" | "settings" | "logout";
  /** label 表示动作展示文本。 */
  label: string;
  /** tone 表示动作使用的视觉语义。 */
  tone: "default" | "danger";
}

// followerUsersConsoleProductName 描述顶部导航栏左侧显示的产品名称。
export const followerUsersConsoleProductName = "运营后台";

// followerUsersConsolePageTitle 描述当前页面在顶栏和内容区中的主标题。
export const followerUsersConsolePageTitle = "小龙虾用户查询";

// followerUsersConsoleNavItems 描述本次设计稿中固定展示的三个一级菜单。
export const followerUsersConsoleNavItems: readonly FollowerUsersConsoleNavItem[] = [
  {
    key: "little-follower-query",
    label: "小龙虾用户查询",
    active: true,
  },
  {
    key: "big-follower-query",
    label: "大龙虾用户查询",
    active: false,
  },
  {
    key: "big-follower-key-query",
    label: "大龙虾 Key 查询",
    active: false,
  },
];

// followerUsersConsoleLanguageOptions 描述顶部语言切换器支持的语言项。
export const followerUsersConsoleLanguageOptions: readonly FollowerUsersConsoleLanguageOption[] = [
  {
    key: "zh",
    label: "中文",
  },
  {
    key: "en",
    label: "EN",
  },
];

// followerUsersConsoleUserActions 描述头像菜单中的个人操作项。
export const followerUsersConsoleUserActions: readonly FollowerUsersConsoleUserAction[] = [
  {
    key: "profile",
    label: "个人资料",
    tone: "default",
  },
  {
    key: "settings",
    label: "账户设置",
    tone: "default",
  },
  {
    key: "logout",
    label: "退出登录",
    tone: "danger",
  },
];

// followerUsersNavigationItems 描述小龙虾用户页面使用的侧边栏导航。
export const followerUsersNavigationItems: readonly AdminNavigationItem[] = [
  {
    key: "overview",
    label: "总览",
    href: "/",
    description: "返回后台前端骨架总览与通用工作台。",
    status: "ready",
  },
  {
    key: "follower-users",
    label: "小龙虾用户",
    href: "/follower-users",
    description: "根据接口文档管理独立表中的小龙虾用户 CRUD。",
    status: "ready",
  },
];

// followerUsersBreadcrumbItems 描述小龙虾用户页面顶部的路径层级。
export const followerUsersBreadcrumbItems: readonly AdminBreadcrumbItem[] = [
  {
    label: "运营后台",
    href: "/",
  },
  {
    label: "小龙虾用户",
  },
];

// followerUsersQuickActions 描述顶部工具栏提供的快捷入口。
export const followerUsersQuickActions: readonly AdminQuickAction[] = [
  {
    label: "回到总览",
    href: "/",
    tone: "primary",
  },
  {
    label: "接口文档",
    href: "/docs/operator-follower-user-custom-table-crud.md",
    tone: "secondary",
  },
  {
    label: "v0 工作流",
    href: "/docs/frontend-v0-workflow.md",
    tone: "secondary",
  },
];

// accountStatusOptions 描述账户状态筛选和表单可用的选项。
export const accountStatusOptions: readonly FollowerUserOption<FollowerUserAccountStatus>[] = [
  {
    value: "active",
    label: "启用",
  },
  {
    value: "disabled",
    label: "停用",
  },
];

// strategyStatusOptions 描述策略状态筛选和表单可用的选项。
export const strategyStatusOptions: readonly FollowerUserOption<FollowerUserStrategyStatus>[] = [
  {
    value: "enabled",
    label: "启用",
  },
  {
    value: "disabled",
    label: "停用",
  },
];

// bindingStatusOptions 描述绑定状态筛选和表单可用的选项。
export const bindingStatusOptions: readonly FollowerUserOption<FollowerUserBindingStatus>[] = [
  {
    value: "pending",
    label: "待绑定",
  },
  {
    value: "bound",
    label: "已绑定",
  },
  {
    value: "unbound",
    label: "未绑定",
  },
];

// responsibilityDomainOptions 描述责任域字段的常见选项，供筛选与输入提示使用。
export const responsibilityDomainOptions: readonly FollowerUserOption<string>[] = [
  {
    value: "risk",
    label: "风险控制",
  },
  {
    value: "ops",
    label: "运营执行",
  },
  {
    value: "growth",
    label: "增长投放",
  },
  {
    value: "liquidity",
    label: "流动性管理",
  },
  {
    value: "compliance",
    label: "合规保障",
  },
];

// accountStatusMetaMap 用于把账户状态转换成中文标签和色彩语义。
export const accountStatusMetaMap: Record<FollowerUserAccountStatus, FollowerUserStatusMeta> = {
  active: {
    label: "启用",
    tone: "accent",
  },
  disabled: {
    label: "停用",
    tone: "muted",
  },
};

// strategyStatusMetaMap 用于把策略状态转换成中文标签和色彩语义。
export const strategyStatusMetaMap: Record<FollowerUserStrategyStatus, FollowerUserStatusMeta> = {
  enabled: {
    label: "已启用",
    tone: "accent",
  },
  disabled: {
    label: "已停用",
    tone: "muted",
  },
};

// bindingStatusMetaMap 用于把绑定状态转换成中文标签和色彩语义。
export const bindingStatusMetaMap: Record<FollowerUserBindingStatus, FollowerUserStatusMeta> = {
  pending: {
    label: "待绑定",
    tone: "warning",
  },
  bound: {
    label: "已绑定",
    tone: "success",
  },
  unbound: {
    label: "未绑定",
    tone: "muted",
  },
};

// defaultFollowerUserFilters 描述列表页首次进入时使用的默认筛选条件。
export const defaultFollowerUserFilters: FollowerUserListQuery = {
  pageNum: 1,
  pageSize: 10,
  userId: "",
  accountStatus: "",
  strategyStatus: "",
  bindingStatus: "",
  responsibilityDomain: "",
  updatedFrom: "",
  updatedTo: "",
};

// defaultFollowerUserCreateForm 描述新建表单的默认值。
export const defaultFollowerUserCreateForm: FollowerUserFormState = {
  userId: "",
  accountStatus: "active",
  strategyStatus: "disabled",
  bindingStatus: "unbound",
  responsibilityDomain: "risk",
  rowVersion: 1,
};

// followerUsersPreviewRecords 描述页面设计预览时使用的示例数据。
export const followerUsersPreviewRecords: readonly FollowerUserRecord[] = [
  {
    userId: "U-1201",
    accountStatus: "active",
    strategyStatus: "enabled",
    bindingStatus: "bound",
    responsibilityDomain: "risk",
    updatedAt: "2026-04-01 18:42:12",
    rowVersion: 6,
  },
  {
    userId: "U-1208",
    accountStatus: "active",
    strategyStatus: "enabled",
    bindingStatus: "pending",
    responsibilityDomain: "growth",
    updatedAt: "2026-04-01 17:58:09",
    rowVersion: 3,
  },
  {
    userId: "U-1220",
    accountStatus: "disabled",
    strategyStatus: "disabled",
    bindingStatus: "unbound",
    responsibilityDomain: "ops",
    updatedAt: "2026-04-01 16:32:44",
    rowVersion: 2,
  },
  {
    userId: "U-1237",
    accountStatus: "active",
    strategyStatus: "disabled",
    bindingStatus: "pending",
    responsibilityDomain: "compliance",
    updatedAt: "2026-04-01 15:06:31",
    rowVersion: 4,
  },
  {
    userId: "U-1265",
    accountStatus: "active",
    strategyStatus: "enabled",
    bindingStatus: "bound",
    responsibilityDomain: "liquidity",
    updatedAt: "2026-04-01 13:21:07",
    rowVersion: 8,
  },
  {
    userId: "U-1293",
    accountStatus: "disabled",
    strategyStatus: "disabled",
    bindingStatus: "unbound",
    responsibilityDomain: "risk",
    updatedAt: "2026-04-01 12:18:55",
    rowVersion: 1,
  },
];

// followerUsersPreviewListData 描述页面在无实时后端时展示的默认分页结果。
export const followerUsersPreviewListData: FollowerUserListData = {
  list: followerUsersPreviewRecords,
  total: followerUsersPreviewRecords.length,
  pageNum: 1,
  pageSize: 10,
  pages: 1,
};
