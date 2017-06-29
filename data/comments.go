package data

// Comments describe code comments
type Comments struct {
	Leading         string
	Trailing        string
	LeadingDetached []string
}

func (c *Comments) String() string {
	return c.Leading
}
