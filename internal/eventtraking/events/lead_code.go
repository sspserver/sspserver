//
// @project GeniusRabbit rotator 2017 - 2019
// @author Dmitry Ponomarev <demdxx@gmail.com> 2017 - 2019
//

package events

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"hash/crc32"

	"geniusrabbit.dev/sspserver/internal/msgpack/thrift"
)

const (
	signConst = "Gei4xish1foojahco"
	signSize  = 4
	magic1    = 1323
	magic2    = 0x01f2f4ac
)

var (
	errInvalidCode    = errors.New("lead code: invalid code")
	errInvalidSign    = errors.New("lead code: invalid sign")
	errInvalidMessage = errors.New("lead code: invalid message")
)

// LeadCode object
type LeadCode struct {
	AuctionID         string `thrift:",1" json:"id"`  // Internal Auction ID
	ImpAdID           string `thrift:",2" json:"iid"` // Specific ID for paticular AD impression
	SourceID          uint64 `thrift:",3" json:"sid,omitempty"`
	ProjectID         uint64 `thrift:",4" json:"pid,omitempty"`
	PublisherCompany  uint64 `thrift:",5" json:"pub,omitempty"`
	AdvertiserCompany uint64 `thrift:",6" json:"adv,omitempty"`
	CampaignID        uint64 `thrift:",7" json:"cid,omitempty"`
	AdID              uint64 `thrift:",8" json:"aid,omitempty"`
	Price             int64  `thrift:",9" json:"p,omitempty"`
	Timestamp         int64  `thrift:",10" json:"ts"` // In seconds
}

// Pack structure to bytes
func (lc *LeadCode) Pack() Code {
	var (
		buff bytes.Buffer
		err  error
	)

	// Sign the record
	if _, err = buff.Write(lc.Sign()); err == nil {
		// Encode message
		if err = thrift.NewEncoder(&buff).Encode(lc); err != nil {
			return CodeObj(nil, err)
		}
	}

	return CodeObj(buff.Bytes(), err)
}

// Unpack data
func (lc *LeadCode) Unpack(code []byte) (err error) {
	if len(code) < signSize+3 {
		return errInvalidMessage
	}

	var (
		sign = code[:signSize]
		dec  = thrift.NewDecoder(nil, code[signSize:])
	)

	if err = dec.Decode(&lc); err != nil {
		return err
	}

	// If sign isn't correct
	if bytes.Compare(sign, lc.Sign()) != 0 {
		return errInvalidSign
	}

	return
}

// String implementation of fmt.Stringer interface
func (lc LeadCode) String() string {
	return string(lc.Sign())
}

// Sign of object (length 10bytes)
func (lc LeadCode) Sign() (sign []byte) {
	// h := sha1.New()
	// fmt.Fprintf(h, "%s%s%d%d%d%d%d%d%d"+signConst,
	// 	lc.AuctionID, lc.ImpAdID, lc.ProjectID, lc.SourceID,
	// 	lc.CampaignID, lc.AdID, lc.Timestamp, (lc.AdID+uint64(lc.Timestamp))/magic1,
	// 	lc.Timestamp-magic2)
	// return h.Sum(nil)

	var buff bytes.Buffer
	fmt.Fprintf(&buff, "%s%s%d%d%d%d%d%d%d"+signConst,
		lc.AuctionID, lc.ImpAdID, lc.ProjectID, lc.SourceID,
		lc.CampaignID, lc.AdID, lc.Timestamp, (lc.AdID+uint64(lc.Timestamp))/magic1,
		lc.Timestamp-magic2)

	sign = make([]byte, 4)
	binary.BigEndian.PutUint32(
		sign,
		crc32.Checksum(buff.Bytes(), crc32.IEEETable),
	)
	return sign
}
