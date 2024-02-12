# aiseg2-influxdb

パナソニック製の AiSEG2 から利用できる Web インターフェースから各情報をスクレイピングして influxdb に投入するツールです。

後述の通りすべての環境で動作を保証していないので自己責任です。

## 動作環境

このツールを使用するには、動作環境にて Node.js のインストール及び、Web操作が可能な AiSEG2端末が必要です。

下記環境で動作確認をしています。作者自宅の機材でしか動作確認していないため、それ以外の環境の動作は保証できません。

- 実行環境
  - Ubuntu `22.04`
  - Node.js `20.11.0`
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

1. ツールを動かすホスト環境に Git 及び Node.js を導入してください。
2. AiSEG2 の IPアドレスを固定化してください。
3. influxdb に格納先となる Bucket を用意してください。

### ツールのインストール

ホスト環境の適当な作業ディレクトリで本リポジトリをクローンします。

```sh
git clone https://github.com/shimosyan/aiseg2-influxdb.git
```

リポジトリ内に入ります。

```sh
cd ./aiseg2-influxdb
```

設定ファイル（`.env`）をサンプルファイルコピーしてご利用の環境に合わせて設定値を入れます。

```sh
cp .env.sample .env
```

### ツールの起動

以下のコマンドで起動することができます。

```sh
npm run script
```

### ツールのデーモン化について

下記コマンドで Node.js の [forever](https://www.npmjs.com/package/forever) ライブラリを使ったデーモン化ができます

```sh
# 開始
npm start

# 停止
npm stop
```

ただし、環境によってはうまく動かないので `npm run script` を `systemd` 化するなど適宜対応してください。

#### `systemd` のサンプルファイル

`/etc/systemd/system/aiseg2-influxdb.service`

```ini
[Unit]
Description=aiseg2-influxdb
After=syslog.target network.target

[Service]
Type=simple
ExecStart=/usr/bin/npm run script
WorkingDirectory=/home/ec2-user/aiseg2-influxdb
KillMode=process
Restart=always
User=ec2-user
Group=ec2-user

[Install]
WantedBy=multi-user.target
```
