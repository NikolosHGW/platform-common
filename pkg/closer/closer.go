package closer

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var privateGlobalCloser = New()

// Add - функция-обёртка над closer.add().
// Добавляет колбэк для отложенного до завершения программы вызова.
func Add(f func() error) {
	privateGlobalCloser.add(f)
}

// Wait - функция-обёртка над close.wait().
// Заставляет процесс ждать, пока не завершатся все отложенные функции.
func Wait() {
	privateGlobalCloser.wait()
}

// CloseAll - функция-обёртка над close.closeAll().
// Завершает все отложенные функции.
func CloseAll() {
	privateGlobalCloser.closeAll()
}

type closer struct {
	mu     sync.Mutex
	once   sync.Once
	funcs  []func() error
	doneCh chan struct{}
}

var (
	closerInstance *closer
	once           sync.Once
)

// New - конструктор для closer.
func New() *closer {
	once.Do(func() {
		closerInstance = &closer{doneCh: make(chan struct{})}

		go func() {
			signalChan := make(chan os.Signal, 1)
			signal.Notify(signalChan, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

			<-signalChan
			log.Print("Получен сигнал завершения, приложение завершается...")
			closerInstance.closeAll()
		}()

	})

	return closerInstance
}

func (c *closer) add(f func() error) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f)
	c.mu.Unlock()
}

func (c *closer) wait() {
	<-c.doneCh
}

func (c *closer) closeAll() {
	c.once.Do(func() {
		defer close(c.doneCh)

		c.mu.Lock()
		defer c.mu.Unlock()

		var wg sync.WaitGroup

		errsCh := make(chan error, len(c.funcs))
		for _, f := range c.funcs {
			wg.Add(1)
			go func(f func() error) {
				defer wg.Done()
				defer func() {
					if r := recover(); r != nil {
						errsCh <- fmt.Errorf("panic при выполнении функции закрытия: %v", r)
					}
				}()

				errsCh <- f()
			}(f)
		}

		go func() {
			wg.Wait()
			close(errsCh)
		}()

		for err := range errsCh {
			if err != nil {
				log.Println("Ошибка при закрытии: ", err)
			}
		}
	})
}
