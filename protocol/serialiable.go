package protocol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
)

var NullRead = bytes.NewReader([]byte{})

var (
	ErrPasswordNotMatch   = errors.New("password isn't match?")
	ErrReadingUnavailable = errors.New("reading unavailable.")
	ErrIgnoreSensor       = errors.New("ignore sensor - entity is not present or disabled.")
	ErrInsufficientBytes  = errors.New("Insufficient Bytes.")
)

//var ErrShortPacket = errors.New("Short Packet.")

type Writer struct {
	buf    *bytes.Buffer
	cached [16]byte
	err    error
}

func (w *Writer) Init(bs []byte) *Writer {
	w.buf = bytes.NewBuffer(bs)
	w.err = nil
	return w
}

func (w *Writer) Reset() *Writer {
	w.buf.Reset()
	w.err = nil
	return w
}

func (w *Writer) ToWriter() io.Writer {
	if nil != w.err {
		return ioutil.Discard
	}
	return w.buf
}

func (r *Writer) Truncate(n int) {
	r.buf.Truncate(n)
}

func (r *Writer) Len() int {
	return r.buf.Len()
}

func (w *Writer) Bytes() []byte {
	return w.buf.Bytes()
}

func (w *Writer) SetError(e error) {
	if nil != w.err {
		return
	}
	w.err = e
}

func (w *Writer) Err() error {
	return w.err
}

func (w *Writer) WriteBytes(bs []byte) {
	if nil != w.err {
		return
	}
	_, w.err = w.buf.Write(bs)
}

func (w *Writer) WriteUint8(v uint8) {
	if nil != w.err {
		return
	}

	w.err = w.buf.WriteByte(byte(v))
}

func (w *Writer) WriteUint16(v uint16) {
	if nil != w.err {
		return
	}
	binary.LittleEndian.PutUint16(w.cached[:], v)
	_, w.err = w.buf.Write(w.cached[:2])
}

func (w *Writer) WriteUint32(v uint32) {
	if nil != w.err {
		return
	}

	binary.LittleEndian.PutUint32(w.cached[:], v)
	_, w.err = w.buf.Write(w.cached[:4])
}

func (w *Writer) WriteUint64(v uint64) {
	if nil != w.err {
		return
	}
	binary.LittleEndian.PutUint64(w.cached[:], v)
	_, w.err = w.buf.Write(w.cached[:8])
}

func (w *Writer) WriteInt8(v int8) {
	if nil != w.err {
		return
	}

	w.err = w.buf.WriteByte(byte(v))
}

func (w *Writer) WriteInt16(v int16) {
	if nil != w.err {
		return
	}
	binary.LittleEndian.PutUint16(w.cached[:], uint16(v))
	_, w.err = w.buf.Write(w.cached[:2])
}

func (w *Writer) WriteInt32(v int32) {
	if nil != w.err {
		return
	}

	binary.LittleEndian.PutUint32(w.cached[:], uint32(v))
	_, w.err = w.buf.Write(w.cached[:4])
}

func (w *Writer) WriteInt64(v int64) {
	if nil != w.err {
		return
	}
	binary.LittleEndian.PutUint64(w.cached[:], uint64(v))
	_, w.err = w.buf.Write(w.cached[:8])
}

func (w *Writer) WriteUint16WithOrder(v uint16, order binary.ByteOrder) {
	if nil != w.err {
		return
	}
	order.PutUint16(w.cached[:], v)
	_, w.err = w.buf.Write(w.cached[:2])
}

func (w *Writer) WriteUint32WithOrder(v uint32, order binary.ByteOrder) {
	if nil != w.err {
		return
	}

	order.PutUint32(w.cached[:], v)
	_, w.err = w.buf.Write(w.cached[:4])
}

func (w *Writer) WriteUint64WithOrder(v uint64, order binary.ByteOrder) {
	if nil != w.err {
		return
	}
	order.PutUint64(w.cached[:], v)
	_, w.err = w.buf.Write(w.cached[:8])
}

func NewWriter(bs []byte) *Writer {
	return &Writer{buf: bytes.NewBuffer(bs)}
}

type Reader struct {
	buf *bytes.Buffer
	err error
}

func (r *Reader) Init(bs []byte) *Reader {
	r.buf = bytes.NewBuffer(bs)
	r.err = nil
	return r
}

func (r *Reader) ToReader() io.Reader {
	if nil != r.err {
		return NullRead
	}
	return r.buf
}

