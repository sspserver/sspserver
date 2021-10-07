package optimizer

import (
	"bytes"
	"encoding/json"
	"sync"
)

type optimizerItem struct {
	mx       sync.RWMutex
	FormatID uint
	Data     map[keyType]*atomicFloatItem
}

func newOptimizerItem(formatID uint) *optimizerItem {
	return &optimizerItem{
		FormatID: formatID,
		Data:     map[keyType]*atomicFloatItem{},
	}
}

func (optItem *optimizerItem) getOrCreate(countryID byte, languageID, deviceID, osID, browserID uint) *atomicFloatItem {
	key := keyType{
		CountryID:  countryID,
		LanguageID: languageID,
		DeviceID:   deviceID,
		OsID:       osID,
		BrowserID:  browserID,
	}

	optItem.mx.RLock()
	if item := optItem.Data[key]; item != nil {
		optItem.mx.RUnlock()
		return item
	}
	optItem.mx.RUnlock()

	optItem.mx.Lock()
	item := new(atomicFloatItem)
	item.Set(initialState)
	optItem.Data[key] = item
	optItem.mx.Unlock()

	return item
}

// MarshalJSON implements json.Marshaler interface
func (optItem *optimizerItem) MarshalJSON() ([]byte, error) {
	optItem.mx.RLock()
	defer optItem.mx.RUnlock()

	var (
		idx    = 0
		buff   bytes.Buffer
		writer = json.NewEncoder(&buff)
	)

	buff.WriteByte('[')

	for key, vl := range optItem.Data {
		val := vl.Get()
		if val <= 0 {
			continue
		}

		keyv := keyValType{
			CountryID:  key.CountryID,
			LanguageID: key.LanguageID,
			DeviceID:   key.DeviceID,
			OsID:       key.OsID,
			BrowserID:  key.BrowserID,
			Value:      val,
		}

		if idx > 0 {
			buff.WriteByte(',')
		}

		if err := writer.Encode(keyv); err != nil {
			return nil, err
		}

		idx++
	}

	buff.WriteByte(']')
	return buff.Bytes(), nil
}

// UnmarshalJSON implements json.Unmarshaler interface
func (optItem *optimizerItem) UnmarshalJSON(data []byte) error {
	var (
		list []keyValType
		body = map[keyType]*atomicFloatItem{}
	)

	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}

	for _, item := range list {
		if item.Value == 0 {
			continue
		}

		key := keyType{
			CountryID:  item.CountryID,
			LanguageID: item.LanguageID,
			DeviceID:   item.DeviceID,
			OsID:       item.OsID,
			BrowserID:  item.BrowserID,
		}

		itemVal := new(atomicFloatItem)
		itemVal.Set(item.Value)
		body[key] = itemVal
	}

	optItem.mx.Lock()
	defer optItem.mx.Unlock()
	optItem.Data = body

	return nil
}
