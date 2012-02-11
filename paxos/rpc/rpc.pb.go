// Code generated by protoc-gen-go from "rpc/rpc.proto"
// DO NOT EDIT!

package rpc

import proto "code.google.com/p/goprotobuf/proto"
import "math"

import "net"
import "net/rpc"
import "github.com/kylelemons/go-rpcgen/codec"
import "net/url"
import "net/http"
import "github.com/kylelemons/go-rpcgen/webrpc"

// Reference proto and math imports to suppress error if they are not otherwise used.
var _ = proto.GetString
var _ = math.Inf

type Paxos int32

const (
	Paxos_MULTI Paxos = 0
	Paxos_FAST  Paxos = 1
)

var Paxos_name = map[int32]string{
	0: "MULTI",
	1: "FAST",
}
var Paxos_value = map[string]int32{
	"MULTI": 0,
	"FAST":  1,
}

func NewPaxos(x Paxos) *Paxos {
	e := Paxos(x)
	return &e
}
func (x Paxos) String() string {
	return proto.EnumName(Paxos_name, int32(x))
}

type Proposal_Type int32

const (
	Proposal_FILE   Proposal_Type = 0
	Proposal_MASTER Proposal_Type = 1
)

var Proposal_Type_name = map[int32]string{
	0: "FILE",
	1: "MASTER",
}
var Proposal_Type_value = map[string]int32{
	"FILE":   0,
	"MASTER": 1,
}

func NewProposal_Type(x Proposal_Type) *Proposal_Type {
	e := Proposal_Type(x)
	return &e
}
func (x Proposal_Type) String() string {
	return proto.EnumName(Proposal_Type_name, int32(x))
}

type Empty struct {
	XXX_unrecognized []byte `json:",omitempty"`
}

func (this *Empty) Reset()         { *this = Empty{} }
func (this *Empty) String() string { return proto.CompactTextString(this) }

type Proposal struct {
	Type             *Proposal_Type `protobuf:"varint,1,opt,name=type,enum=rpc.Proposal_Type" json:"type,omitempty"`
	SequenceId       *uint64        `protobuf:"varint,2,opt,name=sequence_id" json:"sequence_id,omitempty"`
	MasterId         *uint64        `protobuf:"varint,3,opt,name=master_id" json:"master_id,omitempty"`
	RequestId        *uint64        `protobuf:"varint,4,opt,name=request_id" json:"request_id,omitempty"`
	FileName         *string        `protobuf:"bytes,6,opt,name=file_name" json:"file_name,omitempty"`
	FileRev          *string        `protobuf:"bytes,7,opt,name=file_rev" json:"file_rev,omitempty"`
	FilePaxos        *Paxos         `protobuf:"varint,8,opt,name=file_paxos,enum=rpc.Paxos" json:"file_paxos,omitempty"`
	ValueSha1        *string        `protobuf:"bytes,13,opt,name=value_sha1" json:"value_sha1,omitempty"`
	Data             []byte         `protobuf:"bytes,14,opt,name=data" json:"data,omitempty"`
	XXX_unrecognized []byte         `json:",omitempty"`
}

func (this *Proposal) Reset()         { *this = Proposal{} }
func (this *Proposal) String() string { return proto.CompactTextString(this) }

type Promise struct {
	Ack              *bool   `protobuf:"varint,1,opt,name=ack" json:"ack,omitempty"`
	IgnoreSeqBefore  *uint64 `protobuf:"varint,2,opt,name=ignore_seq_before" json:"ignore_seq_before,omitempty"`
	PastSequenceId   *uint64 `protobuf:"varint,3,opt,name=past_sequence_id" json:"past_sequence_id,omitempty"`
	IgnoreRevBefore  *uint64 `protobuf:"varint,6,opt,name=ignore_rev_before" json:"ignore_rev_before,omitempty"`
	PastRev          *uint64 `protobuf:"varint,7,opt,name=past_rev" json:"past_rev,omitempty"`
	PastValueSha1    *string `protobuf:"bytes,10,opt,name=past_value_sha1" json:"past_value_sha1,omitempty"`
	PastData         []byte  `protobuf:"bytes,11,opt,name=past_data" json:"past_data,omitempty"`
	XXX_unrecognized []byte  `json:",omitempty"`
}

