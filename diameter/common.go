package diameter

import (
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
)

func Serve(c diam.Conn, identity, realm, vendorID, productName datatype.Type) {
	var (
		err error
	)

	err = Cer(c, identity, realm, vendorID, productName)
	if err != nil {
		log.Println(err.Error())
	}

	go Watchdog(c, identity, realm)

	log.Println(<-diam.ErrorReports())
	log.Println(<-c.(diam.CloseNotifier).CloseNotify())
	log.Println("Server disconnected.")
}

func Cer(c diam.Conn, identity, realm, vendorID, productName datatype.Type) (err error) {
	m := diam.NewRequest(diam.CapabilitiesExchange, 0, nil)
	m.NewAVP(avp.OriginHost, avp.Mbit, 0, identity)
	m.NewAVP(avp.OriginRealm, avp.Mbit, 0, realm)
	laddr := c.LocalAddr()
	ip, _, _ := net.SplitHostPort(laddr.String())
	m.NewAVP(avp.HostIPAddress, avp.Mbit, 0, datatype.Address(net.ParseIP(ip)))
	m.NewAVP(avp.VendorID, avp.Mbit, 0, vendorID)
	m.NewAVP(avp.ProductName, 0, 0, productName)
	m.NewAVP(avp.SupportedVendorID, avp.Mbit, 0, datatype.Unsigned32(0))
	m.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4))
	m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(0))
	m.NewAVP(avp.AcctApplicationID, avp.Mbit, 0, datatype.Unsigned32(4))
	m.NewAVP(avp.FirmwareRevision, avp.Mbit, 0, datatype.Unsigned32(1))

	log.Printf("Sending message to %s", c.RemoteAddr().String())

	if _, err = m.WriteTo(c); err != nil {
		log.Fatal("Write failed:", err)
		return err
	}

	return
}

func Watchdog(c diam.Conn, identity, realm datatype.Type) {
	for {
		time.Sleep(10 * time.Second)
		m := diam.NewRequest(diam.DeviceWatchdog, 0, nil)
		m.NewAVP(avp.OriginHost, avp.Mbit, 0, identity)
		m.NewAVP(avp.OriginRealm, avp.Mbit, 0, realm)
		m.NewAVP(avp.OriginStateID, avp.Mbit, 0, datatype.Unsigned32(rand.Uint32()))
		log.Printf("Sending message to %s", c.RemoteAddr().String())

		if _, err := m.WriteTo(c); err != nil {
			log.Fatal("Write failed:", err)
		}
	}
}

func OnCEA(c diam.Conn, m *diam.Message) {
	rc, err := m.FindAVP(avp.ResultCode)
	if err != nil {
		log.Fatal(err)
		return
	}
	if v, _ := rc.Data.(datatype.Unsigned32); v != diam.Success {
		log.Fatal("Unexpected response:", rc)
		return
	}

	log.Printf("Receiving message from %s", c.RemoteAddr().String())
}

func OnMSG(c diam.Conn, m *diam.Message) {
	log.Printf("-Receiving message from %s", c.RemoteAddr().String())
	// log.Println(m)
}
