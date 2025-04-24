.PHONY: swag

swag:
	@echo "🔄 Генерация Swagger документации..."
	swag init -g cmd/api/main.go --output docs
	@echo "✅ Документация сгенерирована в папке docs/"