package client

func (c *PClient) Index() (string, error) {
	return c.GetPage("https://na.preva.com/exerciser-api//exerciser-account")
}
