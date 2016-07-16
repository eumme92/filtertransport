package filtertransport

import (
	"errors"
	"net"
	"testing"
)

func TestFilterPrivate(t *testing.T) {
	for _, c := range []struct {
		ip       net.IP
		fileterd bool
	}{
		{net.ParseIP("169.253.255.255"), false},
		{net.ParseIP("169.254.0.1"), true},
		{net.ParseIP("169.254.255.255"), true},
		{net.ParseIP("169.255.0.1"), false},

		{net.ParseIP("172.15.255.255"), false},
		{net.ParseIP("172.16.0.1"), true},
		{net.ParseIP("172.31.255.255"), true},
		{net.ParseIP("172.32.0.1"), false},

		{net.ParseIP("192.167.255.255"), false},
		{net.ParseIP("192.168.0.1"), true},
		{net.ParseIP("192.168.255.255"), true},
		{net.ParseIP("192.169.0.1"), false},

		{net.ParseIP("9.255.255.255"), false},
		{net.ParseIP("10.0.0.1"), true},
		{net.ParseIP("10.255.255.255"), true},
		{net.ParseIP("11.0.0.1"), false},

		{net.ParseIP("126.255.255.255"), false},
		{net.ParseIP("127.0.0.1"), true},
		{net.ParseIP("127.255.255.255"), true},
		{net.ParseIP("128.0.0.1"), false},

		{net.ParseIP("::1"), true},
		{net.ParseIP("::2"), false},

		{net.ParseIP("fb00::1"), false},
		{net.ParseIP("fc00::1"), true},
		{net.ParseIP("fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff"), true},
		{net.ParseIP("fe00::1"), false},
	} {
		if err := FilterPrivate(net.TCPAddr{IP: c.ip}); (err != nil) != c.fileterd {
			t.Errorf("%v should be %t", c.ip, c.fileterd)
		}
	}
}

func TestFilterDial(t *testing.T) {
	FilterDial(
		"tcp", "1.2.3.4:1234",
		func(addr net.TCPAddr) error { return errors.New("") },
		func(network string, address string) (net.Conn, error) {
			t.Errorf("dail should not be called on filter error")
			return nil, nil
		},
	)

	called := false
	FilterDial(
		"tcp", "1.2.3.4:1234",
		func(addr net.TCPAddr) error { return nil },
		func(network string, address string) (net.Conn, error) {
			called = true
			return nil, nil
		},
	)
	if !called {
		t.Errorf("dail should be called on no filter error")
	}
}
