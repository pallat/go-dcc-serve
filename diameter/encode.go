package diameter

import (
	"fmt"
	"github.com/fiorix/go-diameter/diam"
	"github.com/fiorix/go-diameter/diam/avp"
	"github.com/fiorix/go-diameter/diam/datatype"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type Encoder interface {
	Encode(v interface{}) *diam.Message
}

type AVP struct{}

func (a *AVP) Encode(v interface{}) *diam.Message {
	st := reflect.TypeOf(v)
	val := reflect.ValueOf(v)

	r := diam.NewRequest(diam.CreditControl, 4, nil)

	for i := 0; i < st.NumField(); i++ {
		field := st.Field(i)
		tag := field.Tag.Get("cbs")
		attr := strings.Split(tag, ",")
		code, err := strconv.Atoi(attr[0])
		if err != nil {
			log.Fatal(err)
			continue
		}
		switch attr[1] {
		case "OctetString":
			r.NewAVP(code, avp.Mbit, 0, datatype.OctetString(val.Field(i).String()))
		case "UTF8String":
			r.NewAVP(code, avp.Mbit, 0, datatype.UTF8String(val.Field(i).String()))
		case "DiameterIdentity":
			r.NewAVP(code, avp.Mbit, 0, datatype.DiameterIdentity(val.Field(i).String()))
		case "Unsigned32":
			r.NewAVP(code, avp.Mbit, 0, datatype.Unsigned32(val.Field(i).Int()))
		case "Integer32":
			r.NewAVP(code, avp.Mbit, 0, datatype.Integer32(val.Field(i).Int()))
		case "Time":
			r.NewAVP(code, avp.Mbit, 0, datatype.Time(val.Field(i).Interface().(time.Time)))
		case "GroupedAVP":
		}

	}

	return r
}

func Dcode(f reflect.StructField) []string {
	return strings.Split(f.Tag.Get("dcode"), ">")
}

func Dtype(f reflect.StructField) []string {
	return strings.Split(f.Tag.Get("dtype"), ",")
}

func Encode(v interface{}) []*diam.AVP {

	ret := []*diam.AVP{}

	var typ reflect.Type
	var val, inf reflect.Value

	fmt.Println(reflect.ValueOf(v).Kind())

	switch reflect.ValueOf(v).Kind() {
	case reflect.Ptr:
		inf = reflect.ValueOf(v).Elem()
		typ = inf.Elem().Type()
		val = inf.Elem()
	default:
		inf = reflect.ValueOf(v)
		typ = inf.Type()
		val = inf
	}

	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		dcode := Dcode(field)
		dtype := Dtype(field)

		code, err := strconv.Atoi(dcode[0])
		if err != nil {
			log.Fatal(err)
			continue
		}
		switch dtype[0] {
		case "OctetString":
		case "UTF8String":
			if val.Field(i).String() == "" {
				continue
			}
			ret = append(ret, diam.NewAVP(uint32(code), avp.Mbit, 0, datatype.UTF8String(val.Field(i).String())))
		case "DiameterIdentity":
		case "Unsigned32":
		case "Integer32":
		case "Time":
		case "GroupedAVP":
		}

	}

	return ret
}
