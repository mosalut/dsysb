// dsysb

package main

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
)

func (task *task_T) opCheck(length int, args ...int) error {
	vdLength := len(task.vData)
	for k, p := range args {
		limit := p + length
		if vdLength < limit {
			fmt.Println(vdLength, limit)
			return errors.New(fmt.Sprintf("vData error at p%d: %d", k, p))
		}
	}

	return nil
}

/* ------ mov ------ */
func (task *task_T) movsb(p0, p1, length int) error {
	copy(task.vData[p1:p1 + length], task.vData[p0:p0 + length])

	return nil
}

func (task *task_T) mov8(p0, p1 int) error {
	task.vData[p1] = task.vData[p0]

	return nil
}

func (task *task_T) mov16(p0, p1 int) error {
	copy(task.vData[p1:p1 + 2], task.vData[p0:p0 + 2])

	return nil
}

func (task *task_T) mov32(p0, p1 int) error {
	copy(task.vData[p1:p1 + 4], task.vData[p0:p0 + 4])

	return nil
}

func (task *task_T) mov64(p0, p1 int) error {
	copy(task.vData[p1:p1 + 8], task.vData[p0:p0 + 8])

	return nil
}

/* ------ add ------ */
func (task *task_T) add8(p0, p1, p2 int) error {
	task.vData[p2] = byte(int8(task.vData[p0]) + int8(task.vData[p1]))

	return nil
}

func (task *task_T) add16(p0, p1, p2 int) error {
	x := int16(binary.LittleEndian.Uint16(task.vData[p0:p0 + 2]))
	y := int16(binary.LittleEndian.Uint16(task.vData[p1:p1 + 2]))
	binary.LittleEndian.PutUint16(task.vData[p2:p2 + 2], uint16(x + y))

	return nil
}

func (task *task_T) add32(p0, p1, p2 int) error {
	x := int32(binary.LittleEndian.Uint32(task.vData[p0:p0 + 4]))
	y := int32(binary.LittleEndian.Uint32(task.vData[p1:p1 + 4]))
	binary.LittleEndian.PutUint32(task.vData[p2:p2 + 4], uint32(x + y))

	return nil
}

func (task *task_T) add64(p0, p1, p2 int) error {
	x := int64(binary.LittleEndian.Uint64(task.vData[p0:p0 + 8]))
	y := int64(binary.LittleEndian.Uint64(task.vData[p1:p1 + 8]))
	binary.LittleEndian.PutUint64(task.vData[p2:p2 + 8], uint64(x + y))

	return nil
}

func (task *task_T) add8u(p0, p1, p2 int) error {
	task.vData[p2] = byte(uint8(task.vData[p0]) + uint8(task.vData[p1]))

	return nil
}

func (task *task_T) add16u(p0, p1, p2 int) error {
	x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
	y := binary.LittleEndian.Uint16(task.vData[p1:p1 + 2])
	binary.LittleEndian.PutUint16(task.vData[p2:p2 + 2], x + y)

	return nil
}

func (task *task_T) add32u(p0, p1, p2 int) error {
	x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
	y := binary.LittleEndian.Uint32(task.vData[p1:p1 + 4])
	binary.LittleEndian.PutUint32(task.vData[p2:p2 + 4], x + y)

	return nil
}

func (task *task_T) add64u(p0, p1, p2 int) error {
	x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
	y := binary.LittleEndian.Uint64(task.vData[p1:p1 + 8])
	binary.LittleEndian.PutUint64(task.vData[p2:p2 + 8], x + y)

	return nil
}

/* ------ sub ------ */
func (task *task_T) sub8(p0, p1, p2 int) error {
	task.vData[p2] = byte(int8(task.vData[p0]) - int8(task.vData[p1]))

	return nil
}

func (task *task_T) sub16(p0, p1, p2 int) error {
	x := int16(binary.LittleEndian.Uint16(task.vData[p0:p0 + 2]))
	y := int16(binary.LittleEndian.Uint16(task.vData[p1:p1 + 2]))
	binary.LittleEndian.PutUint16(task.vData[p2:p2 + 2], uint16(x - y))

	return nil
}

