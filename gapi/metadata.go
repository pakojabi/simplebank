package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserHeader = "grpcgateway-user-agent"
	userAgentHeader = "user-agent"
	xForwardedForHeader = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP string
}

func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}

	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// grpc gateway
		if userAgents := md.Get(grpcGatewayUserHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}
		// grpc
		if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}
		// grpc gateway
		if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
			mtdt.ClientIP = clientIPs[0]
		}
	}
	
	// grpc
	if peer, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = peer.Addr.String()
	}

	return mtdt
}