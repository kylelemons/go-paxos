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
		File    string
		Propose uint64
		Ack     bool
	}{
		{
			Desc: "test",
			File: "/test/file",
			Propose: 1,
			Ack: true,
		},
		{
			Desc: "test more",
			File: "/test/file",
			Propose: 42,
			Ack: true,
		},
		{
			Desc: "test dup",
			File: "/test/file",
			Propose: 42,
			Ack: false,
		},
		{
			Desc: "oops",
			File: "/test/file",
			Propose: 12,
			Ack: false,
		},
		{
			Desc: "above break",
			File: "/test/file",
			Propose: 42 + breakpoint + 5,
			Ack: false,
		},
		{
			Desc: "below break",
			File: "/test/file",
			Propose: 42 + breakpoint - 5,
			Ack: true,
		},
		{
			Desc: "after wrap",
			File: "/test/file",
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

	a = &Acceptor{filePromises:map[string]uint64{}}
	last, promised = 0, 0

	for idx, test := range tests {
		prop := &rpc.Proposal{
			LeaderId: proto.Uint64(0),
			Origin: proto.Uint64(uint64(idx)),
			Type:   rpc.NewProposal_Type(rpc.Proposal_FILE),
			RequestId: proto.Uint64(uint64(idx)),
			ProposalId: proto.Uint64(uint64(idx+1)),
			File: []*rpc.File{{
				Name: &test.File,
				Revision: &test.Propose,
			}},
		}
		prom := &rpc.Promise{}
		if err := a.Propose(prop, prom); err != nil {
			t.Errorf("%d. Propose(%q, %d) returned %q", idx, test.File, test.Propose, err)
			continue
		}

		if test.Ack {
			promised = test.Propose
		}

		want := &rpc.Promise{
			File: []*rpc.File{{
				Name: &test.File,
				Revision: &test.Propose,
			}},
		}
		want.Ack = &test.Ack
		if idx > 0 {
			want.LastPromised = proto.Uint64(uint64(idx))
			want.IPromise = proto.Uint64(uint64(idx+1))
		}
		if last > 0 {
			want.File[0].LastPromised = &last
		}
		if test.Ack {
			want.File[0].IPromise = &promised
		}
		want.IPromise = proto.Uint64(uint64(idx+1))

		if got, want := prom.String(), want.String(); got != want {
			t.Errorf("%d. Propose(%q, %d) = %s, want %s", idx, test.File, test.Propose, got, want)
		}
		last = promised
	}
}

func TestElection(t *testing.T) {
	tests := []struct{
		LeaderId uint64
		Ack      bool
	}{
		{0, false},
		{1, true},
		//{10, true},
		//{5, false},
		//{10, false},
		//{11, true},
	}

	a := &Acceptor{}

	var active uint64
	for idx, test := range tests {
		prop := &rpc.Proposal{
			Origin: proto.Uint64(1),
			Type:   rpc.NewProposal_Type(rpc.Proposal_ELECTION),
			LeaderId: &active,
			RequestId: proto.Uint64(uint64(idx)),
			ProposalId: proto.Uint64(uint64(idx+1)),
		}
		if test.LeaderId > 0 {
			prop.NewLeaderId = proto.Uint64(test.LeaderId)
		}

		prom := &rpc.Promise{}
		if err := a.Propose(prop, prom); err != nil {
			t.Errorf("%d. Propose(%d) returned %s", idx, test.LeaderId, err)
			continue
		}

		if got, want := prom.Ack != nil && *prom.Ack, test.Ack; got != want {
			t.Errorf("%d. Propose(%d) promise.ack = %v, want %v", idx, test.LeaderId, got, want)
		}
	}
}
