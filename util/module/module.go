package module

import (
	"context"
	"utilware/logger"

	"os"
	"sync"
	"time"
)

var (
	defaultBootTimeout                 = 1000 * time.Millisecond
	bootWaitTimeout                    = 5000 * time.Millisecond
	forceExit                          = false
	disableModules     map[string]bool = make(map[string]bool, 0)
)

type Module func(ctx context.Context) error

type ModuleAgent struct {
	name       string
	module     Module
	after      bool
	background bool
	front      bool
	timeout    time.Duration
}

type ModuleAgents []ModuleAgent

func (m Module) InitWithContext(ctx context.Context) error {
	return m(ctx)
}

func SetDefaultBootTimeout(t time.Duration) {
	defaultBootTimeout = t
}

func SetBootWaitTimeout(t time.Duration) {
	bootWaitTimeout = t
}

func SetForceExit(b bool) {
	forceExit = b
}

func DisableModule(module []string) {
	for i := 0; i < len(module); i++ {
		disableModules[module[i]] = true
	}
}

func NewModuleAgent(name string, module Module) ModuleAgent {
	return ModuleAgent{
		name:    name,
		module:  module,
		timeout: defaultBootTimeout,
	}
}

func (m ModuleAgent) After() ModuleAgent {
	m.after = true
	return m
}

func (m ModuleAgent) Background() ModuleAgent {
	m.background = true
	return m
}

func (m ModuleAgent) Front() ModuleAgent {
	m.front = true
	return m
}

func (m ModuleAgent) Timeout(t time.Duration) ModuleAgent {
	m.timeout = t * time.Millisecond
	return m
}

func Injecte(modules ModuleAgents) byte {
	wg := &sync.WaitGroup{}
	afterModules := ModuleAgents{}

	for i := 0; i < len(modules); i++ {
		// 检测是否禁用模块
		if _, ok := disableModules[modules[i].name]; ok {
			logger.Debug("[%s] module is disabled", modules[i].name)
			continue
		}

		// 延迟执行模块
		if modules[i].after {
			afterModules = append(afterModules, modules[i])
			continue
		}

		if modules[i].front {
			loadModule(modules[i])
			continue
		}

		if !modules[i].background {
			wg.Add(1)
		}
		go func(i int) {
			if !modules[i].background {
				defer wg.Done()
			}
			loadModule(modules[i])
		}(i)
	}

	if bootWaitTimeout > 0 {
		t := time.AfterFunc(bootWaitTimeout, func() {
			logger.Fatal("server start timeout!")
			os.Exit(0)
		})

		defer t.Stop()
	}

	wg.Wait()

	// 延迟执行模块
	for i := 0; i < len(afterModules); i++ {
		loadModule(afterModules[i])
	}

	return 0
}

func loadModule(moduleAgent ModuleAgent) {
	logger.Debug("[%s] module is starting", moduleAgent.name)

	ctx, cancel := context.Background(), func() {}

	if moduleAgent.timeout > 0 {
		t := time.AfterFunc(moduleAgent.timeout, func() {
			if forceExit {
				logger.Fatal("[%s] module is timeout", moduleAgent.name)
			}
			logger.Error("[%s] module is timeout", moduleAgent.name)
		})

		defer t.Stop()
		ctx, cancel = context.WithTimeout(ctx, moduleAgent.timeout)
	}

	defer cancel()

	if e := moduleAgent.module.
		InitWithContext(ctx); e != nil {
		logger.Fatal("[%s] module init error: %v", moduleAgent.name, e)
	} else {
		logger.Debug("[%s] module init success", moduleAgent.name)
	}

}
