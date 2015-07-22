package ipfsld

import (
	"bytes"
	"io"

	codec "github.com/ugorji/go/codec"

	ld "github.com/ipfs/go-ipfsld"
)

type Encoder interface {
	Encode(doc *ld.Doc) error
}

type Decoder interface {
	Decode(doc *ld.Doc) error
}

type encoder struct {
	enc *codec.Encoder
}

func NewEncoder(w io.Writer) *encoder {
	h := &codec.CborHandle{}
	h.Canonical = true
	enc := codec.NewEncoder(w, h)
	return &encoder{enc}
}

func (c *encoder) Encode(doc *ld.Doc) error {
	return c.enc.Encode(&doc.Data)
}

type decoder struct {
	dec *codec.Decoder
}

func NewDecoder(r io.Reader) *decoder {
	h := &codec.CborHandle{}
	h.Canonical = true
	dec := codec.NewDecoder(r, h)
	return &decoder{dec}
}

func (c *decoder) Decode(doc *ld.Doc) error {
	return c.dec.Decode(&doc.Data)
}

// Marshal serializes an ipfs-ld document to a []byte.
func Marshal(doc *ld.Doc) ([]byte, error) {
	var buf bytes.Buffer
	err := MarshalTo(&buf, doc)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// MarshalTo serializes an ipfs-ld document to a writer.
func MarshalTo(w io.Writer, doc *ld.Doc) error {
	return NewEncoder(w).Encode(doc)
}

// Unmarshal deserializes an ipfs-ld document to a []byte.
func Unmarshal(buf []byte) (*ld.Doc, error) {
	doc := new(ld.Doc)
	err := UnmarshalFrom(bytes.NewBuffer(buf), doc)
	if err != nil {
		return nil, err
	}

	// have to call NewDoc so the initial parsing (schema) takes place.
	return ld.NewDoc(doc.Data), nil
}

// UnmarshalFrom deserializes an ipfs-ld document from a reader.
func UnmarshalFrom(r io.Reader, doc *ld.Doc) error {
	return NewDecoder(r).Decode(doc)
}
