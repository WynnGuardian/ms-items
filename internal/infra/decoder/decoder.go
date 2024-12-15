package decoder

import (
	"errors"
	"fmt"
	"math"
)

type IntReader struct {
	Data  []int
	Index int
}

func (br *IntReader) peek() int {
	if br.Index >= len(br.Data) {
		return -1
	}
	return br.Data[br.Index]
}

func (br *IntReader) read() int {
	if br.Index >= len(br.Data) {
		return -1
	}
	b := br.Data[br.Index]
	br.Index++
	return b
}

type DecodedItem struct {
	Name            string
	Version         int
	Identifications map[int]int
}

type ItemDecoder struct {
	reader *IntReader
}

func NewItemDecoder(utf16 string) *ItemDecoder {
	return &ItemDecoder{
		reader: &IntReader{Data: FromUtf16String(utf16).Bytes, Index: 0},
	}
}

func (d *ItemDecoder) Decode() (*DecodedItem, error) {

	decoded := &DecodedItem{}

	if d.reader.read() != 0 {
		return nil, errors.New("invalid item")
	}

	version, err := d.versionBlock()
	if err != nil {
		return nil, err
	}
	decoded.Version = version

	if err := d.itemTypeBlock(); err != nil {
		return nil, err
	}

	name, err := d.nameBlock()
	if err != nil {
		return nil, err
	}

	decoded.Name = name

	ids, err := d.identificationBlock()
	if err != nil {
		return nil, err
	}
	decoded.Identifications = ids

	return decoded, nil
}

func (d *ItemDecoder) versionBlock() (int, error) {
	version := d.reader.read()
	if version == -1 || version != 0 {
		return 0, fmt.Errorf("unsuported encoding version: %d", version)
	}
	return version, nil
}

func (d *ItemDecoder) itemTypeBlock() error {
	_ = d.reader.read() // 0 header
	t := d.reader.read()
	if t == -1 {
		return errors.New("invalid item")
	}
	if t != 0 {
		return fmt.Errorf("invalid item type: %d", t)
	}
	return nil
}

func (d *ItemDecoder) nameBlock() (string, error) {
	_ = d.reader.read() // 0 header
	bytes := make([]byte, 0)
	for d.reader.peek() != 0 {
		bytes = append(bytes, byte(d.reader.read()))
	}
	return string(bytes), nil
}

func (d *ItemDecoder) identificationBlock() (map[int]int, error) {
	ids := make(map[int]int, 0)
	_, _ = d.reader.read(), d.reader.read() // read the 0 start block id and header
	length := d.reader.read()
	if length >= 255 {
		return nil, fmt.Errorf("yoo many identifications: %d", length)
	}
	extended := d.reader.read() == 1
	preIdedCount := int(0)
	if extended {
		preIdedCount = d.reader.read()
	}
	for i := 0; i < int(preIdedCount)+int(length); i++ {
		id := d.reader.read()
		if i < int(preIdedCount) {
			_ = d.reader.DecodeFirstVSI()
			continue
		}
		base := d.reader.DecodeFirstVSI()
		roll := d.reader.read()
		ids[id] = int(math.Floor((float64(base)*(float64(roll)/100) + 0.5)))
		if (id >= 0 && id <= 3) || (id >= 37 && id <= 40) { // reverse costs
			ids[id] = -ids[id]
		}
	}
	return ids, nil
}

func (r *IntReader) DecodeFirstVSI() int {
	val, numBytes := int(0), 0

	for r.peek() != -1 && (r.peek()&0x80) != 0 {
		val |= int(r.read() & 0x7F << (7 * numBytes))
		numBytes++
	}
	val |= int(r.read() & 0x7F << (7 * numBytes))

	return int(uint(val)>>1) ^ -(val & 1)
}
