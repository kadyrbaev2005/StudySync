package api

// StressManager интерфейс для управления стресс-тестом
type StressManager interface {
	StartStress()
	StopStress()
}

// GlobalStressManager глобальный менеджер стресс-теста
var GlobalStressManager StressManager

// SetStressManager устанавливает глобальный менеджер стресс-теста
func SetStressManager(manager StressManager) {
	GlobalStressManager = manager
}
