# マルチステージビルドを使用して軽量なイメージを作成
FROM golang:1.25-alpine AS builder

# 必要なパッケージをインストール
RUN apk add --no-cache git ca-certificates

# 作業ディレクトリを設定
WORKDIR /app

# go.modとgo.sumをコピーして依存関係をダウンロード
COPY go.mod go.sum ./
RUN go mod download

# ソースコードをコピー
COPY . .

# バイナリをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o aiseg2-forwarder ./cmd/aiseg2-forwarder

# 実行用の軽量イメージ
FROM alpine:latest

# 必要なパッケージをインストール
RUN apk --no-cache add ca-certificates tzdata

# タイムゾーンを設定（必要に応じて変更）
ENV TZ=Asia/Tokyo

# 非rootユーザーを作成
RUN adduser -D -s /bin/sh appuser

# 作業ディレクトリを設定
WORKDIR /app

# ビルドしたバイナリをコピー
COPY --from=builder /app/aiseg2-forwarder .

# バイナリに実行権限を付与
RUN chmod +x ./aiseg2-forwarder

# ユーザーを切り替え
USER appuser

# ヘルスチェック用のスクリプト（プロセス確認）
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD pgrep -f aiseg2-forwarder || exit 1

# アプリケーションを実行
CMD ["./aiseg2-forwarder"]
