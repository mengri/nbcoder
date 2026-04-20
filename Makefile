.PHONY: all build frontend backend clean run help

# 变量定义
FRONTEND_DIR := web
BACKEND_DIR := backend
FRONTEND_BUILD_DIR := $(FRONTEND_DIR)/dist
EMBEDDED_DIR := $(FRONTEND_DIR)/embedded/dist
BACKEND_BUILD_DIR := $(BACKEND_DIR)/bin
BINARY_NAME := nbcoder-server

# 默认目标
all: build

# 前端构建
frontend:
	@echo "Building frontend..."
	cd $(FRONTEND_DIR) && npx vite build
	@echo "Copying frontend to embedded directory..."
	cp -r $(FRONTEND_BUILD_DIR) $(EMBEDDED_DIR)
	@echo "Frontend build complete"

# 后端构建
backend:
	@echo "Building backend..."
	cd $(BACKEND_DIR) && mkdir -p $(BACKEND_BUILD_DIR)
	cd $(BACKEND_DIR) && go build -o $(BACKEND_BUILD_DIR)/$(BINARY_NAME) ./cmd/server
	@echo "Backend build complete"

# 完整构建（前端 + 后端）
build: frontend backend
	@echo "Build complete!"

# 清理构建产物
clean:
	@echo "Cleaning..."
	rm -rf $(FRONTEND_BUILD_DIR)
	rm -rf $(EMBEDDED_DIR)
	rm -rf $(BACKEND_BUILD_DIR)
	@echo "Clean complete"

# 运行服务器
run:
	@echo "Starting server..."
	cd $(BACKEND_DIR) && ./bin/$(BINARY_NAME)

# 开发模式（仅后端，假设前端已构建）
dev-backend:
	@echo "Starting backend in development mode..."
	cd $(BACKEND_DIR) && go run ./cmd/server

# 安装前端依赖
install-frontend:
	@echo "Installing frontend dependencies..."
	cd $(FRONTEND_DIR) && npm install

# 前端开发模式
dev-frontend:
	@echo "Starting frontend in development mode..."
	cd $(FRONTEND_DIR) && npm run dev

# 帮助信息
help:
	@echo "Available targets:"
	@echo "  all              - Build both frontend and backend"
	@echo "  build            - Build both frontend and backend (same as all)"
	@echo "  frontend         - Build frontend only"
	@echo "  backend          - Build backend only"
	@echo "  clean            - Clean all build artifacts"
	@echo "  run              - Run the server"
	@echo "  dev-backend      - Run backend in development mode"
	@echo "  dev-frontend     - Run frontend in development mode"
	@echo "  install-frontend - Install frontend dependencies"
	@echo "  help             - Show this help message"
