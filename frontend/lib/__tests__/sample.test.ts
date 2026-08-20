import { describe, it, expect } from 'vitest';
import {
  createTask,
  toggleTaskStatus,
  filterTasksByStatus,
  calculateTaskStats,
  Task,
} from '../sample';

/**
 * @description sample.ts の仕様テスト（Code-as-Docs）
 */
describe('sample.ts - タスク管理ドメインロジック仕様', () => {
  describe('createTask() 仕様', () => {
    it('正常系: 有効なタイトルと優先度で初期状態 todo のタスクが生成されること', () => {
      const task = createTask('テストタスク', 'high');

      expect(task.title).toBe('テストタスク');
      expect(task.status).toBe('todo');
      expect(task.priority).toBe('high');
      expect(task.id).toMatch(/^task-/);
      expect(task.createdAt).toBeDefined();
    });

    it('正常系: 優先度を省略した場合はデフォルトで medium になること', () => {
      const task = createTask('デフォルト優先度タスク');
      expect(task.priority).toBe('medium');
    });

    it('正常系: タイトルの前後の空白文字がトリムされること', () => {
      const task = createTask('   余白付きタスク   ');
      expect(task.title).toBe('余白付きタスク');
    });

    it('異常系: 空文字または空白のみのタイトルが渡された場合は Error を送出すること', () => {
      expect(() => createTask('')).toThrowError('Title cannot be empty');
      expect(() => createTask('   ')).toThrowError('Title cannot be empty');
    });
  });

  describe('toggleTaskStatus() 仕様', () => {
    it('todo -> in_progress -> done -> todo の順序でステータスが循環遷移すること', () => {
      const initialTask: Task = {
        id: 't-1',
        title: 'ステータス遷移テスト',
        status: 'todo',
        priority: 'medium',
        createdAt: new Date().toISOString(),
      };

      const step1 = toggleTaskStatus(initialTask);
      expect(step1.status).toBe('in_progress');

      const step2 = toggleTaskStatus(step1);
      expect(step2.status).toBe('done');

      const step3 = toggleTaskStatus(step2);
      expect(step3.status).toBe('todo');
    });

    it('イミュータビリティ: 元のタスクオブジェクトを変更せず新しいオブジェクトを返すこと', () => {
      const original: Task = {
        id: 't-1',
        title: '不変性テスト',
        status: 'todo',
        priority: 'medium',
        createdAt: new Date().toISOString(),
      };

      const updated = toggleTaskStatus(original);
      expect(original.status).toBe('todo');
      expect(updated).not.toBe(original);
    });
  });

  describe('filterTasksByStatus() 仕様', () => {
    const mockTasks: Task[] = [
      { id: '1', title: 'Task 1', status: 'todo', priority: 'low', createdAt: '' },
      { id: '2', title: 'Task 2', status: 'in_progress', priority: 'medium', createdAt: '' },
      { id: '3', title: 'Task 3', status: 'done', priority: 'high', createdAt: '' },
    ];

    it('指定されたステータスのタスクのみをフィルタリングすること', () => {
      const result = filterTasksByStatus(mockTasks, 'in_progress');
      expect(result).toHaveLength(1);
      expect(result[0].id).toBe('2');
    });

    it('all または undefined が指定された場合は全タスクを返すこと', () => {
      expect(filterTasksByStatus(mockTasks, 'all')).toHaveLength(3);
      expect(filterTasksByStatus(mockTasks)).toHaveLength(3);
    });
  });

  describe('calculateTaskStats() 仕様', () => {
    it('タスクが空配列の場合は全件数と完了率が 0 であること', () => {
      const stats = calculateTaskStats([]);
      expect(stats).toEqual({
        total: 0,
        todo: 0,
        inProgress: 0,
        done: 0,
        completionRate: 0,
      });
    });

    it('各種ステータス件数および完了率（%）が正しく計算されること', () => {
      const mockTasks: Task[] = [
        { id: '1', title: 'Task 1', status: 'todo', priority: 'low', createdAt: '' },
        { id: '2', title: 'Task 2', status: 'in_progress', priority: 'medium', createdAt: '' },
        { id: '3', title: 'Task 3', status: 'done', priority: 'high', createdAt: '' },
        { id: '4', title: 'Task 4', status: 'done', priority: 'high', createdAt: '' },
      ];

      const stats = calculateTaskStats(mockTasks);
      expect(stats.total).toBe(4);
      expect(stats.todo).toBe(1);
      expect(stats.inProgress).toBe(1);
      expect(stats.done).toBe(2);
      expect(stats.completionRate).toBe(50.0);
    });
  });
});
