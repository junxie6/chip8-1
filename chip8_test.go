package chip8

import "testing"

var opcodeTests = map[string][]struct {
	op     int
	before func(*testing.T, *CPU)
	check  func(*testing.T, *CPU)
}{
	"2nnn - CALL addr": {
		{
			0x2100,
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "Stack[0]", c.Stack[0], 0x200)
				checkHex(t, "SP", c.SP, 0x1)
				checkHex(t, "PC", c.PC, 0x100)
			},
		},
	},

	"3xkk - SE Vx, byte": {
		{
			0x3123,
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x200)
			},
		},

		{
			0x3103,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x202)
			},
		},
	},

	"4xkk - SNE Vx, byte": {
		{
			0x4123,
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x202)
			},
		},

		{
			0x4103,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x200)
			},
		},
	},

	"5xy0 - SE Vx, Vy": {
		{
			0x5120,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
				c.V[2] = 0x04
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x200)
			},
		},

		{
			0x5120,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
				c.V[2] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x202)
			},
		},
	},

	"6xkk - LD Vx, byte": {
		{
			0x6102,
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x02)
			},
		},
	},

	"7xkk - ADD Vx, byte": {
		{
			0x7102,
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x02)
			},
		},

		{
			0x7102,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x03)
			},
		},
	},

	"8xy0 - LD Vx, Vy": {
		{
			0x8120,
			func(t *testing.T, c *CPU) {
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x01)
			},
		},
	},

	"8xy1 - OR Vx, Vy": {
		{
			0x8121,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x10
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x11)
			},
		},
	},

	"8xy2 - AND Vx, Vy": {
		{
			0x8122,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x10
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x00)
			},
		},
	},

	"8xy3 - XOR Vx, Vy": {
		{
			0x8123,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x00)
			},
		},
	},

	"8xy4 - ADD Vx, Vy": {
		{
			0x8124,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x2)
				checkHex(t, "VF", c.V[0xF], 0x0)
			},
		},

		{
			0x8124,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0xFF
				c.V[2] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x2)
				checkHex(t, "VF", c.V[0xF], 0x1)
			},
		},
	},

	"8xy5 - SUB Vx, Vy": {
		{
			0x8125,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0xFF
				c.V[2] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", c.V[0xF], 0x1)
				checkHex(t, "V[1]", c.V[1], 0xFC)
			},
		},

		{
			0x8125,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x02
				c.V[2] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", c.V[0xF], 0x0)
				checkHex(t, "V[1]", c.V[1], 0xFF)
			},
		},
	},

	"8xy6 - SHR Vx {, Vy}": {
		{
			0x8126,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", c.V[0xF], 0x1)
				checkHex(t, "V[1]", c.V[1], 0x1)
			},
		},

		{
			0x8126,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x02
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", c.V[0xF], 0x0)
				checkHex(t, "V[1]", c.V[1], 0x1)
			},
		},
	},

	"8xy7 - SUBN Vx, Vy": {
		{
			0x8127,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
				c.V[2] = 0xFF
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", c.V[0xF], 0x1)
				checkHex(t, "V[1]", c.V[1], 0xFC)
			},
		},

		{
			0x8127,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x03
				c.V[2] = 0x02
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", c.V[0xF], 0x0)
				checkHex(t, "V[1]", c.V[1], 0xFF)
			},
		},
	},

	"8xyE - SHL Vx {, Vy}": {
		{
			0x812E,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", c.V[0xF], 0x0)
				checkHex(t, "V[1]", c.V[1], 0x2)
			},
		},

		{
			0x812E,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x81
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "VF", c.V[0xF], 0x1)
				checkHex(t, "V[1]", c.V[1], 0x2)
			},
		},
	},

	"9xy0 - SNE Vx, Vy": {
		{
			0x9120,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
				c.V[2] = 0x02
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x202)
			},
		},

		{
			0x9120,
			func(t *testing.T, c *CPU) {
				c.V[1] = 0x01
				c.V[2] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x200)
			},
		},
	},

	"Annn - LD I, addr": {
		{
			0xA100,
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "I", c.I, 0x100)
			},
		},
	},

	"Bnnn - JP V0, addr": {
		{
			0xB210,
			nil,
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x210)
			},
		},

		{
			0xB210,
			func(t *testing.T, c *CPU) {
				c.V[0] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "PC", c.PC, 0x211)
			},
		},
	},

	"Cxkk - RND Vx, byte": {
		{
			0xC110,
			func(t *testing.T, c *CPU) {
				c.randByteFunc = func() byte {
					return 0x01
				}
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "V[1]", c.V[1], 0x11)
			},
		},
	},

	"Dxyn - DRW Vx, Vy, nibble": {
		{
			0xD001,
			func(t *testing.T, c *CPU) {
				c.I = 0x200
				c.Memory[0x200] = 0x01
			},
			func(t *testing.T, c *CPU) {
				checkHex(t, "Pixel", c.Pixels[0], 0x00)
			},
		},
	},
}

func TestCPU_Step(t *testing.T) {
	c := NewCPU(nil)
	c.Memory[0x200] = 0xA1
	c.Memory[0x201] = 0x00

	if err := c.Step(); err != nil {
		t.Fatal(err)
	}

	checkHex(t, "PC", c.PC, 0x202)
}

func TestOpcodes(t *testing.T) {
	for i, tests := range opcodeTests {
		for _, tt := range tests {
			c := NewCPU(nil)
			if tt.before != nil {
				tt.before(t, c)
			}
			c.Dispatch(uint16(tt.op))
			tt.check(t, c)

			if t.Failed() {
				t.Logf("==============")
				t.Logf("Instruction: %s", i)
				t.Logf("Opcode: 0x%04X", tt.op)
				t.Logf("CPU: %v", c)
				t.Logf("==============")
				t.FailNow()
			}
		}
	}
}

func TestCPU_Load(t *testing.T) {
	c := NewCPU(nil)
	p := []byte{0x01, 0x02}

	n, err := c.LoadBytes(p)
	if err != nil {
		t.Fatal(err)
	}

	if n != len(p) {
		t.Fatal(err)
	}

	checkHex(t, "Memory[0x200]", c.Memory[0x200], 0x01)
	checkHex(t, "Memory[0x201]", c.Memory[0x201], 0x02)
}

func TestCPU_op(t *testing.T) {
	c := NewCPU(nil)
	c.Memory[0x200] = 0xA2
	c.Memory[0x201] = 0xF0

	checkHex(t, "op", c.op(), 0xA2F0)
}

func tryUint16(v interface{}) uint16 {
	switch v := v.(type) {
	case byte:
		return uint16(v)
	case uint16:
		return v
	case int:
		return uint16(v)
	case uint32:
		return uint16(v)
	}

	return 0
}

func checkHex(t *testing.T, subject string, got, want interface{}) {
	g := tryUint16(got)
	w := tryUint16(want)

	if g != w {
		t.Errorf("%s => 0x%04X; want 0x%04X", subject, g, w)
	}
}
