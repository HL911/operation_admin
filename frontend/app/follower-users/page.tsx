import type { Metadata } from "next";
import { FollowerUsersConsoleShell } from "@/features/follower-users/components/follower-users-console-shell";
import { FollowerUsersWorkspace } from "@/features/follower-users/components/follower-users-workspace";

// metadata 负责定义小龙虾用户页面在浏览器中的标题和说明。
export const metadata: Metadata = {
  title: "小龙虾用户",
  description: "深色 Web3 风格的运营后台小龙虾用户查询页，包含固定侧边栏、固定顶栏和 modal CRUD。",
};

/**
 * FollowerUsersPage 负责渲染带固定侧边栏与顶栏的运营后台查询页面。
 */
export default function FollowerUsersPage(): React.JSX.Element {
  return (
    <div className="relative min-h-screen overflow-hidden bg-[#04070a] text-zinc-50">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(34,211,238,0.10),transparent_22%),radial-gradient(circle_at_90%_0%,rgba(255,255,255,0.04),transparent_18%),linear-gradient(180deg,#04070a_0%,#060a0f_100%)]" />
      <div className="relative">
        <FollowerUsersConsoleShell>
          <FollowerUsersWorkspace />
        </FollowerUsersConsoleShell>
      </div>
    </div>
  );
}
