import {
  accountStatusOptions,
  bindingStatusOptions,
  responsibilityDomainOptions,
  strategyStatusOptions,
} from "@/features/follower-users/config/follower-user-content";
import type {
  FollowerUserAccountStatus,
  FollowerUserBindingStatus,
  FollowerUserFormState,
  FollowerUserStrategyStatus,
} from "@/features/follower-users/types";

/**
 * FollowerUserFormFieldsProps 描述新增或编辑表单共用字段区域的输入属性。
 */
export interface FollowerUserFormFieldsProps {
  /** formState 表示当前表单绑定的状态对象。 */
  formState: FollowerUserFormState;
  /** mode 表示当前字段用于新增还是编辑场景。 */
  mode: "create" | "edit";
  /** disabled 表示字段区域是否处于禁用状态。 */
  disabled?: boolean;
  /** onChange 用于同步任意单个字段的变更。 */
  onChange: <TKey extends keyof FollowerUserFormState>(
    key: TKey,
    value: FollowerUserFormState[TKey],
  ) => void;
}

// inputClassName 统一描述深色工作台中的输入框和下拉框样式。
const inputClassName =
  "h-11 w-full rounded-2xl border border-white/10 bg-white/[0.03] px-4 text-sm text-zinc-100 outline-none transition-colors placeholder:text-zinc-500 focus:border-cyan-300/30 focus:bg-white/[0.05]";

/**
 * FormFieldProps 描述单个表单字段容器的输入属性。
 */
interface FormFieldProps {
  /** label 表示字段标题。 */
  label: string;
  /** hint 用于展示字段说明。 */
  hint?: string;
  /** children 表示具体字段控件。 */
  children: React.ReactNode;
}

/**
 * FormField 负责渲染统一标题、提示和控件结构的字段容器。
 */
function FormField({ children, hint, label }: FormFieldProps): React.JSX.Element {
  return (
    <label className="grid gap-2">
      <span className="text-xs font-medium uppercase tracking-[0.2em] text-zinc-400">
        {label}
      </span>
      {children}
      {hint ? <span className="text-xs leading-5 text-zinc-500">{hint}</span> : null}
    </label>
  );
}

/**
 * FollowerUserFormFields 负责渲染新增与编辑弹窗复用的字段集合。
 */
export function FollowerUserFormFields({
  disabled = false,
  formState,
  mode,
  onChange,
}: FollowerUserFormFieldsProps): React.JSX.Element {
  return (
    <div className="grid gap-4">
      <FormField
        label="用户 ID"
        hint={
          mode === "create"
            ? "建议使用接口文档中的唯一用户标识格式，例如 `U-1301`。"
            : "编辑模式下用户 ID 不可修改。"
        }
      >
        <input
          className={inputClassName}
          value={formState.userId}
          disabled={disabled || mode === "edit"}
          onChange={(event) => onChange("userId", event.target.value)}
          placeholder="请输入用户 ID"
        />
      </FormField>

      <div className="grid gap-4 md:grid-cols-2">
        <FormField label="账户状态">
          <select
            className={inputClassName}
            value={formState.accountStatus}
            disabled={disabled}
            onChange={(event) =>
              onChange("accountStatus", event.target.value as FollowerUserAccountStatus)
            }
          >
            {accountStatusOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </FormField>

        <FormField label="策略状态">
          <select
            className={inputClassName}
            value={formState.strategyStatus}
            disabled={disabled}
            onChange={(event) =>
              onChange("strategyStatus", event.target.value as FollowerUserStrategyStatus)
            }
          >
            {strategyStatusOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </FormField>
      </div>

      <div className="grid gap-4 md:grid-cols-2">
        <FormField label="绑定状态">
          <select
            className={inputClassName}
            value={formState.bindingStatus}
            disabled={disabled}
            onChange={(event) =>
              onChange("bindingStatus", event.target.value as FollowerUserBindingStatus)
            }
          >
            {bindingStatusOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </FormField>

        <FormField label="责任域" hint="常见值包括 risk、ops、growth、liquidity、compliance。">
          <input
            className={inputClassName}
            list="follower-user-domain-options"
            value={formState.responsibilityDomain}
            disabled={disabled}
            onChange={(event) => onChange("responsibilityDomain", event.target.value)}
            placeholder="请输入责任域"
          />
          <datalist id="follower-user-domain-options">
            {responsibilityDomainOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </datalist>
        </FormField>
      </div>
    </div>
  );
}
