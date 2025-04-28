.PHONY: swag, test

swag:
	@echo "🔄 Генерация Swagger документации..."
	swag init -g cmd/api/main.go --output docs
	@echo "✅ Документация сгенерирована в папке docs/"

test:
	@echo "🔄 Запуск тестов..."
	go test -v ./...
	#go test -count=1 -v ./...
	@echo "✅ Тесты пройдены..."