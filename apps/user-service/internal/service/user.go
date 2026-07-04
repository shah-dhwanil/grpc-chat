package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/shah-dhwanil/grpc-chat/apps/user-service/internal/dto"
	"github.com/shah-dhwanil/grpc-chat/apps/user-service/internal/repository"
	userv1 "github.com/shah-dhwanil/grpc-chat/packages/api/gen/user/v1"
)


type UserService struct {

	userv1.UnimplementedUserServiceServer
	repository *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{
		repository: repo,
	}
}

func (srv *UserService) CreateUser(ctx context.Context, payload *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error){
	res, err :=srv.repository.UserRepository.CreateUser(ctx, &dto.CreateUserRequest{
		Name: payload.DisplayName,
		PrimaryEmail: payload.PrimaryEmail,
	})
	if err != nil {
		return nil, err
	}
	return &userv1.CreateUserResponse{
		User: mapToUser(res),
	},nil
}
func (srv *UserService) GetUser(ctx context.Context, payload *userv1.GetUserRequest) (*userv1.GetUserResponse, error){
	uid,err := uuid.Parse(payload.UserId)
	if err != nil {
		return nil, err
	}
	res,err:= srv.repository.UserRepository.GetUserByID(ctx,uid)
	if err != nil {
		return nil, err
	}
	return &userv1.GetUserResponse{
		User: mapToUser(res),
	},nil
}
func (srv *UserService) ListUser(ctx context.Context, payload *userv1.ListUserRequest) (*userv1.ListUserResponse, error){
	uids := make([]uuid.UUID, len(payload.UserIds))
	for i, id := range payload.UserIds {
		uid, err := uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		uids[i] = uid
	}
	res,err:= srv.repository.UserRepository.GetUsers(ctx,uids)
	if err != nil {
		return nil, err
	}
	users := make([]*userv1.UserListItem, len(res))
	for i, user := range res {
		users[i] = &userv1.UserListItem{
			Id: user.ID.String(),
			DisplayName: user.Name,
		}
	}
	return &userv1.ListUserResponse{
		Users: users,
	},nil
}
func (srv *UserService) UpdateUser(ctx context.Context, payload *userv1.UpdateUserRequest) (*userv1.UpdateUserResponse, error){
	uid,err := uuid.Parse(ctx.Value("user_id").(string))
	if err != nil {
		return nil, err
	}
	res,err:= srv.repository.UserRepository.UpdateUser(ctx,uid,&dto.UpdateUserRequest{
		Name: payload.DisplayName,
	})
	if err != nil {
		return nil, err
	}
	return &userv1.UpdateUserResponse{
		User: mapToUser(res),
	},nil
}
func (srv *UserService) SetUserPrimaryEmail(ctx context.Context, payload *userv1.SetUserPrimaryEmailRequest) (*userv1.SetUserPrimaryEmailResponse, error){
	uid,err := uuid.Parse(ctx.Value("user_id").(string))
	if err != nil {
		return nil, err
	}
	res,err:= srv.repository.UserRepository.UpdateUser(ctx,uid,&dto.UpdateUserRequest{
		PrimaryEmail: &payload.EmailId,
	})
	if err != nil {
		return nil, err
	}
	return &userv1.SetUserPrimaryEmailResponse{
		User: mapToUser(res),
	},nil
}
func (srv *UserService) DeleteUser(ctx context.Context, _ *userv1.DeleteUserRequest) (*userv1.DeleteUserResponse, error){
	uid,err := uuid.Parse(ctx.Value("user_id").(string))
	if err != nil {
		return nil, err
	}
	err = srv.repository.UserRepository.DeleteUser(context.Background(),uid)
	if err != nil {
		return nil, err
	}
	return &userv1.DeleteUserResponse{},nil
}


func mapToUser(user *dto.User) *userv1.User{
	return &userv1.User{
		Id: user.ID.String(),
		DisplayName: user.Name,
		PrimaryEmail: user.PrimaryEmail,
	}
}