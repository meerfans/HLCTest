// Copyright 2017-2018 The nox developers
package types

import (
	"Nox-DAG-test/script/tool/ecc"
	"time"
	"net"
	"Nox-DAG-test/script/tool/protocol"
	"encoding/binary"
	s "Nox-DAG-test/script/tool/serialization"
	"io"
)
// MaxNetAddressPayload returns the max payload size for the NetAddress
// based on the protocol version.
func MaxNetAddressPayload(pver uint32) uint32 {
	// Services 8 bytes + ip 16 bytes + port 2 bytes.
	plen := uint32(26)

	// Timestamp 4 bytes.
	plen += 4

	return plen
}

// writeNetAddress serializes a NetAddress to w depending on the protocol
// version and whether or not the timestamp is included per ts.  Some messages
// like version do not include the timestamp.
func WriteNetAddress(w io.Writer, pver uint32, na *NetAddress, ts bool) error {
	// TODO fix time ambiguous
	// NOTE: The protocol uses a uint32 for the timestamp so it will
	// stop working somewhere around 2106.  Also timestamp wasn't added until
	// until protocol version >= NetAddressTimeVersion.
	if ts {
		err := s.WriteElements(w, uint32(na.Timestamp.Unix()))
		if err != nil {
			return err
		}
	}

	// Ensure to always write 16 bytes even if the ip is nil.
	var ip [16]byte
	if na.IP != nil {
		copy(ip[:], na.IP.To16())
	}
	err := s.WriteElements(w, na.Services, ip)
	if err != nil {
		return err
	}

	// TODO unify endian
	// Sigh.  protocol mixes little and big endian.
	return binary.Write(w, binary.BigEndian, na.Port)
}
// readNetAddress reads an encoded NetAddress from r depending on the protocol
// version and whether or not the timestamp is included per ts.  Some messages
// like version do not include the timestamp.
func ReadNetAddress(r io.Reader, pver uint32, na *NetAddress, ts bool) error {
	var ip [16]byte

	// TODO fix time ambiguous
	// NOTE: The protocol uses a uint32 for the timestamp so it will
	// stop working somewhere around 2106.  Also timestamp wasn't added until
	// protocol version >= NetAddressTimeVersion
	if ts {
		err := s.ReadElements(r, (*s.Uint32Time)(&na.Timestamp))
		if err != nil {
			return err
		}
	}

	err := s.ReadElements(r, &na.Services, &ip)
	if err != nil {
		return err
	}

	// TODO unify endian
	// Sigh. protocol mixes little and big endian.
	port, err := s.BinarySerializer.Uint16(r, binary.BigEndian)
	if err != nil {
		return err
	}

	*na = NetAddress{
		Timestamp: na.Timestamp,
		Services:  na.Services,
		IP:        net.IP(ip[:]),
		Port:      port,
	}
	return nil
}
// NetAddress defines information about a peer on the network including the time
// it was last seen, the services it supports, its IP address, and port.
type NetAddress struct {
	// Last time the address was seen.  This is, unfortunately, encoded as a
	// uint32 on the wire and therefore is limited to 2106.  This field is
	// not present in the version message (MsgVersion) nor was it
	// added until protocol version >= NetAddressTimeVersion.
	Timestamp time.Time

	// Bitfield which identifies the services supported by the address.
	Services protocol.ServiceFlag

	// IP address of the peer.
	IP net.IP

	// Port the peer is using.  This is encoded in big endian on the wire
	// which differs from most everything else.
	Port uint16
}
type Address interface{
	// String returns the string encoding of the transaction output
	// destination.
	//
	// Please note that String differs subtly from EncodeAddress: String
	// will return the value as a string without any conversion, while
	// EncodeAddress may convert destination types (for example,
	// converting pubkeys to P2PKH addresses) before encoding as a
	// payment address string.
	String() 		string

	// with encode
	Encode()        string

	// Hash160 returns the Hash160(data) where data is the data normally
	// hashed to 160 bits from the respective address type.
	Hash160()       *[20]byte

	EcType()        ecc.EcType

	// raw byte in script, aka the hash in the most case
	ScriptAddress() []byte

	// TODO, revisit the design of address type decision
	// IsForNetwork returns whether or not the address is associated with the
	// passed network.
	// IsForNetwork(hashID [2]byte) bool
}

type AddressType byte

const (
	LegerAddress AddressType = 0x01
	ContractAddress AddressType = 0x02
)


