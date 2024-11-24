package smartContext

import (
	"test_task2/infrastructure/logger"

	"gorm.io/gorm"
)

type SmartContext struct {
	db        *gorm.DB
	logger    *logger.Logger
	apiServer string
}

// NewSmartContext создает новый контекст приложения
func NewSmartContext(db *gorm.DB, logger *logger.Logger, apiServer string) *SmartContext {
	return &SmartContext{db: db, logger: logger, apiServer: apiServer}
}

// GetDB возвращает подключение к базе данных
func (c *SmartContext) GetDB() *gorm.DB {
	return c.db
}

// Infof логирует информацию
func (c *SmartContext) Infof(msg string, keysAndValues ...interface{}) {
	c.logger.Infof(msg, keysAndValues...)
}

// Debug логирует отладочные сообщения
func (c *SmartContext) Debugf(msg string, keysAndValues ...interface{}) {
	c.logger.Debugf(msg, keysAndValues...)
}

// Error логирует ошибки
func (c *SmartContext) Errorf(msg string, keysAndValues ...interface{}) {
	c.logger.Errorf(msg, keysAndValues...)
}

// GetAPIServer возвращает URL API-сервера
func (c *SmartContext) GetAPIServer() string {
	return c.apiServer
}
