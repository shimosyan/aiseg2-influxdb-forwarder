# AiSEG2 InfluxDB Forwarder

Panasonic製スマートホームコントローラー「AiSEG2」から電力メトリクスを取得し、InfluxDBに転送するGoアプリケーションです。

後述の通りすべての環境で動作を保証していないので自己責任です。

## 特徴

- **高い安定性**: Go言語による長期間稼働に適した実装（旧Node.js版のメモリリーク問題を解決）
- **軽量**: 単一バイナリによる軽量な実行環境
- **Docker対応**: コンテナ化による簡単なデプロイメント
- **構造化ログ**: JSON形式のログによる監視性向上

## 動作環境

このツールを使用するには、以下のいずれかの環境が必要です：

### Docker環境（推奨）

- Docker & Docker Compose
- Web操作が可能なAiSEG2端末

### ローカル実行環境

- Go 1.23以降
- Web操作が可能なAiSEG2端末

下記環境で動作確認をしています。作者自宅の機材でしか動作確認していないため、それ以外の環境の動作は保証できません。

- AiSEG2
  - 本体型番 `MKN713`
  - ファームウェア `Ver.2.97I-01`

## 機能概要

本ツールでは以下の機能をサポートしています。

- AiSEG2 から取得したメトリクスを指定の influxdb へ保存
- AiSEG2 から取得できる項目は以下のとおりです。
  - 消費電力の合計
  - 発電電力の合計
  - 売電電力（消費電力の合計と発電電力の合計の差）
  - AiSEG2 が認識している発電機器ごとの発電量（最大3つまで）
  - AiSEG2 の計測回路に設定された回路ごとの消費電力

## 導入方法

### 注意点

AiSEG2 はそのままでは HTTP しか喋らないため、LAN など境界内でツールを展開してください。

境界外からアクセスする場合は別途トンネリングやリバースプロキシなどをつかって経路の暗号化を実施してください。

### 事前準備

1. AiSEG2 の IPアドレスを固定化してください。
2. influxdb に格納先となる Bucket を用意してください。
3. Docker環境（推奨）またはGo言語実行環境を用意してください。

### 1. 環境変数の設定

`.env.example`をコピーして`.env`ファイルを作成し、環境に合わせて設定してください：

```bash
cp .env.example .env
```

`.env`ファイルの設定項目：

```bash
# AiSEG2 設定
AISEG2_HOST=192.168.1.100
AISEG2_USER=your_username
AISEG2_PASSWORD=your_password
AISEG2_USE_HTTPS=0

# InfluxDB 設定
INFLUXDB_HOST=192.168.1.200:8086
INFLUXDB_TOKEN=your_influxdb_token
INFLUXDB_ORG=your_organization
INFLUXDB_BUCKET=your_bucket
INFLUXDB_USE_HTTPS=0
```

**Docker環境でホストマシンのInfluxDBに接続する場合：**

ホストマシンで動作しているInfluxDBに接続する場合は、`INFLUXDB_HOST`を以下のように設定してください：

```bash
# Docker環境からホストマシンのInfluxDBに接続
INFLUXDB_HOST=host.docker.internal:8086
```

### 2. Docker Composeでの実行（推奨）

リポジトリをクローンします。

```bash
git clone https://github.com/shimosyan/aiseg2-influxdb-forwarder.git
cd aiseg2-influxdb-forwarder
```

Docker Composeで起動します。

```bash
# ビルドと起動
docker compose up -d
# または
docker-compose up -d

# ログ確認
docker compose logs -f
# または
docker-compose logs -f

# 停止
docker compose down
# または
docker-compose down
```

### 3. ローカルビルドでの実行

リポジトリをクローンします。

```bash
git clone https://github.com/shimosyan/aiseg2-influxdb-forwarder.git
cd aiseg2-influxdb-forwarder
```

依存関係をインストールしてビルドします。

```bash
# 依存関係のインストール
go mod download

# ビルド
go build -o aiseg2-forwarder ./cmd/aiseg2-forwarder

# 実行
./aiseg2-forwarder
```

## 監視とトラブルシューティング

### ログ確認

```bash
# Docker Composeの場合
docker compose logs -f aiseg2-forwarder
# または
docker-compose logs -f aiseg2-forwarder

# Dockerの場合
docker logs -f aiseg2-forwarder
```

### ヘルスチェック

```bash
# コンテナの状態確認
docker compose ps
# または
docker-compose ps
```

### よくある問題

1. **AiSEG2への接続エラー**
   - ネットワーク接続を確認
   - ホスト名とポート番号の確認
   - Digest認証の認証情報確認

2. **InfluxDBへの書き込みエラー**
   - InfluxDBサーバーの稼働状況確認
   - トークンと権限の確認
   - ネットワーク接続の確認
   - Docker環境の場合：
     - ホストマシンのInfluxDBに接続する場合は `INFLUXDB_HOST=host.docker.internal:8086` を使用
     - Linux環境で `host.docker.internal` が解決できない場合は、ホストマシンのIPアドレスを直接指定
     - `docker-compose.yml` に `extra_hosts` が設定されていることを確認

3. **コンテナが起動しない**
   - `.env`ファイルの設定確認
   - Docker Composeログの確認: `docker compose logs` または `docker-compose logs`

### Docker環境での特記事項

**ホストマシンへの接続について：**

Docker環境からホストマシン上のサービス（InfluxDBやAiSEG2）に接続する場合、以下の点にご注意ください：

1. **host.docker.internal の使用**

   ```bash
   # ホストマシンのInfluxDBに接続する場合
   INFLUXDB_HOST=host.docker.internal:8086
   ```

2. **Linux環境での対応**
   - `host.docker.internal` が解決できない場合は、ホストマシンのIPアドレスを確認：

   ```bash
   # ホストマシンのIPアドレスを確認
   ip route show default | awk '/default/ {print $3}'
   # または
   hostname -I | awk '{print $1}'
   ```

   - 確認したIPアドレスを直接指定：

   ```bash
   INFLUXDB_HOST=172.17.0.1:8086  # 例
   ```

3. **docker-compose.ymlの設定**
   - 本プロジェクトの `docker-compose.yml` には既に `extra_hosts` が設定済み
   - これにより `host.docker.internal` が正しく解決されます

## 旧バージョン（Node.js版）からの移行

このプロジェクトはNode.js + TypeScriptからGo言語に移行されました。主な変更点：

- **安定性の向上**: メモリリーク問題の解決
- **デプロイの簡素化**: 単一バイナリとコンテナ化
- **パフォーマンス向上**: 高速起動と低メモリ使用量
- **監視性の向上**: 構造化ログとヘルスチェック

旧バージョンの設定ファイル（.env）はそのまま利用可能です。
