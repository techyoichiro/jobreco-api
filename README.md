<div id="top"></div>
<img width="958" alt="headder" src="https://github.com/user-attachments/assets/21a55949-07bb-403e-ae2a-ba912e0298d2">

## 使用技術一覧

<!-- シールド一覧 -->
<p style="display: inline">
  <!-- フロントエンドのフレームワーク・ライブラリ一覧 -->
<!--   <img src="https://github.com/user-attachments/assets/c0d80d44-1c5c-4e61-884d-b31f11f037f5">
  <img src="https://github.com/user-attachments/assets/a5edb846-818e-4b0c-a659-f7c9993ed82c">
  <img src="https://github.com/user-attachments/assets/fb5866fb-a538-496b-8565-c728b1c570dd">
  <!-- フロントエンド言語一覧 -->
<!--   <img src="https://github.com/user-attachments/assets/a64bf638-dd8b-4af6-8474-c828f0af07ae"> -->
  <!-- バックエンドのフレームワーク一覧 -->
  <img src="https://github.com/user-attachments/assets/dcc5c9ce-b233-489d-aa39-b3b61fd8a8a4">
  <!-- バックエンド言語 -->
  <img src="https://github.com/user-attachments/assets/24936e52-d8fc-4617-a3c5-b8353005a121">
  <!-- DB -->
  <img src="[https://img.shields.io/badge/-MySQL-4479A1.svg?logo=mysql&style=for-the-badge&logoColor=white](https://github.com/user-attachments/assets/940060f8-2d18-43fe-a0f6-78012a0105ea)">
  <!-- インフラ一覧 -->
  <img src="https://github.com/user-attachments/assets/7ea80ddd-c773-4032-b740-82c44e80eec1">
  <img src="https://github.com/user-attachments/assets/2e9dc12a-5073-4663-aec9-f4fb821f480a">
  <img src="https://github.com/user-attachments/assets/6288c698-5232-4eb7-86a9-8089f7388967">
 
</p>

## 目次

1. [プロジェクトについて](#プロジェクトについて)
2. [環境](#環境)
3. [今後の展望](#今後の展望)

<!-- プロジェクト名を記載 -->

<!-- プロジェクトの概要を記載 -->

<!-- プロジェクトについて -->

## プロジェクトについて
飲食店の勤怠管理アプリ用API
mainブランチにマージされることでGitHub ActionsでDocekrイメージがビルド、Cloud Runにプッシュされる

<img src="https://github.com/user-attachments/assets/68e0cf5a-5ac1-4396-bf10-49e1b2446b89">


## 技術選定理由
APIについてはGo+ginで構築した。
個人開発ということで一度触れたことのあるGoとginでAPIを構成した。
サーバーはCloud Run、DBはSupabaseを使用している。

### バックエンド
フレームワークはGoとEchoで迷ったが、フレームワークの中で一番スター(※世の中で使われる/使う機会)の多いGinを選択した。

### 気づいたこと/工夫したこと
一つの建物に店舗が2つあり、1日にその両方で勤務する従業員がいる。
そのことを記録として残しつつ、テーブルのカラムとロジックをできるだけシンプルに構成した。テーブル数を増やして管理したりする方法も考えたが、システムの規模や保守性を考えてこの構成とした。

工夫したこととしては、フロントエンドとバックエンドを完全に分離することでAPIを新しいフレームワークや技術にリプレイス（学習）したい時に変更を容易にできるようにした。
また小規模な利用ではあるがパスワードを平文のままDBに保存するのはセキュリティ的にも良くないと考え、cryptoライブラリを使用してハッシュ化して保存・復号化して比較を行っている。

## 環境

<!-- 言語、フレームワーク、ミドルウェア、インフラの一覧とバージョンを記載 -->

| 言語・フレームワーク・ライブラリ  | バージョン |
| --------------------- | ---------- |
| Go               | 1.23    |
| Gin                | 1.10.0      |

その他のパッケージのバージョンは go.mod と package.json を参照してください

<!-- コンテナの作成方法、パッケージのインストール方法など、開発環境構築に必要な情報を記載 -->

## 今後の展望
今後は給与明細の出力とシフトの管理機能を追加していく予定。

<p align="right">(<a href="#top">トップへ</a>)</p>