func (r *Reader) Fork(l int) *Reader {
	if nil != r.err {
		return r
	}
	if r.buf.Len() < l {
		return &Reader{buf: nil, err: ErrInsufficientBytes}
	}
	bs := r.buf.Bytes()
	return &Reader{buf: bytes.NewBuffer(bs[:l])}
}

func (r *Reader) Len() int {
	return r.buf.Len()
}

func (r *Reader) Bytes() []byte {
	return r.buf.Bytes()
}

func (r *Reader) SetError(e error) {
	if nil != r.err {
		return
	}
	if io.EOF == e {
		r.err = ErrInsufficientBytes
	} else {
		r.err = e
	}
}

func (r *Reader) Err() error {
	return r.err
}

func (r *Reader) ReadCopy(n int) []byte {
	if nil != r.err {
		return nil
	}
	if 0 == n {
		return nil
	}

	bs := r.buf.Next(n)
	if nil == bs {
		return bs
	}

	copyed := make([]byte, len(bs))
	copy(copyed, bs)
	return copyed
}

func (r *Reader) ReadByte() byte {
	if nil != r.err {
		return 0
	}
	b, err := r.buf.ReadByte()
	if nil != err {
		r.SetError(err)
		return 0
	}
	return b
}

func (r *Reader) ReadBytes(n int) []byte {
	if nil != r.err {
		return nil
	}
	return r.buf.Next(n)
}

func (r *Reader) ReadUint8() uint8 {
	if nil != r.err {
		return 0
	}
	if r.buf.Len() >= 1 {
		b, _ := r.buf.ReadByte()
		return uint8(b)
	}
	r.err = ErrInsufficientBytes
	return 0
}

func (r *Reader) ReadUint16() uint16 {
	if nil != r.err {
		return 0
	}
	if r.buf.Len() >= 2 {
		return binary.LittleEndian.Uint16(r.buf.Next(2))
	}
	r.err = ErrInsufficientBytes
	return 0
}

func (r *Reader) ReadUint32() uint32 {
	if nil != r.err {
		return 0
	}
	if r.buf.Len() >= 4 {
		return binary.LittleEndian.Uint32(r.buf.Next(4))
	}
	r.err = ErrInsufficientBytes
	return 0
}

func (r *Reader) ReadUint64() uint64 {
	if nil != r.err {
		return 0
	}
	if r.buf.Len() >= 8 {
		return binary.LittleEndian.Uint64(r.buf.Next(8))
	}
	r.err = ErrInsufficientBytes
	return 0
}

func (r *Reader) ReadInt8() int8 {
	if nil != r.err {
		return 0
	}
	if r.buf.Len() >= 1 {
		b, _ := r.buf.ReadByte()
		return int8(b)
	}
	r.err = ErrInsufficientBytes
	return 0
}

func (r *Reader) ReadInt16() int16 {
	if nil != r.err {
		return 0
	}
	if r.buf.Len() >= 2 {
		return int16(binary.LittleEndian.Uint16(r.buf.Next(2)))
	}
	r.err = ErrInsufficientBytes
	return 0
}

func (r *Reader) ReadInt32() int32 {
	if nil != r.err {
		return 0
	}
	if r.buf.Len() >= 4 {
		return int32(binary.LittleEndian.Uint32(r.buf.Next(4)))
	}
	r.err = ErrInsufficientBytes
	return 0
}

func (r *Reader) ReadInt64() int64 {
	if nil != r.err {
		return 0
	}
	if r.buf.Len() >= 8 {
		return int64(binary.LittleEndian.Uint64(r.buf.Next(8)))
	}
	r.err = ErrInsufficientBytes
	return 0
}

func (r *Reader) Read(o interface{}) {
	if rb, ok := o.(Readable); ok {
		rb.ReadBytes(r)
	} else {
		err := binary.Read(r.ToReader(),
			binary.LittleEndian, o)
		r.SetError(err)
	}
}

func NewReader(bs []byte) *Reader {
	return &Reader{buf: bytes.NewBuffer(bs)}
}

type Writable interface {
	WriteBytes(w *Writer)
}

type Readable interface {
	ReadBytes(r *Reader)
}

type Serialiable interface {
	Readable
	Writable
}

func ToBytes(o Writable) ([]byte, error) {
	w := NewWriter(make([]byte, 0, 256))
	o.WriteBytes(w)
	return w.Bytes(), w.Err()
}

func FromBytes(o Readable, bs []byte) error {
	r := NewReader(bs)
	o.ReadBytes(r)
	return r.Err()
}
