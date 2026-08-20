'use client';

import React, { useState } from 'react';
import {
  Task,
  TaskPriority,
  createTask,
  toggleTaskStatus,
  calculateTaskStats,
} from '@/lib/sample';

/**
 * 開発検証用サンプルページ
 * 単体テストおよび Playwright E2E ブラウザテストの検証対象コンポーネント
 */
export default function DevSamplePage() {
  const [tasks, setTasks] = useState<Task[]>([
    {
      id: 'task-initial-1',
      title: '初期タスク 1',
      status: 'todo',
      priority: 'medium',
      createdAt: new Date().toISOString(),
    },
  ]);
  const [inputTitle, setInputTitle] = useState('');
  const [inputPriority, setInputPriority] = useState<TaskPriority>('medium');
  const [errorMessage, setErrorMessage] = useState('');

  const stats = calculateTaskStats(tasks);

  const handleAddTask = (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const newTask = createTask(inputTitle, inputPriority);
      setTasks((prev) => [...prev, newTask]);
      setInputTitle('');
      setErrorMessage('');
    } catch (err: unknown) {
      if (err instanceof Error) {
        setErrorMessage(err.message);
      }
    }
  };

  const handleToggle = (id: string) => {
    setTasks((prev) =>
      prev.map((t) => (t.id === id ? toggleTaskStatus(t) : t))
    );
  };

  return (
    <main className="min-h-screen p-8 bg-gray-50 text-gray-900">
      <div className="max-w-2xl mx-auto bg-white p-6 rounded-lg shadow">
        <h1 className="text-2xl font-bold mb-4" data-testid="page-title">
          開発テスト用サンプルページ (Dev Sample)
        </h1>
        <p className="text-sm text-gray-600 mb-6">
          このページはテスト環境（単体テスト・E2E結合テスト）の動作検証用ページです。
        </p>

        {/* 統計エリア */}
        <section className="mb-6 p-4 bg-blue-50 rounded" data-testid="stats-section">
          <h2 className="font-semibold text-blue-900 mb-2">タスク集計情報</h2>
          <div className="grid grid-cols-4 gap-2 text-center">
            <div>
              <span className="text-xs text-gray-500 block">合計</span>
              <span className="text-lg font-bold" data-testid="stat-total">
                {stats.total}
              </span>
            </div>
            <div>
              <span className="text-xs text-gray-500 block">未着手</span>
              <span className="text-lg font-bold text-yellow-600" data-testid="stat-todo">
                {stats.todo}
              </span>
            </div>
            <div>
              <span className="text-xs text-gray-500 block">進行中</span>
              <span className="text-lg font-bold text-blue-600" data-testid="stat-in-progress">
                {stats.inProgress}
              </span>
            </div>
            <div>
              <span className="text-xs text-gray-500 block">完了</span>
              <span className="text-lg font-bold text-green-600" data-testid="stat-done">
                {stats.done}
              </span>
            </div>
          </div>
          <div className="mt-2 text-right text-sm text-gray-600">
            完了率: <span data-testid="stat-completion-rate">{stats.completionRate}%</span>
          </div>
        </section>

        {/* タスク作成フォーム */}
        <form onSubmit={handleAddTask} className="mb-6 space-y-3" data-testid="task-form">
          <div>
            <label htmlFor="task-title" className="block text-sm font-medium mb-1">
              タスク名
            </label>
            <input
              id="task-title"
              data-testid="task-title-input"
              type="text"
              value={inputTitle}
              onChange={(e) => setInputTitle(e.target.value)}
              placeholder="新しいタスクを入力"
              className="w-full border rounded px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>

          <div>
            <label htmlFor="task-priority" className="block text-sm font-medium mb-1">
              優先度
            </label>
            <select
              id="task-priority"
              data-testid="task-priority-select"
              value={inputPriority}
              onChange={(e) => setInputPriority(e.target.value as TaskPriority)}
              className="border rounded px-3 py-2 text-sm"
            >
              <option value="low">Low</option>
              <option value="medium">Medium</option>
              <option value="high">High</option>
            </select>
          </div>

          {errorMessage && (
            <p className="text-red-600 text-sm" data-testid="error-message">
              {errorMessage}
            </p>
          )}

          <button
            type="submit"
            data-testid="add-task-button"
            className="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded text-sm font-medium transition-colors"
          >
            タスクを追加
          </button>
        </form>

        {/* タスク一覧 */}
        <section data-testid="task-list-section">
          <h2 className="font-semibold mb-3">タスク一覧</h2>
          <ul className="divide-y border rounded" data-testid="task-list">
            {tasks.map((task) => (
              <li
                key={task.id}
                data-testid="task-item"
                className="p-3 flex items-center justify-between hover:bg-gray-50"
              >
                <div>
                  <span
                    className={`font-medium ${
                      task.status === 'done' ? 'line-through text-gray-400' : ''
                    }`}
                    data-testid="task-item-title"
                  >
                    {task.title}
                  </span>
                  <span className="ml-2 text-xs px-2 py-0.5 rounded bg-gray-200 text-gray-700">
                    {task.priority}
                  </span>
                </div>
                <button
                  type="button"
                  data-testid="toggle-status-button"
                  onClick={() => handleToggle(task.id)}
                  className="text-xs px-3 py-1 rounded border hover:bg-gray-100"
                >
                  状態: <span data-testid="task-status">{task.status}</span> (変更)
                </button>
              </li>
            ))}
          </ul>
        </section>
      </div>
    </main>
  );
}
