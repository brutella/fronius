package fronius

import (
	"fmt"
	"net/http"
	"net/http/httptest"
)

type Simulator interface {
	URL() string
	Stop()
}

type simulator struct {
	ts *httptest.Server
}

func NewSymoSimulator() Simulator {
	fn := func(rsp http.ResponseWriter, req *http.Request) {
		fmt.Fprint(rsp, `{
	"Head" : {
		"RequestArguments" : {
			"DataCollection" : "",
			"Scope" : "System"
		},
		"Status" : {
			"Code" : 0,
			"Reason" : "",
			"UserMessage" : ""
		},
		"Timestamp" : "2015-05-23T10:42:29+02:00"
	},
	"Body" : {
		"Data" : {
			"PAC" : {
				"Unit" : "W",
				"Values" : {
					"1" : 766
				}
			},
			"DAY_ENERGY" : {
				"Unit" : "Wh",
				"Values" : {
					"1" : 1622
				}
			},
			"YEAR_ENERGY" : {
				"Unit" : "Wh",
				"Values" : {
					"1" : 46146
				}
			},
			"TOTAL_ENERGY" : {
				"Unit" : "Wh",
				"Values" : {
					"1" : 46146
				}
			}
		}
	}
}`)
	}
	ts := httptest.NewServer(http.HandlerFunc(fn))
	return &simulator{ts}
}

func (s *simulator) URL() string {
	return s.ts.URL
}

func (s *simulator) Stop() {
	s.ts.Close()
}