func (this *Promise) Reset()         { *this = Promise{} }
func (this *Promise) String() string { return proto.CompactTextString(this) }

type Ack struct {
	Ack              *bool  `protobuf:"varint,1,opt,name=ack" json:"ack,omitempty"`
	XXX_unrecognized []byte `json:",omitempty"`
}

func (this *Ack) Reset()         { *this = Ack{} }
func (this *Ack) String() string { return proto.CompactTextString(this) }

func init() {
	proto.RegisterEnum("rpc.Paxos", Paxos_name, Paxos_value)
	proto.RegisterEnum("rpc.Proposal_Type", Proposal_Type_name, Proposal_Type_value)
}

// Acceptor is an interface satisfied by the generated client and
// which must be implemented by the object wrapped by the server.
type Acceptor interface {
	Propose(in *Proposal, out *Promise) error
	Accept(in *Proposal, out *Ack) error
}

// internal wrapper for type-safe RPC calling
type rpcAcceptorClient struct {
	*rpc.Client
}

func (this rpcAcceptorClient) Propose(in *Proposal, out *Promise) error {
	return this.Call("Acceptor.Propose", in, out)
}
func (this rpcAcceptorClient) Accept(in *Proposal, out *Ack) error {
	return this.Call("Acceptor.Accept", in, out)
}

// NewAcceptorClient returns an *rpc.Client wrapper for calling the methods of
// Acceptor remotely.
func NewAcceptorClient(conn net.Conn) Acceptor {
	return rpcAcceptorClient{rpc.NewClientWithCodec(codec.NewClientCodec(conn))}
}

// ServeAcceptor serves the given Acceptor backend implementation on conn.
func ServeAcceptor(conn net.Conn, backend Acceptor) error {
	srv := rpc.NewServer()
	if err := srv.RegisterName("Acceptor", backend); err != nil {
		return err
	}
	srv.ServeCodec(codec.NewServerCodec(conn))
	return nil
}

// DialAcceptor returns a Acceptor for calling the Acceptor servince at addr (TCP).
func DialAcceptor(addr string) (Acceptor, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewAcceptorClient(conn), nil
}

// ListenAndServeAcceptor serves the given Acceptor backend implementation
// on all connections accepted as a result of listening on addr (TCP).
func ListenAndServeAcceptor(addr string, backend Acceptor) error {
	clients, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv := rpc.NewServer()
	if err := srv.RegisterName("Acceptor", backend); err != nil {
		return err
	}
	for {
		conn, err := clients.Accept()
		if err != nil {
			return err
		}
		go srv.ServeCodec(codec.NewServerCodec(conn))
	}
	panic("unreachable")
}

// Learner is an interface satisfied by the generated client and
// which must be implemented by the object wrapped by the server.
type Learner interface {
	Learn(in *Proposal, out *Ack) error
}

// internal wrapper for type-safe RPC calling
type rpcLearnerClient struct {
	*rpc.Client
}

func (this rpcLearnerClient) Learn(in *Proposal, out *Ack) error {
	return this.Call("Learner.Learn", in, out)
}

// NewLearnerClient returns an *rpc.Client wrapper for calling the methods of
// Learner remotely.
func NewLearnerClient(conn net.Conn) Learner {
	return rpcLearnerClient{rpc.NewClientWithCodec(codec.NewClientCodec(conn))}
}

// ServeLearner serves the given Learner backend implementation on conn.
func ServeLearner(conn net.Conn, backend Learner) error {
	srv := rpc.NewServer()
	if err := srv.RegisterName("Learner", backend); err != nil {
		return err
	}
	srv.ServeCodec(codec.NewServerCodec(conn))
	return nil
}

// DialLearner returns a Learner for calling the Learner servince at addr (TCP).
func DialLearner(addr string) (Learner, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewLearnerClient(conn), nil
}

