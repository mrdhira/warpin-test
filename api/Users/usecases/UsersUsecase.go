package usecases

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/mrdhira/warpin-test/api/Users/entities"
	"github.com/mrdhira/warpin-test/api/Users/infrastructures/database"
	"github.com/mrdhira/warpin-test/api/Users/infrastructures/repositories"
	"github.com/mrdhira/warpin-test/pkg"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

// IUsersUsecases interface
type IUsersUsecases interface {
	Register(ctx context.Context, Data *entities.RegisterRequest) (Response *pkg.JSONResponse, err error)
	Login(ctx context.Context, Data *entities.LoginRequest) (Response *pkg.JSONResponse, err error)
	Profile(ctx context.Context, Data *entities.ProfileRequest) (Response *pkg.JSONResponse, err error)
	UpdateProfile(ctx context.Context, Data *entities.UpdateProfileRequest) (Response *pkg.JSONResponse, err error)
	UpdatePassword(ctx context.Context, Data *entities.UpdatePasswordRequest) (Response *pkg.JSONResponse, err error)
}

// UsersUsecases struct
type UsersUsecases struct {
	UsersRepository repositories.IUsersRepository
}

// InitUsersUsecases func
func InitUsersUsecases() *UsersUsecases {
	// Init Repositories
	usersRepository := new(repositories.UsersRepository)
	usersRepository.PG = &database.PostgresConnection{}
	usersRepository.Redis = &database.RedisConnection{}

	return &UsersUsecases{
		UsersRepository: usersRepository,
	}
}

// Register usecases
func (u *UsersUsecases) Register(ctx context.Context, Data *entities.RegisterRequest) (Response *pkg.JSONResponse, err error) {
	CheckUsers, err := u.UsersRepository.UsersFindByEmail(ctx, Data.Email)
	if err != nil {
		return
	}

	if CheckUsers != nil {
		return &pkg.JSONResponse{
			Code:    422,
			Message: "Anda sudah terdaftar dengan email " + Data.Email,
		}, nil
	}

	Tx, err := u.UsersRepository.Tx()
	if err != nil {
		return
	}
	defer Tx.RollbackUnlessCommitted()

	Users := &entities.Users{
		Email:       Data.Email,
		PhoneNumber: Data.PhoneNumber,
		FullName:    Data.FullName,
		Gender:      Data.Gender,
		Role:        Data.Role,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	Hash, err := bcrypt.GenerateFromPassword([]byte(Data.Password), bcrypt.DefaultCost)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	Users.Password = string(Hash)

	Users.ID, err = u.UsersRepository.UsersStore(ctx, Tx, Users)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	UsersLog := &entities.UsersLog{
		UserID:      Users.ID,
		Email:       Data.Email,
		PhoneNumber: Data.PhoneNumber,
		FullName:    Data.FullName,
		Gender:      Data.Gender,
		Role:        Data.Role,
		Password:    Users.Password,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	UsersLog.ID, err = u.UsersRepository.UsersLogStore(ctx, Tx, UsersLog)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	defer Tx.Commit()

	return &pkg.JSONResponse{
		Code:    200,
		Message: "Berhasil terdaftar",
		Data:    Users,
	}, nil
}

// Login usecases
func (u *UsersUsecases) Login(ctx context.Context, Data *entities.LoginRequest) (Response *pkg.JSONResponse, err error) {
	Users, err := u.UsersRepository.UsersFindByEmail(ctx, Data.Email)
	if err != nil {
		return
	}

	if Users == nil {
		return &pkg.JSONResponse{
			Code:    404,
			Message: "users tidak ditemukan",
		}, nil
	}

	err = bcrypt.CompareHashAndPassword([]byte(Users.Password), []byte(Data.Password))
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when compare hash and password",
		}).Error(err)

		return &pkg.JSONResponse{
			Code:    403,
			Message: "Password yang anda masukkan salah",
			Error:   err.Error(),
		}, nil
	}

	TokenData := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), &entities.TokenClaim{
		UserID:   Users.ID,
		UserRole: Users.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
		},
	})

	TokenString, err := TokenData.SignedString([]byte("secret"))
	if err != nil {
		log.WithFields(log.Fields{
			"event": "error when signed string for jwt token",
		}).Error(err)

		return &pkg.JSONResponse{
			Code:    400,
			Message: "Terjadi kesalahan sistem",
			Error:   err.Error(),
		}, nil
	}

	return &pkg.JSONResponse{
		Code:    200,
		Message: "Berhasil login",
		Data:    TokenString,
	}, nil
}

