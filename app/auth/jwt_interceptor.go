package auth

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)
var (
	ErrMissingMetadata = errors.New("missing metadata")
	ErrMissingAuthorizationToken = errors.New("missing authorization token")
	ErrInvalidAuthorizationHeader = errors.New("invalid authorization header format")
	ErrUnexpectedSigningMethod = errors.New("unexpected signing method")
	ErrInvalidToken = errors.New("invalid token")
)
func JWTInterceptor(secretKey []byte) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		methodsToSkip := map[string]bool{
			"/client.service.company.v1.CompanyService.Get": true,
		}
		if methodsToSkip[info.FullMethod] {
			log.Println(info.FullMethod)
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, ErrMissingMetadata
		}

		authorization := md["authorization"]
		if len(authorization) == 0 {
			return nil, ErrMissingAuthorizationToken
		}

		authParts := strings.SplitN(authorization[0]," ",2)
		if len(authParts) < 2 {
			return nil, ErrInvalidAuthorizationHeader
		}
		if !strings.EqualFold(authParts[0],"Bearer"){
			return nil, ErrInvalidAuthorizationHeader
		}

		claims := &jwt.MapClaims{}
		token, err := jwt.ParseWithClaims(authParts[1], claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, ErrUnexpectedSigningMethod
			}
			return secretKey, nil
		})
		if err != nil || !token.Valid {
			return nil, ErrInvalidToken
		}

		newCtx := context.WithValue(ctx, "userClaims", claims)

		return handler(newCtx, req)
	}
}
