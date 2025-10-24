package main

type Command interface {
	Name() string
	Run(args []string) error
}

func GetName(c Command) string {
	return c.Name()
}

func Runnable(c Command, args []string) error {
	return c.Run(args)
}