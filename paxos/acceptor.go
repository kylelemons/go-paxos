package paxos

import (
	"fmt"
	"sync"

	"code.google.com/p/goprotobuf/proto"
	"github.com/kylelemons/go-paxos/paxos/rpc"
)

var _ rpc.Acceptor = &Acceptor{}

type acceptedValue struct {
	sequenceId uint64
	value      []byte
}

type Acceptor struct {
	lock sync.Mutex

	accepted map[string]*acceptedValue

	ignoreSeqBefore uint64
	acceptRequestId uint64
}

func (a *Acceptor) Propose(prop *rpc.Proposal, prom *rpc.Promise) error {
	if prop.SequenceId == nil {
		return fmt.Errorf("missing sequence_id")
	}
	if prop.RequestId == nil {
		return fmt.Errorf("missing request_id")
	}
	if prop.FileName == nil {
		return fmt.Errorf("missing file_name")
	}

	a.lock.Lock()
	defer a.lock.Unlock()

	if *prop.SequenceId <= a.ignoreSeqBefore {
		prom.Ack = proto.Bool(false)
		prom.IgnoreSeqBefore = proto.Uint64(a.ignoreSeqBefore)
		return nil
	}

	if val, ok := a.accepted[*prop.FileName]; ok {
		prom.PastSequenceId = proto.Uint64(val.sequenceId)
		prom.PastData = append([]byte{}, val.value...)
	}

	a.ignoreSeqBefore = *prop.SequenceId
	a.acceptRequestId = a.ignoreSeqBefore
	prom.IgnoreSeqBefore = prop.SequenceId
	prom.Ack = proto.Bool(true)
	return nil
}

func (a *Acceptor) Accept(prop *rpc.Proposal, ack *rpc.Ack) error {
	return nil
}
