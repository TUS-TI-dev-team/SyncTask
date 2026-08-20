import React from 'react';
import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import DevSamplePage from '../page';

/**
 * @description DevSamplePage コンポーネントの単体・インタラクションテスト（Code-as-Docs）
 */
describe('DevSamplePage UI コンポーネント仕様', () => {
  it('初期レンダリング: タイトル、初期タスク、初期統計が正しく表示されること', () => {
    render(<DevSamplePage />);

    expect(screen.getByTestId('page-title')).toHaveTextContent('開発テスト用サンプルページ (Dev Sample)');
    expect(screen.getByTestId('stat-total')).toHaveTextContent('1');
    expect(screen.getByTestId('stat-todo')).toHaveTextContent('1');
    expect(screen.getByTestId('stat-done')).toHaveTextContent('0');
    expect(screen.getByTestId('task-item-title')).toHaveTextContent('初期タスク 1');
  });

  it('タスク追加: フォームに入力して追加ボタンを押すとタスク一覧と統計が更新されること', async () => {
    const user = userEvent.setup();
    render(<DevSamplePage />);

    const input = screen.getByTestId('task-title-input');
    const addButton = screen.getByTestId('add-task-button');

    await user.type(input, '新規コンポーネントテストタスク');
    await user.click(addButton);

    // 一覧に2件あること
    const taskItems = screen.getAllByTestId('task-item');
    expect(taskItems).toHaveLength(2);
    expect(taskItems[1]).toHaveTextContent('新規コンポーネントテストタスク');

    // 統計の合計が2になること
    expect(screen.getByTestId('stat-total')).toHaveTextContent('2');
    expect(screen.getByTestId('stat-todo')).toHaveTextContent('2');

    // 入力欄がクリアされていること
    expect(input).toHaveValue('');
  });

  it('バリデーションエラー: タイトル空欄で追加しようとするとエラーメッセージが表示されること', async () => {
    const user = userEvent.setup();
    render(<DevSamplePage />);

    const addButton = screen.getByTestId('add-task-button');
    await user.click(addButton);

    expect(screen.getByTestId('error-message')).toHaveTextContent('Title cannot be empty');
  });

  it('ステータストグル: 変更ボタンを押すとステータスが遷移し統計が反映されること', async () => {
    const user = userEvent.setup();
    render(<DevSamplePage />);

    const toggleButton = screen.getByTestId('toggle-status-button');
    expect(screen.getByTestId('task-status')).toHaveTextContent('todo');

    // 1回目クリック: todo -> in_progress
    await user.click(toggleButton);
    expect(screen.getByTestId('task-status')).toHaveTextContent('in_progress');
    expect(screen.getByTestId('stat-todo')).toHaveTextContent('0');
    expect(screen.getByTestId('stat-in-progress')).toHaveTextContent('1');

    // 2回目クリック: in_progress -> done
    await user.click(toggleButton);
    expect(screen.getByTestId('task-status')).toHaveTextContent('done');
    expect(screen.getByTestId('stat-in-progress')).toHaveTextContent('0');
    expect(screen.getByTestId('stat-done')).toHaveTextContent('1');
    expect(screen.getByTestId('stat-completion-rate')).toHaveTextContent('100%');
  });
});
