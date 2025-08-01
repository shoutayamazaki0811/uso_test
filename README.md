# uso - Premium Entertainment Platform

高級エンターテイメントマッチングプラットフォーム

## 🚀 デプロイ方法

### 1. GitHub Pages（フロントエンドのみ）

1. このリポジトリをGitHubにプッシュ
2. Settings → Pages → Source を "GitHub Actions" に設定
3. `https://[username].github.io/uso/` でアクセス可能

### 2. Vercel（フロントエンド + サーバーレス関数）

```bash
npm install -g vercel
vercel
```

### 3. Railway（フルスタック・データベース付き）

1. [Railway](https://railway.app)にサインアップ
2. GitHubリポジトリを接続
3. PostgreSQLを追加
4. 環境変数を設定

### 4. Render（フルスタック・無料プラン有り）

1. [Render](https://render.com)にサインアップ
2. Web ServiceとPostgreSQLを作成
3. 環境変数を設定

## 🛠️ ローカル開発

### フロントエンドのみ
```bash
# 任意のHTTPサーバーで起動
python -m http.server 8080
# または
npx serve
```

### フルスタック
```bash
# PostgreSQLを起動
docker run -p 5432:5432 -e POSTGRES_PASSWORD=password postgres

# データベース作成
createdb uso_db
psql uso_db < migrations/001_initial_schema.sql

# 環境変数設定
cp .env.example .env
# .envファイルを編集

# サーバー起動
go run cmd/server/main.go
```

## 📱 機能

- ✅ ユーザー登録・ログイン
- ✅ キャスト検索・フィルタリング
- ✅ 予約システム（Stripe決済）
- ✅ メッセージング
- ✅ レビュー・評価
- ✅ 管理画面
- ✅ レスポンシブデザイン

## 🔧 技術スタック

- **Backend**: Go, Gin, PostgreSQL
- **Frontend**: Vanilla JavaScript, CSS3
- **Payment**: Stripe
- **Email**: Resend
- **Auth**: JWT