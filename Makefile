.PHONY: help up down backend frontend logs clean install restart rebuild stop

# デフォルトターゲット
help:
	@echo "使用可能なコマンド:"
	@echo "  make up          - バックエンドとフロントエンドを起動"
	@echo "  make down        - すべてのコンテナを停止"
	@echo "  make stop        - すべてのコンテナを強制停止"
	@echo "  make restart     - すべてのサービスを再起動"
	@echo "  make rebuild     - すべてを停止・再ビルド・起動"
	@echo "  make backend     - バックエンドのみを起動"
	@echo "  make frontend    - フロントエンドのみを起動"
	@echo "  make logs        - すべてのログを表示"
	@echo "  make logs-backend  - バックエンドのログを表示"
	@echo "  make logs-frontend - フロントエンドのログを表示"
	@echo "  make clean       - すべてのコンテナとボリュームを削除"
	@echo "  make install     - フロントエンドの依存関係をインストール"
	@echo "  make status      - サービスのステータスを確認"

# すべてのサービスを起動
up:
	@echo "Starting all services..."
	@cd backend && docker compose up -d
	@cd frontend && docker compose up -d
	@echo "Services started successfully!"
	@echo "Backend: http://localhost:8080"
	@echo "Frontend: http://localhost:5173"

# すべてのサービスを停止
down:
	@echo "Stopping all services..."
	@cd backend && docker compose down
	@cd frontend && docker compose down
	@echo "Services stopped successfully!"

# すべてのコンテナを強制停止
stop:
	@echo "Force stopping all services..."
	@cd backend && docker compose down --remove-orphans
	@cd frontend && docker compose down --remove-orphans
	@echo "All services stopped!"

# すべてのサービスを再起動（停止→起動）
restart:
	@echo "Restarting all services..."
	@$(MAKE) down
	@$(MAKE) up

# すべてを停止・再ビルド・起動
rebuild:
	@echo "Rebuilding all services..."
	@$(MAKE) stop
	@echo "Building backend..."
	@cd backend && docker compose build --no-cache
	@echo "Building frontend..."
	@cd frontend && docker compose build --no-cache
	@echo "Starting all services..."
	@$(MAKE) up
	@echo "Rebuild completed!"

# バックエンドのみを起動
backend:
	@echo "Starting backend..."
	@cd backend && docker compose up -d
	@echo "Backend started: http://localhost:8080"

# フロントエンドのみを起動
frontend:
	@echo "Starting frontend..."
	@cd frontend && docker compose up -d
	@echo "Frontend started: http://localhost:5173"

# すべてのログを表示
logs:
	@echo "Showing logs for all services (Ctrl+C to exit)..."
	@docker compose -f backend/compose.yml -f frontend/compose.yml logs -f

# バックエンドのログを表示
logs-backend:
	@echo "Showing backend logs (Ctrl+C to exit)..."
	@cd backend && docker compose logs -f

# フロントエンドのログを表示
logs-frontend:
	@echo "Showing frontend logs (Ctrl+C to exit)..."
	@cd frontend && docker compose logs -f

# すべてのコンテナとボリュームを削除
clean:
	@echo "Cleaning up all containers and volumes..."
	@cd backend && docker compose down -v
	@cd frontend && docker compose down -v
	@echo "Cleanup completed!"

# フロントエンドの依存関係をインストール
install:
	@echo "Installing frontend dependencies..."
	@cd frontend && docker compose exec frontend npm install
	@echo "Dependencies installed successfully!"

# バックエンドを再起動
restart-backend:
	@echo "Restarting backend..."
	@cd backend && docker compose restart
	@echo "Backend restarted!"

# フロントエンドを再起動
restart-frontend:
	@echo "Restarting frontend..."
	@cd frontend && docker compose restart
	@echo "Frontend restarted!"

# バックエンドのビルド
build-backend:
	@echo "Building backend..."
	@cd backend && docker compose build
	@echo "Backend built successfully!"

# フロントエンドのビルド
build-frontend:
	@echo "Building frontend..."
	@cd frontend && docker compose build
	@echo "Frontend built successfully!"

# すべてのサービスをビルド
build: build-backend build-frontend

# DynamoDBのテーブルを確認
dynamodb-tables:
	@echo "Listing DynamoDB tables..."
	@docker exec -it go-google-auth-dynamodb aws dynamodb list-tables --endpoint-url http://localhost:8000 --region ap-northeast-1

# ステータス確認
status:
	@echo "=== Backend Services ==="
	@cd backend && docker compose ps
	@echo ""
	@echo "=== Frontend Services ==="
	@cd frontend && docker compose ps
