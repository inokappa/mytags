default: ## ヘルプを表示する
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

depend: ## 依存パッケージの導入
	@dep init

test: ## test テストの実行
	@gom test -v

build: build ## バイナリをビルドする
	@./build.sh mytags.go

release: release ## バイナリをリリースする. 引数に `_VER=バージョン番号` を指定する.
	@ghr -u inokappa -r tagCtrl v${_VER} ./pkg/
