package enode

import (
	"net"
	"reflect"
	"strings"
	"testing"
)

var parseNodeWithResolveTests = []struct {
	rawurl     string
	wantError  string
	wantResult *Node
	host       string
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
		wantError: `invalid domain or IP address`,
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
		rawurl:    "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@localhost:3",
		wantResult: NewV4(
			hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
			net.IP{0x7f, 0x0, 0x0, 0x1},
			3,
			3,
		),
	},
	{
		rawurl:    "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@localhost",
		wantResult: NewV4(
			hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
			net.IP{0x7f, 0x0, 0x0, 0x1},
			30303,
			30303,
		),
	},
	{
		rawurl:    "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@gdns.oogle.com:3",
		wantResult: NewV4(
			hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
			net.IP{0xd8, 0xef, 0x20, 0x3b},
			3,
			3,
		),
		host: "any-in-",
	},
	{
		rawurl:    "enode://1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439@dns.google.com",
		wantResult: NewV4(
			hexPubkey("1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
			net.IP{0xd8, 0xef, 0x20, 0x3b},
			30303,
			30303,
		),
		host: "any-in-",
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
		wantError: `parse ://foo: missing protocol scheme`,
	},
}

func TestParseNodeWithDomainResolution(t *testing.T) {
	for _, test := range parseNodeWithResolveTests {
		n, err := ParseV4WithResolve(test.rawurl)
		if test.wantError != "" {
			if err == nil {
				t.Errorf("test %q:\n  got nil error, expected %#q", test.rawurl, test.wantError)
				continue
			} else if !strings.Contains(err.Error(), test.wantError) {
				t.Errorf("test %q:\n  got error %#q, expected %#q\n%v", test.rawurl, err.Error(), test.wantError, n)
				continue
			}
		} else {
			if err != nil {
				t.Errorf("test %q:\n  unexpected error: %v", test.rawurl, err)
				continue
			}
			if test.host != "" {
				hosts, err := net.LookupAddr(test.wantResult.IP().String())
				if err != nil {
					t.Errorf("failed to resolve host: %s", err)
				}

				if len(hosts) == 0 {
					t.Errorf("failed to resolve host. Hosts list is empty")
				}

				if !strings.HasPrefix(hosts[0], test.host) {
					t.Errorf("got wrong host. got: %#v\nwant: %#v\n", hosts[0], test.host)
				}
			} else if !reflect.DeepEqual(n, test.wantResult) {
				t.Errorf("test %q:\n result mismatch:\n" +
					"got PubKey: %#v\nwant: %#v\n\n" +
					"got IP: %#v\nwant: %#v(equal %v)\n\n" +
					"got TCP: %#v\nwant: %#v\n\n" +
					"got UPD: %#v\nwant: %#v\n",
					test.rawurl,
					n.Pubkey(), test.wantResult.Pubkey(),
					n.IP().String(), test.wantResult.IP().String(),
					n.IP().Equal(test.wantResult.IP()),
					n.TCP(), test.wantResult.TCP(),
					n.UDP(), test.wantResult.UDP())
			}
		}
	}
}