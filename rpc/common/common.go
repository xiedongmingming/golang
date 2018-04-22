package common

import "errors"

type Args struct {
	A, B int
}
type Quotient struct {
	Quo, Rem int
}
type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {

	if args.B == 0 {
		return errors.New("divide by zero")
	}

	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B

	return nil
}



type Mpmsg struct {
}

func (mpMsg *Mpmsg) Handle(args *map[string]interface{}, reply *map[string]interface{}) error {

	*reply = map[string]interface{}{}

	return nil
}