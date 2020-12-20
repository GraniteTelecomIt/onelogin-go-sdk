package users

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/onelogin/onelogin-go-sdk/pkg/services"
	"github.com/onelogin/onelogin-go-sdk/pkg/services/olhttp"
)

const errUsersV2Context = "users v2 service"

// V2Service holds the information needed to interface with a repository
type V2Service struct {
	Endpoint, ErrorContext string
	Repository             services.Repository
	Host                   string
}

// New creates the new svc service v2.
func New(repo services.Repository, host string) *V2Service {
	return &V2Service{
		Endpoint:     fmt.Sprintf("%s/api/2/users", host),
		Repository:   repo,
		ErrorContext: errUsersV2Context,
		Host:         host,
	}
}

// Query retrieves all the users from the repository that meet the query criteria passed in the
// request payload. If an empty payload is given, it will retrieve all users
func (svc *V2Service) Query(query *UserQuery) ([]User, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    query,
	})
	if err != nil {
		return nil, err
	}

	var users []User
	json.Unmarshal(resp, &users)
	return users, nil
}

// GetOne retrieves the user by id and returns it
func (svc *V2Service) GetOne(id int32) (*User, error) {
	resp, err := svc.Repository.Read(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	})
	if err != nil {
		return nil, err
	}
	var user User
	json.Unmarshal(resp, &user)
	return &user, nil
}

// Create takes a user without an id and attempts to use the parameters to create it
// in the API. Modifies the user in place, or returns an error if one occurs
func (svc *V2Service) Create(user *User) error {
	resp, err := svc.Repository.Create(olhttp.OLHTTPRequest{
		URL:        svc.Endpoint,
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    user,
	})
	if err != nil {
		return err
	}
	json.Unmarshal(resp, user)
	return nil
}

// Update takes a user and an id and attempts to use the parameters to update it
// in the API. Modifies the user in place, or returns an error if one occurs
func (svc *V2Service) Update(user *User) error {
	if user.ID == nil {
		return errors.New("No ID Given")
	}
	resp, err := svc.Repository.Update(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d", svc.Endpoint, *user.ID),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
		Payload:    user,
	})
	if err != nil {
		return err
	}
	json.Unmarshal(resp, user)
	return nil
}

// Destroy deletes the user with the given id, and if successful, it returns nil
func (svc *V2Service) Destroy(id int32) error {
	if _, err := svc.Repository.Destroy(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/%d", svc.Endpoint, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	}); err != nil {
		return err
	}
	return nil
}

// Logout logs out the user from OneLogin
func (svc *V2Service) Logout(id int32) error {
	if _, err := svc.Repository.Update(olhttp.OLHTTPRequest{
		URL:        fmt.Sprintf("%s/api/1/users/%d/logout", svc.Host, id),
		Headers:    map[string]string{"Content-Type": "application/json"},
		AuthMethod: "bearer",
	}); err != nil {
		return err
	}
	return nil
}
