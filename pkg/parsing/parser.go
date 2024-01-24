package parsing

type Engine struct {
	adblock *AdblockEngine
}

func NewEngine() (*Engine, error) {
	adblock, err := NewAdblockEngine()
	if err != nil {
		return nil, err
	}
	engine := Engine{
		adblock: adblock,
	}
	return &engine, nil
}

func (engine *Engine) Close() {
	engine.adblock.Close()
}
