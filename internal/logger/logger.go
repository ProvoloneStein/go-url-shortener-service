package logger

import (
	"fmt"
	"go.uber.org/zap"
)

// Initialize инициализирует синглтон логера с необходимым уровнем логирования.
func Initialize(level string) (*zap.Logger, error) {
	var zl *zap.Logger
	// преобразуем текстовый уровень логирования в zap.AtomicLevel
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return zl, fmt.Errorf("Ошибка при иницизилации логгера: %s", err)
	}
	// создаём новую конфигурацию логера
	cfg := zap.NewProductionConfig()
	// устанавливаем уровень
	cfg.Level = lvl
	// создаём логер на основе конфигурации
	zl, err = cfg.Build()
	if err != nil {
		return zl, fmt.Errorf("Ошибка при иницизилации логгера: %s", err)
	}
	return zl, nil
}