func (task *task_T) sub32(p0, p1, p2 int) error {
	x := int32(binary.LittleEndian.Uint32(task.vData[p0:p0 + 4]))
	y := int32(binary.LittleEndian.Uint32(task.vData[p1:p1 + 4]))
	binary.LittleEndian.PutUint32(task.vData[p2:p2 + 4], uint32(x - y))

	return nil
}

func (task *task_T) sub64(p0, p1, p2 int) error {
	x := int64(binary.LittleEndian.Uint64(task.vData[p0:p0 + 8]))
	y := int64(binary.LittleEndian.Uint64(task.vData[p1:p1 + 8]))
	binary.LittleEndian.PutUint64(task.vData[p2:p2 + 8], uint64(x - y))

	return nil
}

func (task *task_T) sub8u(p0, p1, p2 int) error {
	task.vData[p2] = byte(uint8(task.vData[p0]) - uint8(task.vData[p1]))

	return nil
}

func (task *task_T) sub16u(p0, p1, p2 int) error {
	x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
	y := binary.LittleEndian.Uint16(task.vData[p1:p1 + 2])
	binary.LittleEndian.PutUint16(task.vData[p2:p2 + 2], x - y)

	return nil
}

func (task *task_T) sub32u(p0, p1, p2 int) error {
	x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
	y := binary.LittleEndian.Uint32(task.vData[p1:p1 + 4])
	binary.LittleEndian.PutUint32(task.vData[p2:p2 + 4], x - y)

	return nil
}

func (task *task_T) sub64u(p0, p1, p2 int) error {
	x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
	y := binary.LittleEndian.Uint64(task.vData[p1:p1 + 8])
	binary.LittleEndian.PutUint64(task.vData[p2:p2 + 8], x - y)

	return nil
}

/* ------ mul ------ */
func (task *task_T) mul8(p0, p1, p2 int) error {
	task.vData[p2] = byte(int8(task.vData[p0]) * int8(task.vData[p1]))

	return nil
}

func (task *task_T) mul16(p0, p1, p2 int) error {
	x := int16(binary.LittleEndian.Uint16(task.vData[p0:p0 + 2]))
	y := int16(binary.LittleEndian.Uint16(task.vData[p1:p1 + 2]))
	binary.LittleEndian.PutUint16(task.vData[p2:p2 + 2], uint16(x * y))

	return nil
}

func (task *task_T) mul32(p0, p1, p2 int) error {
	x := int32(binary.LittleEndian.Uint32(task.vData[p0:p0 + 4]))
	y := int32(binary.LittleEndian.Uint32(task.vData[p1:p1 + 4]))
	binary.LittleEndian.PutUint32(task.vData[p2:p2 + 4], uint32(x * y))

	return nil
}

func (task *task_T) mul64(p0, p1, p2 int) error {
	x := int64(binary.LittleEndian.Uint64(task.vData[p0:p0 + 8]))
	y := int64(binary.LittleEndian.Uint64(task.vData[p1:p1 + 8]))
	binary.LittleEndian.PutUint64(task.vData[p2:p2 + 8], uint64(x * y))

	return nil
}

func (task *task_T) mul8u(p0, p1, p2 int) error {
	task.vData[p2] = byte(uint8(task.vData[p0]) * uint8(task.vData[p1]))

	return nil
}

func (task *task_T) mul16u(p0, p1, p2 int) error {
	x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
	y := binary.LittleEndian.Uint16(task.vData[p1:p1 + 2])
	binary.LittleEndian.PutUint16(task.vData[p2:p2 + 2], x * y)

	return nil
}

func (task *task_T) mul32u(p0, p1, p2 int) error {
	x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
	y := binary.LittleEndian.Uint32(task.vData[p1:p1 + 4])
	binary.LittleEndian.PutUint32(task.vData[p2:p2 + 4], x * y)

	return nil
}

func (task *task_T) mul64u(p0, p1, p2 int) error {
	x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
	y := binary.LittleEndian.Uint64(task.vData[p1:p1 + 8])
	binary.LittleEndian.PutUint64(task.vData[p2:p2 + 8], x * y)

	return nil
}

