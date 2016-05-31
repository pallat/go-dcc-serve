package diameter

import (
	"reflect"
	"testing"
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
)

type balanceReq struct {
	SessionID           string    `dcode:"263" dtype:"UTF8String"`
	AuthApplicationID   int       `dcode:"258" dtype:"Unsigned32"`
	DestinationRealm    string    `dcode:"283" dtype:"DiameterIdentity"`
	OriginHost          string    `dcode:"264" dtype:"OctetString"`
	OriginRealm         string    `dcode:"296" dtype:"OctetString"`
	CCRequestType       int       `dcode:"416" dtype:"Integer32"`
	SubscriptionIDType  int       `dcode:"443>450" dtype:"Integer32"`
	SubscriptionIDData  string    `dcode:"443>444" dtype:"UTF8String"`
	ServiceContextID    string    `dcode:"461" dtype:"UTF8String"`
	RequestedAction     int       `dcode:"436" dtype:"Integer32"`
	EventTimestamp      time.Time `dcode:"55" dtype:"Time"`
	ServiceIdentifier   int       `dcode:"439" dtype:"Unsigned32"`
	CCRequestNumber     int       `dcode:"415" dtype:"Unsigned32"`
	RouteRecord         string    `dcode:"282" dtype:"OctetString"`
	DestinationHost     string    `dcode:"293" dtype:"OctetString"`
	CallingPartyAddress string    `dcode:"873>21100>20336" dtype:"UTF8String"`
	AccessMethod        int       `dcode:"873>21100>20340" dtype:"Unsigned32"`
	AccountQueryMethod  int       `dcode:"873>21100>20346" dtype:"Unsigned32"`
	SSPTime             time.Time `dcode:"873>21100>20386" dtype:"Time"`
}

func TestGetTagDcode(t *testing.T) {
	var req interface{} = balanceReq{
		SessionID: "dtac.co.th;OMR2014",
	}

	v := reflect.TypeOf(req)

	if Dcode(v.Field(0))[0] != "263" {
		t.Error("It should be 263 but was ", Dcode(v.Field(0))[0])
	}
	if Dtype(v.Field(0))[0] != "UTF8String" {
		t.Error("It should be UTF8String but was ", Dtype(v.Field(0)))
	}

}

func TestGetTagDcodeThreeLayers(t *testing.T) {
	var req interface{} = balanceReq{
		AccountQueryMethod: 9,
	}

	v := reflect.TypeOf(req)

	if !reflect.DeepEqual(Dcode(v.Field(17)), []string{"873", "21100", "20346"}) {
		t.Error("It should be ", []string{"873", "21100", "20346"}, " but was ", Dcode(v.Field(17)))
	}
	if Dtype(v.Field(17))[0] != "Unsigned32" {
		t.Error("It should be Unsigned32 but was ", Dtype(v.Field(17)))
	}
}

func TestEncodeToAVPUTF8String(t *testing.T) {
	var req interface{} = balanceReq{
		SessionID: "dtac.co.th;OMR2014",
	}

	encoded := Encode(&req)

	expected := []*diam.AVP{}
	expected = append(expected, diam.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("dtac.co.th;OMR2014")))

	if !reflect.DeepEqual(encoded, expected) {
		t.Error("Expected: ", expected)
		t.Error("But got: ", encoded)
	}
}

func TestEncodeToAVPUTF8StringWhenPassedNotInterface(t *testing.T) {
	var req interface{} = balanceReq{
		SessionID: "dtac.co.th;OMR2014",
	}

	encoded := Encode(req)

	expected := []*diam.AVP{}
	expected = append(expected, diam.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String("dtac.co.th;OMR2014")))

	if !reflect.DeepEqual(encoded, expected) {
		t.Error("Expected: ", expected)
		t.Error("But got: ", encoded)
	}
}

func xTestEncode(t *testing.T) {
	sessionID := "dtac.co.th;OMR" + time.Now().Format("200601021504050000")

	req := balanceReq{
		SessionID:           sessionID,
		AuthApplicationID:   4,
		DestinationRealm:    "www.huawei.com",
		OriginHost:          "jenkin13_OMR_TEST01",
		OriginRealm:         "dtac.co.th",
		CCRequestType:       4,
		SubscriptionIDType:  0,
		SubscriptionIDData:  "66947451960",
		ServiceContextID:    "QueryBalance@huawei.com",
		RequestedAction:     2,
		EventTimestamp:      time.Now(),
		ServiceIdentifier:   0,
		CCRequestNumber:     0,
		RouteRecord:         "10.89.111.40",
		DestinationHost:     "cbp211",
		CallingPartyAddress: "66947451960",
		AccessMethod:        9,
		AccountQueryMethod:  1,
		SSPTime:             time.Now(),
	}

	davp := []*diam.AVP{}
	davp = append(davp, diam.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sessionID)))
	davp = append(davp, diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4)))
	davp = append(davp, diam.NewAVP(avp.DestinationRealm, avp.Mbit, 0, datatype.OctetString("www.huawei.com"))) //peerRealm.Data)
	davp = append(davp, diam.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.OctetString("jenkin13_OMR_TEST01")))  //identity)
	davp = append(davp, diam.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.OctetString("dtac.co.th")))          //realm)
	davp = append(davp, diam.NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Integer32(4)))
	davp = append(davp, diam.NewAVP(avp.SubscriptionID, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Integer32(0)),
			diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String("66947451960")),
		},
	}))
	davp = append(davp, diam.NewAVP(avp.ServiceContextID, avp.Mbit, 0, datatype.UTF8String("QueryBalance@huawei.com")))
	davp = append(davp, diam.NewAVP(avp.RequestedAction, avp.Mbit, 0, datatype.Integer32(2)))
	davp = append(davp, diam.NewAVP(avp.EventTimestamp, avp.Mbit, 0, datatype.Time(time.Now())))
	davp = append(davp, diam.NewAVP(avp.ServiceIdentifier, avp.Mbit, 0, datatype.Unsigned32(0)))
	davp = append(davp, diam.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, datatype.Unsigned32(0)))
	davp = append(davp, diam.NewAVP(avp.RouteRecord, avp.Mbit, 0, datatype.OctetString("10.89.111.40")))
	davp = append(davp, diam.NewAVP(avp.DestinationHost, avp.Mbit, 0, datatype.OctetString("cbp211")))
	davp = append(davp, diam.NewAVP(avp.ServiceInformation, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(21100, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(20336, avp.Mbit, 0, datatype.UTF8String("66947451960")),
					diam.NewAVP(20340, avp.Mbit, 0, datatype.Unsigned32(9)),
					diam.NewAVP(20346, avp.Mbit, 0, datatype.Unsigned32(1)),
					diam.NewAVP(20386, avp.Mbit, 0, datatype.Time(time.Now())),
				},
			}),
		},
	}))

	expected := davp

	encoded := Encode(req)

	if !reflect.DeepEqual(expected, encoded) {
		t.Error("Expected: ", expected)
		t.Error("But Got: ", encoded)
	}

}
