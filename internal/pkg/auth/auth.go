package auth

import (
	"context"
	"fmt"
	"github.com/gidyon/antibug/internal/pkg/errs"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	// Pharmacist ...
	Pharmacist = "PHARMARCY"
	// Physician ...
	Physician = "PHYSICIAN"
	// Researcher ...
	Researcher = "RESEARCHER"
	// LabTechnician ...
	LabTechnician = "LAB_TECHNICIAN"
	// Admin ...
	Admin = "ADMIN"
	// Super Admin ...
)

// AuthenticateCtx authenticates ctx in request
func AuthenticateCtx(ctx context.Context) error {
	_, err := ParseFromCtx(ctx)
	if err != nil {
		return err
	}
	return nil
}

// AuthorizeGroup authenticates whether token belongs to member of a particular group
func AuthorizeGroup(ctx context.Context, allowedGroups ...string) (*Payload, error) {
	claims, err := ParseFromCtx(ctx)
	if err != nil {
		return nil, err
	}

	err = matchGroup(claims.Payload.Group, allowedGroups)
	if err != nil {
		return nil, err
	}

	return claims.Payload, nil
}

// AuthorizeGroupAndID authenticates member of a particular group and ad having given ID
func AuthorizeGroupAndID(ctx context.Context, ID string, allowedGroups ...string) error {
	claims, err := ParseFromCtx(ctx)
	if err != nil {
		return err
	}

	err = matchGroup(claims.Payload.Group, allowedGroups)
	if err != nil {
		return err
	}

	if claims.ID != ID {
		return errs.TokenCredentialNotMatching("ID")
	}

	return nil
}

// AuthorizeGroupFromToken authenticates member of a particular group from token
func AuthorizeGroupFromToken(token string, allowedGroups ...string) (*Payload, error) {
	claims, err := ParseToken(token)
	if err != nil {
		return nil, err
	}

	err = matchGroup(claims.Payload.Group, allowedGroups)
	if err != nil {
		return nil, err
	}

	return claims.Payload, nil
}

// AuthorizeGroupAndIDFromToken authenticates group and ID of a particular group
func AuthorizeGroupAndIDFromToken(token, ID string, allowedGroups ...string) error {
	claims, err := AuthorizeGroupFromToken(token, allowedGroups...)
	if err != nil {
		return err
	}
	if claims.ID != ID {
		return errs.TokenCredentialNotMatching("ID")
	}
	return nil
}

// AddGroupAndIDMD creates a metadata context with the group ID and token ID
func AddGroupAndIDMD(ctx context.Context, group, ID string) context.Context {
	payload := &Payload{ID: ID, Group: group}
	token, err := GenToken(ctx, payload, group, 0)
	if err != nil {
		panic(err)
	}

	return addTokenMD(ctx, token)
}

func addTokenMD(ctx context.Context, token string) context.Context {
	return metadata.NewIncomingContext(
		ctx, metadata.Pairs("authorization", fmt.Sprintf("Bearer %s", token)),
	)
}

func matchGroup(claimGroup string, allowedGroups []string) error {
	for _, group := range allowedGroups {
		if claimGroup == group {
			return nil
		}
	}
	return status.Errorf(codes.PermissionDenied, "permission denied for group %s", claimGroup)
}
