package youtube

import (
	"bytes"
	"encoding/base64"
	"net/url"
)

type ProtoBuilder struct {
	byteBuffer bytes.Buffer
}

func (pb *ProtoBuilder) ToBytes() []byte {
	return pb.byteBuffer.Bytes()
}

func (pb *ProtoBuilder) ToUrlEncodedBase64() string {
	b64 := base64.URLEncoding.EncodeToString(pb.ToBytes())
	return url.QueryEscape(b64)
}

func (pb *ProtoBuilder) writeVarint(val int64) error {
	if val == 0 {
		_, err := pb.byteBuffer.Write([]byte{0})
		return err
	}
	for {
		b := byte(val & 0x7F)
		val >>= 7
		if val != 0 {
			b |= 0x80
		}
		_, err := pb.byteBuffer.Write([]byte{b})
		if err != nil {
			return err
		}
		if val == 0 {
			break
		}
	}
	return nil
}

func (pb *ProtoBuilder) field(field int, wireType byte) error {
	val := int64(field<<3) | int64(wireType&0x07)
	return pb.writeVarint(val)
}

func (pb *ProtoBuilder) Varint(field int, val int64) error {
	err := pb.field(field, 0)
	if err != nil {
		return err
	}
	return pb.writeVarint(val)
}

func (pb *ProtoBuilder) String(field int, stringVal string) error {
	strBts := []byte(stringVal)
	return pb.Bytes(field, strBts)
}

func (pb *ProtoBuilder) Bytes(field int, bytesVal []byte) error {
	if err := pb.field(field, 2); err != nil {
		return err
	}

	if err := pb.writeVarint(int64(len(bytesVal))); err != nil {
		return err
	}

	_, err := pb.byteBuffer.Write(bytesVal)
	return err
}