/* ------ quo ------ */
func (task *task_T) quo8(p0, p1, p2, p3 int) error {
	d := int8(task.vData[p1])
	if d == 0 {
		return errors.New("quo8 p1, divisor is zero")
	}
	task.vData[p2] = byte(int8(task.vData[p0]) / d)
	task.vData[p3] = byte(int8(task.vData[p0]) % d)

	return nil
}

func (task *task_T) quo16(p0, p1, p2, p3 int) error {
	x := int16(binary.LittleEndian.Uint16(task.vData[p0:p0 + 2]))
	y := int16(binary.LittleEndian.Uint16(task.vData[p1:p1 + 2]))
	if y == 0 {
		return errors.New("quo16 p1, divisor is zero")
	}
	binary.LittleEndian.PutUint16(task.vData[p2:p2 + 2], uint16(x / y))
	binary.LittleEndian.PutUint16(task.vData[p3:p3 + 2], uint16(x % y))

	return nil
}

func (task *task_T) quo32(p0, p1, p2, p3 int) error {
	x := int32(binary.LittleEndian.Uint32(task.vData[p0:p0 + 4]))
	y := int32(binary.LittleEndian.Uint32(task.vData[p1:p1 + 4]))
	if y == 0 {
		return errors.New("quo32 p1, divisor is zero")
	}
	binary.LittleEndian.PutUint32(task.vData[p2:p2 + 4], uint32(x / y))
	binary.LittleEndian.PutUint32(task.vData[p3:p3 + 4], uint32(x % y))

	return nil
}

func (task *task_T) quo64(p0, p1, p2, p3 int) error {
	x := int64(binary.LittleEndian.Uint64(task.vData[p0:p0 + 8]))
	y := int64(binary.LittleEndian.Uint64(task.vData[p1:p1 + 8]))
	if y == 0 {
		return errors.New("quo64 p1, divisor is zero")
	}
	binary.LittleEndian.PutUint64(task.vData[p2:p2 + 8], uint64(x / y))
	binary.LittleEndian.PutUint64(task.vData[p3:p3 + 8], uint64(x % y))

	return nil
}

func (task *task_T) quo8u(p0, p1, p2, p3 int) error {
	d := int(task.vData[p1])
	if d == 0 {
		return errors.New("quo8 p1, divisor is zero")
	}
	task.vData[p2] = byte(uint8(task.vData[p0]) / uint8(task.vData[p1]))
	task.vData[p3] = byte(uint8(task.vData[p0]) % uint8(task.vData[p1]))

	return nil
}

func (task *task_T) quo16u(p0, p1, p2, p3 int) error {
	x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
	y := binary.LittleEndian.Uint16(task.vData[p1:p1 + 2])
	if y == 0 {
		return errors.New("quo16 p1, divisor is zero")
	}
	binary.LittleEndian.PutUint16(task.vData[p2:p2 + 2], x / y)
	binary.LittleEndian.PutUint16(task.vData[p3:p3 + 2], x % y)

	return nil
}

func (task *task_T) quo32u(p0, p1, p2, p3 int) error {
	x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
	y := binary.LittleEndian.Uint32(task.vData[p1:p1 + 4])
	if y == 0 {
		return errors.New("quo32 p1, divisor is zero")
	}
	binary.LittleEndian.PutUint32(task.vData[p2:p2 + 4], x / y)
	binary.LittleEndian.PutUint32(task.vData[p3:p3 + 4], x % y)

	return nil
}

func (task *task_T) quo64u(p0, p1, p2, p3 int) error {
	x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
	y := binary.LittleEndian.Uint64(task.vData[p1:p1 + 8])
	if y == 0 {
		return errors.New("quo64 p1, divisor is zero")
	}
	binary.LittleEndian.PutUint64(task.vData[p2:p2 + 8], x / y)
	binary.LittleEndian.PutUint64(task.vData[p3:p3 + 8], x % y)

	return nil
}

/* ------ inc ------ */
func (task *task_T) inc8(p0 int) error {
	task.vData[p0] = byte(int8(task.vData[p0]) + 1)

	return nil
}

