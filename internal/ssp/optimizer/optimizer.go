package optimizer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
)

type keyType struct {
	CountryID  byte
	LanguageID uint
	DeviceID   uint
	OsID       uint
	BrowserID  uint
}

// Optimizer of the platform requests
type Optimizer struct {
	mx   sync.RWMutex
	data []*optimizerItem
}

// Init Optimizer object
func (opt *Optimizer) Init() {
	if opt.data == nil {
		opt.data = make([]*optimizerItem, 0, 10)
	}
}

// Inc incrementation of the filter effectivity
func (opt *Optimizer) Inc(formatID uint, countryID byte, languageID, deviceID, osID, browserID uint, result float64) float64 {
	return opt.getOrCreate(formatID, countryID, languageID, deviceID, osID, browserID).Inc(result)
}

// Test accoring to current state
func (opt *Optimizer) Test(formatID uint, countryID byte, languageID, deviceID, osID, browserID uint, min float64) bool {
	return opt.getOrCreate(formatID, countryID, languageID, deviceID, osID, browserID).Test(min)
}

// Len returns count of elements
func (opt *Optimizer) Len() (cnt int) {
	for _, it := range opt.data {
		cnt += len(it.Data)
	}
	return cnt
}

// Enum all counters
func (opt *Optimizer) Enum(fnk func(formatID uint, countryID byte, languageID, deviceID, osID, browserID uint, min float64)) {
	opt.mx.RLock()
	defer opt.mx.RUnlock()

	for _, item := range opt.data {
		item.mx.RLock()
		for key, it := range item.Data {
			fnk(item.FormatID, key.CountryID, key.LanguageID, key.DeviceID, key.OsID, key.BrowserID, it.Get())
		}
		item.mx.RUnlock()
	}
}

// MarshalJSON implements json.Marshaler interface
func (opt *Optimizer) MarshalJSON() ([]byte, error) {
	var buff bytes.Buffer

	buff.WriteByte('{')
	for i, item := range opt.data {
		if i > 0 {
			buff.WriteByte(',')
		}
		buff.WriteByte('"')
		buff.WriteString(strconv.FormatUint(uint64(item.FormatID), 10))
		buff.WriteByte('"')
		buff.WriteByte(':')

		if err := json.NewEncoder(&buff).Encode(item); err != nil {
			return nil, err
		}
	}
	buff.WriteByte('}')

	return buff.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (opt *Optimizer) UnmarshalJSON(data []byte) error {
	var (
		base = map[string]json.RawMessage{}
		body = []*optimizerItem{}
	)

	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}

	for key, value := range base {
		val, err := strconv.ParseUint(key, 10, 64)
		if err != nil {
			return fmt.Errorf("parse format ID: %s", err.Error())
		}

		item := newOptimizerItem(uint(val))
		if err = json.Unmarshal(value, item); err != nil {
			return fmt.Errorf("item ID %s decode error: %s", key, err.Error())
		}
		body = append(body, item)
	}

	opt.mx.Lock()
	defer opt.mx.Unlock()
	opt.data = body

	return nil
}

// Encode optimizer data
func (opt *Optimizer) Encode() ([]byte, error) {
	return opt.MarshalJSON()
}

// Decode optimizer data
func (opt *Optimizer) Decode(data []byte) error {
	return opt.UnmarshalJSON(data)
}

// getOrCreate cofficient counter
func (opt *Optimizer) getOrCreate(formatID uint, countryID byte, languageID, deviceID, osID, browserID uint) *atomicFloatItem {
	if item := opt.getItem(formatID); item != nil {
		return item.getOrCreate(countryID, languageID, deviceID, osID, browserID)
	}

	opt.mx.Lock()
	item := newOptimizerItem(formatID)
	opt.data = append(opt.data, item)
	opt.mx.Unlock()

	return item.getOrCreate(countryID, languageID, deviceID, osID, browserID)
}

func (opt *Optimizer) getItem(formatID uint) *optimizerItem {
	opt.mx.RLock()
	defer opt.mx.RUnlock()

	for _, item := range opt.data {
		if item.FormatID == formatID {
			return item
		}
	}
	return nil
}
