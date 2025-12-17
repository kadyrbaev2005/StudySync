package stress

import (
	"runtime"
	"sync"

	"github.com/kadyrbayev2005/studysync/internal/services"
)

// GlobalStressWG глобальная WaitGroup для отслеживания активных горутин стресс-теста
var GlobalStressWG sync.WaitGroup

// Глобальные переменные для управления стресс-тестом
var (
	stopChan   chan struct{} // Канал для сигнала остановки
	stopOnce   sync.Once     // Гарантирует однократную отправку сигнала остановки
	isRunning  bool          // Флаг состояния стресс-теста
	runningMux sync.RWMutex  // Mutex для безопасного чтения/записи флага
)

// init инициализирует глобальные переменные пакета
func init() {
	stopChan = make(chan struct{})
}

// HeavyWork выполняет интенсивные вычисления для создания нагрузки на CPU
func HeavyWork() {
	sum := 0
	for i := 0; i < 100_000_000; i++ {
		sum += i
	}
	// Предотвращаем оптимизацию компилятора
	_ = sum
}

// StartStress запускает нагрузку на CPU для профилирования
func StartStress() {
	runningMux.Lock()
	defer runningMux.Unlock()

	if isRunning {
		services.Info("CPU stress test is already running")
		return
	}

	isRunning = true
	numCPU := runtime.NumCPU()
	services.Info("Starting CPU stress test", "goroutines", numCPU)

	// Сбрасываем канал остановки на новый
	stopOnce = sync.Once{}
	stopChan = make(chan struct{})

	for i := 0; i < numCPU; i++ {
		GlobalStressWG.Add(1)
		go func(id int) {
			defer GlobalStressWG.Done()
			services.Info("CPU stress goroutine started", "id", id)

			for {
				select {
				case <-stopChan:
					services.Info("CPU stress goroutine stopped", "id", id)
					return
				default:
					HeavyWork()
				}
			}
		}(i)
	}
}

// StopStress останавливает нагрузку на CPU через канал
func StopStress() {
	runningMux.Lock()
	defer runningMux.Unlock()

	if !isRunning {
		services.Info("CPU stress test is not running")
		return
	}

	isRunning = false
	services.Info("Stopping CPU stress test")

	// Отправляем сигнал остановки всем горутинам
	stopOnce.Do(func() {
		close(stopChan)
	})
}

// IsRunning возвращает true если стресс-тест запущен
func IsRunning() bool {
	runningMux.RLock()
	defer runningMux.RUnlock()
	return isRunning
}
