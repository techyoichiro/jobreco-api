# Goバージョン
FROM golang:1.23-alpine

# コンテナ内の作業ディレクトリを設定
WORKDIR /app

# 必要なパッケージをインストール
RUN apk update && apk add --no-cache git

# タイムゾーン設定
RUN apk add --no-cache tzdata && \
    cp /usr/share/zoneinfo/Asia/Tokyo /etc/localtime && \
    echo "Asia/Tokyo" > /etc/timezone

# go.modとgo.sumを最初にコピーして依存関係をダウンロード
COPY go.mod go.sum ./
RUN go mod download

# .envファイルをコンテナにコピー
COPY .env ./

# ローカルのソースコードをコンテナにコピー
COPY . .

# アプリケーションをビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main ./cmd/main.go

# airをインストール
RUN go install github.com/air-verse/air@latest

# 公開予定のコンテナのポートを明示
EXPOSE 8080

# ホットリロードを実行する（airを使う場合）
CMD ["air", "-c", ".air.toml"]

# バイナリを直接実行する場合
# CMD ["./main"]
