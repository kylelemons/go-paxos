package paxos

import (
	"fmt"
	"sync"

	"code.google.com/p/goprotobuf/proto"
	"github.com/kylelemons/go-paxos/paxos/rpc"
)

var _ rpc.Acceptor = &Acceptor{}

type savedFile struct {
	sequenceId uint64
	value      []byte
}

const breakpoint = 1<<63
func breaksPromise(proposed, promise uint64) bool {
	return promise - proposed < breakpoint
}

func checkProposal(prop *rpc.Proposal) error {
	if prop.Origin == nil {
		return fmt.Errorf("proposal: missing origin")
	}
	if prop.ProposalId == nil {
		return fmt.Errorf("proposal: missing proposal_id")
	} else if *prop.ProposalId <= 0 {
		return fmt.Errorf("proposal: proposal_id must be greater than zero")
	}
	if prop.RequestId == nil {
		return fmt.Errorf("proposal: missing request_id")
	}
	if prop.LeaderId == nil {
		return fmt.Errorf("proposal: missing leader_id")
	}
	return nil
}

func checkFile(file *rpc.File) error {
	if file.Name == nil {
		return fmt.Errorf("file: missing name")
	}
	if file.Revision == nil {
		return fmt.Errorf("file: missing revision")
	}
	return nil
}

type Acceptor struct {
	lock sync.Mutex

	promised uint64
}

func (a *Acceptor) Propose(prop *rpc.Proposal, prom *rpc.Promise) error {
	if err := checkProposal(prop); err != nil {
		return err
	}
	for idx, file := range prop.File {
		if err := checkFile(file); err != nil {
			return fmt.Errorf("proposal.file[%d]: %s", idx, err)
		}
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	switch {
	case prop.Type == nil:
	case *prop.Type == rpc.Proposal_FILE:
	case *prop.Type == rpc.Proposal_ELECTION:
	}

	// Send the last promise we made if we have one
	if a.promised > 0 {
		prom.LastPromised = proto.Uint64(a.promised)
	}

	if breaksPromise(*prop.ProposalId, a.promised) {
		// Nack if the proposal would break our promise
		prom.Ack = proto.Bool(false)
	} else {
		// Make the promise if we can
		a.promised = *prop.ProposalId
		prom.IPromise = prop.ProposalId
	}

	// If nothing failed, we can ack
	if prom.Ack == nil {
		prom.Ack = proto.Bool(true)
	}
	return nil
}

func (a *Acceptor) Accept(prop *rpc.Proposal, ack *rpc.Ack) error {
	return nil
}
