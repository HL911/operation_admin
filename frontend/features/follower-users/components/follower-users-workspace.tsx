"use client";

import { useEffect, useState } from "react";
import {
  ChevronLeft,
  ChevronRight,
  LoaderCircle,
  PencilLine,
  Plus,
  RefreshCcw,
  Search,
  Trash2,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Dialog } from "@/components/ui/dialog";
import { FollowerUserFormFields } from "@/features/follower-users/components/follower-user-form-fields";
import { FollowerUserStatusPill } from "@/features/follower-users/components/follower-user-status-pill";
import {
  accountStatusMetaMap,
  defaultFollowerUserCreateForm,
  defaultFollowerUserFilters,
  followerUsersConsolePageTitle,
  followerUsersPreviewRecords,
  responsibilityDomainOptions,
  strategyStatusMetaMap,
  bindingStatusMetaMap,
} from "@/features/follower-users/config/follower-user-content";
import {
  createFollowerUser,
  deleteFollowerUser,
  fetchFollowerUserDetail,
  fetchFollowerUserList,
  updateFollowerUser,
} from "@/features/follower-users/lib/follower-user-api";
import type {
  FollowerUserAccountStatus,
  FollowerUserBindingStatus,
  FollowerUserCreatePayload,
  FollowerUserFormState,
  FollowerUserListData,
  FollowerUserListQuery,
  FollowerUserRecord,
  FollowerUserStrategyStatus,
  FollowerUserUpdatePayload,
} from "@/features/follower-users/types";
import { cn } from "@/lib/utils";

/**
 * ModalMode 描述当前页面可能打开的模态框类型。
 */
type ModalMode = "create" | "edit" | "delete";

/**
 * NoticeState 描述页面顶部反馈条的状态。
 */
interface NoticeState {
  /** tone 表示反馈条对应的视觉语义。 */
  tone: "success" | "danger";
  /** message 表示要展示给用户的反馈文案。 */
  message: string;
}

/**
 * FilterFieldProps 描述筛选字段容器需要的输入属性。
 */
interface FilterFieldProps {
  /** label 表示字段标题。 */
  label: string;
  /** children 表示要渲染的输入控件。 */
  children: React.ReactNode;
}

// filterControlClassName 统一描述搜索框和筛选控件的基础样式。
const filterControlClassName =
  "h-11 w-full rounded-2xl border border-white/10 bg-white/[0.03] px-4 text-sm text-zinc-100 outline-none transition-colors placeholder:text-zinc-500 focus:border-cyan-300/30 focus:bg-white/[0.05]";

// subtleButtonClassName 统一描述筛选区和表格中的轻量按钮样式。
const subtleButtonClassName =
  "!rounded-2xl !border-white/10 !bg-white/[0.03] !text-zinc-200 hover:!bg-white/[0.08]";

// tableHeaders 描述表格中需要展示的列标题。
const tableHeaders = ["用户ID", "账户状态", "策略状态", "绑定状态", "责任域", "更新时间", "操作"];

/**
 * normalizeRecord 负责补齐记录对象中缺失的 rowVersion。
 */
function normalizeRecord(record: FollowerUserRecord): FollowerUserRecord {
  return {
    ...record,
    rowVersion: record.rowVersion ?? 1,
  };
}

/**
 * mapRecordToFormState 负责把记录对象转换成表单需要的结构。
 */
function mapRecordToFormState(record: FollowerUserRecord): FollowerUserFormState {
  return {
    userId: record.userId,
    accountStatus: record.accountStatus,
    strategyStatus: record.strategyStatus,
    bindingStatus: record.bindingStatus,
    responsibilityDomain: record.responsibilityDomain,
    rowVersion: record.rowVersion ?? 1,
  };
}

/**
 * buildCreatePayload 负责从表单状态提取创建接口所需的请求体。
 */
function buildCreatePayload(formState: FollowerUserFormState): FollowerUserCreatePayload {
  return {
    userId: formState.userId.trim(),
    accountStatus: formState.accountStatus,
    strategyStatus: formState.strategyStatus,
    bindingStatus: formState.bindingStatus,
    responsibilityDomain: formState.responsibilityDomain.trim(),
  };
}

/**
 * buildUpdatePayload 负责从表单状态提取更新接口所需的请求体。
 */
function buildUpdatePayload(formState: FollowerUserFormState): FollowerUserUpdatePayload {
  return {
    accountStatus: formState.accountStatus,
    strategyStatus: formState.strategyStatus,
    bindingStatus: formState.bindingStatus,
    responsibilityDomain: formState.responsibilityDomain.trim(),
    rowVersion: formState.rowVersion,
  };
}

/**
 * validateFollowerUserForm 负责校验新增或编辑表单中的必填字段。
 */
function validateFollowerUserForm(
  formState: FollowerUserFormState,
  options: {
    /** includeUserId 表示当前校验是否要求必须填写用户 ID。 */
    includeUserId: boolean;
  },
): string | null {
  if (options.includeUserId && !formState.userId.trim()) {
    return "用户 ID 不能为空。";
  }

  if (!formState.responsibilityDomain.trim()) {
    return "责任域不能为空。";
  }

  if (formState.rowVersion < 1) {
    return "版本号必须大于等于 1。";
  }

  return null;
}

