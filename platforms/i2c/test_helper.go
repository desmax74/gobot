package i2c

var rgb = map[string]interface{}{
	"red":   1.0,
	"green": 1.0,
	"blue":  1.0,
}

func castColor(color string) byte {
	return byte(rgb[color].(float64))
}

var red = castColor("red")
var green = castColor("green")
var blue = castColor("blue")

type i2cTestAdaptor struct {
	name        string
	i2cReadImpl func() []byte
}

func (t *i2cTestAdaptor) I2cStart(byte) (err error) { return nil }
func (t *i2cTestAdaptor) I2cRead(uint) (data []byte, err error) {
	return t.i2cReadImpl(), nil
}
func (t *i2cTestAdaptor) I2cWrite([]byte) (err error) { return nil }
func (t *i2cTestAdaptor) Name() string                { return t.name }
func (t *i2cTestAdaptor) Connect() (errs []error)     { return }
func (t *i2cTestAdaptor) Finalize() (errs []error)    { return }

func newI2cTestAdaptor(name string) *i2cTestAdaptor {
	return &i2cTestAdaptor{
		name: name,
		i2cReadImpl: func() []byte {
			return []byte{}
		},
	}
}
