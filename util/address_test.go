// Copyright (c) 2013-2017 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package util_test

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/daglabs/btcd/util"
	"golang.org/x/crypto/ripemd160"
)

func TestAddresses(t *testing.T) {
	tests := []struct {
		name    string
		addr    string
		encoded string
		valid   bool
		result  util.Address
		f       func() (util.Address, error)
		prefix  util.Bech32Prefix
	}{
		// Positive P2PKH tests.
		{
			name:    "mainnet p2pkh",
			addr:    "dagcoin:qr35ennsep3hxfe7lnz5ee7j5jgmkjswss74as46gy",
			encoded: "dagcoin:qr35ennsep3hxfe7lnz5ee7j5jgmkjswss74as46gy",
			valid:   true,
			result: util.TstAddressPubKeyHash(
				util.Bech32PrefixDAGCoin,
				[ripemd160.Size]byte{
					0xe3, 0x4c, 0xce, 0x70, 0xc8, 0x63, 0x73, 0x27, 0x3e, 0xfc,
					0xc5, 0x4c, 0xe7, 0xd2, 0xa4, 0x91, 0xbb, 0x4a, 0x0e, 0x84}),
			f: func() (util.Address, error) {
				pkHash := []byte{
					0xe3, 0x4c, 0xce, 0x70, 0xc8, 0x63, 0x73, 0x27, 0x3e, 0xfc,
					0xc5, 0x4c, 0xe7, 0xd2, 0xa4, 0x91, 0xbb, 0x4a, 0x0e, 0x84}
				return util.NewAddressPubKeyHash(pkHash, util.Bech32PrefixDAGCoin)
			},
			prefix: util.Bech32PrefixDAGCoin,
		},
		{
			name:    "mainnet p2pkh 2",
			addr:    "dagcoin:qq80qvqs0lfxuzmt7sz3909ze6camq9d4gwzqeljga",
			encoded: "dagcoin:qq80qvqs0lfxuzmt7sz3909ze6camq9d4gwzqeljga",
			valid:   true,
			result: util.TstAddressPubKeyHash(
				util.Bech32PrefixDAGCoin,
				[ripemd160.Size]byte{
					0x0e, 0xf0, 0x30, 0x10, 0x7f, 0xd2, 0x6e, 0x0b, 0x6b, 0xf4,
					0x05, 0x12, 0xbc, 0xa2, 0xce, 0xb1, 0xdd, 0x80, 0xad, 0xaa}),
			f: func() (util.Address, error) {
				pkHash := []byte{
					0x0e, 0xf0, 0x30, 0x10, 0x7f, 0xd2, 0x6e, 0x0b, 0x6b, 0xf4,
					0x05, 0x12, 0xbc, 0xa2, 0xce, 0xb1, 0xdd, 0x80, 0xad, 0xaa}
				return util.NewAddressPubKeyHash(pkHash, util.Bech32PrefixDAGCoin)
			},
			prefix: util.Bech32PrefixDAGCoin,
		},
		{
			name:    "testnet p2pkh",
			addr:    "dagtest:qputx94qseratdmjs0j395mq8u03er0x3ucluj5qam",
			encoded: "dagtest:qputx94qseratdmjs0j395mq8u03er0x3ucluj5qam",
			valid:   true,
			result: util.TstAddressPubKeyHash(
				util.Bech32PrefixDAGTest,
				[ripemd160.Size]byte{
					0x78, 0xb3, 0x16, 0xa0, 0x86, 0x47, 0xd5, 0xb7, 0x72, 0x83,
					0xe5, 0x12, 0xd3, 0x60, 0x3f, 0x1f, 0x1c, 0x8d, 0xe6, 0x8f}),
			f: func() (util.Address, error) {
				pkHash := []byte{
					0x78, 0xb3, 0x16, 0xa0, 0x86, 0x47, 0xd5, 0xb7, 0x72, 0x83,
					0xe5, 0x12, 0xd3, 0x60, 0x3f, 0x1f, 0x1c, 0x8d, 0xe6, 0x8f}
				return util.NewAddressPubKeyHash(pkHash, util.Bech32PrefixDAGTest)
			},
			prefix: util.Bech32PrefixDAGTest,
		},

		// Negative P2PKH tests.
		{
			name:  "p2pkh wrong hash length",
			addr:  "",
			valid: false,
			f: func() (util.Address, error) {
				pkHash := []byte{
					0x00, 0x0e, 0xf0, 0x30, 0x10, 0x7f, 0xd2, 0x6e, 0x0b, 0x6b,
					0xf4, 0x05, 0x12, 0xbc, 0xa2, 0xce, 0xb1, 0xdd, 0x80, 0xad,
					0xaa}
				return util.NewAddressPubKeyHash(pkHash, util.Bech32PrefixDAGCoin)
			},
			prefix: util.Bech32PrefixDAGCoin,
		},
		{
			name:   "p2pkh bad checksum",
			addr:   "dagcoin:qr35ennsep3hxfe7lnz5ee7j5jgmkjswss74as46gx",
			valid:  false,
			prefix: util.Bech32PrefixDAGCoin,
		},

		// Positive P2SH tests.
		{
			// Taken from transactions:
			// output: 3c9018e8d5615c306d72397f8f5eef44308c98fb576a88e030c25456b4f3a7ac
			// input:  837dea37ddc8b1e3ce646f1a656e79bbd8cc7f558ac56a169626d649ebe2a3ba
			name:    "mainnet p2sh",
			addr:    "dagcoin:pruptvpkmxamee0f72sq40gm70wfr624zq8mc2ujcn",
			encoded: "dagcoin:pruptvpkmxamee0f72sq40gm70wfr624zq8mc2ujcn",
			valid:   true,
			result: util.TstAddressScriptHash(
				util.Bech32PrefixDAGCoin,
				[ripemd160.Size]byte{
					0xf8, 0x15, 0xb0, 0x36, 0xd9, 0xbb, 0xbc, 0xe5, 0xe9, 0xf2,
					0xa0, 0x0a, 0xbd, 0x1b, 0xf3, 0xdc, 0x91, 0xe9, 0x55, 0x10}),
			f: func() (util.Address, error) {
				script := []byte{
					0x52, 0x41, 0x04, 0x91, 0xbb, 0xa2, 0x51, 0x09, 0x12, 0xa5,
					0xbd, 0x37, 0xda, 0x1f, 0xb5, 0xb1, 0x67, 0x30, 0x10, 0xe4,
					0x3d, 0x2c, 0x6d, 0x81, 0x2c, 0x51, 0x4e, 0x91, 0xbf, 0xa9,
					0xf2, 0xeb, 0x12, 0x9e, 0x1c, 0x18, 0x33, 0x29, 0xdb, 0x55,
					0xbd, 0x86, 0x8e, 0x20, 0x9a, 0xac, 0x2f, 0xbc, 0x02, 0xcb,
					0x33, 0xd9, 0x8f, 0xe7, 0x4b, 0xf2, 0x3f, 0x0c, 0x23, 0x5d,
					0x61, 0x26, 0xb1, 0xd8, 0x33, 0x4f, 0x86, 0x41, 0x04, 0x86,
					0x5c, 0x40, 0x29, 0x3a, 0x68, 0x0c, 0xb9, 0xc0, 0x20, 0xe7,
					0xb1, 0xe1, 0x06, 0xd8, 0xc1, 0x91, 0x6d, 0x3c, 0xef, 0x99,
					0xaa, 0x43, 0x1a, 0x56, 0xd2, 0x53, 0xe6, 0x92, 0x56, 0xda,
					0xc0, 0x9e, 0xf1, 0x22, 0xb1, 0xa9, 0x86, 0x81, 0x8a, 0x7c,
					0xb6, 0x24, 0x53, 0x2f, 0x06, 0x2c, 0x1d, 0x1f, 0x87, 0x22,
					0x08, 0x48, 0x61, 0xc5, 0xc3, 0x29, 0x1c, 0xcf, 0xfe, 0xf4,
					0xec, 0x68, 0x74, 0x41, 0x04, 0x8d, 0x24, 0x55, 0xd2, 0x40,
					0x3e, 0x08, 0x70, 0x8f, 0xc1, 0xf5, 0x56, 0x00, 0x2f, 0x1b,
					0x6c, 0xd8, 0x3f, 0x99, 0x2d, 0x08, 0x50, 0x97, 0xf9, 0x97,
					0x4a, 0xb0, 0x8a, 0x28, 0x83, 0x8f, 0x07, 0x89, 0x6f, 0xba,
					0xb0, 0x8f, 0x39, 0x49, 0x5e, 0x15, 0xfa, 0x6f, 0xad, 0x6e,
					0xdb, 0xfb, 0x1e, 0x75, 0x4e, 0x35, 0xfa, 0x1c, 0x78, 0x44,
					0xc4, 0x1f, 0x32, 0x2a, 0x18, 0x63, 0xd4, 0x62, 0x13, 0x53,
					0xae}
				return util.NewAddressScriptHash(script, util.Bech32PrefixDAGCoin)
			},
			prefix: util.Bech32PrefixDAGCoin,
		},
		{
			// Taken from transactions:
			// output: b0539a45de13b3e0403909b8bd1a555b8cbe45fd4e3f3fda76f3a5f52835c29d
			// input: (not yet redeemed at time test was written)
			name:    "mainnet p2sh 2",
			addr:    "dagcoin:pr5vxqxg0xrwl2zvxlq9rxffqx00sm44ksj47shjr6",
			encoded: "dagcoin:pr5vxqxg0xrwl2zvxlq9rxffqx00sm44ksj47shjr6",
			valid:   true,
			result: util.TstAddressScriptHash(
				util.Bech32PrefixDAGCoin,
				[ripemd160.Size]byte{
					0xe8, 0xc3, 0x00, 0xc8, 0x79, 0x86, 0xef, 0xa8, 0x4c, 0x37,
					0xc0, 0x51, 0x99, 0x29, 0x01, 0x9e, 0xf8, 0x6e, 0xb5, 0xb4}),
			f: func() (util.Address, error) {
				hash := []byte{
					0xe8, 0xc3, 0x00, 0xc8, 0x79, 0x86, 0xef, 0xa8, 0x4c, 0x37,
					0xc0, 0x51, 0x99, 0x29, 0x01, 0x9e, 0xf8, 0x6e, 0xb5, 0xb4}
				return util.NewAddressScriptHashFromHash(hash, util.Bech32PrefixDAGCoin)
			},
			prefix: util.Bech32PrefixDAGCoin,
		},
		{
			// Taken from bitcoind base58_keys_valid.
			name:    "testnet p2sh",
			addr:    "dagtest:przhjdpv93xfygpqtckdc2zkzuzqeyj2pg6ghunlhx",
			encoded: "dagtest:przhjdpv93xfygpqtckdc2zkzuzqeyj2pg6ghunlhx",
			valid:   true,
			result: util.TstAddressScriptHash(
				util.Bech32PrefixDAGTest,
				[ripemd160.Size]byte{
					0xc5, 0x79, 0x34, 0x2c, 0x2c, 0x4c, 0x92, 0x20, 0x20, 0x5e,
					0x2c, 0xdc, 0x28, 0x56, 0x17, 0x04, 0x0c, 0x92, 0x4a, 0x0a}),
			f: func() (util.Address, error) {
				hash := []byte{
					0xc5, 0x79, 0x34, 0x2c, 0x2c, 0x4c, 0x92, 0x20, 0x20, 0x5e,
					0x2c, 0xdc, 0x28, 0x56, 0x17, 0x04, 0x0c, 0x92, 0x4a, 0x0a}
				return util.NewAddressScriptHashFromHash(hash, util.Bech32PrefixDAGTest)
			},
			prefix: util.Bech32PrefixDAGTest,
		},

		// Negative P2SH tests.
		{
			name:  "p2sh wrong hash length",
			addr:  "",
			valid: false,
			f: func() (util.Address, error) {
				hash := []byte{
					0x00, 0xf8, 0x15, 0xb0, 0x36, 0xd9, 0xbb, 0xbc, 0xe5, 0xe9,
					0xf2, 0xa0, 0x0a, 0xbd, 0x1b, 0xf3, 0xdc, 0x91, 0xe9, 0x55,
					0x10}
				return util.NewAddressScriptHashFromHash(hash, util.Bech32PrefixDAGCoin)
			},
			prefix: util.Bech32PrefixDAGCoin,
		},
	}

	for _, test := range tests {
		// Decode addr and compare error against valid.
		decoded, err := util.DecodeAddress(test.addr, test.prefix)
		if (err == nil) != test.valid {
			t.Errorf("%v: decoding test failed: %v", test.name, err)
			return
		}

		if err == nil {
			// Ensure the stringer returns the same address as the
			// original.
			if decodedStringer, ok := decoded.(fmt.Stringer); ok {
				addr := test.addr

				if addr != decodedStringer.String() {
					t.Errorf("%v: String on decoded value does not match expected value: %v != %v",
						test.name, test.addr, decodedStringer.String())
					return
				}
			}

			// Encode again and compare against the original.
			encoded := decoded.EncodeAddress()
			if test.encoded != encoded {
				t.Errorf("%v: decoding and encoding produced different addressess: %v != %v",
					test.name, test.encoded, encoded)
				return
			}

			// Perform type-specific calculations.
			var saddr []byte
			switch d := decoded.(type) {
			case *util.AddressPubKeyHash:
				saddr = util.TstAddressSAddr(encoded)

			case *util.AddressScriptHash:
				saddr = util.TstAddressSAddr(encoded)

			case *util.AddressPubKey:
				// Ignore the error here since the script
				// address is checked below.
				saddr, _ = hex.DecodeString(d.String())
			}

			// Check script address, as well as the Hash160 method for P2PKH and
			// P2SH addresses.
			if !bytes.Equal(saddr, decoded.ScriptAddress()) {
				t.Errorf("%v: script addresses do not match:\n%x != \n%x",
					test.name, saddr, decoded.ScriptAddress())
				return
			}
			switch a := decoded.(type) {
			case *util.AddressPubKeyHash:
				if h := a.Hash160()[:]; !bytes.Equal(saddr, h) {
					t.Errorf("%v: hashes do not match:\n%x != \n%x",
						test.name, saddr, h)
					return
				}

			case *util.AddressScriptHash:
				if h := a.Hash160()[:]; !bytes.Equal(saddr, h) {
					t.Errorf("%v: hashes do not match:\n%x != \n%x",
						test.name, saddr, h)
					return
				}
			}

			// Ensure the address is for the expected network.
			if !decoded.IsForPrefix(test.prefix) {
				t.Errorf("%v: calculated network does not match expected",
					test.name)
				return
			}
		}

		if !test.valid {
			// If address is invalid, but a creation function exists,
			// verify that it returns a nil addr and non-nil error.
			if test.f != nil {
				_, err := test.f()
				if err == nil {
					t.Errorf("%v: address is invalid but creating new address succeeded",
						test.name)
					return
				}
			}
			continue
		}

		// Valid test, compare address created with f against expected result.
		addr, err := test.f()
		if err != nil {
			t.Errorf("%v: address is valid but creating new address failed with error %v",
				test.name, err)
			return
		}

		if !reflect.DeepEqual(addr, test.result) {
			t.Errorf("%v: created address does not match expected result",
				test.name)
			return
		}
	}
}