// ListenAndServeLearner serves the given Learner backend implementation
// on all connections accepted as a result of listening on addr (TCP).
func ListenAndServeLearner(addr string, backend Learner) error {
	clients, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	srv := rpc.NewServer()
	if err := srv.RegisterName("Learner", backend); err != nil {
		return err
	}
	for {
		conn, err := clients.Accept()
		if err != nil {
			return err
		}
		go srv.ServeCodec(codec.NewServerCodec(conn))
	}
	panic("unreachable")
}

// AcceptorWeb is the web-based RPC version of the interface which
// must be implemented by the object wrapped by the webrpc server.
type AcceptorWeb interface {
	Propose(r *http.Request, in *Proposal, out *Promise) error
	Accept(r *http.Request, in *Proposal, out *Ack) error
}

// internal wrapper for type-safe webrpc calling
type rpcAcceptorWebClient struct {
	remote   *url.URL
	protocol webrpc.Protocol
}

func (this rpcAcceptorWebClient) Propose(in *Proposal, out *Promise) error {
	return webrpc.Post(this.protocol, this.remote, "/Acceptor/Propose", in, out)
}

func (this rpcAcceptorWebClient) Accept(in *Proposal, out *Ack) error {
	return webrpc.Post(this.protocol, this.remote, "/Acceptor/Accept", in, out)
}

// Register a AcceptorWeb implementation with the given webrpc ServeMux.
// If mux is nil, the default webrpc.ServeMux is used.
func RegisterAcceptorWeb(this AcceptorWeb, mux webrpc.ServeMux) error {
	if mux == nil {
		mux = webrpc.DefaultServeMux
	}
	if err := mux.Handle("/Acceptor/Propose", func(c *webrpc.Call) error {
		in, out := new(Proposal), new(Promise)
		if err := c.ReadRequest(in); err != nil {
			return err
		}
		if err := this.Propose(c.Request, in, out); err != nil {
			return err
		}
		return c.WriteResponse(out)
	}); err != nil {
		return err
	}
	if err := mux.Handle("/Acceptor/Accept", func(c *webrpc.Call) error {
		in, out := new(Proposal), new(Ack)
		if err := c.ReadRequest(in); err != nil {
			return err
		}
		if err := this.Accept(c.Request, in, out); err != nil {
			return err
		}
		return c.WriteResponse(out)
	}); err != nil {
		return err
	}
	return nil
}

// NewAcceptorWebClient returns a webrpc wrapper for calling the methods of Acceptor
// remotely via the web.  The remote URL is the base URL of the webrpc server.
func NewAcceptorWebClient(pro webrpc.Protocol, remote *url.URL) Acceptor {
	return rpcAcceptorWebClient{remote, pro}
}

// LearnerWeb is the web-based RPC version of the interface which
// must be implemented by the object wrapped by the webrpc server.
type LearnerWeb interface {
	Learn(r *http.Request, in *Proposal, out *Ack) error
}

// internal wrapper for type-safe webrpc calling
type rpcLearnerWebClient struct {
	remote   *url.URL
	protocol webrpc.Protocol
}

func (this rpcLearnerWebClient) Learn(in *Proposal, out *Ack) error {
	return webrpc.Post(this.protocol, this.remote, "/Learner/Learn", in, out)
}

// Register a LearnerWeb implementation with the given webrpc ServeMux.
// If mux is nil, the default webrpc.ServeMux is used.
func RegisterLearnerWeb(this LearnerWeb, mux webrpc.ServeMux) error {
	if mux == nil {
		mux = webrpc.DefaultServeMux
	}
	if err := mux.Handle("/Learner/Learn", func(c *webrpc.Call) error {
		in, out := new(Proposal), new(Ack)
		if err := c.ReadRequest(in); err != nil {
			return err
		}
		if err := this.Learn(c.Request, in, out); err != nil {
			return err
		}
		return c.WriteResponse(out)
	}); err != nil {
		return err
	}
	return nil
}

// NewLearnerWebClient returns a webrpc wrapper for calling the methods of Learner
// remotely via the web.  The remote URL is the base URL of the webrpc server.
func NewLearnerWebClient(pro webrpc.Protocol, remote *url.URL) Learner {
	return rpcLearnerWebClient{remote, pro}
}
