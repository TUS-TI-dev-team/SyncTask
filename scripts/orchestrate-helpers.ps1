# ==============================================================================
# SyncTask Orchestration Helper Script (PowerShell)
# ==============================================================================
param (
    [Parameter(Position = 0, Mandatory = $true)]
    [ValidateSet("create-worktree", "cleanup-worktree", "list-worktrees")]
    [string]$Command,

    [Parameter(Position = 1)]
    [string]$IssueNum,

    [Parameter(Position = 2)]
    [string]$Param2
)

$RepoRoot = (Get-Item -Path $PSScriptRoot).Parent.FullName

function Create-Worktree {
    param (
        [string]$Issue,
        [string]$Endpoint
    )

    $WorktreePath = Join-Path $RepoRoot ".worktrees\issue-$Issue"
    $BranchName = "feature/issue-$Issue-$Endpoint"

    Write-Host "[orchestrate] Creating git worktree at $WorktreePath..." -ForegroundColor Cyan
    $WorktreesDir = Join-Path $RepoRoot ".worktrees"
    if (-not (Test-Path $WorktreesDir)) {
        New-Item -ItemType Directory -Path $WorktreesDir | Out-Null
    }

    git -C $RepoRoot worktree add -b $BranchName $WorktreePath main

    Write-Host "[orchestrate] Creating Herdr workspace..." -ForegroundColor Cyan
    $CreateRes = herdr workspace create --cwd $WorktreePath --label "ep-issue-$Issue" --no-focus
    $JsonObj = $CreateRes | ConvertFrom-Json

    $WsId = $JsonObj.result.workspace.workspace_id
    $RootPaneId = $JsonObj.result.root_pane.pane_id

    Write-Host "[orchestrate] Workspace created: ID=$WsId, RootPane=$RootPaneId" -ForegroundColor Green

    $SupervisorName = "sup-issue-$Issue"
    Write-Host "[orchestrate] Starting Endpoint Supervisor agent ($SupervisorName)..." -ForegroundColor Cyan
    herdr agent start $SupervisorName --kind agy --pane $RootPaneId

    Write-Host "[orchestrate] Initializing supervisor with prompt..." -ForegroundColor Cyan
    $Prompt = "/endpoint-supervisor Issue #$Issue`: $Endpoint"
    herdr agent prompt $SupervisorName $Prompt

    Write-Host "[orchestrate] Setup completed for Issue #$Issue. Supervisor is running." -ForegroundColor Green
    [PSCustomObject]@{
        issue_num    = $Issue
        workspace_id = $WsId
        root_pane_id = $RootPaneId
        supervisor   = $SupervisorName
    } | ConvertTo-Json -Compress
}

function Cleanup-Worktree {
    param (
        [string]$Issue,
        [string]$WsId
    )

    $WorktreePath = Join-Path $RepoRoot ".worktrees\issue-$Issue"

    Write-Host "[orchestrate] Closing Herdr workspace $WsId..." -ForegroundColor Yellow
    herdr workspace close $WsId

    Write-Host "[orchestrate] Removing git worktree $WorktreePath..." -ForegroundColor Yellow
    git -C $RepoRoot worktree remove $WorktreePath --force
    git -C $RepoRoot worktree prune

    Write-Host "[orchestrate] Cleanup completed for Issue #$Issue." -ForegroundColor Green
}

function List-Worktrees {
    Write-Host "=== Git Worktrees ===" -ForegroundColor Cyan
    git -C $RepoRoot worktree list
    Write-Host "`n=== Herdr Workspaces ===" -ForegroundColor Cyan
    herdr workspace list
}

switch ($Command) {
    "create-worktree" {
        if (-not $IssueNum -or -not $Param2) {
            Write-Error "Usage: .\scripts\orchestrate-helpers.ps1 create-worktree <issue_num> <endpoint_name>"
            exit 1
        }
        Create-Worktree -Issue $IssueNum -Endpoint $Param2
    }
    "cleanup-worktree" {
        if (-not $IssueNum -or -not $Param2) {
            Write-Error "Usage: .\scripts\orchestrate-helpers.ps1 cleanup-worktree <issue_num> <workspace_id>"
            exit 1
        }
        Cleanup-Worktree -Issue $IssueNum -WsId $Param2
    }
    "list-worktrees" {
        List-Worktrees
    }
}
