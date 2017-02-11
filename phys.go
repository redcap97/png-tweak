/*
Copyright 2017 Akira Midorikawa

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"bytes"
	"encoding/binary"
	"hash/crc32"
)

type PhysChunk struct {
	X    uint32
	Y    uint32
	Unit uint8
}

func (self *PhysChunk) GenerateChunk() *Chunk {
	buf := bytes.NewBuffer(make([]byte, 0))
	u32 := make([]byte, 4)

	buf.WriteString("pHYs")

	binary.BigEndian.PutUint32(u32, self.X)
	buf.Write(u32)

	binary.BigEndian.PutUint32(u32, self.Y)
	buf.Write(u32)

	buf.WriteByte(self.Unit)

	data := buf.Bytes()
	return &Chunk{9, data, crc32.ChecksumIEEE(data)}
}