func (task *task_T) inc16(p0 int) error {
	x := int16(binary.LittleEndian.Uint16(task.vData[p0:p0 + 2]))
	binary.LittleEndian.PutUint16(task.vData[p0:p0 + 2], uint16(x + 1))

	return nil
}

func (task *task_T) inc32(p0 int) error {
	x := int32(binary.LittleEndian.Uint32(task.vData[p0:p0 + 4]))
	binary.LittleEndian.PutUint32(task.vData[p0:p0 + 4], uint32(x + 1))

	return nil
}

func (task *task_T) inc64(p0 int) error {
	x := int64(binary.LittleEndian.Uint64(task.vData[p0:p0 + 8]))
	binary.LittleEndian.PutUint64(task.vData[p0:p0 + 8], uint64(x + 1))

	return nil
}

func (task *task_T) inc8u(p0 int) error {
	task.vData[p0] = byte(uint8(task.vData[p0]) + 1)

	return nil
}

func (task *task_T) inc16u(p0 int) error {
	x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
	binary.LittleEndian.PutUint16(task.vData[p0:p0 + 2], x + 1)

	return nil
}

func (task *task_T) inc32u(p0 int) error {
	x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
	binary.LittleEndian.PutUint32(task.vData[p0:p0 + 4], x + 1)

	return nil
}

func (task *task_T) inc64u(p0 int) error {
	x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
	binary.LittleEndian.PutUint64(task.vData[p0:p0 + 8], x + 1)

	return nil
}

/* ------ dec ------ */
func (task *task_T) dec8(p0 int) error {
	task.vData[p0] = byte(int8(task.vData[p0]) + 1)

	return nil
}

func (task *task_T) dec16(p0 int) error {
	x := int16(binary.LittleEndian.Uint16(task.vData[p0:p0 + 2]))
	binary.LittleEndian.PutUint16(task.vData[p0:p0 + 2], uint16(x + 1))

	return nil
}

func (task *task_T) dec32(p0 int) error {
	x := int32(binary.LittleEndian.Uint32(task.vData[p0:p0 + 4]))
	binary.LittleEndian.PutUint32(task.vData[p0:p0 + 4], uint32(x + 1))

	return nil
}

func (task *task_T) dec64(p0 int) error {
	x := int64(binary.LittleEndian.Uint64(task.vData[p0:p0 + 8]))
	binary.LittleEndian.PutUint64(task.vData[p0:p0 + 8], uint64(x + 1))

	return nil
}

func (task *task_T) dec8u(p0 int) error {
	task.vData[p0] = byte(uint8(task.vData[p0]) - 1)

	return nil
}

func (task *task_T) dec16u(p0 int) error {
	x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
	binary.LittleEndian.PutUint16(task.vData[p0:p0 + 2], x - 1)

	return nil
}

func (task *task_T) dec32u(p0 int) error {
	x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
	binary.LittleEndian.PutUint32(task.vData[p0:p0 + 4], x - 1)

	return nil
}

func (task *task_T) dec64u(p0 int) error {
	x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
	binary.LittleEndian.PutUint64(task.vData[p0:p0 + 8], x - 1)

	return nil
}

/* ------ write ------ */
func (task *task_T) writeUint8(reg *reg_T, p0 int) error {
	task.vData[p0] = reg.vUint8

	return nil
}

func (task *task_T) writeUint16(reg *reg_T, p0 int) error {
	binary.LittleEndian.PutUint16(task.vData[p0:p0 + 2], reg.vUint16)

	return nil
}

func (task *task_T) writeUint32(reg *reg_T, p0 int) error {
	binary.LittleEndian.PutUint32(task.vData[p0:p0 + 4], reg.vUint32)

	return nil
}

func (task *task_T) writeUint64(reg *reg_T, p0 int) error {
	binary.LittleEndian.PutUint64(task.vData[p0:p0 + 8], reg.vUint64)

	return nil
}

/* ------ read ------ */
func (task *task_T) readUint8(reg *reg_T, p0 int) error {
	reg.vUint8 = task.vData[p0]

	return nil
}

