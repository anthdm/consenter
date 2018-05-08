package codec

import (
	"encoding/binary"
	"io"

	pb "github.com/anthdm/consenter/pkg/protos"
	"github.com/golang/protobuf/proto"
)

// DecodeProto decodes msg to r.
func DecodeProto(r io.Reader, msg *pb.Message) error {
	var n uint32
	if err := binary.Read(r, binary.LittleEndian, &n); err != nil {
		return err
	}
	buf := make([]byte, n)
	_, err := io.ReadFull(r, buf)
	if err != nil {
		return err
	}
	return proto.Unmarshal(buf, msg)
}

// EncodeProto encodes msg to w.
func EncodeProto(w io.Writer, msg *pb.Message) error {
	b, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, uint32(len(b)))
	_, err = w.Write(buf)
	_, err = w.Write(b)
	return err
}
