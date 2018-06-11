package core

import (
	"ground-x/go-gxplatform/common"
	"fmt"
)

func (c *core) handleFinalCommitted() error {
	fmt.Printf("#### istanbul.FinalCommitted... \n")
	logger := c.logger.New("state", c.state)
	logger.Trace("Received a final committed proposal")
	c.startNewRound(common.Big0)
	return nil
}
