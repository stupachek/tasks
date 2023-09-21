package database

import (
	"bytes"
	"encoding/binary"
	"epam/model"
	"io"
	"math"
	"os"
	"strings"
	"unsafe"

	"golang.org/x/exp/slices"
)

const (
	INSERT byte = 0
	DELETE byte = 1
)

func Insert(data model.User, fileName string) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0200)
	if err != nil {
		return err
	}
	row := make([]byte, 0)
	// row = append(row, INSERT)
	name := []byte(data.Name)
	row = append(row, byte(uint8(len(name))))
	row = append(row, name...)
	var activeAge uint64 = data.Age
	if data.Active {
		bitmask := uint64(math.Exp2(64))
		activeAge |= bitmask
	}
	binaryAge := make([]byte, 8)
	binary.BigEndian.PutUint64(binaryAge, activeAge)
	row = append(row, binaryAge...)
	binaryMass := make([]byte, 8)
	binary.BigEndian.PutUint64(binaryMass, math.Float64bits(data.Mass))
	row = append(row, binaryMass...)
	books := []byte(strings.Join(data.Books, ","))
	row = append(row, byte(uint8(len(books))))
	row = append(row, books...)
	_, err = file.Write(row)
	if err != nil {
		return err
	}
	err = file.Sync()
	if err != nil {
		return err
	}
	err = file.Close()
	return err
}

const (
	maxUsers int = 8
)

func Users(r io.Reader) []model.User {
	out := []model.User{}
	for i := 0; i < maxUsers; i++ {
		user, err := ReadUser(r)
		if err != nil || user.Name == "" {
			break
		}
		out = append(out, user)
	}
	return out
}

func ReadUser(r io.Reader) (model.User, error) {
	var nameLength uint8
	if err := binary.Read(r, binary.BigEndian, &nameLength); err != nil {
		return model.User{}, err
	}
	name := make([]byte, nameLength)
	if err := binary.Read(r, binary.BigEndian, &name); err != nil {
		return model.User{}, err
	}
	var ActiveAge uint64
	if err := binary.Read(r, binary.BigEndian, &ActiveAge); err != nil {
		return model.User{}, err
	}
	bitmask := uint64(math.Exp2(64))
	activeUser := uint8(ActiveAge & bitmask >> (64 - 1))
	age := ActiveAge << 1
	age = age >> 1
	var mass float64
	if err := binary.Read(r, binary.BigEndian, &mass); err != nil {
		return model.User{}, err
	}
	var bookLength uint8
	if err := binary.Read(r, binary.BigEndian, &bookLength); err != nil {
		return model.User{}, err
	}
	book := make([]byte, bookLength)
	if err := binary.Read(r, binary.BigEndian, &book); err != nil {
		return model.User{}, err
	}
	books := strings.Split(string(book), ",")
	return model.User{
		Name:   string(name),
		Age:    age,
		Active: *(*bool)(unsafe.Pointer(&activeUser)),
		Mass:   mass,
		Books:  books,
	}, nil
}

