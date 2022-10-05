package options

import (
	"fmt"
	"math"
	"strconv"
)

type portOption struct {
	port uint
}

func NewPortOption() *portOption {
	return &portOption{
		port: 8080,
	}
}

func (o *portOption) Value() uint {
	return o.port
}

func (o *portOption) Set(s string) error {
	v, err := strconv.ParseUint(s, 10, 16)
	if err != nil {
		return fmt.Errorf("invalid number or outside")
	}

	if v > math.MaxUint16 {
		return fmt.Errorf("outside valid range for port")
	}
	if v <= 1024 {
		return fmt.Errorf("reserved port number")
	}

	o.port = uint(v)
	return nil
}

func (h *portOption) Type() string {
	return "Port"
}

func (h *portOption) String() string {
	return fmt.Sprint(h.port)
}
