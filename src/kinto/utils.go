package kinto

import (
	"crypto/x509/pkix"
	"net/http"
	"reflect"
	"regexp"
	"sort"
)

func newRequest() *http.Request {
	return &http.Request{
		Method: "GET",
		URL:    kinto,
		Header: map[string][]string{
			"X-Automated-Tool": {"chris"},
		},
	}
}

func comp(a, b pkix.RDNSequence) bool {
	aSet := make(pkix.RelativeDistinguishedNameSET, 0)
	bSet := make(pkix.RelativeDistinguishedNameSET, 0)
	for _, i := range a {
		for _, j := range i {
			aSet = append(aSet, j)
		}
	}
	for _, i := range b {
		for _, j := range i {
			bSet = append(bSet, j)
		}
	}
	return compset(aSet, bSet)
	//needed := len(a)
	//matched := 0
	//for _, i := range a {
	//	for _, j := range b {
	//		if compset(i, j) {
	//			matched += 1
	//			goto CONTINUE
	//		}
	//	}
	//	return false
	//CONTINUE:
	//}
	//return matched == needed
}

func compset(a, b pkix.RelativeDistinguishedNameSET) bool {
	if len(a) != len(b) {
		return false
	}
	sort.Slice(a, func(i, j int) bool {
		if len(a[i].Type) < len(a[j].Type) {
			return true
		}
		for index, value := range a[i].Type {
			if value < a[j].Type[index] {
				return true
			}
		}
		return false
	})
	sort.Slice(b, func(i, j int) bool {
		if len(b[i].Type) < len(b[j].Type) {
			return true
		}
		for index, value := range b[i].Type {
			if value < b[j].Type[index] {
				return true
			}
		}
		return false
	})
	return reflect.DeepEqual(a, b)
}

var r = regexp.MustCompilePOSIX(`^0{2}*`)

func StripPadding(s string) string {
	return r.ReplaceAllString(s, "")
}
