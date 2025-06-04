// dsysb

package main

import (
	"encoding/binary"
)

/* ------ mov ------ */
func (task *task_T) movsb(p0, p1, length int) {
	copy(task.vData[p1:p1 + length], task.vData[p0:p0 + length])
}

func (task *task_T) mov8(p0, p1 int) {
	task.vData[p1] = task.vData[p0]
}

func (task *task_T) mov16(p0, p1 int) {
	copy(task.vData[p1:p1 + 2], task.vData[p0:p0 + 2])
}

func (task *task_T) mov32(p0, p1 int) {
	copy(task.vData[p1:p1 + 4], task.vData[p0:p0 + 4])
}

func (task *task_T) mov64(p0, p1 int) {
	copy(task.vData[p1:p1 + 8], task.vData[p0:p0 + 8])
}

/* ------ add ------ */
func (task *task_T) add8(p0, p1, p2 int) {
	task.vData[p2] = byte(int8(task.vData[p0]) + int8(task.vData[p1]))
}

func (task *task_T) add16(p0, p1, p2 int) {
	x := int16(binary.LittleEndian.Uint16(task.vData[p0:p0 + 2]))
	y := int16(binary.LittleEndian.Uint16(task.vData[p1:p1 + 2]))
	binary.LittleEndian.PutUint16(task.vData[p2:p2 + 2], uint16(x + y))
}

func (task *task_T) add32(p0, p1, p2 int) {
	x := int32(binary.LittleEndian.Uint32(task.vData[p0:p0 + 4]))
	y := int32(binary.LittleEndian.Uint32(task.vData[p1:p1 + 4]))
	binary.LittleEndian.PutUint32(task.vData[p2:p2 + 4], uint32(x + y))
}

func (task *task_T) add64(p0, p1, p2 int) {
	x := int64(binary.LittleEndian.Uint64(task.vData[p0:p0 + 8]))
	y := int64(binary.LittleEndian.Uint64(task.vData[p1:p1 + 8]))
	binary.LittleEndian.PutUint64(task.vData[p2:p2 + 8], uint64(x + y))
}

func (task *task_T) add8u(p0, p1, p2 int) {
	task.vData[p2] = byte(uint8(task.vData[p0]) + uint8(task.vData[p1]))
}

func (task *task_T) add16u(p0, p1, p2 int) {
	x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
	y := binary.LittleEndian.Uint16(task.vData[p1:p1 + 2])
	binary.LittleEndian.PutUint16(task.vData[p2:p2 + 2], x + y)
}

func (task *task_T) add32u(p0, p1, p2 int) {
	x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
	y := binary.LittleEndian.Uint32(task.vData[p1:p1 + 4])
	binary.LittleEndian.PutUint32(task.vData[p2:p2 + 4], x + y)
}

func (task *task_T) add64u(p0, p1, p2 int) {
	x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
	y := binary.LittleEndian.Uint64(task.vData[p1:p1 + 8])
	binary.LittleEndian.PutUint64(task.vData[p2:p2 + 8], x + y)
}

/* ------ sub ------ */
func (task *task_T) sub8(p0, p1, p2 int) {
	task.vData[p2] = byte(int8(task.vData[p0]) - int8(task.vData[p1]))
}

func (task *task_T) sub16(p0, p1, p2 int) {
	x := int16(binary.LittleEndian.Uint16(task.vData[p0:p0 + 2]))
	y := int16(binary.LittleEndian.Uint16(task.vData[p1:p1 + 2]))
	binary.LittleEndian.PutUint16(task.vData[p2:p2 + 2], uint16(x - y))
}

func (task *task_T) sub32(p0, p1, p2 int) {
	x := int32(binary.LittleEndian.Uint32(task.vData[p0:p0 + 4]))
	y := int32(binary.LittleEndian.Uint32(task.vData[p1:p1 + 4]))
	binary.LittleEndian.PutUint32(task.vData[p2:p2 + 4], uint32(x - y))
}

func (task *task_T) sub64(p0, p1, p2 int) {
	x := int64(binary.LittleEndian.Uint64(task.vData[p0:p0 + 8]))
	y := int64(binary.LittleEndian.Uint64(task.vData[p1:p1 + 8]))
	binary.LittleEndian.PutUint64(task.vData[p2:p2 + 8], uint64(x - y))
}

func (task *task_T) sub8u(p0, p1, p2 int) {
	task.vData[p2] = byte(uint8(task.vData[p0]) - uint8(task.vData[p1]))
}

func (task *task_T) sub16u(p0, p1, p2 int) {
	x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
	y := binary.LittleEndian.Uint16(task.vData[p1:p1 + 2])
	binary.LittleEndian.PutUint16(task.vData[p2:p2 + 2], x - y)
}

