## notion-meeting-app
### できること
1. notionに議事録のテンプレートページを定期的に自動作成できる
2. 作成するテンプレートページの内容・場所・間隔を簡単にカスタマイズできる
3. いつでも自動作成するのを停止できます。自動作成の軌道も同様に簡単にできる
### 使い方
1. ショートカット`Register the notion info`を選択する
2. 表示されたmodalの項目を埋める
```
1. notionのtoken
2. notionのページを作成する場所
3. notionのテンプレートページの内容
4. テンプレートページを作成する曜日
```
3. `/start`のslash commandでスケジューラを起動させる

※スケジューラを停止させたい場合は``/stop``のslash commandで停止させる。
### notion情報登録の際の、ページ内容の作り方
- 1行目はページのタイトルに該当し、2行目以降はページの中身に該当する。
- 文章・heading1・heading2・heading3・リスト・番号付きリスト・todoリスト・toggleリストなどのテキストタイプが使える。
- テキストタイプを変更する場合には `Enterキー`で改行する必要がある。
- ページ内容を記入する際には以下の形式である必要がある。
```
`テキストタイプ` `テキスト`...
```
以下にテキストタイプの種類を示す。
|text type| about|
|:---|:---|
| paragraph | 文章 |
| heading_1 | 1番大きいheader |
| heading_2 | 2番大きいheader |
| heading_3 | 3番大きいheader |
| bulleted_list_item | リスト|
| numbered_list_item | 番号つきリスト |
| to_do | todoリスト |
| toggle | toggleリスト |

以下にページ内容の登録例と実際のnotionの画面を示す。
```
title ミーティング
heading_2 概要
タイトル:
heading_2 内容
numbered_list_item sample
numbered_list_item sample
heading_2 次回
次回の日程:
書記: sample, リード: sample, 計測係: sample
```