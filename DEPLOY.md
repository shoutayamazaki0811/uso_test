# USO デプロイガイド

このガイドでは、USOアプリケーションを無料でホスティングする方法を説明します。

## 🚀 デプロイオプション

### 1. GitHub Pages（フロントエンドのみ）

最も簡単な方法です。静的なフロントエンドのみをホスティングします。

**手順:**
1. GitHubにリポジトリを作成
2. コードをプッシュ
3. Settings → Pages → Source で "GitHub Actions" を選択
4. `https://[username].github.io/uso/` でアクセス可能

### 2. Vercel（フロントエンド推奨）

**手順:**
1. [Vercel](https://vercel.com)にサインアップ
2. GitHubアカウントを連携
3. リポジトリをインポート
4. デプロイ完了！

**特徴:**
- 自動HTTPS
- カスタムドメイン対応
- 自動デプロイ

### 3. Netlify（フロントエンド）

**手順:**
1. [Netlify](https://netlify.com)にサインアップ
2. GitHubリポジトリを連携
3. ビルド設定はデフォルトのまま
4. デプロイ完了！

**特徴:**
- 月100GBの帯域幅無料
- フォーム機能
- 関数機能（制限あり）

### 4. Railway（フルスタック推奨）

**手順:**
1. [Railway](https://railway.app)にサインアップ
2. New Project → Deploy from GitHub repo
3. PostgreSQLを追加
4. 環境変数を設定:
   ```
   DATABASE_URL=（自動設定）
   JWT_SECRET=your-secret-key
   ADMIN_PASSWORD=your-admin-pass
   ```
5. デプロイ完了！

**特徴:**
- 月$5の無料クレジット
- PostgreSQL付き
- 自動スケーリング

### 5. Render（フルスタック）

**手順:**
1. [Render](https://render.com)にサインアップ
2. New → Web Service
3. GitHubリポジトリを接続
4. 環境変数を設定
5. New → PostgreSQL でデータベース作成
6. データベースURLを環境変数に追加

**特徴:**
- 無料プランあり（スリープあり）
- PostgreSQL無料
- 自動デプロイ

## 🔧 環境変数

すべてのフルスタックデプロイで必要な環境変数:

```env
DATABASE_URL=postgres://...
JWT_SECRET=your-super-secret-key
ADMIN_PASSWORD=secure-admin-password
PORT=8080

# オプション（メール通知を使う場合）
RESEND_API_KEY=re_xxxxx
FROM_EMAIL=noreply@uso.app

# オプション（決済を使う場合）
STRIPE_SECRET_KEY=sk_test_xxxxx
STRIPE_WEBHOOK_SECRET=whsec_xxxxx
```

## 📱 デプロイ後の確認

1. トップページが表示されるか
2. スワイプ機能が動作するか
3. ナビゲーションが機能するか
4. レスポンシブデザインが適用されているか

## 🆘 トラブルシューティング

### ビルドエラー
- Go modulesが正しくダウンロードされているか確認
- go.modファイルが存在するか確認

### データベース接続エラー
- DATABASE_URLが正しく設定されているか確認
- SSLモードの設定を確認（`?sslmode=require`）

### 画像が表示されない
- 画像パスが正しいか確認
- 静的ファイルの配信設定を確認

## 🎉 完了！

デプロイが完了したら、URLを共有して使い始めましょう！

質問がある場合は、Issueを作成してください。