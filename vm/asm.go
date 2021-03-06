/*
 * Copyright 2018 De-labtory
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package vm

import (
	"github.com/DE-labtory/koa/opcode"
)

var opCodes = map[opcode.Type]opCode{
	// 0x0 range
	opcode.Add: add{},
	opcode.Mul: mul{},
	opcode.Sub: sub{},
	opcode.Div: div{},
	opcode.Mod: mod{},
	opcode.And: and{},
	opcode.Or:  or{},

	// 0x10 range
	opcode.LT:  lt{},
	opcode.LTE: lte{},
	opcode.GT:  gt{},
	opcode.GTE: gte{},
	opcode.EQ:  eq{},
	opcode.NOT: not{},

	// 0x20 range
	opcode.Pop:       pop{},
	opcode.Push:      push{},
	opcode.Mload:     mload{},
	opcode.Mstore:    mstore{},
	opcode.Msize:     msize{},
	opcode.LoadFunc:  loadfunc{},
	opcode.LoadArgs:  loadargs{},
	opcode.Returning: returning{},
	opcode.Jump:      jump{},
	opcode.JumpDst:   jumpDst{},

	// 0x30 range
	opcode.Jumpi: jumpi{},
	opcode.DUP:   dup{},
	opcode.SWAP:  swap{},
	opcode.Exit:  exit{},
}

// Converts rawByteCode to assembly code.
func disassemble(rawByteCode []byte) (*asm, error) {
	asm := newAsm()

	for i := 0; i < len(rawByteCode); i++ {
		op, ok := opCodes[opcode.Type(rawByteCode[i])]

		if !ok {
			return nil, ErrInvalidOpcode
		}

		switch op.hex()[0] {
		case uint8(opcode.Push):
			body := make([]uint8, 0)
			body = append(body, rawByteCode[i+1:i+9]...)

			asm.code = append(asm.code, op)
			asm.code = append(asm.code, Data{Body: body})
			i += 8
		default:
			asm.code = append(asm.code, op)
		}
	}

	return asm, nil
}

// Do some analysis step (calculating the cost of running the code)
func analysis() {
}

// Assemble Reader read assembly codes and can jump to certain assembly code
type asmReader interface {
	next() hexer
	jump(i uint64)
	validateJumpDst(i uint64) bool
}

type Data struct {
	Body []uint8
}

func (d Data) hex() []uint8 {
	return d.Body
}

type hexer interface {
	hex() []uint8
}

type asm struct {
	code []hexer
	cost uint64
	pc   uint64
}

func newAsm() *asm {
	return &asm{
		code: make([]hexer, 0),
		cost: 0,
		pc:   0,
	}
}

func (a *asm) next() hexer {
	if a.pc+1 == uint64(len(a.code)) {
		return nil
	}

	code := a.code[a.pc+1]
	a.pc += 1
	return code
}

func (a *asm) jump(pc uint64) {
	if pc > uint64(len(a.code))-1 {
		panic("Access to invalid program counter!")
	}
	a.pc = pc
}

func (a *asm) validateJumpDst(pc uint64) bool {
	if opcode.Type(a.code[pc].hex()[0]) != opcode.JumpDst {
		return false
	}
	return true
}

func (a *asm) print() {

}