func (task *task_T) sub32u(p0, p1, p2 int) {
	x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
	y := binary.LittleEndian.Uint32(task.vData[p1:p1 + 4])
	binary.LittleEndian.PutUint32(task.vData[p2:p2 + 4], x - y)
}

func (task *task_T) sub64u(p0, p1, p2 int) {
	x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
	y := binary.LittleEndian.Uint64(task.vData[p1:p1 + 8])
	binary.LittleEndian.PutUint64(task.vData[p2:p2 + 8], x - y)
}

/* ------ mul ------ */
func (task *task_T) mul8(p0, p1, p2 int) {
	task.vData[p2] = byte(int8(task.vData[p0]) * int8(task.vData[p1]))
}

func (task *task_T) mul16(p0, p1, p2 int) {
	x := int16(binary.LittleEndian.Uint16(task.vData[p0:p0 + 2]))
	y := int16(binary.LittleEndian.Uint16(task.vData[p1:p1 + 2]))
	binary.LittleEndian.PutUint16(task.vData[p2:p2 + 2], uint16(x * y))
}

func (task *task_T) mul32(p0, p1, p2 int) {
	x := int32(binary.LittleEndian.Uint32(task.vData[p0:p0 + 4]))
	y := int32(binary.LittleEndian.Uint32(task.vData[p1:p1 + 4]))
	binary.LittleEndian.PutUint32(task.vData[p2:p2 + 4], uint32(x * y))
}

func (task *task_T) mul64(p0, p1, p2 int) {
	x := int64(binary.LittleEndian.Uint64(task.vData[p0:p0 + 8]))
	y := int64(binary.LittleEndian.Uint64(task.vData[p1:p1 + 8]))
	binary.LittleEndian.PutUint64(task.vData[p2:p2 + 8], uint64(x * y))
}

func (task *task_T) mul8u(p0, p1, p2 int) {
	task.vData[p2] = byte(uint8(task.vData[p0]) * uint8(task.vData[p1]))
}

func (task *task_T) mul16u(p0, p1, p2 int) {
	x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
	y := binary.LittleEndian.Uint16(task.vData[p1:p1 + 2])
	binary.LittleEndian.PutUint16(task.vData[p2:p2 + 2], x * y)
}

func (task *task_T) mul32u(p0, p1, p2 int) {
	x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
	y := binary.LittleEndian.Uint32(task.vData[p1:p1 + 4])
	binary.LittleEndian.PutUint32(task.vData[p2:p2 + 4], x * y)
}

func (task *task_T) mul64u(p0, p1, p2 int) {
	x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
	y := binary.LittleEndian.Uint64(task.vData[p1:p1 + 8])
	binary.LittleEndian.PutUint64(task.vData[p2:p2 + 8], x * y)
}

/* ------ quo ------ */
func (task *task_T) quo8(p0, p1, p2, p3 int) {
	task.vData[p2] = byte(int8(task.vData[p0]) / int8(task.vData[p1]))
	task.vData[p3] = byte(int8(task.vData[p0]) % int8(task.vData[p1]))
}

func (task *task_T) quo16(p0, p1, p2, p3 int) {
	x := int16(binary.LittleEndian.Uint16(task.vData[p0:p0 + 2]))
	y := int16(binary.LittleEndian.Uint16(task.vData[p1:p1 + 2]))
	binary.LittleEndian.PutUint16(task.vData[p2:p2 + 2], uint16(x / y))
	binary.LittleEndian.PutUint16(task.vData[p3:p3 + 2], uint16(x % y))
}

func (task *task_T) quo32(p0, p1, p2, p3 int) {
	x := int32(binary.LittleEndian.Uint32(task.vData[p0:p0 + 4]))
	y := int32(binary.LittleEndian.Uint32(task.vData[p1:p1 + 4]))
	binary.LittleEndian.PutUint32(task.vData[p2:p2 + 4], uint32(x / y))
	binary.LittleEndian.PutUint32(task.vData[p3:p3 + 4], uint32(x % y))
}

func (task *task_T) quo64(p0, p1, p2, p3 int) {
	x := int64(binary.LittleEndian.Uint64(task.vData[p0:p0 + 8]))
	y := int64(binary.LittleEndian.Uint64(task.vData[p1:p1 + 8]))
	binary.LittleEndian.PutUint64(task.vData[p2:p2 + 8], uint64(x / y))
	binary.LittleEndian.PutUint64(task.vData[p3:p3 + 8], uint64(x % y))
}

func (task *task_T) quo8u(p0, p1, p2, p3 int) {
	task.vData[p2] = byte(uint8(task.vData[p0]) / uint8(task.vData[p1]))
	task.vData[p3] = byte(uint8(task.vData[p0]) % uint8(task.vData[p1]))
}

