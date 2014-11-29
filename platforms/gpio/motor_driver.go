package gpio

import (
	"github.com/hybridgroup/gobot"
)

var _ gobot.Driver = (*MotorDriver)(nil)

// Represents a Motor
type MotorDriver struct {
	name             string
	connection       gobot.Connection
	SpeedPin         string
	SwitchPin        string
	DirectionPin     string
	ForwardPin       string
	BackwardPin      string
	CurrentState     byte
	CurrentSpeed     byte
	CurrentMode      string
	CurrentDirection string
}

// NewMotorDriver return a new MotorDriver given a PwmDigitalWriter, name and pin
func NewMotorDriver(a PwmDigitalWriter, name string) *MotorDriver {
	return &MotorDriver{
		name:             name,
		connection:       a.(gobot.Adaptor),
		CurrentState:     0,
		CurrentSpeed:     0,
		CurrentMode:      "digital",
		CurrentDirection: "forward",
	}
}

func (m *MotorDriver) Name() string                 { return m.name }
func (m *MotorDriver) Connection() gobot.Connection { return m.connection }

func (m *MotorDriver) adaptor() PwmDigitalWriter {
	return m.Connection().(PwmDigitalWriter)
}

// Start starts the MotorDriver. Returns true on successful start of the driver
func (m *MotorDriver) Start() (errs []error) { return }

// Halt halts the MotorDriver. Returns true on successful halt of the driver
func (m *MotorDriver) Halt() (errs []error) { return }

// Off turns the motor off or sets the motor to a 0 speed
func (m *MotorDriver) Off() (err error) {
	if m.isDigital() {
		err = m.changeState(0)
	} else {
		err = m.Speed(0)
	}
	return
}

// On turns the motor on or sets the motor to a maximum speed
func (m *MotorDriver) On() (err error) {
	if m.isDigital() {
		err = m.changeState(1)
	} else {
		if m.CurrentSpeed == 0 {
			m.CurrentSpeed = 255
		}
		err = m.Speed(m.CurrentSpeed)
	}
	return
}

// Min sets the motor to the minimum speed
func (m *MotorDriver) Min() (err error) {
	return m.Off()
}

// Max sets the motor to the maximum speed
func (m *MotorDriver) Max() (err error) {
	return m.Speed(255)
}

// InOn returns true if the motor is on
func (m *MotorDriver) IsOn() bool {
	if m.isDigital() {
		return m.CurrentState == 1
	}
	return m.CurrentSpeed > 0
}

// InOff returns true if the motor is off
func (m *MotorDriver) IsOff() bool {
	return !m.IsOn()
}

// Toggle sets the motor to the opposite of it's current state
func (m *MotorDriver) Toggle() (err error) {
	if m.IsOn() {
		err = m.Off()
	} else {
		err = m.On()
	}
	return
}

// Speed sets the speed of the motor
func (m *MotorDriver) Speed(value byte) (err error) {
	m.CurrentMode = "analog"
	m.CurrentSpeed = value
	return m.adaptor().PwmWrite(m.SpeedPin, value)
}

// Forward sets the forward pin to the specified speed
func (m *MotorDriver) Forward(speed byte) (err error) {
	err = m.Direction("forward")
	if err != nil {
		return
	}
	err = m.Speed(speed)
	if err != nil {
		return
	}
	return
}

// Backward sets the backward pin to the specified speed
func (m *MotorDriver) Backward(speed byte) (err error) {
	err = m.Direction("backward")
	if err != nil {
		return
	}
	err = m.Speed(speed)
	if err != nil {
		return
	}
	return
}

// Direction sets the direction pin to the specified speed
func (m *MotorDriver) Direction(direction string) (err error) {
	m.CurrentDirection = direction
	if m.DirectionPin != "" {
		var level byte
		if direction == "forward" {
			level = 1
		} else {
			level = 0
		}
		err = m.adaptor().DigitalWrite(m.DirectionPin, level)
	} else {
		var forwardLevel, backwardLevel byte
		switch direction {
		case "forward":
			forwardLevel = 1
			backwardLevel = 0
		case "backward":
			forwardLevel = 0
			backwardLevel = 1
		case "none":
			forwardLevel = 0
			backwardLevel = 0
		}
		err = m.adaptor().DigitalWrite(m.ForwardPin, forwardLevel)
		if err != nil {
			return
		}
		err = m.adaptor().DigitalWrite(m.BackwardPin, backwardLevel)
		if err != nil {
			return
		}
	}
	return
}

func (m *MotorDriver) isDigital() bool {
	if m.CurrentMode == "digital" {
		return true
	}
	return false
}

func (m *MotorDriver) changeState(state byte) (err error) {
	m.CurrentState = state
	if state == 1 {
		m.CurrentSpeed = 0
	} else {
		m.CurrentSpeed = 255
	}
	if m.ForwardPin != "" {
		if state == 0 {
			err = m.Direction(m.CurrentDirection)
			if err != nil {
				return
			}
			if m.SpeedPin != "" {
				err = m.Speed(m.CurrentSpeed)
				if err != nil {
					return
				}
			}
		} else {
			err = m.Direction("none")
		}
	} else {
		err = m.adaptor().DigitalWrite(m.SpeedPin, state)
	}

	return
}
