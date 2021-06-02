package client

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"

	spb "github.com/openconfig/gribi/v1/proto/service"
)

func TestHandleParams(t *testing.T) {
	tests := []struct {
		desc      string
		inOpts    []ClientOpt
		wantState *clientState
		wantErr   bool
	}{{
		desc:   "client with default parameters",
		inOpts: nil,
		wantState: &clientState{
			SessParams: &spb.SessionParameters{},
		},
	}, {
		desc: "ALL_PRIMARY client",
		inOpts: []ClientOpt{
			AllPrimaryClients(),
		},
		wantState: &clientState{
			SessParams: &spb.SessionParameters{
				Redundancy: spb.SessionParameters_ALL_PRIMARY,
			},
		},
	}, {
		desc: "SINGLE_PRIMARY client",
		inOpts: []ClientOpt{
			ElectedPrimaryClient(&spb.Uint128{High: 0, Low: 1}),
		},
		wantState: &clientState{
			SessParams: &spb.SessionParameters{
				Redundancy: spb.SessionParameters_SINGLE_PRIMARY,
			},
			ElectionID: &spb.Uint128{High: 0, Low: 1},
		},
	}, {
		desc: "SINGLE_PRIMARY and ALL_PRIMARY both included",
		inOpts: []ClientOpt{
			ElectedPrimaryClient(&spb.Uint128{High: 0, Low: 1}),
			AllPrimaryClients(),
		},
		wantErr: true,
	}}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			got, err := handleParams(tt.inOpts...)
			if (err != nil) != tt.wantErr {
				t.Fatalf("did not get expected error, wanted error? %v got error: %v", tt.wantErr, err)
			}
			if diff := cmp.Diff(tt.wantState, got, protocmp.Transform()); diff != "" {
				t.Fatalf("did not get expected state, diff(-want,+got):\n%s", diff)
			}
		})
	}
}
