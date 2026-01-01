.PHONY: up down build logs restart clean env help test test-coverage test-coverage-html coverage-check coverage-summary clean-coverage

# Docker Composeを起動
up:
	docker compose up -d

# Docker Composeを停止
down:
	docker compose down

# Docker Composeをビルドして起動
build:
	docker compose up -d --build

# ログを表示
logs:
	docker compose logs -f app

# 再起動
restart:
	docker compose restart app

# コンテナとボリュームを完全削除
clean:
	docker compose down -v

# .envファイルを生成（既存の場合はスキップ）- Docker経由で実行
env:
	@docker run --rm -v $(PWD):/app -w /app alpine sh -c '\
		if [ ! -f .env ]; then \
			cp .env.example .env; \
			echo ".env file created from .env.example"; \
			echo "Please update the values in .env with your actual credentials"; \
		else \
			echo ".env file already exists. Remove it first if you want to regenerate."; \
		fi'

# .envファイルを強制的に再生成 - Docker経由で実行
env-force:
	@docker run --rm -v $(PWD):/app -w /app alpine sh -c 'cp .env.example .env && echo ".env file created from .env.example"'

# テストを実行
test:
	@echo "Running unit tests..."
	@go test -v ./internal/...

# テストをカバレッジ付きで実行
test-coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=coverage.out -covermode=atomic ./internal/... 2>&1 | grep -v "no such tool" || true
	@echo "\nCoverage by package:"
	@go test -cover ./internal/... 2>&1 | grep -E "^(ok|coverage:)" || true
	@echo "\nTotal coverage:"
	@go tool cover -func=coverage.out | grep total:

# HTMLカバレッジレポートを生成
test-coverage-html: test-coverage
	@echo "Generating HTML coverage report..."
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# カバレッジが80%以上かチェック
coverage-check: test-coverage
	@echo "Checking coverage threshold (80%)..."
	@COVERAGE=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	echo "Current coverage: $$COVERAGE%"; \
	if [ $$(echo "$$COVERAGE < 80" | bc -l) -eq 1 ]; then \
		echo "❌ Coverage $$COVERAGE% is below threshold 80%"; \
		exit 1; \
	else \
		echo "✅ Coverage $$COVERAGE% meets threshold 80%"; \
	fi

# カバレッジ情報を表示
coverage-summary: test-coverage
	@echo "\n=== Coverage Summary ==="
	@echo "By Package:"
	@go test -cover ./internal/... 2>&1 | grep -E "^(ok|coverage:)" || true
	@echo "\n=== Detailed Coverage ==="
	@go tool cover -func=coverage.out | grep -E "(config|handlers|middleware|services)" | sort -k3 -n
	@echo "\n=== Total Coverage ==="
	@go tool cover -func=coverage.out | grep total:

# カバレッジファイルを削除
clean-coverage:
	@echo "Cleaning coverage files..."
	@rm -f coverage.out coverage.html
	@echo "Coverage files removed"

# ヘルプ
help:
	@echo "使用可能なコマンド:"
	@echo ""
	@echo "Docker コマンド:"
	@echo "  make up                 - Docker Composeを起動"
	@echo "  make down               - Docker Composeを停止"
	@echo "  make build              - Docker Composeをビルドして起動"
	@echo "  make logs               - アプリケーションのログを表示"
	@echo "  make restart            - アプリケーションを再起動"
	@echo "  make clean              - コンテナとボリュームを完全削除"
	@echo ""
	@echo "テストコマンド:"
	@echo "  make test               - ユニットテストを実行"
	@echo "  make test-coverage      - カバレッジ付きでテストを実行"
	@echo "  make test-coverage-html - HTMLカバレッジレポートを生成"
	@echo "  make coverage-check     - カバレッジが80%以上かチェック"
	@echo "  make coverage-summary   - 詳細なカバレッジ情報を表示"
	@echo "  make clean-coverage     - カバレッジファイルを削除"
	@echo ""
	@echo "環境設定:"
	@echo "  make env                - .envファイルを生成（Docker経由、既存の場合はスキップ）"
	@echo "  make env-force          - .envファイルを強制的に再生成（Docker経由）"
	@echo ""
	@echo "  make help               - このヘルプを表示"