func (task *task_T) readUint16(reg *reg_T, p0 int) error {
	reg.vUint16 = binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])

	return nil
}

func (task *task_T) readUint32(reg *reg_T, p0 int) error {
	reg.vUint32 = binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])

	return nil
}

func (task *task_T) readUint64(reg *reg_T, p0 int) error {
	reg.vUint64 = binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])

	return nil
}

/* ------ compare ------ */
func (task *task_T) eq(reg *reg_T, p0, p1, p2, p3 int, ip *int) error {
	flag := uint8(task.vData[p3])

	switch flag {
	case 0:
		x := task.vData[p0]
		y := task.vData[p1]
		reg.vBool = x == y
	case 1:
		x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
		y := binary.LittleEndian.Uint16(task.vData[p1:p1 + 2])
		reg.vBool = x == y
	case 2:
		x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
		y := binary.LittleEndian.Uint32(task.vData[p1:p1 + 4])
		reg.vBool = x == y
	case 3:
		x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
		y := binary.LittleEndian.Uint64(task.vData[p1:p1 + 8])
		reg.vBool = x == y
	default:
		return errors.New("Wrong type of task op eq")
	}

	if reg.vBool {
		*ip = p2
	}

	return nil
}

func (task *task_T) gt(reg *reg_T, p0, p1, p2, p3 int, ip *int) error {
	flag := uint8(task.vData[p3])

	switch flag {
	case 0:
		x := task.vData[p0]
		y := task.vData[p1]
		reg.vBool = x > y
	case 1:
		x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
		y := binary.LittleEndian.Uint16(task.vData[p1:p1 + 2])
		reg.vBool = x > y
	case 2:
		x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
		y := binary.LittleEndian.Uint32(task.vData[p1:p1 + 4])
		reg.vBool = x > y
	case 3:
		x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
		y := binary.LittleEndian.Uint64(task.vData[p1:p1 + 8])
		reg.vBool = x > y
	case 4:
		x := int8(task.vData[p0])
		y := int8(task.vData[p1])
		reg.vBool = x > y
	case 5:
		x := int16(binary.LittleEndian.Uint16(task.vData[p0:p0 + 2]))
		y := int16(binary.LittleEndian.Uint16(task.vData[p1:p1 + 2]))
		reg.vBool = x > y
	case 6:
		x := int32(binary.LittleEndian.Uint32(task.vData[p0:p0 + 4]))
		y := int32(binary.LittleEndian.Uint32(task.vData[p1:p1 + 4]))
		reg.vBool = x > y
	case 7:
		x := int64(binary.LittleEndian.Uint64(task.vData[p0:p0 + 8]))
		y := int64(binary.LittleEndian.Uint64(task.vData[p1:p1 + 8]))
		reg.vBool = x > y
	default:
		return errors.New("Wrong type of task op gt")
	}

	if reg.vBool {
		*ip = p2
	}

	return nil
}

func (task *task_T) lt(reg *reg_T, p0, p1, p2, p3 int, ip *int) error {
	flag := uint8(task.vData[p3])

	switch flag {
	case 0:
		x := task.vData[p0]
		y := task.vData[p1]
		reg.vBool = x < y
	case 1:
		x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
		y := binary.LittleEndian.Uint16(task.vData[p1:p1 + 2])
		reg.vBool = x < y
	case 2:
		x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
		y := binary.LittleEndian.Uint32(task.vData[p1:p1 + 4])
		reg.vBool = x < y
	case 3:
		x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
		y := binary.LittleEndian.Uint64(task.vData[p1:p1 + 8])
		reg.vBool = x < y
	case 4:
		x := int8(task.vData[p0])
		y := int8(task.vData[p1])
		reg.vBool = x < y
	case 5:
		x := int16(binary.LittleEndian.Uint16(task.vData[p0:p0 + 2]))
		y := int16(binary.LittleEndian.Uint16(task.vData[p1:p1 + 2]))
		reg.vBool = x < y
	case 6:
		x := int32(binary.LittleEndian.Uint32(task.vData[p0:p0 + 4]))
		y := int32(binary.LittleEndian.Uint32(task.vData[p1:p1 + 4]))
		reg.vBool = x < y
	case 7:
		x := int64(binary.LittleEndian.Uint64(task.vData[p0:p0 + 8]))
		y := int64(binary.LittleEndian.Uint64(task.vData[p1:p1 + 8]))
		reg.vBool = x < y
	default:
		return errors.New("Wrong type of task op lt")
	}

	if reg.vBool {
		*ip = p2
	}

	return nil
}