/**
 * matchesPreviewFilters 负责判断示例数据是否命中当前筛选条件。
 */
function matchesPreviewFilters(
  record: FollowerUserRecord,
  query: FollowerUserListQuery,
): boolean {
  // normalizedUserIdKeyword 用于统一处理用户 ID 搜索词。
  const normalizedUserIdKeyword = query.userId.trim().toLowerCase();
  // normalizedDomainKeyword 用于统一处理责任域搜索词。
  const normalizedDomainKeyword = query.responsibilityDomain.trim().toLowerCase();

  if (normalizedUserIdKeyword && !record.userId.toLowerCase().includes(normalizedUserIdKeyword)) {
    return false;
  }

  if (query.accountStatus && record.accountStatus !== query.accountStatus) {
    return false;
  }

  if (query.strategyStatus && record.strategyStatus !== query.strategyStatus) {
    return false;
  }

  if (query.bindingStatus && record.bindingStatus !== query.bindingStatus) {
    return false;
  }

  if (
    normalizedDomainKeyword &&
    !record.responsibilityDomain.toLowerCase().includes(normalizedDomainKeyword)
  ) {
    return false;
  }

  return true;
}

/**
 * buildPreviewListData 负责根据当前筛选条件构造假数据模式下的分页结果。
 */
function buildPreviewListData(
  records: readonly FollowerUserRecord[],
  query: FollowerUserListQuery,
): FollowerUserListData {
  // matchedRecords 用于保存命中筛选条件的假数据记录。
  const matchedRecords = records.filter((record) => matchesPreviewFilters(record, query));
  // pages 表示假数据结果对应的总页数。
  const pages = Math.max(1, Math.ceil(matchedRecords.length / query.pageSize));
  // safePageNum 用于把当前页码约束在有效区间内。
  const safePageNum = Math.min(query.pageNum, pages);
  // startIndex 用于计算当前页起始索引。
  const startIndex = (safePageNum - 1) * query.pageSize;

  return {
    list: matchedRecords.slice(startIndex, startIndex + query.pageSize),
    total: matchedRecords.length,
    pageNum: safePageNum,
    pageSize: query.pageSize,
    pages,
  };
}

/**
 * formatPreviewTimestamp 负责把当前时间格式化成假数据记录使用的更新时间文本。
 */
function formatPreviewTimestamp(date: Date): string {
  // year 表示当前年份文本。
  const year = date.getFullYear();
  // month 表示补零后的月份文本。
  const month = String(date.getMonth() + 1).padStart(2, "0");
  // day 表示补零后的日期文本。
  const day = String(date.getDate()).padStart(2, "0");
  // hour 表示补零后的小时文本。
  const hour = String(date.getHours()).padStart(2, "0");
  // minute 表示补零后的分钟文本。
  const minute = String(date.getMinutes()).padStart(2, "0");
  // second 表示补零后的秒钟文本。
  const second = String(date.getSeconds()).padStart(2, "0");

  return `${year}-${month}-${day} ${hour}:${minute}:${second}`;
}

/**
 * renderNoticeClassName 负责根据反馈条语义返回对应样式。
 */
function renderNoticeClassName(tone: NoticeState["tone"]): string {
  return tone === "success"
    ? "border-emerald-300/15 bg-emerald-300/10 text-emerald-100"
    : "border-red-300/15 bg-red-400/10 text-red-100";
}

/**
 * FilterField 负责渲染带标题的统一筛选字段容器。
 */
function FilterField({ children, label }: FilterFieldProps): React.JSX.Element {
  return (
    <label className="grid gap-2">
      <span className="text-[11px] font-medium uppercase tracking-[0.2em] text-zinc-500">
        {label}
      </span>
      {children}
    </label>
  );
}

/**
 * FollowerUsersWorkspace 负责渲染小龙虾用户查询页的主内容区与 modal CRUD 交互。
 */
