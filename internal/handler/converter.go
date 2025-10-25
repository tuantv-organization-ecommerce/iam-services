package handler

import (
	"github.com/tvttt/iam-services/internal/domain"
	pb "github.com/tvttt/iam-services/pkg/proto"
)

// domainRoleToPB converts domain.Role to pb.Role
func domainRoleToPB(role *domain.Role) *pb.Role {
	if role == nil {
		return nil
	}

	pbRole := &pb.Role{
		Id:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}

	if len(role.Permissions) > 0 {
		pbRole.Permissions = make([]*pb.Permission, len(role.Permissions))
		for i, perm := range role.Permissions {
			pbRole.Permissions[i] = domainPermissionToPB(&perm)
		}
	}

	return pbRole
}

// domainPermissionToPB converts domain.Permission to pb.Permission
func domainPermissionToPB(perm *domain.Permission) *pb.Permission {
	if perm == nil {
		return nil
	}

	return &pb.Permission{
		Id:          perm.ID,
		Name:        perm.Name,
		Resource:    perm.Resource,
		Action:      perm.Action,
		Description: perm.Description,
		CreatedAt:   perm.CreatedAt.Format("2006-01-02T15:04:05Z"),
		UpdatedAt:   perm.UpdatedAt.Format("2006-01-02T15:04:05Z"),
	}
}