func (task *task_T) gteq(reg *reg_T, p0, p1, p2, p3 int, ip *int) error {
	flag := uint8(task.vData[p3])

	switch flag {
	case 0:
		x := task.vData[p0]
		y := task.vData[p1]
		reg.vBool = x >= y
	case 1:
		x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
		y := binary.LittleEndian.Uint16(task.vData[p1:p1 + 2])
		reg.vBool = x >= y
	case 2:
		x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
		y := binary.LittleEndian.Uint32(task.vData[p1:p1 + 4])
		reg.vBool = x >= y
	case 3:
		x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
		y := binary.LittleEndian.Uint64(task.vData[p1:p1 + 8])
		reg.vBool = x >= y
	case 4:
		x := int8(task.vData[p0])
		y := int8(task.vData[p1])
		reg.vBool = x >= y
	case 5:
		x := int16(binary.LittleEndian.Uint16(task.vData[p0:p0 + 2]))
		y := int16(binary.LittleEndian.Uint16(task.vData[p1:p1 + 2]))
		reg.vBool = x >= y
	case 6:
		x := int32(binary.LittleEndian.Uint32(task.vData[p0:p0 + 4]))
		y := int32(binary.LittleEndian.Uint32(task.vData[p1:p1 + 4]))
		reg.vBool = x >= y
	case 7:
		x := int64(binary.LittleEndian.Uint64(task.vData[p0:p0 + 8]))
		y := int64(binary.LittleEndian.Uint64(task.vData[p1:p1 + 8]))
		reg.vBool = x >= y
	default:
		return errors.New("Wrong type of task op gteq")
	}

	if reg.vBool {
		*ip = p2
	}

	return nil
}

func (task *task_T) lteq(reg *reg_T, p0, p1, p2, p3 int, ip *int) error {
	flag := uint8(task.vData[p3])

	switch flag {
	case 0:
		x := task.vData[p0]
		y := task.vData[p1]
		reg.vBool = x <= y
	case 1:
		x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
		y := binary.LittleEndian.Uint16(task.vData[p1:p1 + 2])
		reg.vBool = x <= y
	case 2:
		x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
		y := binary.LittleEndian.Uint32(task.vData[p1:p1 + 4])
		reg.vBool = x <= y
	case 3:
		x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
		y := binary.LittleEndian.Uint64(task.vData[p1:p1 + 8])
		reg.vBool = x <= y
	case 4:
		x := int8(task.vData[p0])
		y := int8(task.vData[p1])
		reg.vBool = x <= y
	case 5:
		x := int16(binary.LittleEndian.Uint16(task.vData[p0:p0 + 2]))
		y := int16(binary.LittleEndian.Uint16(task.vData[p1:p1 + 2]))
		reg.vBool = x <= y
	case 6:
		x := int32(binary.LittleEndian.Uint32(task.vData[p0:p0 + 4]))
		y := int32(binary.LittleEndian.Uint32(task.vData[p1:p1 + 4]))
		reg.vBool = x <= y
	case 7:
		x := int64(binary.LittleEndian.Uint64(task.vData[p0:p0 + 8]))
		y := int64(binary.LittleEndian.Uint64(task.vData[p1:p1 + 8]))
		reg.vBool = x <= y
	default:
		return errors.New("task op call:Wrong type of task op lteq")
	}

	if reg.vBool {
		*ip = p2
	}

	return nil
}

func (task *task_T) eqBytes(reg *reg_T, p0, p1, p2, length int, ip *int) error {
	reg.vBool = hex.EncodeToString(task.vData[p0:p0 + length]) == hex.EncodeToString(task.vData[p1:p1 + length])

	if reg.vBool {
		*ip = p2
	}

	return nil
}