export function FollowerUsersWorkspace(): React.JSX.Element {
  // filters 保存当前列表筛选条件。
  const [filters, setFilters] = useState<FollowerUserListQuery>(defaultFollowerUserFilters);
  // listData 保存实时接口返回的列表结果；为空时回退到假数据模式。
  const [listData, setListData] = useState<FollowerUserListData | null>(null);
  // previewRecords 保存假数据模式下允许本地增删改的记录集合。
  const [previewRecords, setPreviewRecords] = useState<readonly FollowerUserRecord[]>(
    followerUsersPreviewRecords,
  );
  // listError 保存实时列表请求失败时的错误说明。
  const [listError, setListError] = useState<string>("");
  // isListLoading 表示列表区域是否正在同步实时数据。
  const [isListLoading, setIsListLoading] = useState<boolean>(true);
  // activeDialog 表示当前打开的模态框类型。
  const [activeDialog, setActiveDialog] = useState<ModalMode | null>(null);
  // activeRecord 保存当前准备编辑或删除的用户记录。
  const [activeRecord, setActiveRecord] = useState<FollowerUserRecord | null>(null);
  // createForm 保存新增弹窗绑定的表单状态。
  const [createForm, setCreateForm] = useState<FollowerUserFormState>(defaultFollowerUserCreateForm);
  // editForm 保存编辑或删除弹窗绑定的表单状态。
  const [editForm, setEditForm] = useState<FollowerUserFormState>(defaultFollowerUserCreateForm);
  // isRecordLoading 表示编辑或删除前是否正在刷新该用户详情。
  const [isRecordLoading, setIsRecordLoading] = useState<boolean>(false);
  // isSubmitting 表示当前是否正在提交新增、编辑或删除动作。
  const [isSubmitting, setIsSubmitting] = useState<boolean>(false);
  // createError 保存新增动作失败时的错误说明。
  const [createError, setCreateError] = useState<string>("");
  // actionError 保存编辑或删除动作失败时的错误说明。
  const [actionError, setActionError] = useState<string>("");
  // notice 保存页面顶部的成功或错误反馈条。
  const [notice, setNotice] = useState<NoticeState | null>(null);

  useEffect(() => {
    void (async () => {
      setIsListLoading(true);
      setListError("");

      try {
        // initialListData 表示页面首次加载时拿到的实时列表结果。
        const initialListData = await fetchFollowerUserList(defaultFollowerUserFilters);

        setListData(initialListData);
      } catch (error) {
        setListError(error instanceof Error ? error.message : "实时列表加载失败。");
      } finally {
        setIsListLoading(false);
      }
    })();
  }, []);

  /**
   * reloadList 负责按当前筛选条件重新同步实时列表。
   */
  async function reloadList(nextFilters: FollowerUserListQuery): Promise<void> {
    setIsListLoading(true);
    setListError("");

    try {
      // nextListData 表示最新一次从接口返回的分页结果。
      const nextListData = await fetchFollowerUserList(nextFilters);

      setListData(nextListData);
      setFilters(nextFilters);
    } catch (error) {
      setFilters(nextFilters);
      setListError(error instanceof Error ? error.message : "实时列表加载失败。");
    } finally {
      setIsListLoading(false);
    }
  }

  /**
   * hydrateRecordForAction 负责在打开编辑或删除弹窗后补齐当前记录的最新详情。
   */
  async function hydrateRecordForAction(record: FollowerUserRecord): Promise<void> {
    // normalizedRecord 表示已经补齐 rowVersion 的当前行记录。
    const normalizedRecord = normalizeRecord(record);

    setActiveRecord(normalizedRecord);
    setEditForm(mapRecordToFormState(normalizedRecord));
    setActionError("");

    if (listData === null) {
      setIsRecordLoading(false);
      return;
    }

    setIsRecordLoading(true);

    try {
      // detailRecord 表示详情接口返回的最新记录对象。
      const detailRecord = normalizeRecord(await fetchFollowerUserDetail(record.userId));

      setActiveRecord(detailRecord);
      setEditForm(mapRecordToFormState(detailRecord));
    } catch (error) {
      setActionError(
        error instanceof Error
          ? `${error.message} 当前先使用列表中的版本信息继续展示。`
          : "未能刷新实时详情，当前先使用列表中的版本信息继续展示。",
      );
    } finally {
      setIsRecordLoading(false);
    }
  }

  /**
   * handleDialogOpenChange 负责统一处理模态框的打开与关闭。
   */
  function handleDialogOpenChange(open: boolean): void {
    if (open) {
      return;
    }

    setActiveDialog(null);
    setActiveRecord(null);
    setCreateError("");
    setActionError("");
    setIsRecordLoading(false);
  }

  /**
   * handleFilterInputChange 负责更新筛选条件中的单个字段。
   */
  function handleFilterInputChange<TKey extends keyof FollowerUserListQuery>(
    key: TKey,
    value: FollowerUserListQuery[TKey],
  ): void {
    setFilters((currentFilters) => ({
      ...currentFilters,
      [key]: value,
    }));
  }

  /**
   * handleCreateFormChange 负责更新新增表单中的单个字段。
   */
  function handleCreateFormChange<TKey extends keyof FollowerUserFormState>(
    key: TKey,
    value: FollowerUserFormState[TKey],
  ): void {
    setCreateForm((currentForm) => ({
      ...currentForm,
      [key]: value,
    }));
  }

  /**
   * handleEditFormChange 负责更新编辑表单中的单个字段。
   */
  function handleEditFormChange<TKey extends keyof FollowerUserFormState>(
    key: TKey,
    value: FollowerUserFormState[TKey],
  ): void {
    setEditForm((currentForm) => ({
      ...currentForm,
      [key]: value,
    }));
  }

  /**
   * handleSearchSubmit 负责根据当前筛选条件刷新列表。
   */
  async function handleSearchSubmit(event: React.FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    setNotice(null);

    await reloadList({
      ...filters,
      pageNum: 1,
    });
  }

  /**
   * handleResetFilters 负责恢复默认筛选条件并重新查询列表。
   */
  async function handleResetFilters(): Promise<void> {
    setNotice(null);

    await reloadList(defaultFollowerUserFilters);
  }

  /**
   * handleOpenCreateDialog 负责打开新增用户弹窗。
   */
  function handleOpenCreateDialog(): void {
    setCreateForm(defaultFollowerUserCreateForm);
    setCreateError("");
    setActiveDialog("create");
  }

  /**
   * handleOpenEditDialog 负责打开编辑弹窗并准备当前记录详情。
   */
  async function handleOpenEditDialog(record: FollowerUserRecord): Promise<void> {
    setActiveDialog("edit");
    await hydrateRecordForAction(record);
  }

  /**
   * handleOpenDeleteDialog 负责打开删除确认弹窗并准备当前记录详情。
   */
  async function handleOpenDeleteDialog(record: FollowerUserRecord): Promise<void> {
    setActiveDialog("delete");
    await hydrateRecordForAction(record);
  }

  /**
   * handleCreateSubmit 负责提交新增请求，并在成功后刷新当前列表。
   */
  async function handleCreateSubmit(event: React.FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    setCreateError("");
    setNotice(null);

    // validationMessage 保存新增表单当前的校验结果。
    const validationMessage = validateFollowerUserForm(createForm, {
      includeUserId: true,
    });

    if (validationMessage) {
      setCreateError(validationMessage);
      return;
    }

    setIsSubmitting(true);

    try {
      if (listData === null) {
        // createdPreviewRecord 表示假数据模式下新建的本地记录。
        const createdPreviewRecord = normalizeRecord({
          ...buildCreatePayload(createForm),
          updatedAt: formatPreviewTimestamp(new Date()),
          rowVersion: 1,
        });

        setPreviewRecords((currentRecords) => [createdPreviewRecord, ...currentRecords]);
        setFilters((currentFilters) => ({
          ...currentFilters,
          pageNum: 1,
        }));
        setNotice({
          tone: "success",
          message: `已在假数据模式下创建用户 ${createdPreviewRecord.userId}。`,
        });
        setCreateForm(defaultFollowerUserCreateForm);
        setActiveDialog(null);
        return;
      }

      // createdRecord 表示实时接口创建成功后返回的最新用户对象。
      const createdRecord = await createFollowerUser(buildCreatePayload(createForm));

      setNotice({
        tone: "success",
        message: `已创建用户 ${createdRecord.userId}。`,
      });
      setCreateForm(defaultFollowerUserCreateForm);
      setActiveDialog(null);
      await reloadList({
        ...filters,
        pageNum: 1,
      });
    } catch (error) {
      setCreateError(error instanceof Error ? error.message : "创建失败。");
    } finally {
      setIsSubmitting(false);
    }
  }

  /**
   * handleUpdateSubmit 负责提交编辑请求，并在成功后刷新当前列表。
   */
  async function handleUpdateSubmit(event: React.FormEvent<HTMLFormElement>): Promise<void> {
    event.preventDefault();
    setActionError("");
    setNotice(null);

    if (!activeRecord) {
      setActionError("请先选择一个用户，再进行修改。");
      return;
    }

    // validationMessage 保存编辑表单当前的校验结果。
    const validationMessage = validateFollowerUserForm(editForm, {
      includeUserId: false,
    });

    if (validationMessage) {
      setActionError(validationMessage);
      return;
    }

    setIsSubmitting(true);

    try {
      if (listData === null) {
        // updatedPreviewRecord 表示假数据模式下编辑后的最新记录。
        const updatedPreviewRecord = normalizeRecord({
          ...activeRecord,
          ...buildUpdatePayload(editForm),
          updatedAt: formatPreviewTimestamp(new Date()),
          rowVersion: editForm.rowVersion + 1,
        });

        setPreviewRecords((currentRecords) =>
          currentRecords.map((record) =>
            record.userId === activeRecord.userId ? updatedPreviewRecord : record,
          ),
        );
        setActiveRecord(updatedPreviewRecord);
        setEditForm(mapRecordToFormState(updatedPreviewRecord));
        setNotice({
          tone: "success",
          message: `已在假数据模式下更新用户 ${activeRecord.userId}。`,
        });
        setActiveDialog(null);
        return;
      }

      // updatedRecord 表示实时接口更新成功后返回的最新记录对象。
      const updatedRecord = await updateFollowerUser(
        activeRecord.userId,
        buildUpdatePayload(editForm),
      );
      // nextRecord 表示补齐 rowVersion 后用于前端状态同步的对象。
      const nextRecord = normalizeRecord({
        ...updatedRecord,
        rowVersion: updatedRecord.rowVersion ?? editForm.rowVersion + 1,
      });

      setActiveRecord(nextRecord);
      setEditForm(mapRecordToFormState(nextRecord));
      setNotice({
        tone: "success",
        message: `已更新用户 ${activeRecord.userId}。`,
      });
      setActiveDialog(null);
      await reloadList(filters);
    } catch (error) {
      setActionError(error instanceof Error ? error.message : "修改失败。");
    } finally {
      setIsSubmitting(false);
    }
  }

  /**
   * handleDeleteSubmit 负责提交删除请求，并在成功后刷新当前列表。
   */
  async function handleDeleteSubmit(): Promise<void> {
    setActionError("");
    setNotice(null);

    if (!activeRecord) {
      setActionError("请先选择一个用户，再执行删除。");
      return;
    }

    setIsSubmitting(true);

    try {
      if (listData === null) {
        setPreviewRecords((currentRecords) =>
          currentRecords.filter((record) => record.userId !== activeRecord.userId),
        );
        setNotice({
          tone: "success",
          message: `已在假数据模式下删除用户 ${activeRecord.userId}。`,
        });
        setActiveDialog(null);
        setActiveRecord(null);
        return;
      }

      // deleteResult 表示实时接口返回的删除执行结果。
      const deleteResult = await deleteFollowerUser(activeRecord.userId, {
        rowVersion: editForm.rowVersion,
      });

      setNotice({
        tone: "success",
        message: deleteResult.success
          ? `已删除用户 ${deleteResult.userId}。`
          : "删除请求已提交，但接口未明确返回成功标记。",
      });
      setActiveDialog(null);
      setActiveRecord(null);
      await reloadList(filters);
    } catch (error) {
      setActionError(error instanceof Error ? error.message : "删除失败。");
    } finally {
      setIsSubmitting(false);
    }
  }

  /**
   * handlePageChange 负责切换当前页码并重新刷新列表。
   */
  async function handlePageChange(nextPageNum: number): Promise<void> {
    await reloadList({
      ...filters,
      pageNum: nextPageNum,
    });
  }

  // previewListData 表示当前筛选条件下的假数据分页结果。
  const previewListData = buildPreviewListData(previewRecords, filters);
  // displayListData 表示当前真正用于渲染的分页数据。
  const displayListData = listData ?? previewListData;
  // displayList 表示当前页要渲染的记录集合。
  const displayList = displayListData.list;
  // isPreviewMode 表示当前是否处于假数据模式。
  const isPreviewMode = listData === null;
  // shouldShowListError 表示当前是否需要展示实时同步失败提示。
  const shouldShowListError = Boolean(listError);

  return (
    <div className="space-y-5">
      <div className="flex flex-col gap-4 rounded-[28px] border border-white/10 bg-[rgba(8,11,16,0.76)] px-5 py-5 backdrop-blur-xl sm:flex-row sm:items-center sm:justify-between sm:px-6">
        <div className="min-w-0">
          <h1 className="truncate text-2xl font-semibold tracking-[0.02em] text-zinc-50 sm:text-[1.75rem]">
            {followerUsersConsolePageTitle}
          </h1>
          <p className="mt-2 text-xs uppercase tracking-[0.22em] text-zinc-500">
            Query / Filter / Table / Modal CRUD
          </p>
        </div>

        <div className="flex items-center gap-3">
          <Button
            type="button"
            variant="secondary"
            className={cn(subtleButtonClassName, "hidden sm:inline-flex")}
            disabled={isListLoading}
            onClick={() => {
              void reloadList(filters);
            }}
          >
            {isListLoading ? (
              <LoaderCircle className="mr-2 size-4 animate-spin" />
            ) : (
              <RefreshCcw className="mr-2 size-4" />
            )}
            刷新
          </Button>
          <Button
            type="button"
            className="rounded-2xl bg-cyan-300/12 px-4 text-cyan-50 shadow-[0_0_0_1px_rgba(125,211,252,0.12)] hover:bg-cyan-300/18"
            onClick={handleOpenCreateDialog}
          >
            <Plus className="mr-2 size-4" />
            新增用户
          </Button>
        </div>
      </div>

      {notice ? (
        <div
          className={cn(
            "rounded-[22px] border px-4 py-3 text-sm font-medium",
            renderNoticeClassName(notice.tone),
          )}
        >
          {notice.message}
        </div>
      ) : null}

      {shouldShowListError ? (
        <div className="rounded-[22px] border border-white/10 bg-white/[0.03] px-4 py-3 text-sm text-zinc-300">
          {isPreviewMode
            ? "实时接口暂不可用，当前以假数据模式继续预览和操作页面。"
            : `最近一次同步失败：${listError}`}
        </div>
      ) : null}

      <section className="rounded-[28px] border border-white/10 bg-[rgba(8,11,16,0.76)] px-5 py-5 backdrop-blur-xl sm:px-6">
        <form className="space-y-4" onSubmit={handleSearchSubmit}>
          <div className="grid gap-4 xl:grid-cols-[minmax(0,1.25fr)_repeat(4,minmax(0,1fr))]">
            <FilterField label="用户ID">
              <div className="relative">
                <Search className="pointer-events-none absolute left-4 top-1/2 size-4 -translate-y-1/2 text-zinc-500" />
                <input
                  className={cn(filterControlClassName, "pl-11")}
                  value={filters.userId}
                  onChange={(event) => handleFilterInputChange("userId", event.target.value)}
                  placeholder="搜索用户 ID"
                />
              </div>
            </FilterField>

            <FilterField label="账户状态">
              <select
                className={filterControlClassName}
                value={filters.accountStatus}
                onChange={(event) =>
                  handleFilterInputChange(
                    "accountStatus",
                    event.target.value as "" | FollowerUserAccountStatus,
                  )
                }
              >
                <option value="">全部</option>
                <option value="active">启用</option>
                <option value="disabled">停用</option>
              </select>
            </FilterField>

            <FilterField label="策略状态">
              <select
                className={filterControlClassName}
                value={filters.strategyStatus}
                onChange={(event) =>
                  handleFilterInputChange(
                    "strategyStatus",
                    event.target.value as "" | FollowerUserStrategyStatus,
                  )
                }
              >
                <option value="">全部</option>
                <option value="enabled">启用</option>
                <option value="disabled">停用</option>
              </select>
            </FilterField>

            <FilterField label="绑定状态">
              <select
                className={filterControlClassName}
                value={filters.bindingStatus}
                onChange={(event) =>
                  handleFilterInputChange(
                    "bindingStatus",
                    event.target.value as "" | FollowerUserBindingStatus,
                  )
                }
              >
                <option value="">全部</option>
                <option value="pending">待绑定</option>
                <option value="bound">已绑定</option>
                <option value="unbound">未绑定</option>
              </select>
            </FilterField>

            <FilterField label="责任域">
              <input
                className={filterControlClassName}
                list="follower-user-filter-domain-options"
                value={filters.responsibilityDomain}
                onChange={(event) =>
                  handleFilterInputChange("responsibilityDomain", event.target.value)
                }
                placeholder="输入责任域"
              />
            </FilterField>
          </div>

          <datalist id="follower-user-filter-domain-options">
            {responsibilityDomainOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </datalist>

          <div className="flex flex-wrap gap-3">
            <Button
              type="submit"
              className="rounded-2xl bg-cyan-300/12 px-4 text-cyan-50 shadow-[0_0_0_1px_rgba(125,211,252,0.12)] hover:bg-cyan-300/18"
              disabled={isListLoading}
            >
              {isListLoading ? (
                <LoaderCircle className="mr-2 size-4 animate-spin" />
              ) : (
                <Search className="mr-2 size-4" />
              )}
              查询
            </Button>

            <Button
              type="button"
              variant="secondary"
              className={subtleButtonClassName}
              disabled={isListLoading}
              onClick={() => {
                void handleResetFilters();
              }}
            >
              重置
            </Button>
          </div>
        </form>
      </section>

      <section className="overflow-hidden rounded-[28px] border border-white/10 bg-[rgba(8,11,16,0.76)] backdrop-blur-xl">
        <div className="hidden overflow-x-auto md:block">
          <table className="min-w-full border-collapse">
            <thead>
              <tr className="border-b border-white/10 text-left text-[11px] font-medium uppercase tracking-[0.2em] text-zinc-500">
                {tableHeaders.map((headerLabel) => (
                  <th key={headerLabel} className="px-6 py-4">
                    {headerLabel}
                  </th>
                ))}
              </tr>
            </thead>
            <tbody>
              {displayList.length > 0 ? (
                displayList.map((record) => {
                  // normalizedRecord 表示当前行中已经补齐版本号的记录对象。
                  const normalizedRecord = normalizeRecord(record);

                  return (
                    <tr
                      key={normalizedRecord.userId}
                      className="border-b border-white/6 text-sm transition-colors last:border-b-0 hover:bg-white/[0.03]"
                    >
                      <td className="px-6 py-4 font-medium text-zinc-50">
                        {normalizedRecord.userId}
                      </td>
                      <td className="px-6 py-4">
                        <FollowerUserStatusPill
                          label={accountStatusMetaMap[normalizedRecord.accountStatus].label}
                          tone={accountStatusMetaMap[normalizedRecord.accountStatus].tone}
                        />
                      </td>
                      <td className="px-6 py-4">
                        <FollowerUserStatusPill
                          label={strategyStatusMetaMap[normalizedRecord.strategyStatus].label}
                          tone={strategyStatusMetaMap[normalizedRecord.strategyStatus].tone}
                        />
                      </td>
                      <td className="px-6 py-4">
                        <FollowerUserStatusPill
                          label={bindingStatusMetaMap[normalizedRecord.bindingStatus].label}
                          tone={bindingStatusMetaMap[normalizedRecord.bindingStatus].tone}
                        />
                      </td>
                      <td className="px-6 py-4 text-zinc-300">
                        {normalizedRecord.responsibilityDomain}
                      </td>
                      <td className="px-6 py-4 text-zinc-400">{normalizedRecord.updatedAt}</td>
                      <td className="px-6 py-4">
                        <div className="flex items-center justify-end gap-2">
                          <Button
                            type="button"
                            variant="secondary"
                            className={cn(subtleButtonClassName, "!h-9 !px-3")}
                            onClick={() => {
                              void handleOpenEditDialog(normalizedRecord);
                            }}
                          >
                            <PencilLine className="mr-2 size-3.5" />
                            修改
                          </Button>
                          <Button
                            type="button"
                            variant="secondary"
                            className="!h-9 !rounded-2xl !border-red-400/18 !bg-red-400/[0.08] !px-3 !text-red-100 hover:!bg-red-400/[0.14]"
                            onClick={() => {
                              void handleOpenDeleteDialog(normalizedRecord);
                            }}
                          >
                            <Trash2 className="mr-2 size-3.5" />
                            删除
                          </Button>
                        </div>
                      </td>
                    </tr>
                  );
                })
              ) : (
                <tr>
                  <td colSpan={tableHeaders.length} className="px-6 py-16 text-center">
                    <p className="text-base font-medium text-zinc-100">暂无匹配记录</p>
                    <p className="mt-2 text-sm text-zinc-500">可以尝试调整筛选条件，或者直接新增用户。</p>
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>

        <div className="grid gap-3 px-4 py-4 md:hidden">
          {displayList.length > 0 ? (
            displayList.map((record) => {
              // normalizedRecord 表示当前移动卡片使用的记录对象。
              const normalizedRecord = normalizeRecord(record);

              return (
                <article
                  key={normalizedRecord.userId}
                  className="rounded-[24px] border border-white/10 bg-white/[0.03] p-4"
                >
                  <div className="flex items-start justify-between gap-3">
                    <div>
                      <p className="text-sm font-medium text-zinc-50">{normalizedRecord.userId}</p>
                      <p className="mt-2 text-xs text-zinc-500">{normalizedRecord.updatedAt}</p>
                    </div>
                    <FollowerUserStatusPill
                      label={accountStatusMetaMap[normalizedRecord.accountStatus].label}
                      tone={accountStatusMetaMap[normalizedRecord.accountStatus].tone}
                    />
                  </div>

                  <div className="mt-4 flex flex-wrap gap-2">
                    <FollowerUserStatusPill
                      label={strategyStatusMetaMap[normalizedRecord.strategyStatus].label}
                      tone={strategyStatusMetaMap[normalizedRecord.strategyStatus].tone}
                    />
                    <FollowerUserStatusPill
                      label={bindingStatusMetaMap[normalizedRecord.bindingStatus].label}
                      tone={bindingStatusMetaMap[normalizedRecord.bindingStatus].tone}
                    />
                    <FollowerUserStatusPill label={normalizedRecord.responsibilityDomain} tone="muted" />
                  </div>

                  <div className="mt-4 flex gap-2">
                    <Button
                      type="button"
                      variant="secondary"
                      className={cn(subtleButtonClassName, "flex-1 !h-9")}
                      onClick={() => {
                        void handleOpenEditDialog(normalizedRecord);
                      }}
                    >
                      <PencilLine className="mr-2 size-3.5" />
                      修改
                    </Button>
                    <Button
                      type="button"
                      variant="secondary"
                      className="!h-9 flex-1 !rounded-2xl !border-red-400/18 !bg-red-400/[0.08] !px-3 !text-red-100 hover:!bg-red-400/[0.14]"
                      onClick={() => {
                        void handleOpenDeleteDialog(normalizedRecord);
                      }}
                    >
                      <Trash2 className="mr-2 size-3.5" />
                      删除
                    </Button>
                  </div>
                </article>
              );
            })
          ) : (
            <div className="rounded-[24px] border border-white/10 bg-white/[0.03] px-4 py-10 text-center">
              <p className="text-base font-medium text-zinc-100">暂无匹配记录</p>
              <p className="mt-2 text-sm text-zinc-500">可以尝试调整筛选条件，或者直接新增用户。</p>
            </div>
          )}
        </div>

        <div className="flex flex-col gap-3 border-t border-white/10 px-5 py-4 text-sm text-zinc-500 sm:flex-row sm:items-center sm:justify-between sm:px-6">
          <p>
            第 {displayListData.pageNum} / {displayListData.pages} 页，共 {displayListData.total} 条记录
          </p>

          <div className="flex items-center gap-2">
            <Button
              type="button"
              variant="secondary"
              className={subtleButtonClassName}
              disabled={displayListData.pageNum <= 1 || isListLoading}
              onClick={() => {
                void handlePageChange(displayListData.pageNum - 1);
              }}
            >
              <ChevronLeft className="mr-2 size-4" />
              上一页
            </Button>
            <Button
              type="button"
              variant="secondary"
              className={subtleButtonClassName}
              disabled={displayListData.pageNum >= displayListData.pages || isListLoading}
              onClick={() => {
                void handlePageChange(displayListData.pageNum + 1);
              }}
            >
              下一页
              <ChevronRight className="ml-2 size-4" />
            </Button>
          </div>
        </div>
      </section>

      <Dialog
        open={activeDialog === "create"}
        onOpenChange={handleDialogOpenChange}
        title="新增用户"
        description="通过模态框直接完成新增，不跳转页面。"
      >
        <form className="space-y-5" onSubmit={handleCreateSubmit}>
          {createError ? (
            <div className="rounded-[20px] border border-red-300/15 bg-red-400/10 px-4 py-3 text-sm text-red-100">
              {createError}
            </div>
          ) : null}

          <FollowerUserFormFields
            formState={createForm}
            mode="create"
            disabled={isSubmitting}
            onChange={handleCreateFormChange}
          />

          <div className="flex flex-col-reverse gap-3 pt-2 sm:flex-row sm:justify-end">
            <Button
              type="button"
              variant="secondary"
              className={subtleButtonClassName}
              onClick={() => handleDialogOpenChange(false)}
            >
              取消
            </Button>
            <Button
              type="submit"
              className="rounded-2xl bg-cyan-300/12 px-4 text-cyan-50 shadow-[0_0_0_1px_rgba(125,211,252,0.12)] hover:bg-cyan-300/18"
              disabled={isSubmitting}
            >
              {isSubmitting ? (
                <LoaderCircle className="mr-2 size-4 animate-spin" />
              ) : (
                <Plus className="mr-2 size-4" />
              )}
              确认新增
            </Button>
          </div>
        </form>
      </Dialog>

      <Dialog
        open={activeDialog === "edit"}
        onOpenChange={handleDialogOpenChange}
        title={activeRecord ? `修改 ${activeRecord.userId}` : "修改用户"}
        description="编辑动作同样通过 modal 完成，不跳转页面。"
      >
        <form className="space-y-5" onSubmit={handleUpdateSubmit}>
          {isRecordLoading ? (
            <div className="rounded-[20px] border border-white/10 bg-white/[0.04] px-4 py-3 text-sm text-zinc-300">
              <span className="inline-flex items-center gap-2">
                <LoaderCircle className="size-4 animate-spin" />
                正在同步当前用户的最新详情。
              </span>
            </div>
          ) : null}

          {actionError ? (
            <div className="rounded-[20px] border border-red-300/15 bg-red-400/10 px-4 py-3 text-sm text-red-100">
              {actionError}
            </div>
          ) : null}

          <FollowerUserFormFields
            formState={editForm}
            mode="edit"
            disabled={isSubmitting}
            onChange={handleEditFormChange}
          />

          <div className="flex flex-col-reverse gap-3 pt-2 sm:flex-row sm:justify-end">
            <Button
              type="button"
              variant="secondary"
              className={subtleButtonClassName}
              onClick={() => handleDialogOpenChange(false)}
            >
              取消
            </Button>
            <Button
              type="submit"
              className="rounded-2xl bg-cyan-300/12 px-4 text-cyan-50 shadow-[0_0_0_1px_rgba(125,211,252,0.12)] hover:bg-cyan-300/18"
              disabled={isSubmitting}
            >
              {isSubmitting ? (
                <LoaderCircle className="mr-2 size-4 animate-spin" />
              ) : (
                <PencilLine className="mr-2 size-4" />
              )}
              保存修改
            </Button>
          </div>
        </form>
      </Dialog>

      <Dialog
        open={activeDialog === "delete"}
        onOpenChange={handleDialogOpenChange}
        title="确认删除"
        description="删除动作通过确认弹窗完成，不跳转页面。"
        className="max-w-xl"
      >
        <div className="space-y-5">
          {isRecordLoading ? (
            <div className="rounded-[20px] border border-white/10 bg-white/[0.04] px-4 py-3 text-sm text-zinc-300">
              <span className="inline-flex items-center gap-2">
                <LoaderCircle className="size-4 animate-spin" />
                正在刷新删除前的版本信息。
              </span>
            </div>
          ) : null}

          {actionError ? (
            <div className="rounded-[20px] border border-red-300/15 bg-red-400/10 px-4 py-3 text-sm text-red-100">
              {actionError}
            </div>
          ) : null}

          <div className="rounded-[24px] border border-white/10 bg-white/[0.03] p-5">
            <p className="text-[11px] font-medium uppercase tracking-[0.22em] text-zinc-500">
              用户ID
            </p>
            <p className="mt-3 text-2xl font-semibold tracking-[0.02em] text-zinc-50">
              {activeRecord?.userId ?? "--"}
            </p>
            <p className="mt-3 text-sm leading-6 text-zinc-400">
              确认删除后，该用户会从当前列表中移除。整个过程保持在当前页面内完成。
            </p>
          </div>

          <div className="flex flex-col-reverse gap-3 pt-2 sm:flex-row sm:justify-end">
            <Button
              type="button"
              variant="secondary"
              className={subtleButtonClassName}
              onClick={() => handleDialogOpenChange(false)}
            >
              取消
            </Button>
            <Button
              type="button"
              variant="secondary"
              className="!rounded-2xl !border-red-400/18 !bg-red-400/[0.08] !px-4 !text-red-100 hover:!bg-red-400/[0.14]"
              disabled={isSubmitting}
              onClick={() => {
                void handleDeleteSubmit();
              }}
            >
              {isSubmitting ? (
                <LoaderCircle className="mr-2 size-4 animate-spin" />
              ) : (
                <Trash2 className="mr-2 size-4" />
              )}
              确认删除
            </Button>
          </div>
        </div>
      </Dialog>
    </div>
  );
}
