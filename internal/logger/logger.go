package logger

import (
	"fmt"

	"go.uber.org/zap"
)

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize(level string) (*zap.Logger, error) {
	var zl *zap.Logger
	// Преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return zl, fmt.Errorf("Logger: Ошибка при иницизилации логгера: %w", err)
	}
	// Создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// Устанавливаем уровень
	cfg.Level = lvl
	// Создаём логер на основе конфигурации
	zl, err = cfg.Build()
	if err != nil {
		return zl, fmt.Errorf("Logger: Ошибка при иницизилации логгера: %w", err)
	}
	return zl, nil
}
