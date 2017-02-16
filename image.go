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
	"bufio"
	"bytes"
	"container/list"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
	"os"
)

var Signature = []byte{0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a}

type Image struct {
	ChunkList *list.List
	Trailer   []byte
}

func Parse(data []byte) (*Image, error) {
	buf := bytes.NewBuffer(data)

	if !bytes.Equal(buf.Next(8), Signature) {
		return nil, errors.New("Incorrect signature")
	}

	chunkList := list.New()

	for {
		u32 := buf.Next(4)

		if len(u32) == 0 {
			break
		}

		if len(u32) != 4 {
			return nil, errors.New("Broken structure of chunk")
		}

		length := binary.BigEndian.Uint32(u32)
		if int(length) < 0 || int(length+8) < 0 || buf.Len() < int(length+8) {
			return nil, errors.New("Broken structure of chunk")
		}

		data := buf.Next(int(length + 4))
		crc := binary.BigEndian.Uint32(buf.Next(4))

		chunk := &Chunk{length, data, crc}
		chunkList.PushBack(chunk)

		if chunk.Type() == "IEND" {
			break
		}
	}

	trailer := buf.Next(buf.Len())
	return &Image{chunkList, trailer}, nil
}

func Load(path string) (*Image, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	image, err := Parse(data)
	if err != nil {
		return nil, errors.New(path + ": " + err.Error())
	}

	return image, nil
}

func (self *Image) dump(w io.Writer) error {
	u32 := make([]byte, 4)
	writer := &WriterWithError{writer: w}

	writer.Write(Signature)

	for e := self.ChunkList.Front(); e != nil; e = e.Next() {
		chunk := e.Value.(*Chunk)

		binary.BigEndian.PutUint32(u32, chunk.Length)
		writer.Write(u32)

		writer.Write(chunk.Data)

		binary.BigEndian.PutUint32(u32, chunk.Crc)
		writer.Write(u32)
	}

	writer.Write(self.Trailer)
	return writer.Error
}

func (self *Image) Write(path string) (retErr error) {
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); retErr == nil {
			retErr = err
		}
	}()

	w := bufio.NewWriterSize(file, 64*1024)
	if err := self.dump(w); err != nil {
		return err
	}

	if err := w.Flush(); err != nil {
		return err
	}

	return nil
}

func (self *Image) SetPhysChunk(phys *PhysChunk) error {
	physChunk := phys.GenerateChunk()

	for e := self.ChunkList.Front(); e != nil; e = e.Next() {
		chunk := e.Value.(*Chunk)

		if chunk.Type() == "pHYs" {
			self.ChunkList.InsertBefore(physChunk, e)
			self.ChunkList.Remove(e)
			return nil
		}
	}

	for e := self.ChunkList.Front(); e != nil; e = e.Next() {
		chunk := e.Value.(*Chunk)

		if chunk.Type() == "IDAT" {
			self.ChunkList.InsertBefore(physChunk, e)
			return nil
		}
	}

	return errors.New("IDAT chunk not found")
}
