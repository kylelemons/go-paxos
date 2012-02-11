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
			Error:    "missing sequence_id",
		},
		{
			Proposal: &rpc.Proposal{
				SequenceId: proto.Uint64(0),
			},
			Error:    "missing request_id",
		},
		{
			Proposal: &rpc.Proposal{
				SequenceId: proto.Uint64(0),
				RequestId: proto.Uint64(12),
			},
			Error:    "missing file_name",
		},
		{
			Proposal: &rpc.Proposal{
				SequenceId: proto.Uint64(0),
				RequestId: proto.Uint64(12),
				FileName: proto.String(""),
			},
			Error:    "<nil>",
		},
	}

	a := &Acceptor{
		accepted: map[string]*acceptedValue{},
	}

	for _, test := range tests {
		err := a.Propose(test.Proposal, &rpc.Promise{})
		if got, want := fmt.Sprint(err), test.Error; got != want {
			t.Errorf("Promise(%s) returned %q, want %q", test.Proposal, got, want)
		}
	}
}

func TestPropose(t *testing.T) {
	tests := []struct{
		Proposal *rpc.Proposal
		Promise  *rpc.Promise
	}{
		{
			Proposal: &rpc.Proposal{
				SequenceId: proto.Uint64(1),
				RequestId: proto.Uint64(12),
				FileName: proto.String("/p/test"),
			},
			Promise: &rpc.Promise{
				Ack: proto.Bool(true),
				IgnoreSeqBefore: proto.Uint64(1),
			},
		},
		{
			Proposal: &rpc.Proposal{
				SequenceId: proto.Uint64(42),
				RequestId: proto.Uint64(13),
				FileName: proto.String("/p/test/more"),
			},
			Promise: &rpc.Promise{
				Ack: proto.Bool(true),
				IgnoreSeqBefore: proto.Uint64(42),
			},
		},
		{
			Proposal: &rpc.Proposal{
				SequenceId: proto.Uint64(42),
				RequestId: proto.Uint64(14),
				FileName: proto.String("/p/test/dup"),
			},
			Promise: &rpc.Promise{
				Ack: proto.Bool(false),
				IgnoreSeqBefore: proto.Uint64(42),
			},
		},
		{
			Proposal: &rpc.Proposal{
				SequenceId: proto.Uint64(12),
				RequestId: proto.Uint64(14),
				FileName: proto.String("/p/oops"),
			},
			Promise: &rpc.Promise{
				Ack: proto.Bool(false),
				IgnoreSeqBefore: proto.Uint64(42),
			},
		},
		{
			Proposal: &rpc.Proposal{
				SequenceId: proto.Uint64(43),
				RequestId: proto.Uint64(15),
				FileName: proto.String("/p/test/next"),
			},
			Promise: &rpc.Promise{
				Ack: proto.Bool(true),
				IgnoreSeqBefore: proto.Uint64(43),
			},
		},
	}

	a := &Acceptor{
		accepted: map[string]*acceptedValue{},
	}

	for idx, test := range tests {
		prom := &rpc.Promise{}
		if err := a.Propose(test.Proposal, prom); err != nil {
			t.Errorf("%d. Promise(%q) returned %q", idx, test.Proposal, err)
			continue
		}
		if got, want := prom.String(), test.Promise.String(); got != want {
			t.Errorf("%d. Promise(%s) = %s, want %s", idx, test.Proposal, got, want)
		}
	}
}