func TestDecodeAddressErrorConditions(t *testing.T) {
	tests := []struct {
		address      string
		prefix       util.Bech32Prefix
		errorMessage string
	}{
		{
			"bitcoincash:qpzry9x8gf2tvdw0s3jn54khce6mua7lcw20ayyn",
			util.Bech32PrefixUnknown,
			"decoded address's prefix could not be parsed",
		},
		{
			"dagreg:qpm2qsznhks23z7629mms6s4cwef74vcwvtmvqeszh",
			util.Bech32PrefixDAGTest,
			"decoded address is of wrong network",
		},
		{
			"dagreg:raskzctpv9skzctpv9skzctpv9skzctpvyd070wnqg",
			util.Bech32PrefixDAGReg,
			"unknown address type",
		},
		{
			"dagreg:raskzcg5egs6nnj",
			util.Bech32PrefixDAGReg,
			"decoded address is of unknown size",
		},
	}

	for _, test := range tests {
		_, err := util.DecodeAddress(test.address, test.prefix)
		if err == nil {
			t.Errorf("decodeAddress unexpectedly succeeded")
		} else if !strings.Contains(err.Error(), test.errorMessage) {
			t.Errorf("received mismatched error. Expected %s but got %s",
				test.errorMessage, err)
		}
	}
}

