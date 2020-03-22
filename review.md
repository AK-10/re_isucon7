# isucon7 振り返り

## 点数遷移
初期 4000 ~ 5500
最終 110000 ~ 150000程度


## やったこと
- alp導入
- messageJsonifyで発生しているN+1を修正
- dbにある画像をファイルの読み出しに変更，またnginx側でキャッシュするようにした
- historyのN+1を解消
- indexをはった
- 画像をgo-cacheでキャッシュするようにした(不採用)
- deploy scriptを書いた

## やり損ねたこと
- redisを使う
- 
