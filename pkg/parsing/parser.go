package parsing

import "log"

type Engine struct {
	log *log.Logger
}

func NewEngine(log *log.Logger) (*Engine, error) {

	engine := Engine{log: log}
	return &engine, nil
}

func (engine *Engine) Close() {

}
