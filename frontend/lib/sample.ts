/**
 * @file sample.ts
 * @description SyncTask フロントエンドのサンプルビジネスロジックモジュール
 * Code-as-Docs 思想に基づき、各関数の事前条件・事後条件および仕様を明記しています。
 */

/**
 * タスクのステータス定義
 */
export type TaskStatus = 'todo' | 'in_progress' | 'done';

/**
 * タスクの優先度定義
 */
export type TaskPriority = 'low' | 'medium' | 'high';

/**
 * タスクエンティティ
 */
export interface Task {
  /** 一意なタスクID */
  id: string;
  /** タスクのタイトル（必須・1文字以上） */
  title: string;
  /** タスクの進捗状況 */
  status: TaskStatus;
  /** タスクの優先度 */
  priority: TaskPriority;
  /** 作成日時 (ISO 8601 文字列) */
  createdAt: string;
}

/**
 * タスク統計情報の集計結果
 */
export interface TaskStats {
  total: number;
  todo: number;
  inProgress: number;
  done: number;
  completionRate: number; // 0.0 - 100.0 (%)
}

/**
 * 新規タスクオブジェクトを作成します。
 *
 * @spec
 * 1. title が空文字列（トリム後）の場合、エラー（Error: Title cannot be empty）を送出する。
 * 2. status は初期値として 'todo' が設定される。
 * 3. id はタイムスタンプまたはランダム文字列から一意に生成される。
 * 4. priority のデフォルトは 'medium' とする。
 *
 * @param title タスクのタイトル
 * @param priority 優先度（省略時は 'medium'）
 * @returns 作成された Task オブジェクト
 */
export function createTask(title: string, priority: TaskPriority = 'medium'): Task {
  const trimmedTitle = title.trim();
  if (trimmedTitle.length === 0) {
    throw new Error('Title cannot be empty');
  }

  return {
    id: `task-${Date.now()}-${Math.random().toString(36).substring(2, 7)}`,
    title: trimmedTitle,
    status: 'todo',
    priority,
    createdAt: new Date().toISOString(),
  };
}

/**
 * タスクのステータスを次の状態に遷移させます。
 * 状態遷移: 'todo' -> 'in_progress' -> 'done' -> 'todo' (ループ)
 *
 * @spec
 * 1. 'todo' の場合は 'in_progress' に変更した新しい Task を返す。
 * 2. 'in_progress' の場合は 'done' に変更した新しい Task を返す。
 * 3. 'done' の場合は 'todo' に変更した新しい Task を返す。
 * 4. 引数のオブジェクトは変更せず（イミュータブル）、新しいオブジェクトを返す。
 *
 * @param task 遷移前のタスク
 * @returns 状態が更新された新しい Task オブジェクト
 */
export function toggleTaskStatus(task: Task): Task {
  let nextStatus: TaskStatus;
  switch (task.status) {
    case 'todo':
      nextStatus = 'in_progress';
      break;
    case 'in_progress':
      nextStatus = 'done';
      break;
    case 'done':
      nextStatus = 'todo';
      break;
  }
  return { ...task, status: nextStatus };
}

/**
 * 指定したステータスに一致するタスクのみを抽出します。
 *
 * @spec
 * 1. status が指定された場合、その status を持つタスク配列を返す。
 * 2. status が 'all' または undefined の場合は全タスクをそのまま返す。
 *
 * @param tasks タスク配列
 * @param status 絞り込み条件 ('all' または TaskStatus)
 * @returns 絞り込み後のタスク配列
 */
export function filterTasksByStatus(tasks: Task[], status?: TaskStatus | 'all'): Task[] {
  if (!status || status === 'all') {
    return [...tasks];
  }
  return tasks.filter((t) => t.status === status);
}

/**
 * タスク一覧から各種ステータスごとの件数および完了率（%）を算出します。
 *
 * @spec
 * 1. tasks が空配列の場合、total: 0, todo: 0, inProgress: 0, done: 0, completionRate: 0 を返す。
 * 2. completionRate は (done / total) * 100 で計算され、小数点第1位に四捨五入される。
 *
 * @param tasks タスク配列
 * @returns TaskStats 集計結果オブジェクト
 */
export function calculateTaskStats(tasks: Task[]): TaskStats {
  const total = tasks.length;
  if (total === 0) {
    return {
      total: 0,
      todo: 0,
      inProgress: 0,
      done: 0,
      completionRate: 0,
    };
  }

  let todo = 0;
  let inProgress = 0;
  let done = 0;

  for (const task of tasks) {
    if (task.status === 'todo') todo++;
    else if (task.status === 'in_progress') inProgress++;
    else if (task.status === 'done') done++;
  }

  const completionRate = Math.round((done / total) * 1000) / 10;

  return {
    total,
    todo,
    inProgress,
    done,
    completionRate,
  };
}
