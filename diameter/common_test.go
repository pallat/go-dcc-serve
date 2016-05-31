package diameter_test

import (
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/diamtest"

	"server/diameter"
)

const (
	identity    = datatype.DiameterIdentity("jenkin13_OMR_TEST01")
	realm       = datatype.DiameterIdentity("dtac.co.th")
	vendorID    = datatype.Unsigned32(0)
	productName = datatype.UTF8String("omr")
)

func TestCer(t *testing.T) {
	errc := make(chan error, 1)

	smux := diam.NewServeMux()
	smux.Handle("CER", handleCER(errc))

	srv := diamtest.NewServer(smux, nil)
	defer srv.Close()

	wait := make(chan struct{})
	cmux := diam.NewServeMux()
	cmux.Handle("CEA", handleCEA(errc, wait))

	cli, err := diam.Dial(srv.Address, cmux, nil)
	if err != nil {
		t.Fatal(err)
	}

	diameter.Cer(cli, identity, realm, vendorID, productName)

	select {
	case <-wait:
	case err := <-errc:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("Timed out: no CER or CEA received")
	}
}

func TestWatchdog(t *testing.T) {
	errc := make(chan error, 1)

	smux := diam.NewServeMux()
	smux.Handle("CER", handleCER(errc))

	srv := diamtest.NewServer(smux, nil)
	defer srv.Close()

	wait := make(chan struct{})
	cmux := diam.NewServeMux()
	cmux.Handle("CEA", handleCEA(errc, wait))

	cli, err := diam.Dial(srv.Address, cmux, nil)
	if err != nil {
		t.Fatal(err)
	}

	diameter.Cer(cli, identity, realm, vendorID, productName)
	go diameter.Watchdog(cli, identity, realm)

	select {
	case <-wait:
	case err := <-errc:
		t.Fatal(err)
	case <-time.After(time.Second):
		t.Fatal("Timed out: no CER or CEA received")
	}
}

func BenchmarkCer(b *testing.B) {
	// run the Fib function b.N times
	for n := 0; n < b.N; n++ {
		errc := make(chan error, 1)

		smux := diam.NewServeMux()
		smux.Handle("CER", handleCER(errc))

		srv := diamtest.NewServer(smux, nil)
		defer srv.Close()

		wait := make(chan struct{})
		cmux := diam.NewServeMux()
		cmux.Handle("CEA", handleCEA(errc, wait))

		cli, err := diam.Dial(srv.Address, cmux, nil)
		if err != nil {
			b.Fatal(err)
		}

		diameter.Cer(cli, identity, realm, vendorID, productName)
		go diameter.Watchdog(cli, identity, realm)

		select {
		case <-wait:
		case err := <-errc:
			b.Fatal(err)
		case <-time.After(time.Second):
			b.Fatal("Timed out: no CER or CEA received")
		}
	}
}

func handleCER(errc chan error) diam.HandlerFunc {
	type CER struct {
		OriginHost        string    `avp:"Origin-Host"`
		OriginRealm       string    `avp:"Origin-Realm"`
		VendorID          int       `avp:"Vendor-Id"`
		ProductName       string    `avp:"Product-Name"`
		OriginStateID     *diam.AVP `avp:"Origin-State-Id"`
		AcctApplicationID *diam.AVP `avp:"Acct-Application-Id"`
	}
	return func(c diam.Conn, m *diam.Message) {
		go func() {
			<-c.(diam.CloseNotifier).CloseNotify()
			//log.Println("Client", c.RemoteAddr(), "disconnected")
		}()
		var req CER
		err := m.Unmarshal(&req)
		if err != nil {
			errc <- err
			return
		}
		if req.OriginHost != "jenkin13_OMR_TEST01" {
			errc <- fmt.Errorf("Unexpected OriginHost. Want cli, have %q", req.OriginHost)
			return
		}
		if req.OriginRealm != "dtac.co.th" {
			errc <- fmt.Errorf("Unexpected OriginRealm. Want localhost, have %q", req.OriginRealm)
			return
		}
		if req.VendorID != 0 {
			errc <- fmt.Errorf("Unexpected VendorID. Want 99, have %d", req.VendorID)
			return
		}
		if req.ProductName != "omr" {
			errc <- fmt.Errorf("Unexpected ProductName. Want go-diameter, have %q", req.ProductName)
			return
		}

		a := m.Answer(diam.Success)
		_, err = sendCEA(c, a, req.OriginStateID, req.AcctApplicationID)
		if err != nil {
			errc <- err
		}
	}
}

func sendCEA(w io.Writer, m *diam.Message, OriginStateID, AcctApplicationID *diam.AVP) (n int64, err error) {
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.OctetString("srv"))
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.OctetString("localhost"))
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP("127.0.0.1")))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(99))
	m.NewAVP(avp.ProductName, avp.Mbit, 0, datatype.UTF8String("go-diameter"))
	m.AddAVP(OriginStateID)
	m.AddAVP(AcctApplicationID)
	return m.WriteTo(w)
}

func handleCEA(errc chan error, wait chan struct{}) diam.HandlerFunc {
	type CEA struct {
		OriginHost        string `avp:"Origin-Host"`
		OriginRealm       string `avp:"Origin-Realm"`
		VendorID          int    `avp:"Vendor-Id"`
		ProductName       string `avp:"Product-Name"`
		OriginStateID     int    `avp:"Origin-State-Id"`
		AcctApplicationID int    `avp:"Acct-Application-Id"`
	}
	return func(c diam.Conn, m *diam.Message) {
		defer close(wait)
		fmt.Println("Receiving message from %s", c.RemoteAddr().String())
		// var resp CEA
		// err := m.Unmarshal(&resp)
		// if err != nil {
		// 	errc <- err
		// 	return
		// }
		// if resp.OriginHost != "srv" {
		// 	errc <- fmt.Errorf("Unexpected OriginHost. Want srv, have %q", resp.OriginHost)
		// 	return
		// }
		// if resp.OriginRealm != "localhost" {
		// 	errc <- fmt.Errorf("Unexpected OriginRealm. Want localhost, have %q", resp.OriginRealm)
		// 	return
		// }
		// if resp.VendorID != 99 {
		// 	errc <- fmt.Errorf("Unexpected VendorID. Want 99, have %d", resp.VendorID)
		// 	return
		// }
		// if resp.ProductName != "go-diameter" {
		// 	errc <- fmt.Errorf("Unexpected ProductName. Want go-diameter, have %q", resp.ProductName)
		// 	return
		// }
		// if resp.OriginStateID != 1234 {
		// 	errc <- fmt.Errorf("Unexpected OriginStateID. Want 1234, have %d", resp.OriginStateID)
		// 	return
		// }
		// if resp.AcctApplicationID != 1 {
		// 	errc <- fmt.Errorf("Unexpected AcctApplicationID. Want 1, have %d", resp.AcctApplicationID)
		// 	return
		// }
	}
}
