package youtube

type PlayerCache struct {
	Version            string
	SignatureTimestamp string
	Operations         []Operation
}

type Operation struct {
	Name  int
	Value int
}

const (
	OpReverse = iota
	OpSplice
	OpSwap
)

func (c *PlayerCache) setSts(sts string) *PlayerCache {
	c.SignatureTimestamp = sts
	return c
}

func (c *PlayerCache) getSts(version string) (string, bool) {
	if c.Version != version {
		return "", false
	}
	return c.SignatureTimestamp, true
}

func (c *PlayerCache) addOps(op ...Operation) *PlayerCache {
	c.Operations = append(c.Operations, op...)
	return c
}

func (c *PlayerCache) getOps(version string) ([]Operation, bool) {
	if c.Version != version {
		return nil, false
	}
	return c.Operations, true
}
