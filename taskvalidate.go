// dsysb

package main

import (
	"encoding/binary"
	"errors"
	"fmt"
)

func (task *task_T) validate() error {
	// variable ip int for instructs

	instructsLength := len(task.instructs)

	insStartPositions := []uint16{}
	p2sOfCompare := []uint16{}

	for ip := 0; ip < instructsLength; {
		ipx := ip
		ip++
		if instructsLength == ip {
			return nil
		}

		insStartPositions = append(insStartPositions, uint16(ipx))

		var aip int
		var err error
		switch task.instructs[ipx] {
		case ins_movsb:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // source
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // distination
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // length
			ip += 2
			err = task.opCheck(p2, p0, p1)
		case ins_mov8:
			aip = ip + 4
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // source
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // distination
			ip += 2
			err = task.opCheck(1, p0, p1)
		case ins_mov16:
			aip = ip + 4
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // source
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // distination
			ip += 2
			err = task.opCheck(2, p0, p1)
		case ins_mov32:
			aip = ip + 4
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // source
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // distination
			ip += 2
			err = task.opCheck(4, p0, p1)
		case ins_mov64:
			aip = ip + 4
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // source
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // distination
			ip += 2
			err = task.opCheck(8, p0, p1)
		case ins_add8:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // sum
			ip += 2
			err = task.opCheck(1, p0, p1, p2)
		case ins_add16:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // sum
			ip += 2
			err = task.opCheck(2, p0, p1, p2)
		case ins_add32:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // sum
			ip += 2
			err = task.opCheck(4, p0, p1, p2)
		case ins_add64:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // sum
			ip += 2
			err = task.opCheck(8, p0, p1, p2)
		case ins_add8u:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // sum
			ip += 2
			err = task.opCheck(1, p0, p1, p2)
		case ins_add16u:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // sum
			ip += 2
			err = task.opCheck(2, p0, p1, p2)
		case ins_add32u:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // sum
			ip += 2
			err = task.opCheck(4, p0, p1, p2)
		case ins_add64u:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // adder
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // sum
			ip += 2
			err = task.opCheck(8, p0, p1, p2)
		case ins_sub8:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(1, p0, p1, p2)
		case ins_sub16:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(2, p0, p1, p2)
		case ins_sub32:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(4, p0, p1, p2)
		case ins_sub64:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(8, p0, p1, p2)
		case ins_sub8u:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(1, p0, p1, p2)
		case ins_sub16u:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(2, p0, p1, p2)
		case ins_sub32u:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(4, p0, p1, p2)
		case ins_sub64u:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(8, p0, p1, p2)
		case ins_mul8:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(1, p0, p1, p2)
		case ins_mul16:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(2, p0, p1, p2)
		case ins_mul32:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(4, p0, p1, p2)
		case ins_mul64:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(8, p0, p1, p2)
		case ins_mul8u:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(1, p0, p1, p2)
		case ins_mul16u:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(2, p0, p1, p2)
		case ins_mul32u:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(4, p0, p1, p2)
		case ins_mul64u:
			aip = ip + 6
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(8, p0, p1, p2)
		case ins_quo8:
			aip = ip + 8
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(1, p0, p1, p2, p3)
		case ins_quo16:
			aip = ip + 8
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(2, p0, p1, p2, p3)
		case ins_quo32:
			aip = ip + 8
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(4, p0, p1, p2, p3)
		case ins_quo64:
			aip = ip + 8
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(8, p0, p1, p2, p3)
		case ins_quo8u:
			aip = ip + 8
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(1, p0, p1, p2, p3)
		case ins_quo16u:
			aip = ip + 8
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(2, p0, p1, p2, p3)
		case ins_quo32u:
			aip = ip + 8
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(4, p0, p1, p2, p3)
		case ins_quo64u:
			aip = ip + 8
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(8, p0, p1, p2, p3)
		case ins_inc8:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(1, p0)
		case ins_inc16:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(2, p0)
		case ins_inc32:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(4, p0)
		case ins_inc64:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(8, p0)
		case ins_inc8u:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(1, p0)
		case ins_inc16u:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(2, p0)
		case ins_inc32u:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(4, p0)
		case ins_inc64u:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(8, p0)
		case ins_dec8:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(1, p0)
		case ins_dec16:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(2, p0)
		case ins_dec32:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(4, p0)
		case ins_dec64:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(8, p0)
		case ins_dec8u:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(1, p0)
		case ins_dec16u:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(2, p0)
		case ins_dec32u:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(4, p0)
		case ins_dec64u:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(8, p0)
		case ins_write_uint8:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // distination 
			ip += 2
			err = task.opCheck(1, p0)
		case ins_write_uint16:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(2, p0)
		case ins_write_uint32:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(4, p0)
		case ins_write_uint64:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(8, p0)
		case ins_read_uint8:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) //  source
			ip += 2
			err = task.opCheck(1, p0)
		case ins_read_uint16:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(2, p0)
		case ins_read_uint32:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(4, p0)
		case ins_read_uint64:
			aip = ip + 2
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			err = task.opCheck(8, p0)
		case ins_eq:
			aip = ip + 8
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			if p2 < ip {
				return errors.New(fmt.Sprintf("jump error at ip:%d, invalid index", aip))
			}
			p2sOfCompare = append(p2sOfCompare, uint16(p2))
			err = task.opCheckInnerA(aip, p0, p1, p2, p3)
		case ins_gt:
			aip = ip + 8
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			if p2 < ip {
				return errors.New(fmt.Sprintf("jump error at ip:%d, invalid index", aip))
			}
			p2sOfCompare = append(p2sOfCompare, uint16(p2))
			err = task.opCheckInnerB(aip, p0, p1, p2, p3)
		case ins_lt:
			aip = ip + 8
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			if p2 < ip {
				return errors.New(fmt.Sprintf("jump error at ip:%d, invalid index", aip))
			}
			p2sOfCompare = append(p2sOfCompare, uint16(p2))
			err = task.opCheckInnerB(aip, p0, p1, p2, p3)
		case ins_gteq:
			aip = ip + 8
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			if p2 < ip {
				return errors.New(fmt.Sprintf("jump error at ip:%d, invalid index", aip))
			}
			p2sOfCompare = append(p2sOfCompare, uint16(p2))
			err = task.opCheckInnerB(aip, p0, p1, p2, p3)
		case ins_lteq:
			aip = ip + 8
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			if p2 < ip {
				return errors.New(fmt.Sprintf("jump error at ip:%d, invalid index", aip))
			}
			p2sOfCompare = append(p2sOfCompare, uint16(p2))
			err = task.opCheckInnerB(aip, p0, p1, p2, p3)
		case ins_eq_bytes:
			aip = ip + 8
			if instructsLength < aip {
				return errors.New(fmt.Sprintf("Instruction error at ip:%d", aip))
			}
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			p3 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2]))
			ip += 2
			if p2 < ip {
				return errors.New(fmt.Sprintf("jump error at ip:%d, invalid index or jump backward", aip))
			}
			p2sOfCompare = append(p2sOfCompare, uint16(p2))
			err = task.opCheck(p3, p0, p1, p2)
		case ins_height:
		case ins_transfer_dsb_from_caller:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // amount
			ip += 2
			err = task.opCheck(8, p0)
		case ins_transfer_dsb_to_caller:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // amount
			ip += 2
			err = task.opCheck(8, p0)
		case ins_transfer_from_caller:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // asset id
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // amount
			ip += 2
			err = task.opCheck(32, p0)
			if err != nil {
				return err
			}
			err = task.opCheck(8, p1)
		case ins_transfer_to_caller:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // asset id
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // amount
			ip += 2
			err = task.opCheck(32, p0)
			if err != nil {
				return err
			}
			err = task.opCheck(8, p1)
		case ins_pushsb:
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // vdata position
			ip += 2
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // length
			ip += 2
			err = task.opCheck(p2, p1)
		case ins_push8:
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // vdata position
			ip += 2
			err = task.opCheck(1, p1)
		case ins_push16:
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // vdata position
			ip += 2
			err = task.opCheck(2, p1)
		case ins_push32:
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // vdata position
			ip += 2
			err = task.opCheck(4, p1)
		case ins_push64:
			ip += 2
			p1 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // vdata position
			ip += 2
			err = task.opCheck(8, p1)
		default:
			return errors.New("Invalid instruction")
		}

		if err != nil {
			return err
		}
	}

	for _, p2 := range p2sOfCompare {
		var ok bool
		for _, index := range insStartPositions {
			if p2 == index {
				ok = true
				break
			}
		}
		if !ok {
			return errors.New(fmt.Sprintf("Jmp to an invalid index: %d", p2))
		}
	}

	return nil
}

