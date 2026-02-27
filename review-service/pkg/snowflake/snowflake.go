package snowflake

import (
	"errors"
	"fmt"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var (
	InvalidInitParamErr = errors.New("Invalid init parameter")
	InvalidTimeFormatErr = errors.New("Invalid time format")
)

var node *sf.Node

func Init(startTime string, nodeID int64) error {
	if len(startTime) == 0 || nodeID < 0 || nodeID > 1023 {
		return InvalidInitParamErr
	}
	var st time.Time
	var err error
	st, err = time.Parse("2006-01-02 15:04:05", startTime)
	if err != nil {
		return fmt.Errorf("Failed to parse start time: %v", err)
	}
	sf.Epoch = st.UnixNano() / 1e6
	node, err = sf.NewNode(nodeID)
	if err != nil {
		return fmt.Errorf("Failed to init snowflake node: %v", err)
	}
	return nil
}

func GenerateID() int64 {
	return node.Generate().Int64()
}