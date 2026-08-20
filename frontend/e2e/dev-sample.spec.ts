import { test, expect } from '@playwright/test';

/**
 * @description DevSample ページのブラウザ実動作・結合テスト（Code-as-Docs）
 */
test.describe('DevSample ブラウザ結合テスト (E2E)', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/dev-sample');
  });

  test('シナリオ1: ページが正常に表示され、初期要素が確認できること', async ({ page }) => {
    // ページタイトルの確認
    const title = page.getByTestId('page-title');
    await expect(title).toBeVisible();
    await expect(title).toHaveText('開発テスト用サンプルページ (Dev Sample)');

    // 初期タスクの存在確認
    const initialItem = page.getByTestId('task-item').first();
    await expect(initialItem).toContainText('初期タスク 1');
    await expect(page.getByTestId('stat-total')).toHaveText('1');
  });

  test('シナリオ2: 新規タスクの追加フローがブラウザ上で正常に動作すること', async ({ page }) => {
    // フォームに入力
    const input = page.getByTestId('task-title-input');
    const select = page.getByTestId('task-priority-select');
    const submitBtn = page.getByTestId('add-task-button');

    await input.fill('Playwright E2E テストタスク');
    await select.selectOption('high');
    await submitBtn.click();

    // タスク一覧に新しいタスクが表示されること
    const taskItems = page.getByTestId('task-item');
    await expect(taskItems).toHaveCount(2);

    const newTask = taskItems.nth(1);
    await expect(newTask).toContainText('Playwright E2E テストタスク');
    await expect(newTask).toContainText('high');

    // 統計値が更新されること
    await expect(page.getByTestId('stat-total')).toHaveText('2');
    await expect(page.getByTestId('stat-todo')).toHaveText('2');
  });

  test('シナリオ3: タスクのステータス更新操作と完了率計算の反映', async ({ page }) => {
    const toggleButton = page.getByTestId('toggle-status-button').first();

    // 1回クリック: todo -> in_progress
    await toggleButton.click();
    await expect(page.getByTestId('task-status').first()).toHaveText('in_progress');
    await expect(page.getByTestId('stat-in-progress')).toHaveText('1');

    // 2回クリック: in_progress -> done
    await toggleButton.click();
    await expect(page.getByTestId('task-status').first()).toHaveText('done');
    await expect(page.getByTestId('stat-done')).toHaveText('1');
    await expect(page.getByTestId('stat-completion-rate')).toHaveText('100%');
  });
});