func TestParsePrefix(t *testing.T) {
	tests := []struct {
		prefixStr      string
		expectedPrefix util.Bech32Prefix
		expectedError  bool
	}{
		{"dagcoin", util.Bech32PrefixDAGCoin, false},
		{"dagreg", util.Bech32PrefixDAGReg, false},
		{"dagtest", util.Bech32PrefixDAGTest, false},
		{"dagsim", util.Bech32PrefixDAGSim, false},
		{"blabla", util.Bech32PrefixUnknown, true},
		{"unknown", util.Bech32PrefixUnknown, true},
		{"", util.Bech32PrefixUnknown, true},
	}

	for _, test := range tests {
		result, err := util.ParsePrefix(test.prefixStr)
		if (err != nil) != test.expectedError {
			t.Errorf("TestParsePrefix: %s: expected error status: %t, but got %t",
				test.prefixStr, test.expectedError, err != nil)
		}

		if result != test.expectedPrefix {
			t.Errorf("TestParsePrefix: %s: expected prefix: %d, but got %d",
				test.prefixStr, test.expectedPrefix, result)
		}
	}
}

func TestPrefixToString(t *testing.T) {
	tests := []struct {
		prefix            util.Bech32Prefix
		expectedPrefixStr string
	}{
		{util.Bech32PrefixDAGCoin, "dagcoin"},
		{util.Bech32PrefixDAGReg, "dagreg"},
		{util.Bech32PrefixDAGTest, "dagtest"},
		{util.Bech32PrefixDAGSim, "dagsim"},
		{util.Bech32PrefixUnknown, ""},
	}

	for _, test := range tests {
		result := test.prefix.String()

		if result != test.expectedPrefixStr {
			t.Errorf("TestPrefixToString: %s: expected string: %s, but got %s",
				test.prefix, test.expectedPrefixStr, result)
		}
	}
}