func (task *task_T) getIndex(reg *reg_T) error {
	var err error
	reg.vUint32, err = getIndex()
	if err != nil {
		return err
	}

	return nil
}

/* ------ transfer ------ */
func (task *task_T) transferDSBFromCaller(state *state_T, address string, p0 int) error {
	amount := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])

	accountFrom, ok := state.accounts[address]
	if !ok {
		return errors.New("task op call:the from address is not found")
	}

	accountTo, ok := state.accounts[task.address]
	if !ok {
		return errors.New("task op call:the to address is not found")
	}

	if accountFrom.balance < amount {
		return errors.New("task op call:not enough more DSBs")
	}

	accountFrom.balance -= amount
	accountTo.balance += amount

	return nil
}

func (task *task_T) transferDSBToCaller(state *state_T, address string, p0 int) error {
	amount := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])

	accountFrom, ok := state.accounts[task.address]
	if !ok {
		return errors.New("task op call:the from address is not found")
	}

	accountTo, ok := state.accounts[address]
	if !ok {
		return errors.New("task op call:the to address is not found")
	}

	if accountFrom.balance < amount {
		return errors.New("task op call:not enough more DSBs")
	}

	accountFrom.balance -= amount
	accountTo.balance += amount

	return nil
}

func (task *task_T) transferFromCaller(state *state_T, address string, p0, p1 int) error {
	id := hex.EncodeToString(task.vData[p0:p0 + 32])
	amount := binary.LittleEndian.Uint64(task.vData[p1:p1 + 8])

	accountFrom, ok := state.accounts[address]
	if !ok {
		return errors.New("task op call:the from account is not found")
	}

	accountTo, ok := state.accounts[task.address]
	if !ok {
		return errors.New("task op call:the to account is not found")
	}

	_, ok = accountFrom.assets[id]
	if !ok {
		return errors.New("task op call:the from account's asset is not found")
	}

	if accountFrom.assets[id] < amount {
		return errors.New("task op call:not enough more tokens")
	}

	accountFrom.assets[id] -= amount
	_, ok = accountTo.assets[id]
	if !ok {
		accountTo.assets[id] = 0
	}
	accountTo.assets[id] += amount

	return nil
}

func (task *task_T) transferToCaller(state *state_T, address string, p0, p1 int) error {
	id := hex.EncodeToString(task.vData[p0:p0 + 32])
	amount := binary.LittleEndian.Uint64(task.vData[p1:p1 + 8])

	accountFrom, ok := state.accounts[task.address]
	if !ok {
		return errors.New("task op call:the from address is not found")
	}

	accountTo, ok := state.accounts[address]
	if !ok {
		return errors.New("task op call:the to address is not found")
	}

	_, ok = accountFrom.assets[id]
	if !ok {
		return errors.New("task op call:the from account's asset is not found")
	}

	if accountFrom.assets[id] < amount {
		return errors.New("task op call:not enough more tokens")
	}

	accountFrom.assets[id] -= amount
	_, ok = accountTo.assets[id]
	if !ok {
		accountTo.assets[id] = 0
	}
	accountTo.assets[id] += amount

	return nil
}

/* ------ push ------ */
func (task *task_T) pushsb(params []byte, p0, p1, length int) error {
	copy(task.vData[p1:p1 + length], params[p0:p0 + length])

	return nil
}

func (task *task_T) push8(params []byte, p0, p1 int) error {
	task.vData[p1] = params[p0]

	return nil
}

func (task *task_T) push16(params []byte, p0, p1 int) error {
	copy(task.vData[p1:p1 + 2], params[p0:p0 + 2])

	return nil
}

func (task *task_T) push32(params []byte, p0, p1 int) error {
	copy(task.vData[p1:p1 + 4], params[p0:p0 + 4])

	return nil
}

func (task *task_T) push64(params []byte, p0, p1 int) error {
	copy(task.vData[p1:p1 + 8], params[p0:p0 + 8])

	return nil
}
