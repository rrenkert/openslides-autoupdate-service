package permission

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/OpenSlides/openslides-autoupdate-service/internal/restrict/permission/perm"
)

func NewTestPermission() *Permission {
	p := &Permission{
		hs: newHandlerStore(),
	}

	for _, con := range fakeCollections() {
		con.Connect(p.hs)
	}

	return p
}

func fakeCollections() []perm.Connecter {
	return []perm.Connecter{
		collectionMock{},
	}
}

type collectionMock struct{}

func (c collectionMock) Connect(s perm.HandlerStore) {
	s.RegisterAction("dummy_allowed", allowedMock(true))
	s.RegisterAction("dummy_not_allowed", allowedMock(false))
	s.RegisterAction("dummy_error", allowErrorMock{errors.New("original error message")})

	s.RegisterRestricter("dummy", allowedMock(false))
}

type allowedMock bool

func (a allowedMock) IsAllowed(ctx context.Context, userID int, data map[string]json.RawMessage) (bool, error) {
	return bool(a), nil
}

func (a allowedMock) RestrictFQFields(ctx context.Context, userID int, fqfields []perm.FQField, result map[string]bool) error {
	if !a {
		return nil
	}

	for _, fqfield := range fqfields {
		result[fqfield.String()] = true
	}
	return nil
}

type allowErrorMock struct {
	err error
}

func (a allowErrorMock) IsAllowed(ctx context.Context, userID int, data map[string]json.RawMessage) (bool, error) {
	return false, a.err
}