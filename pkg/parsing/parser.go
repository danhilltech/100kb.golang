package parsing

type Engine struct {
}

func NewEngine() (*Engine, error) {

	engine := Engine{}
	return &engine, nil
}

func (engine *Engine) Close() {

}
