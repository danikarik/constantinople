package metropolis

import (
	"context"
	"errors"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

const (
	dialTime time.Duration = 3 * time.Second
	timeout  time.Duration = 60 * time.Second
)

// ErrServerAddr raises when server address is not specified.
var ErrServerAddr = errors.New("proto: gRPC server address is empty")

// New return new client instance.
func New(ctx context.Context, addr string) (*grpc.ClientConn, error) {
	if addr == "" {
		return nil, ErrServerAddr
	}
	conn, err := connect(ctx, addr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func connect(ctx context.Context, url string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(ctx, dialTime)
	defer cancel()
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithKeepaliveParams(keepalive.ClientParameters{
		Time:    timeout,
		Timeout: timeout,
	}))
	cc, err := dial(ctx, url, opts...)
	if err != nil {
		return nil, err
	}
	return cc, nil
}

func dial(ctx context.Context, address string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	result := make(chan interface{}, 1)

	writeResult := func(res interface{}) {
		select {
		case result <- res:
		default:
		}
	}

	dialer := func(address string, timeout time.Duration) (net.Conn, error) {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		conn, err := (&net.Dialer{Cancel: ctx.Done()}).Dial("tcp", address)
		if err != nil {
			writeResult(err)
			return nil, err
		}
		return conn, nil
	}

	go func() {
		opts = append(opts,
			grpc.WithBlock(),
			grpc.FailOnNonTempDialError(true),
			grpc.WithDialer(dialer),
			grpc.WithInsecure(),
		)
		conn, err := grpc.DialContext(ctx, address, opts...)
		var res interface{}
		if err != nil {
			res = err
		} else {
			res = conn
		}
		writeResult(res)
	}()

	select {
	case res := <-result:
		if conn, ok := res.(*grpc.ClientConn); ok {
			return conn, nil
		}
		return nil, res.(error)
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
