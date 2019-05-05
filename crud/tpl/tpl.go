// Code generated by go-bindata.
// sources:
// mgo.go.tpl
// mock.go.tpl
// DO NOT EDIT!

package tpl

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data, name string) ([]byte, error) {
	gz, err := gzip.NewReader(strings.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _mgoGoTpl = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xec\x58\x5b\x6f\xdb\x36\x14\x7e\xb6\x7e\xc5\xa9\x52\x04\x52\x6a\x48\xd9\xab\x56\x07\x18\x62\xb7\x33\x16\x24\x45\x93\x62\xc0\x8a\xa2\xd5\xe5\xd8\x61\x2b\x91\x2e\x45\xd7\x2b\x5c\xff\xf7\x81\x17\x59\xb4\x2e\xb6\x9b\x04\xd8\x1e\xe6\x27\x99\x97\x73\x0e\xbf\xef\x23\x0f\x0f\xc3\x10\x2e\x59\x86\x30\x47\x8a\x3c\x16\x98\xfd\x0a\x63\x06\x94\x09\xe0\x58\xb5\x81\xb8\x47\x60\xdf\x90\xaf\x38\x11\x08\xf1\x4c\x20\x07\xcc\x88\x20\x74\x1e\x38\xce\x22\x4e\xbf\xc4\x73\x84\x97\x10\xbc\x31\x9f\x17\x8e\x43\x8a\x05\xe3\x02\x3c\x67\xe0\x0a\x52\xa0\xeb\x38\x03\x77\x4e\xc4\xfd\x32\x09\x52\x56\x84\xf3\x9c\x25\x71\x5e\x92\x39\x0d\x8b\x39\x73\xf7\x75\x86\x49\xc9\xa8\xeb\xf8\x8e\x13\x86\xd2\xc9\xbb\xc5\x02\xf9\xef\xcb\x62\x01\x17\x7f\x12\x71\x3f\x1d\x03\x29\x1b\xed\xb0\x22\xe2\x1e\xa6\x63\x47\x7c\x5f\x60\xf7\x9c\x52\xf0\x65\x2a\x60\xed\x0c\xa6\x63\x90\x1e\x82\x9b\xe4\x33\xa6\x62\x9a\x81\xfc\x7d\x92\x4d\x91\xfb\x91\x64\x43\x56\x10\x81\xc5\x42\x7c\x77\xe1\xb3\x6a\x7c\x09\xc1\x15\x5b\x21\xbf\xa5\xf1\x17\x84\x8b\x8f\x24\x73\x3f\x39\x83\x46\x04\x60\x9b\x19\x12\x9a\x13\x8a\x72\xd8\x25\xc7\x58\xe0\x1d\x29\x10\x24\x2e\x81\xfa\x32\xc3\x52\xd5\xf7\x51\xb6\xb7\xbd\x5a\x9d\xd2\xce\xbb\x45\xd6\x6b\x67\xa9\xfa\x7a\xec\x58\x9d\xee\x27\x67\xd3\x81\xea\x5b\x4c\x19\xcf\x24\xaa\x5c\x7f\xb1\x99\x92\xc0\xee\xa8\x2e\x6c\xcd\xcc\x1d\x6c\x1b\xbf\x26\xd4\x3f\x03\xb4\x0e\xa7\x13\x6f\xe9\xa8\xc7\x74\x9b\xae\x63\x29\x7d\x8b\x29\x52\x61\x47\x7f\xd6\x60\xd9\xb8\xe0\x6a\x60\xdb\xac\x6e\x57\xb4\x2f\x39\xdf\xb5\xd5\x63\x2a\xd5\x03\x3b\xf8\xd7\x1d\x75\x5c\x8a\x71\xf3\xab\x25\xb0\xa3\x3b\xed\xbf\x47\x08\x56\xa7\x15\xa1\x65\xa4\xc7\xa8\x09\xa4\x4f\xa6\x56\xaf\x34\x2b\xe7\x97\x3b\x0a\x20\xbb\x90\x56\x66\xe5\x84\xb2\x6d\x4f\x35\xf7\x08\xf5\x16\xf9\x37\x92\xa2\x54\x6a\x69\x3e\xbb\xa5\x1a\x86\x70\xb2\x88\xc5\x7d\xe4\x86\x0d\xa6\x43\xf7\xa4\x4b\xc8\x95\xe5\x5a\xc9\x59\x52\xf1\x56\xcc\x59\x70\xc9\xf2\x1c\x53\x41\x18\x95\x3d\x46\xf6\xcd\x1e\x1d\xf2\x35\xae\xba\x6d\xeb\x93\x00\x62\xa0\xb8\xea\x76\xef\xcc\x96\x34\xed\x35\xe0\x65\x49\xd3\xa5\x0f\xde\x59\xe7\xd8\x21\x20\xe7\x8c\xfb\x7a\x25\x26\xde\x68\x04\x59\x12\x8c\x63\x11\x27\x71\x89\xc1\xa5\x97\x25\xc1\x75\x5c\x20\xbc\x00\xd7\xec\x35\xd7\xaf\xc7\x07\x13\x5a\x2e\x39\x4e\x69\x86\x7f\x7b\xd2\xaf\xfa\x5a\xff\x81\xdf\x23\x78\xff\xa1\x14\x9c\xd0\xf9\xba\x6b\x27\x6d\x36\xbe\x33\xe0\x28\x96\x9c\xc2\x69\x67\x78\x6b\x67\x30\xc8\x92\xc8\x20\x9c\x25\x43\xf5\x5f\xbb\x8d\xa0\xfa\x1a\x3a\x83\xcd\x10\x28\xc9\x0d\xb2\x5b\x00\x3b\xd8\xe6\x6c\x29\x30\x72\xdf\xdc\xdc\xde\x81\xe4\x58\x21\xe9\x95\xcd\x5d\x67\xfc\xfb\xc6\x96\x57\x45\x6f\xf6\x64\x63\xb4\x0f\x8d\x01\xad\xd4\x11\x9e\xc1\x09\x8d\x0b\xec\x3c\x51\x4e\xe0\x2c\x54\x44\x58\x64\xb4\xec\x8d\xb4\xc5\x6b\x5c\x55\x46\x3d\xdf\x19\x50\xb6\x92\x74\xe9\x2e\xb6\x92\x4d\xd2\xce\x08\xca\x20\x4b\x82\x29\x2d\x91\x0b\xef\xb4\x33\xd1\x49\x6c\xa7\xe3\x08\x9a\x9e\x24\xc6\x2f\x21\xb8\xe1\x64\x4e\x68\x9c\xc3\x45\xa4\x96\x6b\x8d\x91\x23\xea\x74\x15\x01\x65\x2b\xd9\x54\x67\x9e\xaa\x49\xf2\x4b\x66\x6a\x65\xcf\x46\x92\x20\xe5\xd4\x30\xee\xba\x6a\xcd\xce\x60\xb3\x15\x41\x2b\x14\x8b\x54\x6d\xbd\x6f\x13\x57\xb4\xbe\xbb\x83\x70\xdd\x86\x78\x73\x04\xd3\xda\xc1\x13\x10\x79\x50\x2b\xbb\x4c\x9b\x1c\xa1\x50\x8a\x24\x71\xaf\x51\xb4\xa2\xe8\x04\xd2\xb4\x8c\x46\x20\x77\xdd\x84\xf3\x6b\x26\x5e\xb1\x25\xcd\x54\x6f\x05\xaa\x84\x70\x20\x41\xae\x1a\x0c\xe8\x3b\x4a\xd1\x8b\x9f\x66\x2d\xc7\x43\xbd\xfe\xf1\x7a\xed\x3e\x2f\x51\xb8\xc3\xe6\x4e\xd5\x6a\x5a\x6b\xd1\x58\xed\x91\x49\x65\x0d\xdd\xd8\x22\xa9\x55\x2b\xb5\xb2\xd9\xab\x96\x6d\xd0\x55\x1e\xfc\x17\xf1\xd2\x47\xa0\x57\x11\x77\x6a\x42\x0a\x1a\x44\x1f\xb3\x1a\xcb\xa7\x96\xf9\x18\x73\x3c\x28\xf3\xf1\xe4\x6a\x72\x37\x79\xb0\xd2\xb5\x8f\x47\x2b\xfd\xbf\xa9\xe4\xb7\x58\xb0\x6f\x5d\x4a\x3e\x8a\x8f\x1e\x8a\x29\xc9\x1f\x46\xe7\x6b\x14\x07\xb8\x7c\x3d\x79\xf8\x91\xd5\x85\xef\x03\x58\xdc\x7f\x60\xe9\x0d\xde\xcc\x4f\x5f\x35\xc3\x59\x12\xbc\x22\x34\xeb\x86\x5b\x83\xf9\x35\xb8\xa1\xa8\x72\x90\xd5\xbf\x0f\x4d\x4a\xf2\x03\x99\xc1\xce\x0b\x57\xa4\x14\x7b\x6e\x77\x36\xca\x47\x20\x2a\xad\x79\xa5\x88\xb9\xbe\xef\xd6\xe8\xa9\x36\x7d\x79\x35\xe9\x9a\x66\x8d\x22\xab\x1e\x8c\x34\xb3\x87\xb2\xd9\xac\x44\x31\x84\x9c\x14\x44\xc8\x6b\xee\x10\x4a\x59\x03\xeb\xab\x51\x8b\x80\x12\xde\x7f\x38\x8a\x83\x62\x9b\xfa\xc7\xeb\x8d\x02\xf4\xd9\x36\xf4\x60\x5a\xfe\x85\x9c\x79\x3e\xfc\xf8\x01\xcf\x4c\xb0\x75\xa3\x84\xbb\x38\xdf\x9d\xdf\x67\x40\x6d\xc1\xe2\x1c\x46\x10\x2f\x16\x48\x33\xaf\x38\xaf\xb2\x02\x4b\x27\x39\x16\x6b\xf7\xf9\x5c\xa0\x3b\x84\xed\x64\x79\x94\x0f\x2a\x93\x9d\xce\x0f\x19\xcc\x65\xaa\x31\x33\xb7\xd6\x0a\x6b\x4a\x73\x86\x5d\x06\x0f\xa1\x38\x97\x93\x36\x4d\x9d\x7a\x85\x1f\xdc\x7e\x21\x0b\x4f\x53\xe2\x07\x57\x92\x12\x4f\x11\xe3\x07\xb7\x8c\x0b\x4f\x32\x63\x69\xf7\xb7\x3c\x6f\x6a\xb7\x7c\x9c\x78\xcb\x9d\xab\x2a\x5b\xd2\x63\xe5\x9b\xca\xb1\xc7\x5c\x59\xe5\xb8\x27\x16\xb1\x0f\x9e\x72\xaf\xd5\xfb\xbf\x0c\x1f\x2d\xc3\xad\x38\xbe\x06\x9a\x2f\xdf\x28\x42\x97\x34\xfd\xa7\x5a\xf5\xf6\x92\x93\x52\x1c\x95\x47\x42\x53\xb1\x1d\x16\x4e\xed\xfa\x09\x6e\xc1\xad\x43\xaf\x75\xce\x69\x77\x1d\xa7\x9d\xa9\xea\xf6\xc8\xac\xb3\x9c\x6c\xde\xbc\xa7\xe3\x8d\x8d\xbc\x29\x56\x0f\x1e\x03\x7b\xf7\xbe\x89\xf9\x71\x27\x80\x31\x62\x9f\x03\xba\x69\xcf\x69\x50\xd1\xae\x76\xe1\xcf\xf0\x7e\xf4\xb1\x61\x85\xf0\x14\x97\x8a\xa3\x8e\x8b\xc7\xf3\xd8\xbd\x8f\x0e\xad\xd5\x5c\xef\x8e\xbb\xf8\x98\xeb\x7d\xbb\x90\x53\x6b\x92\x4b\x22\xb3\x66\xd1\x37\xea\x92\x85\x52\x84\xc2\xc5\xba\x23\xef\x2e\xe8\x67\x70\x09\x24\x30\xfe\x76\xd9\x07\xaf\xa7\x03\xbe\x7d\xde\x39\xed\xda\x70\xed\x1a\xae\xfd\x3c\x10\xe8\xf7\x01\xf3\x26\x18\xd5\x8f\x75\x55\x59\x56\x77\xea\x1a\xaf\x7e\xdb\xad\x0a\x3d\xfd\xf2\x17\xd9\x2f\x7d\x5a\x2b\x2f\xe0\x17\xd9\xad\x1f\x31\xed\xfe\xc6\x1e\x6c\xd4\x59\xf5\x1c\xcb\x63\x33\xec\xba\xec\x1c\xda\xfb\xd2\x82\xdf\xbc\x95\x68\x8c\x94\x88\xfe\x09\x00\x00\xff\xff\xe3\x2f\x3a\x1b\x03\x19\x00\x00"

func mgoGoTplBytes() ([]byte, error) {
	return bindataRead(
		_mgoGoTpl,
		"mgo.go.tpl",
	)
}

func mgoGoTpl() (*asset, error) {
	bytes, err := mgoGoTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "mgo.go.tpl", size: 6403, mode: os.FileMode(420), modTime: time.Unix(1555396366, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _mockGoTpl = "\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xcc\x55\xc1\x6e\xe2\x3a\x14\x5d\xc7\x5f\x71\x5e\x90\x50\x42\xf3\x48\xbb\xe5\x15\x36\x05\xf5\x21\x55\x6d\xa5\xb6\x7a\x8b\xaa\x6a\xf3\xc8\x0d\x78\x1a\x9c\xc8\x31\x30\x23\xe0\xdf\x47\x76\x42\x0b\xc4\xa5\x1d\x36\x33\xab\xd8\x37\xd7\xe7\x1e\x9f\x63\x5f\x87\x21\x2e\xb2\x98\x30\x26\x41\x32\x52\x14\xff\x83\x7e\x06\x91\x29\x48\xda\xc4\xa0\x26\x84\x6c\x4e\x72\x21\xb9\x22\x44\x89\x22\x09\x8a\xb9\xe2\x62\xdc\x66\x2c\x8f\x46\xaf\xd1\x98\x70\x8e\xf6\x6d\x35\xec\x31\xc6\xa7\x79\x26\x15\x3c\xe6\xb8\x24\x65\x26\x0b\x97\xf9\x8c\x85\xa1\x4e\x7b\xc8\x73\x92\xff\xce\xa6\x39\x7a\xff\x71\x35\x19\xf6\xc1\x8b\xbd\x38\x16\x5c\x4d\x30\xec\x33\xf5\x23\x27\xfb\x9a\x42\xc9\xd9\x48\x61\xc9\x1c\x0d\x20\x14\x5e\xbe\x15\x99\xe8\xb8\xe7\x68\x5f\x65\x0b\x92\x77\x22\x7a\x25\xf4\x9e\x79\x1c\x14\x4a\x72\x31\x76\x5f\x98\x73\x8e\xf6\x8d\xe4\x63\x2e\xa2\x14\x3d\xb6\xb6\x30\xba\x23\x39\xe7\x23\xd2\x94\x8a\x6a\x98\x25\x46\x82\xdd\x3c\xbd\xb2\x91\x47\x6a\xd2\x71\xc3\xbd\x92\xa1\xdb\xb0\xf1\xde\x20\xbf\x13\x8f\x23\x15\x15\x78\x7c\x6a\xd9\x36\x58\xb1\xbb\xa6\x85\x1d\xe6\x42\x92\x36\x27\x82\xa0\x85\xbd\x12\x4b\x66\x62\xf4\x21\x80\xe7\xc3\x6b\x59\xff\x04\x30\x96\xf9\x9a\xa2\x24\x35\x93\x02\x4d\x6b\xe2\x72\x1d\x40\xf0\xb4\x62\xfa\x46\xc8\x22\x94\xcc\x66\x8a\x3a\xee\xed\xcd\xdd\x3d\xb4\x3c\x86\x99\x57\xc0\x4e\xc0\xaf\xb0\xbc\x8d\xb0\xd5\xa1\x68\xed\xd8\xe7\xc3\x23\x29\xb7\xb8\xee\x66\x0f\xfb\xe8\x74\x91\x92\xf0\x8a\xb6\x11\xda\xc7\x09\xce\x4a\xd1\xf5\x9f\xa6\xf5\x54\x2d\x99\xe3\x0c\xfb\x1d\xd4\xb0\x02\xe6\xec\x9e\x9e\x8e\xa1\xb3\x95\x14\x30\x67\xcd\x9c\xaa\x18\xba\x88\xf2\x9c\x44\xbc\xa9\x1e\x40\x7f\xfc\x37\x45\xdf\x75\x7b\xc8\xe3\xcd\x2d\x3b\xa0\xdc\xc3\x3d\xc2\x65\xfd\x6c\xaf\xbf\x20\x66\x59\xc0\xab\x6d\x49\x5f\x9a\xb0\x85\x86\x88\xa6\x64\xbb\x37\x6e\x03\xad\x30\xc0\xaf\x99\xc0\x13\x9c\xa2\xd7\xad\x0b\xb8\x5a\xd5\x63\xbd\x5d\x7f\x56\x2b\x54\xe3\xc7\xfd\xd4\xbf\xcf\x9e\xd0\xed\x6a\xd5\x8c\x45\x95\x88\x65\x6b\x69\x5f\xd3\xc2\x73\x79\x8c\x38\xa3\xc2\x34\x2f\xfa\xce\x0b\xe5\xfa\xc6\x90\xb9\x36\xfb\x00\x2c\x73\xe6\xed\x9d\x2d\xa1\xbb\xef\xac\xc5\xb4\x3e\xa5\xf4\xa9\x69\xfd\xc1\xd5\xe0\x7e\x70\xb4\x6f\x65\x8d\x23\x7d\xfb\xf3\x8d\x39\x88\x68\xa4\xae\xcb\x7e\x49\xea\x13\xcd\x2f\x07\xc7\x5f\x94\x4b\x52\xc7\xab\x6d\xb9\x26\xb5\xf6\x62\x3a\xeb\xef\xf2\x44\xf0\x34\xf8\x92\x31\x55\xfe\x01\xf4\xed\xb6\x7f\xc5\x0b\x75\xe0\x89\xdc\xb6\xe5\x0b\x16\x68\x34\x2f\x4b\x92\x82\x54\x80\x94\x4f\xb9\xd2\x06\xd4\xe4\xfd\xf0\xd5\xdc\x57\x38\x4b\x12\x7d\xff\x4f\x99\x93\xf2\x69\x35\x4a\x32\x89\xe7\x00\xa6\x31\xc8\x48\x8c\x69\xb3\x55\xa3\x16\x4f\x30\xc7\x5f\xef\xea\xe9\x40\x49\x48\x47\x4f\xd1\x6c\x42\xcf\xf5\xa4\x0a\x9b\x2c\x5d\xe9\xe4\xc4\x8c\x46\x99\x50\x5c\xcc\x48\x4f\xd6\x15\x42\xb9\x95\xae\x06\x58\xad\xf4\x4c\xaf\x37\xc1\x72\x75\xca\xa7\xd5\xea\xfd\x8d\xbe\x3d\x24\x7b\x3f\x02\xcc\x7d\xb3\xa0\x44\xd7\xd8\x65\x91\x12\xd0\xf9\x5f\x52\xf4\x6a\x86\x6b\x2b\xab\xf5\xb6\xd7\x35\xec\xd2\xdf\x9f\x01\x00\x00\xff\xff\x49\x66\x66\x21\x1f\x0a\x00\x00"

func mockGoTplBytes() ([]byte, error) {
	return bindataRead(
		_mockGoTpl,
		"mock.go.tpl",
	)
}

func mockGoTpl() (*asset, error) {
	bytes, err := mockGoTplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "mock.go.tpl", size: 2591, mode: os.FileMode(420), modTime: time.Unix(1557020544, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"mgo.go.tpl": mgoGoTpl,
	"mock.go.tpl": mockGoTpl,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"mgo.go.tpl": &bintree{mgoGoTpl, map[string]*bintree{}},
	"mock.go.tpl": &bintree{mockGoTpl, map[string]*bintree{}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}

