#!/usr/bin/env bash
set -euo pipefail

# ==============================================================================
# SyncTask Orchestration Helper Script (Bash)
# ==============================================================================

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

usage() {
  cat <<EOF
Usage:
  $0 create-worktree <issue_num> <endpoint_name>
  $0 cleanup-worktree <issue_num> <workspace_id>
  $0 list-worktrees

Examples:
  $0 create-worktree 66 post-tasks
  $0 cleanup-worktree 66 w2
EOF
  exit 1
}

create_worktree() {
  local issue_num="$1"
  local endpoint_name="$2"
  local worktree_path="${REPO_ROOT}/.worktrees/issue-${issue_num}"
  local branch_name="feature/issue-${issue_num}-${endpoint_name}"

  echo "[orchestrate] Creating git worktree at ${worktree_path}..."
  mkdir -p "${REPO_ROOT}/.worktrees"
  git -C "${REPO_ROOT}" worktree add -b "${branch_name}" "${worktree_path}" main

  echo "[orchestrate] Creating Herdr workspace..."
  local create_res
  create_res=$(herdr workspace create --cwd "${worktree_path}" --label "ep-issue-${issue_num}" --no-focus)
  
  local ws_id root_pane_id
  ws_id=$(printf '%s\n' "${create_res}" | jq -r '.result.workspace.workspace_id')
  root_pane_id=$(printf '%s\n' "${create_res}" | jq -r '.result.root_pane.pane_id')

  echo "[orchestrate] Workspace created: ID=${ws_id}, RootPane=${root_pane_id}"

  local supervisor_name="sup-issue-${issue_num}"
  echo "[orchestrate] Starting Endpoint Supervisor agent (${supervisor_name})..."
  herdr agent start "${supervisor_name}" --kind agy --pane "${root_pane_id}"

  echo "[orchestrate] Initializing supervisor with prompt..."
  local prompt="/endpoint-supervisor Issue #${issue_num}: ${endpoint_name}"
  herdr agent prompt "${supervisor_name}" "${prompt}"

  echo "[orchestrate] Setup completed for Issue #${issue_num}. Supervisor is running."
  printf '{"issue_num":"%s","workspace_id":"%s","root_pane_id":"%s","supervisor":"%s"}\n' \
    "${issue_num}" "${ws_id}" "${root_pane_id}" "${supervisor_name}"
}

cleanup_worktree() {
  local issue_num="$1"
  local ws_id="$2"
  local worktree_path="${REPO_ROOT}/.worktrees/issue-${issue_num}"

  echo "[orchestrate] Closing Herdr workspace ${ws_id}..."
  herdr workspace close "${ws_id}" || true

  echo "[orchestrate] Removing git worktree ${worktree_path}..."
  git -C "${REPO_ROOT}" worktree remove "${worktree_path}" --force || true
  git -C "${REPO_ROOT}" worktree prune || true

  echo "[orchestrate] Cleanup completed for Issue #${issue_num}."
}

list_worktrees() {
  git -C "${REPO_ROOT}" worktree list
  herdr workspace list
}

COMMAND="${1:-}"
case "${COMMAND}" in
  create-worktree)
    if [ "$#" -ne 3 ]; then usage; fi
    create_worktree "$2" "$3"
    ;;
  cleanup-worktree)
    if [ "$#" -ne 3 ]; then usage; fi
    cleanup_worktree "$2" "$3"
    ;;
  list-worktrees)
    list_worktrees
    ;;
  *)
    usage
    ;;
esac
