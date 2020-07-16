// Code generated by protoc-gen-grpc-gateway. DO NOT EDIT.
// source: antibiogram.proto

/*
Package antibug_antibiogram is a reverse proxy.

It translates gRPC into RESTful JSON APIs.
*/
package antibug_antibiogram

import (
	"context"
	"io"
	"net/http"

	"github.com/golang/protobuf/descriptor"
	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

// Suppress "imported and not used" errors
var _ codes.Code
var _ io.Reader
var _ status.Status
var _ = runtime.String
var _ = utilities.NewDoubleArray
var _ = descriptor.ForMessage

var (
	filter_AntibiogramAPI_GenPathogensAntibiogram_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_AntibiogramAPI_GenPathogensAntibiogram_0(ctx context.Context, marshaler runtime.Marshaler, client AntibiogramAPIClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq Filter
	var metadata runtime.ServerMetadata

	if err := req.ParseForm(); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	if err := runtime.PopulateQueryParameters(&protoReq, req.Form, filter_AntibiogramAPI_GenPathogensAntibiogram_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.GenPathogensAntibiogram(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_AntibiogramAPI_GenPathogensAntibiogram_0(ctx context.Context, marshaler runtime.Marshaler, server AntibiogramAPIServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq Filter
	var metadata runtime.ServerMetadata

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_AntibiogramAPI_GenPathogensAntibiogram_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.GenPathogensAntibiogram(ctx, &protoReq)
	return msg, metadata, err

}

var (
	filter_AntibiogramAPI_GenPathogenAntibiogram_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_AntibiogramAPI_GenPathogenAntibiogram_0(ctx context.Context, marshaler runtime.Marshaler, client AntibiogramAPIClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq Filter
	var metadata runtime.ServerMetadata

	if err := req.ParseForm(); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	if err := runtime.PopulateQueryParameters(&protoReq, req.Form, filter_AntibiogramAPI_GenPathogenAntibiogram_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.GenPathogenAntibiogram(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_AntibiogramAPI_GenPathogenAntibiogram_0(ctx context.Context, marshaler runtime.Marshaler, server AntibiogramAPIServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq Filter
	var metadata runtime.ServerMetadata

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_AntibiogramAPI_GenPathogenAntibiogram_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.GenPathogenAntibiogram(ctx, &protoReq)
	return msg, metadata, err

}

var (
	filter_AntibiogramAPI_GenAntimicrobialsAntibiogram_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_AntibiogramAPI_GenAntimicrobialsAntibiogram_0(ctx context.Context, marshaler runtime.Marshaler, client AntibiogramAPIClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq Filter
	var metadata runtime.ServerMetadata

	if err := req.ParseForm(); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	if err := runtime.PopulateQueryParameters(&protoReq, req.Form, filter_AntibiogramAPI_GenAntimicrobialsAntibiogram_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.GenAntimicrobialsAntibiogram(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_AntibiogramAPI_GenAntimicrobialsAntibiogram_0(ctx context.Context, marshaler runtime.Marshaler, server AntibiogramAPIServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq Filter
	var metadata runtime.ServerMetadata

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_AntibiogramAPI_GenAntimicrobialsAntibiogram_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.GenAntimicrobialsAntibiogram(ctx, &protoReq)
	return msg, metadata, err

}

var (
	filter_AntibiogramAPI_GenAntimicrobialAntibiogram_0 = &utilities.DoubleArray{Encoding: map[string]int{}, Base: []int(nil), Check: []int(nil)}
)

func request_AntibiogramAPI_GenAntimicrobialAntibiogram_0(ctx context.Context, marshaler runtime.Marshaler, client AntibiogramAPIClient, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq Filter
	var metadata runtime.ServerMetadata

	if err := req.ParseForm(); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}
	if err := runtime.PopulateQueryParameters(&protoReq, req.Form, filter_AntibiogramAPI_GenAntimicrobialAntibiogram_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := client.GenAntimicrobialAntibiogram(ctx, &protoReq, grpc.Header(&metadata.HeaderMD), grpc.Trailer(&metadata.TrailerMD))
	return msg, metadata, err

}

func local_request_AntibiogramAPI_GenAntimicrobialAntibiogram_0(ctx context.Context, marshaler runtime.Marshaler, server AntibiogramAPIServer, req *http.Request, pathParams map[string]string) (proto.Message, runtime.ServerMetadata, error) {
	var protoReq Filter
	var metadata runtime.ServerMetadata

	if err := runtime.PopulateQueryParameters(&protoReq, req.URL.Query(), filter_AntibiogramAPI_GenAntimicrobialAntibiogram_0); err != nil {
		return nil, metadata, status.Errorf(codes.InvalidArgument, "%v", err)
	}

	msg, err := server.GenAntimicrobialAntibiogram(ctx, &protoReq)
	return msg, metadata, err

}

// RegisterAntibiogramAPIHandlerServer registers the http handlers for service AntibiogramAPI to "mux".
// UnaryRPC     :call AntibiogramAPIServer directly.
// StreamingRPC :currently unsupported pending https://github.com/grpc/grpc-go/issues/906.
func RegisterAntibiogramAPIHandlerServer(ctx context.Context, mux *runtime.ServeMux, server AntibiogramAPIServer) error {

	mux.Handle("GET", pattern_AntibiogramAPI_GenPathogensAntibiogram_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_AntibiogramAPI_GenPathogensAntibiogram_0(rctx, inboundMarshaler, server, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_AntibiogramAPI_GenPathogensAntibiogram_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_AntibiogramAPI_GenPathogenAntibiogram_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_AntibiogramAPI_GenPathogenAntibiogram_0(rctx, inboundMarshaler, server, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_AntibiogramAPI_GenPathogenAntibiogram_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_AntibiogramAPI_GenAntimicrobialsAntibiogram_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_AntibiogramAPI_GenAntimicrobialsAntibiogram_0(rctx, inboundMarshaler, server, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_AntibiogramAPI_GenAntimicrobialsAntibiogram_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_AntibiogramAPI_GenAntimicrobialAntibiogram_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateIncomingContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := local_request_AntibiogramAPI_GenAntimicrobialAntibiogram_0(rctx, inboundMarshaler, server, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_AntibiogramAPI_GenAntimicrobialAntibiogram_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

// RegisterAntibiogramAPIHandlerFromEndpoint is same as RegisterAntibiogramAPIHandler but
// automatically dials to "endpoint" and closes the connection when "ctx" gets done.
func RegisterAntibiogramAPIHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return RegisterAntibiogramAPIHandler(ctx, mux, conn)
}

// RegisterAntibiogramAPIHandler registers the http handlers for service AntibiogramAPI to "mux".
// The handlers forward requests to the grpc endpoint over "conn".
func RegisterAntibiogramAPIHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return RegisterAntibiogramAPIHandlerClient(ctx, mux, NewAntibiogramAPIClient(conn))
}

// RegisterAntibiogramAPIHandlerClient registers the http handlers for service AntibiogramAPI
// to "mux". The handlers forward requests to the grpc endpoint over the given implementation of "AntibiogramAPIClient".
// Note: the gRPC framework executes interceptors within the gRPC handler. If the passed in "AntibiogramAPIClient"
// doesn't go through the normal gRPC flow (creating a gRPC client etc.) then it will be up to the passed in
// "AntibiogramAPIClient" to call the correct interceptors.
func RegisterAntibiogramAPIHandlerClient(ctx context.Context, mux *runtime.ServeMux, client AntibiogramAPIClient) error {

	mux.Handle("GET", pattern_AntibiogramAPI_GenPathogensAntibiogram_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_AntibiogramAPI_GenPathogensAntibiogram_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_AntibiogramAPI_GenPathogensAntibiogram_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_AntibiogramAPI_GenPathogenAntibiogram_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_AntibiogramAPI_GenPathogenAntibiogram_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_AntibiogramAPI_GenPathogenAntibiogram_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_AntibiogramAPI_GenAntimicrobialsAntibiogram_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_AntibiogramAPI_GenAntimicrobialsAntibiogram_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_AntibiogramAPI_GenAntimicrobialsAntibiogram_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	mux.Handle("GET", pattern_AntibiogramAPI_GenAntimicrobialAntibiogram_0, func(w http.ResponseWriter, req *http.Request, pathParams map[string]string) {
		ctx, cancel := context.WithCancel(req.Context())
		defer cancel()
		inboundMarshaler, outboundMarshaler := runtime.MarshalerForRequest(mux, req)
		rctx, err := runtime.AnnotateContext(ctx, mux, req)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}
		resp, md, err := request_AntibiogramAPI_GenAntimicrobialAntibiogram_0(rctx, inboundMarshaler, client, req, pathParams)
		ctx = runtime.NewServerMetadataContext(ctx, md)
		if err != nil {
			runtime.HTTPError(ctx, mux, outboundMarshaler, w, req, err)
			return
		}

		forward_AntibiogramAPI_GenAntimicrobialAntibiogram_0(ctx, mux, outboundMarshaler, w, req, resp, mux.GetForwardResponseOptions()...)

	})

	return nil
}

var (
	pattern_AntibiogramAPI_GenPathogensAntibiogram_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3}, []string{"api", "antibug", "antibiograms", "pathogens"}, "", runtime.AssumeColonVerbOpt(true)))

	pattern_AntibiogramAPI_GenPathogenAntibiogram_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3}, []string{"api", "antibug", "antibiograms", "pathogen"}, "", runtime.AssumeColonVerbOpt(true)))

	pattern_AntibiogramAPI_GenAntimicrobialsAntibiogram_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3}, []string{"api", "antibug", "antibiograms", "antimicrobials"}, "", runtime.AssumeColonVerbOpt(true)))

	pattern_AntibiogramAPI_GenAntimicrobialAntibiogram_0 = runtime.MustPattern(runtime.NewPattern(1, []int{2, 0, 2, 1, 2, 2, 2, 3}, []string{"api", "antibug", "antibiograms", "antimicrobial"}, "", runtime.AssumeColonVerbOpt(true)))
)

var (
	forward_AntibiogramAPI_GenPathogensAntibiogram_0 = runtime.ForwardResponseMessage

	forward_AntibiogramAPI_GenPathogenAntibiogram_0 = runtime.ForwardResponseMessage

	forward_AntibiogramAPI_GenAntimicrobialsAntibiogram_0 = runtime.ForwardResponseMessage

	forward_AntibiogramAPI_GenAntimicrobialAntibiogram_0 = runtime.ForwardResponseMessage
)