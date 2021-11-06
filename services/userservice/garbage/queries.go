package garbage

import (
	"context"
	"fmt"
	"github.com/MarcGrol/userautomation/core/user"
)

type UserDatabase interface {
	UserInSegment(ctx context.Context, whereClause string) (bool, error)
}

type userDatabase struct{}

func NewUserDatabase() UserDatabase {
	return &userDatabase{}
}

func (db *userDatabase) UserInSegment(ctx context.Context, whereClause string) (bool, error) {
	return true, nil
}

func Query(db UserDatabase, segmentQuery string) user.FilterFunc {
	return func(ctx context.Context, u user.User) (bool, error) {
		whereClause := fmt.Sprintf(`user_uid = '%s' AND %s`, u.UID, segmentQuery)
		exists, err := db.UserInSegment(ctx, whereClause)
		if err != nil {
			return false, fmt.Errorf("Error check: %s", err)
		}
		return exists, nil
	}
}