func Reader() io.Reader {
	return bytes.NewBuffer(
		[]byte{
			0x8, 0x4a, 0x6f, 0x68, 0x6e, 0x20, 0x44, 0x6f, 0x65, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1e, 0x40, 0x54, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x11, 0x48, 0x61, 0x72, 0x72, 0x79, 0x20, 0x50, 0x6f, 0x74, 0x74, 0x65, 0x72, 0x2c, 0x31, 0x39, 0x38, 0x34, 0x8, 0x4a, 0x61, 0x6b, 0x65, 0x20, 0x44, 0x6f, 0x65, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x14, 0x40, 0x4e, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xa, 0x20, 0x4a, 0x61, 0x6e, 0x65, 0x20, 0x44, 0x6f, 0x65, 0x20, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x96, 0x3f, 0xe8, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1c, 0x48, 0x61, 0x72, 0x72, 0x79, 0x20, 0x50, 0x6f, 0x74, 0x74, 0x65, 0x72, 0x2c, 0x47, 0x61, 0x6d, 0x65, 0x20, 0x6f, 0x66, 0x20, 0x54, 0x68, 0x72, 0x6f, 0x6e, 0x65, 0x73, 0x1, 0x9, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x5a, 0x40, 0xbf, 0x40, 0x0, 0x0, 0x0, 0x0, 0x0, 0xc, 0x48, 0x61,
			0x72, 0x72, 0x79, 0x20, 0x50, 0x6f, 0x74, 0x74, 0x65, 0x72, 0xd9, 0x56, 0x6d, 0x30, 0x77, 0x65, 0x45, 0x35, 0x47, 0x56, 0x58, 0x68, 0x55, 0x62, 0x6c, 0x4a, 0x57, 0x56, 0x30, 0x64, 0x53, 0x55, 0x46, 0x5a, 0x74, 0x4d, 0x57, 0x39, 0x57, 0x56, 0x6c, 0x6c, 0x33, 0x57, 0x6b, 0x52,
			0x53, 0x56, 0x31, 0x5a, 0x74, 0x65, 0x44, 0x42, 0x61, 0x52, 0x56, 0x59, 0x77, 0x56, 0x6d, 0x73, 0x78, 0x57, 0x47, 0x56, 0x47, 0x57, 0x6c, 0x5a, 0x69, 0x56, 0x45, 0x5a, 0x49, 0x57, 0x56, 0x64, 0x34, 0x53, 0x32, 0x4d, 0x78, 0x54, 0x6e, 0x4e, 0x69, 0x52, 0x30, 0x5a, 0x58, 0x56,
			0x6d, 0x35, 0x43, 0x62, 0x31, 0x5a, 0x73, 0x56, 0x6d, 0x46, 0x57, 0x4d, 0x56, 0x70, 0x57, 0x54, 0x56, 0x56, 0x57, 0x61, 0x47, 0x56, 0x71, 0x51, 0x54, 0x6b, 0x3d, 0xa, 0x56, 0x6d, 0x30, 0x77, 0x65, 0x45, 0x35, 0x47, 0x56, 0x58, 0x68, 0x55, 0x62, 0x6c, 0x4a, 0x57, 0x56, 0x30, 0x64, 0x53, 0x55, 0x46, 0x5a, 0x74, 0x4d, 0x57, 0x39, 0x57, 0x56, 0x6c, 0x6c, 0x33, 0x57, 0x6b, 0x52, 0x53, 0x56, 0x31, 0x5a, 0x74, 0x65, 0x44, 0x42, 0x61, 0x52, 0x56, 0x59, 0x77, 0x56, 0x6d, 0x73, 0x78, 0x57, 0x47, 0x56, 0x47, 0x57, 0x6c, 0x5a, 0x69, 0x56, 0x45, 0x5a, 0x49, 0x57, 0x56, 0x64, 0x34, 0x53, 0x32, 0x4d, 0x78, 0x54, 0x6e, 0x4e, 0x69, 0x52, 0x30, 0x5a, 0x58, 0x56, 0x6d, 0x35, 0x43, 0x62, 0x31, 0x5a, 0x73, 0x56, 0x6d, 0x46, 0x57, 0x4d, 0x56, 0x70, 0x57, 0x54, 0x56, 0x56, 0x57, 0x61, 0x47, 0x56, 0x71, 0x51, 0x54, 0x6b, 0x3d, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x10, 0x54, 0x68, 0x65, 0x20, 0x48, 0x75, 0x6e, 0x67, 0x65, 0x72, 0x20, 0x47, 0x61, 0x6d, 0x65, 0x73, 0x8, 0x0, 0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, 0x80, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x1b, 0x4d, 0x6f, 0x62, 0x79, 0x20, 0x44, 0x69, 0x63, 0x6b, 0x2c, 0x49, 0x74, 0x2c, 0x54, 0x68, 0x65, 0x20, 0x47, 0x72, 0x65, 0x65, 0x6e, 0x20, 0x4d, 0x69, 0x6c, 0x65,
		},
	)
}

func nearestMass(users []model.User, target float64) model.User {
	slices.SortFunc(users, func(a model.User, b model.User) bool {
		return a.Mass < b.Mass
	})
	i, ok := slices.BinarySearchFunc(users, target, func(u model.User, mass float64) int {
		if u.Mass < mass {
			return -1
		}
		if u.Mass > mass {
			return 1
		}
		return 0
	})
	if !ok && i != 0 {
		if math.Abs(target-users[i-1].Mass) < math.Abs(target-users[i].Mass) {
			return users[i-1]
		}
	}
	return users[i]
}
