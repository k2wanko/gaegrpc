//go:generate protoc --go_out=plugins=grpc:. test/test.proto

package gaegrpc

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/aetest"
	"google.golang.org/grpc"

	testpb "github.com/k2wanko/gaegrpc/test"
)

var aeInst aetest.Instance

func TestMain(m *testing.M) {
	var err error
	if aeInst, err = aetest.NewInstance(nil); err != nil {
		panic(fmt.Sprintf("aetestInstance: %v", err))
	}

	code := m.Run()
	aeInst.Close()

	os.Exit(code)
}

func newTestContext() context.Context {
	r, _ := aeInst.NewRequest("GET", "/", nil)
	return appengine.NewContext(r)
}

func newTestRequest(t *testing.T, method, url string, body io.Reader) *http.Request {
	req, err := aeInst.NewRequest(method, url, body)
	if err != nil {
		t.Fatalf("newTestRequest(%q, %q): %v", method, url, err)
	}
	return req
}

func newTestGrpcRequest(t *testing.T, method string, headers http.Header, messages ...proto.Message) *http.Request {
	writer := new(bytes.Buffer)
	for _, msg := range serializeProtoMessages(messages) {
		grpcPreamble := []byte{0, 0, 0, 0, 0}
		binary.BigEndian.PutUint32(grpcPreamble[1:], uint32(len(msg)))
		writer.Write(grpcPreamble)
		writer.Write(msg)
	}
	req := newTestRequest(t, "POST", method, writer)
	if headers != nil {
		req.Header = headers
	}
	req.ProtoMajor = 2
	req.ProtoMinor = 0
	req.Header.Set("content-type", "application/grpc+proto")
	return req
}

func serializeProtoMessages(messages []proto.Message) [][]byte {
	out := [][]byte{}
	for _, m := range messages {
		b, _ := proto.Marshal(m)
		out = append(out, b)
	}
	return out
}

type testService struct{}

func (s *testService) Ping(ctx context.Context, req *testpb.PingRequest) (res *testpb.Pong, err error) {
	res = &testpb.Pong{}
	return
}

func TestPing(t *testing.T) {
	sv := NewServer(grpc.UnaryInterceptor(func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		id := appengine.AppID(ctx)
		t.Logf("AppID: %v", id)
		if want := "testapp"; id != want {
			t.Errorf("id = %v; want = %v", id, want)
		}
		resp, err = handler(ctx, req)
		return
	}))
	testpb.RegisterTestServer(sv, &testService{})
	h := NewWrapHandler(sv)

	r, w := newTestGrpcRequest(t, "/com.github.k2wanko.gaegrpc.test.Test/Ping", nil, &testpb.PingRequest{}), httptest.NewRecorder()
	h.ServeHTTP(w, r)
}
