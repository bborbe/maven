package model

import "fmt"

type Port int

func (p Port) Address() string {
	return fmt.Sprintf(":%d", p)
}
