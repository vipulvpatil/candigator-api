package utilities

type IdGeneratorMockConstant struct {
	Id string
}

func (c *IdGeneratorMockConstant) Generate() string {
	return c.Id
}

type IdGeneratorMockSeries struct {
	Series []string
	index  int
}

func (c *IdGeneratorMockSeries) Generate() string {
	var next string
	if c.index < len(c.Series) {
		next = c.Series[c.index]
		c.index++
	} else {
		next = "undefined"
	}
	return next
}
