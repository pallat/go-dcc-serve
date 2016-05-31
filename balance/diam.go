// Copyright 2013-2015 go-diameter authors.  All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

// Diameter client example. This is by no means a complete client.
// The commands in here are not fully implemented. For that you have
// to read the RFCs (base and credit control) and follow the spec.

package balance

import (
	"log"
	"server/dictionary"

	"github.com/ant0ine/go-json-rest/rest"
	"github.com/fiorix/go-diameter/diam"
	// "github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"github.com/fiorix/go-diameter/diam/dict"

	"server/diameter"
)

const (
	identity    = datatype.DiameterIdentity("jenkin13_OMR_TEST01")
	realm       = datatype.DiameterIdentity("dtac.co.th")
	vendorID    = datatype.Unsigned32(0)
	productName = datatype.UTF8String("omr")

	dtacAddr = "10.89.111.12:6553"
	dtnAddr  = "10.89.111.40:6573"
)

var (
	cdtac             diam.Conn
	cdtn              diam.Conn
	DefaultRestWriter rest.ResponseWriter
)

func Start() {
	var err error

	response = make(map[string]chan BalanceInfo)
	dict.Default = dictionary.Load()
	diam.HandleFunc("CEA", diameter.OnCEA)
	diam.HandleFunc("CCA", OnCCA)

	cdtac, err = diam.Dial(dtacAddr, nil, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	go diameter.Serve(cdtac, identity, realm, vendorID, productName)

	cdtn, err = diam.Dial(dtnAddr, nil, nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	go diameter.Serve(cdtn, identity, realm, vendorID, productName)
}

func OnCCA(c diam.Conn, m *diam.Message) {
	log.Printf("Receiving message from %s", c.RemoteAddr().String())
	if m.Header.CommandCode == 272 {
		var balance BalanceInfo
		m.Unmarshal(&balance)
		response[balance.SessionId] <- balance
	}
}
