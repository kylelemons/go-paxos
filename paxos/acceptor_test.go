package paxos

import (
	"fmt"
	"testing"

	"code.google.com/p/goprotobuf/proto"
	"github.com/kylelemons/go-paxos/paxos/rpc"
)

func TestProposeMissingFields(t *testing.T) {
	tests := []struct{
		Proposal *rpc.Proposal
		Error    string
	}{
		{
			Proposal: &rpc.Proposal{},
			Error:    "proposal: missing origin",
		},
		{
			Proposal: &rpc.Proposal{
				Origin: proto.Uint64(0),
			},
			Error:    "proposal: missing proposal_id",
		},
		{
			Proposal: &rpc.Proposal{
				Origin: proto.Uint64(0),
				ProposalId: proto.Uint64(0),
			},
			Error:    "proposal: proposal_id must be greater than zero",
		},
		{
			Proposal: &rpc.Proposal{
				Origin: proto.Uint64(0),
				ProposalId: proto.Uint64(1),
			},
			Error:    "proposal: missing request_id",
		},
		{
			Proposal: &rpc.Proposal{
				Origin: proto.Uint64(0),
				ProposalId: proto.Uint64(1),
				RequestId: proto.Uint64(0),
			},
			Error:    "proposal: missing leader_id",
		},
		{
			Proposal: &rpc.Proposal{
				Origin: proto.Uint64(0),
				ProposalId: proto.Uint64(1),
				RequestId: proto.Uint64(0),
				LeaderId: proto.Uint64(0),
			},
			Error:    "<nil>",
		},
	}

	a := &Acceptor{}

	for _, test := range tests {
		err := a.Propose(test.Proposal, &rpc.Promise{})
		if got, want := fmt.Sprint(err), test.Error; got != want {
			t.Errorf("Promise(%s) returned %q, want %q", test.Proposal, got, want)
		}
	}
}

func TestProposeWraparound(t *testing.T) {
	tests := []struct{
		Desc    string
		Propose uint64
		Ack     bool
	}{
		{
			Desc: "/p/test",
			Propose: 1,
			Ack: true,
		},
		{
			Desc: "/p/test/more",
			Propose: 42,
			Ack: true,
		},
		{
			Desc: "/p/test/dup",
			Propose: 42,
			Ack: false,
		},
		{
			Desc: "/p/oops",
			Propose: 12,
			Ack: false,
		},
		{
			Desc: "/p/above/break",
			Propose: 42 + breakpoint + 5,
			Ack: false,
		},
		{
			Desc: "/p/below/break",
			Propose: 42 + breakpoint - 5,
			Ack: true,
		},
		{
			Desc: "/p/after/wrap",
			Propose: 32,
			Ack: true,
		},
	}

	a := &Acceptor{}
	var last, promised uint64

	for idx, test := range tests {
		prop := &rpc.Proposal{
			LeaderId: proto.Uint64(0),
			Origin: proto.Uint64(uint64(idx)),
			Type:   rpc.NewProposal_Type(rpc.Proposal_FILE),
			RequestId: proto.Uint64(uint64(idx)),
			ProposalId: &test.Propose,
		}
		prom := &rpc.Promise{}
		if err := a.Propose(prop, prom); err != nil {
			t.Errorf("%d. Propose(%d) returned %q", idx, test.Propose, err)
			continue
		}

		if test.Ack {
			promised = test.Propose
		}

		want := &rpc.Promise{}
		want.Ack = &test.Ack
		if last > 0 {
			want.LastPromised = &last
		}
		if test.Ack {
			want.IPromise = &promised
		}

		if got, want := prom.String(), want.String(); got != want {
			t.Errorf("%d. Propose(%s) = %s, want %s", idx, test.Propose, got, want)
		}
		last = promised
	}
}