func (task *task_T) quo16u(p0, p1, p2, p3 int) {
	x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
	y := binary.LittleEndian.Uint16(task.vData[p1:p1 + 2])
	binary.LittleEndian.PutUint16(task.vData[p2:p2 + 2], x / y)
	binary.LittleEndian.PutUint16(task.vData[p3:p3 + 2], x % y)
}

func (task *task_T) quo32u(p0, p1, p2, p3 int) {
	x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
	y := binary.LittleEndian.Uint32(task.vData[p1:p1 + 4])
	binary.LittleEndian.PutUint32(task.vData[p2:p2 + 4], x / y)
	binary.LittleEndian.PutUint32(task.vData[p3:p3 + 4], x % y)
}

func (task *task_T) quo64u(p0, p1, p2, p3 int) {
	x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
	y := binary.LittleEndian.Uint64(task.vData[p1:p1 + 8])
	binary.LittleEndian.PutUint64(task.vData[p2:p2 + 8], x / y)
	binary.LittleEndian.PutUint64(task.vData[p3:p3 + 8], x % y)
}

/* ------ inc ------ */
func (task *task_T) inc8(p0 int) {
	task.vData[p0] = byte(int8(task.vData[p0]) + 1)
}

func (task *task_T) inc16(p0 int) {
	x := int16(binary.LittleEndian.Uint16(task.vData[p0:p0 + 2]))
	binary.LittleEndian.PutUint16(task.vData[p0:p0 + 2], uint16(x + 1))
}

func (task *task_T) inc32(p0 int) {
	x := int32(binary.LittleEndian.Uint32(task.vData[p0:p0 + 4]))
	binary.LittleEndian.PutUint32(task.vData[p0:p0 + 4], uint32(x + 1))
}

func (task *task_T) inc64(p0 int) {
	x := int64(binary.LittleEndian.Uint64(task.vData[p0:p0 + 8]))
	binary.LittleEndian.PutUint64(task.vData[p0:p0 + 8], uint64(x + 1))
}

func (task *task_T) inc8u(p0 int) {
	task.vData[p0] = byte(uint8(task.vData[p0]) + 1)
}

func (task *task_T) inc16u(p0 int) {
	x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
	binary.LittleEndian.PutUint16(task.vData[p0:p0 + 2], x + 1)
}

func (task *task_T) inc32u(p0 int) {
	x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
	binary.LittleEndian.PutUint32(task.vData[p0:p0 + 4], x + 1)
}

func (task *task_T) inc64u(p0 int) {
	x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
	binary.LittleEndian.PutUint64(task.vData[p0:p0 + 8], x + 1)
}

/* ------ dec ------ */
func (task *task_T) dec8(p0 int) {
	task.vData[p0] = byte(int8(task.vData[p0]) + 1)
}

func (task *task_T) dec16(p0 int) {
	x := int16(binary.LittleEndian.Uint16(task.vData[p0:p0 + 2]))
	binary.LittleEndian.PutUint16(task.vData[p0:p0 + 2], uint16(x + 1))
}

func (task *task_T) dec32(p0 int) {
	x := int32(binary.LittleEndian.Uint32(task.vData[p0:p0 + 4]))
	binary.LittleEndian.PutUint32(task.vData[p0:p0 + 4], uint32(x + 1))
}

func (task *task_T) dec64(p0 int) {
	x := int64(binary.LittleEndian.Uint64(task.vData[p0:p0 + 8]))
	binary.LittleEndian.PutUint64(task.vData[p0:p0 + 8], uint64(x + 1))
}

func (task *task_T) dec8u(p0 int) {
	task.vData[p0] = byte(uint8(task.vData[p0]) - 1)
}

func (task *task_T) dec16u(p0 int) {
	x := binary.LittleEndian.Uint16(task.vData[p0:p0 + 2])
	binary.LittleEndian.PutUint16(task.vData[p0:p0 + 2], x - 1)
}

func (task *task_T) dec32u(p0 int) {
	x := binary.LittleEndian.Uint32(task.vData[p0:p0 + 4])
	binary.LittleEndian.PutUint32(task.vData[p0:p0 + 4], x - 1)
}

func (task *task_T) dec64u(p0 int) {
	x := binary.LittleEndian.Uint64(task.vData[p0:p0 + 8])
	binary.LittleEndian.PutUint64(task.vData[p0:p0 + 8], x - 1)
}

func (task *task_T) getIndex(p0 int) error {
	index, err := getIndex()
	if err != nil {
		return err
	}

	binary.LittleEndian.PutUint32(task.vData[p0:p0 + 4], index)

	return nil
}
