package conv

type Conf struct {
	InputType  string
	OutputType string
}

var conf Conf

func GetConf() *Conf {
	return &conf
}
