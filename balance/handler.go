package balance

import (
	"fmt"
	"time"

	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"

	"github.com/ant0ine/go-json-rest/rest"
)

const (
	BalanceInformation  = 21100
	AccessMethod        = 20340
	AccountQueryMethod  = 20346
	SSPTime             = 20386
	CallingPartyAddress = 20336
)

var response map[string]chan BalanceInfo

type BalanceInfo struct {
	SessionId          string `avp:"Session-Id"`
	ServiceInformation struct {
		BalanceInformation struct {
			FirstActiveDate   string `avp:"First-Active-Date"`
			SubscriberState   int    `avp:"Subscriber-State"`
			ActivePeriod      string `avp:"Active-Period"`
			GracePeriod       string `avp:"Grace-Period"`
			DisablePeriod     string `avp:"Disable-Period"`
			Balance           int    `avp:"Balance"`
			LanguageIVR       int    `avp:"Language-IVR"`
			LanguageSMS       int    `avp:"Language-SMS"`
			LanguageUSSD      int    `avp:"Language-USSD"`
			AccountChangeInfo []struct {
				AccountId             string `avp:"Account-Id"`
				AccountType           int    `avp:"Account-Type"`
				AccountTypeDesc       string `avp:"Account-Type-Desc"`
				AccountBeginDate      string `avp:"Account-Begin-Date"`
				RelatedType           int    `avp:"Related-Type"`
				RelatedObjectID       string `avp:"Related-Object-ID"`
				CurrentAccountBalance int    `avp:"Current-Account-Balance"`
				AccountEndDate        string `avp:"Account-End-Date"`
				MeasureType           int    `avp:"Measure-Type"`
			} `avp:"Account-Change-Info"`
			OfferInformation []struct {
				OfferInfo []struct {
					OfferOrderKey            string `avp:"Offer-Order-Key"`
					EffectiveTime            string `avp:"Effective-Time"`
					Status                   string `avp:"Status"`
					CurrentCycle             int    `avp:"Current-Cycle"`
					TotalCycle               int    `avp:"Total-Cycle"`
					OfferOrderIntegrationKey string `avp:"Offer-Order-Integration-Key"`
					ExternalOfferCode        string `avp:"External-Offer-Code"`
				} `avp:"Offer-Info"`
			} `avp:"Offer-Information,omitempty"`
		} `avp:"Balance-Information"`
	} `avp:"Service-Information"`
}

func Balance(w rest.ResponseWriter, req *rest.Request) {
	corp := req.PathParam("corp")
	subr := req.PathParam("subr")
	r := diam.NewRequest(diam.CreditControl, 4, nil)

	sessionID := "dtac.co.th;OMR" + time.Now().Format("200601021504050000")

	r.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sessionID))
	r.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(4))
	r.NewAVP(avp.DestinationRealm, avp.Mbit, 0, datatype.OctetString("www.huawei.com")) //peerRealm.Data)
	r.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.OctetString("jenkin13_OMR_TEST01"))  //identity)
	r.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.OctetString("dtac.co.th"))          //realm)
	r.NewAVP(avp.CCRequestType, avp.Mbit, 0, datatype.Integer32(4))
	r.NewAVP(avp.SubscriptionID, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.SubscriptionIDType, avp.Mbit, 0, datatype.Integer32(0)),
			diam.NewAVP(avp.SubscriptionIDData, avp.Mbit, 0, datatype.UTF8String(subr)),
		},
	})
	r.NewAVP(avp.ServiceContextID, avp.Mbit, 0, datatype.UTF8String("QueryBalance@huawei.com"))
	r.NewAVP(avp.RequestedAction, avp.Mbit, 0, datatype.Integer32(2))
	r.NewAVP(avp.EventTimestamp, avp.Mbit, 0, datatype.Time(time.Now()))
	r.NewAVP(avp.ServiceIdentifier, avp.Mbit, 0, datatype.Unsigned32(0))
	r.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, datatype.Unsigned32(0))
	r.NewAVP(avp.RouteRecord, avp.Mbit, 0, datatype.OctetString("10.89.111.40"))
	r.NewAVP(avp.DestinationHost, avp.Mbit, 0, datatype.OctetString("cbp211"))
	r.NewAVP(avp.ServiceInformation, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(BalanceInformation, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(CallingPartyAddress, avp.Mbit, 0, datatype.UTF8String(subr)),
					diam.NewAVP(AccessMethod, avp.Mbit, 0, datatype.Unsigned32(9)),
					diam.NewAVP(AccountQueryMethod, avp.Mbit, 0, datatype.Unsigned32(1)),
					diam.NewAVP(SSPTime, avp.Mbit, 0, datatype.Time(time.Now())),
				},
			}),
		},
	})

	var err error

	response[sessionID] = make(chan BalanceInfo)
	if corp == "dtac" {
		_, err = r.WriteTo(cdtac)
	} else {
		_, err = r.WriteTo(cdtn)
	}

	if err != nil {
		fmt.Println(err.Error())
	}

	resp := <-response[sessionID]

	w.WriteJson(resp)
}