func (task *task_T)opCheckInnerA(aip, p0, p1, p2, p3 int) error {
	err := task.opCheck(1, p0, p1, p2, p3)
	if err != nil {
		return err
	}

	flag := uint8(task.vData[p3])

	switch flag {
	case 0:
		return nil
	case 1:
		return task.opCheck(2, p0, p1, p2)
	case 2:
		return task.opCheck(4, p0, p1, p2)
	case 3:
		return task.opCheck(8, p0, p1, p2)
	default:
		return errors.New(fmt.Sprintf("Wrong type of task op eq at ip :%d", aip))
	}
}

func (task *task_T)opCheckInnerB(aip, p0, p1, p2, p3 int) error {
	err := task.opCheck(1, p0, p1, p2, p3)
	if err != nil {
		return err
	}

	flag := uint8(task.vData[p3])

	switch flag {
	case 0:
		return nil
	case 1:
		return task.opCheck(2, p0, p1, p2)
	case 2:
		return task.opCheck(4, p0, p1, p2)
	case 3:
		return task.opCheck(8, p0, p1, p2)
	case 4:
		return nil
	case 5:
		return task.opCheck(2, p0, p1, p2)
	case 6:
		return task.opCheck(4, p0, p1, p2)
	case 7:
		return task.opCheck(8, p0, p1, p2)
	default:
		return errors.New(fmt.Sprintf("Wrong type of task op compare at ip :%d", aip))
	}
}

func (task *task_T) paramsCheck(pLength, length int, args ...int) error {
	for k, p := range args {
		limit := p + length
		if pLength < limit {
			fmt.Println(pLength, limit)
			return errors.New(fmt.Sprintf("params error at p%d: %d", k, p))
		}
	}

	return nil
}

func (task *task_T) validateCall(params []byte) error {
	// variable ip int for instructs

	instructsLength := len(task.instructs)

	for ip := 0; ip < instructsLength; {
		ipx := ip
		ip++
		if instructsLength == ip {
			return nil
		}

		var err error
		switch task.instructs[ipx] {
		case ins_pushsb:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // params position
			ip += 4
			p2 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // length
			ip += 2
			err = task.paramsCheck(p2, p0)
		case ins_push8:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // params position
			ip += 4
			err = task.paramsCheck(1, p0)
		case ins_push16:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // params position
			ip += 4
			err = task.paramsCheck(2, p0)
		case ins_push32:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // params position
			ip += 4
			err = task.paramsCheck(4, p0)
		case ins_push64:
			p0 := int(binary.LittleEndian.Uint16(task.instructs[ip:ip + 2])) // params position
			ip += 4
			err = task.paramsCheck(8, p0)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
