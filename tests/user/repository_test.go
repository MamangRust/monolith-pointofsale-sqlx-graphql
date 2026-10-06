package user_test

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-pointofsale-shared/domain/requests"

	pbrole "github.com/MamangRust/monolith-pointofsale-pb/role"
	pbuserrole "github.com/MamangRust/monolith-pointofsale-pb/user_role"
	tests "github.com/MamangRust/monolith-pointofsale-test"
	"github.com/MamangRust/monolith-pointofsale-user/repository"
	"github.com/stretchr/testify/suite"
)

func ptr[T any](v T) *T {
	return &v
}

type UserRepositoryTestSuite struct {
	tests.BaseTestSuite
	repo   *repository.Repositories
	userID int
}

func (s *UserRepositoryTestSuite) SetupSuite() {
	s.BaseTestSuite.SetupSuite()

	queries := s.SQLxDB()
	s.SetupRoleService()
	s.repo = repository.NewRepositories(&repository.Deps{
		Db:              queries,
		RoleQueryClient: pbrole.NewRoleQueryServiceClient(s.Conns["role"]),
		UserRoleClient:  pbuserrole.NewUserRoleServiceClient(s.Conns["role"]),
	})
}

func (s *UserRepositoryTestSuite) TearDownSuite() {
	s.BaseTestSuite.TearDownSuite()
}

func (s *UserRepositoryTestSuite) Test1_CreateUser() {
	ctx := context.Background()

	req := &requests.CreateUserRequest{
		FirstName:       "John",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "password123",
		ConfirmPassword: "password123",
	}
	created, err := s.repo.UserCommand.CreateUser(ctx, req)
	s.NoError(err)
	s.NotNil(created)
	s.Equal(req.FirstName, created.Firstname)
	s.Equal(req.Email, created.Email)
	s.userID = int(created.UserID)
}

func (s *UserRepositoryTestSuite) Test2_FindById() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	found, err := s.repo.UserQuery.FindById(ctx, s.userID)
	s.NoError(err)
	s.NotNil(found)
	s.Equal(s.userID, int(found.UserID))
}

func (s *UserRepositoryTestSuite) Test3_FindByEmail() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	found, err := s.repo.UserQuery.FindByEmail(ctx, "john.doe@example.com")
	s.NoError(err)
	s.NotNil(found)
	s.Equal("john.doe@example.com", found.Email)
}

func (s *UserRepositoryTestSuite) Test4_FindByEmailWithPassword() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	found, err := s.repo.UserQuery.FindByEmail(ctx, "john.doe@example.com")
	s.NoError(err)
	s.NotNil(found)
	s.NotEmpty(found.Password)
}

func (s *UserRepositoryTestSuite) Test5_FindAllUser() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	all, _, err := s.repo.UserQuery.FindAllUsers(ctx, &requests.FindAllUsers{Search: "", Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(all)
}

func (s *UserRepositoryTestSuite) Test6_FindActive() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	active, _, err := s.repo.UserQuery.FindByActive(ctx, &requests.FindAllUsers{Search: "", Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(active)
}

func (s *UserRepositoryTestSuite) Test7_UpdateUser() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	updateReq := &requests.UpdateUserRequest{
		UserID:          &s.userID,
		FirstName:       "Updated",
		LastName:        "Doe",
		Email:           "john.doe@example.com",
		Password:        "newpassword123",
		ConfirmPassword: "newpassword123",
	}
	updated, err := s.repo.UserCommand.UpdateUser(ctx, updateReq)
	s.NoError(err)
	s.NotNil(updated)
	s.Equal("Updated", updated.Firstname)
}

func (s *UserRepositoryTestSuite) Test8_TrashAndFindTrashed() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	// Trash
	_, err := s.repo.UserCommand.TrashedUser(ctx, s.userID)
	s.NoError(err)

	// FindTrashed
	trashed, _, err := s.repo.UserQuery.FindByTrashed(ctx, &requests.FindAllUsers{Search: "", Page: 1, PageSize: 10})
	s.NoError(err)
	s.NotEmpty(trashed)

	// FindActive — should NOT include trashed user
	active, _, err := s.repo.UserQuery.FindByActive(ctx, &requests.FindAllUsers{Search: "", Page: 1, PageSize: 10})
	s.NoError(err)
	for _, u := range active {
		s.NotEqual(s.userID, int(u.UserID))
	}
}

func (s *UserRepositoryTestSuite) Test9_Restore() {
	s.Require().NotZero(s.userID)
	ctx := context.Background()

	_, err := s.repo.UserCommand.RestoreUser(ctx, s.userID)
	s.NoError(err)

	// Verify it's back in active
	active, _, err := s.repo.UserQuery.FindByActive(ctx, &requests.FindAllUsers{Search: "", Page: 1, PageSize: 10})
	s.NoError(err)
	var found bool
	for _, u := range active {
		if int(u.UserID) == s.userID {
			found = true
			break
		}
	}
	s.True(found, "restored user should appear in active list")
}

func (s *UserRepositoryTestSuite) Test10_DeletePermanent() {
	ctx := context.Background()

	// Create a fresh user (testify runs methods in name order; Test10 runs before Test1)
	u, err := s.repo.UserCommand.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Del", LastName: "Perm", Email: "del.perm@example.com",
		Password: "password123", ConfirmPassword: "password123",
	})
	s.Require().NoError(err)

	// Must trash first
	_, err = s.repo.UserCommand.TrashedUser(ctx, int(u.UserID))
	s.Require().NoError(err)

	success, err := s.repo.UserCommand.DeleteUserPermanent(ctx, int(u.UserID))
	s.NoError(err)
	s.True(success)
}

func (s *UserRepositoryTestSuite) Test11_UpdatePassword() {
	ctx := context.Background()

	// Create a fresh user for password test
	u, err := s.repo.UserCommand.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName: "Pass", LastName: "Test", Email: "password.test@example.com",
		Password: "oldpass123", ConfirmPassword: "oldpass123",
	})
	s.Require().NoError(err)

	updatedPass, err := s.repo.UserCommand.UpdateUser(ctx, &requests.UpdateUserRequest{UserID: ptr(int(u.UserID))})
	s.NoError(err)
	s.NotNil(updatedPass)

	// Cleanup
	s.repo.UserCommand.TrashedUser(ctx, int(u.UserID))
	s.repo.UserCommand.DeleteUserPermanent(ctx, int(u.UserID))
}

func TestUserRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}
	suite.Run(t, new(UserRepositoryTestSuite))
}
