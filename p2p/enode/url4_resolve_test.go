package enode

import (
	"crypto/ecdsa"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/autonity/autonity/p2p/enr"
)

func newV4WithHostNoError(t *testing.T, pubkey *ecdsa.PublicKey, host string, tcp, udp int, resolveFunc func(host string) ([]net.IP, error)) *Node {
	n, err := NewV4WithHost(pubkey, host, tcp, udp, resolveFunc)
	require.NoError(t, err)
	return n
}

func TestParseNodeWithDomainResolution(t *testing.T) {
	var parseNodeWithResolveTests = []struct {
		rawurl     string
		wantError  string
		wantResult *Node
	}{
		{
			rawurl:    "http://foobar",
			wantError: `invalid URL scheme, want "enode"`,
		},
		{
			rawurl:    "enode://01010101@123.124.125.126:3",
			wantError: `invalid public key (wrong length, want 128 hex chars)`,
		},
		// Complete nodes with IP address.
		{
			rawurl:    "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@hostname:3",
			wantError: `lookup hostname`,
		},
		{
			rawurl:    "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@127.0.0.1:foo",
			wantError: `invalid port`,
		},
		{
			rawurl:    "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@127.0.0.1:3?discport=foo",
			wantError: `invalid discport in query`,
		},
		{
			rawurl: "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@localhost:3",
			wantResult: newV4WithHostNoError(
				t,
				hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
				"localhost",
				3,
				3,
				net.LookupIP,
			),
		},
		{
			rawurl: "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@localhost",
			wantResult: newV4WithHostNoError(
				t,
				hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
				"localhost",
				30303,
				30303,
				net.LookupIP,
			),
		},
		{
			rawurl: "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@gdns.oogle.com:3",
			wantResult: newV4WithHostNoError(
				t,
				hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
				"gdns.oogle.com",
				3,
				3,
				net.LookupIP,
			),
		},
		{
			rawurl: "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@dns.google.com",
			wantResult: newV4WithHostNoError(
				t,
				hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
				"dns.google.com",
				30303,
				30303,
				net.LookupIP,
			),
		},
		{
			rawurl: "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@127.0.0.1:52150",
			wantResult: NewV4(
				hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
				net.IP{0x7f, 0x0, 0x0, 0x1},
				52150,
				52150,
			),
		},
		{
			rawurl: "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@[::]:52150",
			wantResult: NewV4(
				hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
				net.ParseIP("::"),
				52150,
				52150,
			),
		},
		{
			rawurl: "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@[2001:db8:3c4d:15::abcd:ef12]:52150",
			wantResult: NewV4(
				hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
				net.ParseIP("2001:db8:3c4d:15::abcd:ef12"),
				52150,
				52150,
			),
		},
		{
			rawurl: "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@127.0.0.1:52150?discport=22334",
			wantResult: NewV4(
				hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
				net.IP{0x7f, 0x0, 0x0, 0x1},
				52150,
				22334,
			),
		},
		// Incomplete nodes with no address.
		{
			rawurl: "1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439",
			wantResult: NewV4(
				hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
				nil, 0, 0,
			),
		},
		{
			rawurl: "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439",
			wantResult: NewV4(
				hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
				nil, 0, 0,
			),
		},
		// Invalid URLs
		{
			rawurl:    "01010101",
			wantError: `invalid public key (wrong length, want 128 hex chars)`,
		},
		{
			rawurl:    "enode://01010101",
			wantError: `invalid public key (wrong length, want 128 hex chars)`,
		},
		{
			// This test checks that errors from url.Parse are handled.
			rawurl:    "://foo",
			wantError: `parse "://foo": missing protocol scheme`,
		},
	}

	for i, test := range parseNodeWithResolveTests {
		n, err := ParseV4(test.rawurl)

		var gotErr string
		if err != nil {
			gotErr = strings.ReplaceAll(err.Error(), "\"", "")
		}

		wantError := strings.ReplaceAll(test.wantError, "\"", "")
		if wantError != "" {
			if err == nil {
				t.Errorf("test %q:\n  got nil error, expected %#q", test.rawurl, wantError)
				continue
			} else if !strings.Contains(gotErr, wantError) {
				t.Errorf("test %q:\n  got error %#q, expected %#q\n%v", test.rawurl, gotErr, wantError, n)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("test %q:\n  unexpected error: %v", test.rawurl, err)
				continue
			} else {
				// Check that an IP was set
				assert.True(t, len(n.IP().String()) > 0)

				// Nullify the IPs since dns resolution can return different
				// results, for example when round robining is setup.
				zeroIP := net.IPv4(0, 0, 0, 0)
				n.r.Set(enr.IP(zeroIP))
				test.wantResult.r.Set(enr.IP(zeroIP))
				n.r.Set(enr.IP(net.IPv6zero))
				test.wantResult.r.Set(enr.IP(net.IPv6zero))

				// Function references are not comparable, so we nullify them
				// to allow comparing the remaining fields.
				n.resolveFunc = nil
				test.wantResult.resolveFunc = nil
				assert.Equal(t, test.wantResult, n, i)
			}
		}
	}
}
