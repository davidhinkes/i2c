// Package i2c provides the T type, allowing SMB I2C writes.
package i2c

import (
	"fmt"
	"os"
)

/*
#include <linux/i2c.h>
#include <linux/i2c-dev.h>
#include <sys/ioctl.h>
#include <stdint.h>

static inline int ioctl_i2c_slave(int fd, uint8_t reg) {
	return ioctl(fd, I2C_SLAVE, reg);
}

static int ioctl_i2c_smbus(int fd, struct i2c_smbus_ioctl_data data) {
	return ioctl(fd, I2C_SMBUS, &data);
}
*/
import "C"

// T is the main type for this package, for which users should create with Make().
type T struct {
	file *os.File
}

func (i *T) Close() error {
	return i.file.Close()
}

func (i *T) setAddr(addr uint8) error {
	if _, err := C.ioctl_i2c_slave(C.int(i.file.Fd()), C.uint8_t(addr)); err != nil {
		return fmt.Errorf("ioctl_i2c_slave: %w", err)
	}
	return nil
}

// Write reads from an I2C device address addr with offset (aka register).
// The result is copied into dst. The size of dst dictates how many bytes
// are read.
// Function will return an error if len(dst) bytes cannot be read.
func (i *T) Read(dst []byte, addr uint8, offset uint8) error {
	if err := i.setAddr(addr); err != nil {
		return err
	}
	// See https://github.com/mozilla-b2g/i2c-tools/blob/master/lib/smbus.c, function i2c_smbus_read_i2c_block_data
	n := len(dst)
	ioctlData := C.struct_i2c_smbus_ioctl_data{
		read_write: C.I2C_SMBUS_READ,
		command:    C.uchar(offset),
		size:       C.I2C_SMBUS_I2C_BLOCK_DATA,
		data:       &C.union_i2c_smbus_data{byte(n)}, // just a [34]byte really
	}

	if _, err := C.ioctl_i2c_smbus(C.int(i.file.Fd()), ioctlData); err != nil {
		return fmt.Errorf("ioctl_i2c_smbus: %w", err)
	}
	if want, got := n, int(ioctlData.data[0]); want != got {
		return fmt.Errorf("wanted read of size %v bytes, actually read %v", want, got)
	}
	copy(dst, ioctlData.data[1:n+1])
	return nil
}

// Write writes src to I2C device address addr with offset (aka register).
// This will return an error if all len(src) bytes cannot be written.example:
func (i *T) Write(src []byte, addr uint8, offset uint8) error {
	if err := i.setAddr(addr); err != nil {
		return err
	}
	// See https://github.com/mozilla-b2g/i2c-tools/blob/master/lib/smbus.c, function i2c_smbus_write_i2c_block_data
	n := len(src)
	ioctlData := C.struct_i2c_smbus_ioctl_data{
		read_write: C.I2C_SMBUS_WRITE,
		command:    C.uchar(offset),
		size:       C.I2C_SMBUS_I2C_BLOCK_DATA,
		data:       &C.union_i2c_smbus_data{byte(n)}, // just a [34]byte really
	}
	copy(ioctlData.data[1:n+1], src)
	if _, err := C.ioctl_i2c_smbus(C.int(i.file.Fd()), ioctlData); err != nil {
		return fmt.Errorf("ioctl_i2c_smbus: %w", err)
	}
	if want, got := n, int(ioctlData.data[0]); want != got {
		return fmt.Errorf("wanted read of size %v bytes, actually read %v", want, got)
	}
	return nil
}

// Make creates a T given a path to the dev file.
func Make(path string) (T, error) {
	file, err := os.OpenFile(path, os.O_RDWR, 0)
	ret := T{
		file: file,
	}
	return ret, err
}