// Profile usecases
func (u *UsersUsecases) Profile(ctx context.Context, Data *entities.ProfileRequest) (Response *pkg.JSONResponse, err error) {
	Profile, err := u.UsersRepository.ProfileByID(ctx, Data.UserID)
	if err != nil {
		return
	}

	return &pkg.JSONResponse{
		Code:    200,
		Message: "OK",
		Data:    Profile,
	}, nil
}

// UpdateProfile usecases
func (u *UsersUsecases) UpdateProfile(ctx context.Context, Data *entities.UpdateProfileRequest) (Response *pkg.JSONResponse, err error) {
	if Data.FullName == "" && Data.PhoneNumber == "" {
		return &pkg.JSONResponse{
			Code:    422,
			Message: "Nama dan Nomor Handphone tidak bisa kosong semua",
		}, nil
	}

	Users, err := u.UsersRepository.UsersFindByID(ctx, Data.UserID)
	if err != nil {
		return
	}

	Tx, err := u.UsersRepository.Tx()
	if err != nil {
		return
	}
	defer Tx.RollbackUnlessCommitted()

	UsersLog := &entities.UsersLog{
		UserID:      Users.ID,
		Email:       Users.Email,
		PhoneNumber: Users.PhoneNumber,
		FullName:    Users.FullName,
		Gender:      Users.Gender,
		Role:        Users.Role,
		Password:    Users.Password,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	UpdatePayload := map[string]interface{}{}
	if Data.FullName != "" {
		UsersLog.FullName = Data.FullName
		UpdatePayload["full_name"] = Data.FullName
	}
	if Data.PhoneNumber != "" {
		UsersLog.PhoneNumber = Data.PhoneNumber
		UpdatePayload["phone_number"] = Data.PhoneNumber
	}

	err = u.UsersRepository.UsersUpdate(ctx, Tx, Data.UserID, UpdatePayload)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	UsersLog.ID, err = u.UsersRepository.UsersLogStore(ctx, Tx, UsersLog)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	defer Tx.Commit()

	return &pkg.JSONResponse{
		Code:    200,
		Message: "Profile berhasil diupdate",
	}, nil
}

// UpdatePassword usecases
func (u *UsersUsecases) UpdatePassword(ctx context.Context, Data *entities.UpdatePasswordRequest) (Response *pkg.JSONResponse, err error) {
	Users, err := u.UsersRepository.UsersFindByID(ctx, Data.UserID)
	if err != nil {
		return
	}

	Tx, err := u.UsersRepository.Tx()
	if err != nil {
		return
	}
	defer Tx.RollbackUnlessCommitted()

	Hash, err := bcrypt.GenerateFromPassword([]byte(Data.Password), bcrypt.DefaultCost)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	UsersLog := &entities.UsersLog{
		UserID:      Users.ID,
		Email:       Users.Email,
		PhoneNumber: Users.PhoneNumber,
		FullName:    Users.FullName,
		Gender:      Users.Gender,
		Role:        Users.Role,
		Password:    string(Hash),
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	UpdatePayload := map[string]interface{}{
		"password": string(Hash),
	}

	err = u.UsersRepository.UsersUpdate(ctx, Tx, Data.UserID, UpdatePayload)
	if err != nil {
		return
	}

	UsersLog.ID, err = u.UsersRepository.UsersLogStore(ctx, Tx, UsersLog)
	if err != nil {
		defer Tx.Rollback()
		return
	}

	defer Tx.Commit()

	return &pkg.JSONResponse{
		Code:    200,
		Message: "Password berhasil diupdate",
	}, nil
}
