package service

import (
	"context"
	"errors"
	"time"

	"github.com/Se7enSe7enSe7en/go-toolkit/pkg/logger"
	"github.com/Se7enSe7enSe7en/tenant-manager/internal/auth"
	repo "github.com/Se7enSe7enSe7en/tenant-manager/internal/database/generated"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, email string, password string, name string) (repo.User, repo.Session, error)
	Login(ctx context.Context, email, password string) (repo.User, repo.Session, error)
	Logout(ctx context.Context, sessionID string) error
	LoginWithGoogle(ctx context.Context, googleSub, email, name string, emailVerified bool) (repo.User, repo.Session, error)
	UserFromSession(ctx context.Context, sessionID string) (repo.User, error)
}

type authService struct {
	db                  *pgxpool.Pool
	queries             *repo.Queries
	dummyHashedPassword string
}

func NewAuthService(db *pgxpool.Pool) *authService {
	dummyHashedPassword, _ := auth.HashPassword("a$$word")

	return new(authService{
		db:                  db,
		queries:             repo.New(db),
		dummyHashedPassword: dummyHashedPassword,
	})
}

func (s *authService) Register(ctx context.Context, email string, password string, name string) (repo.User, repo.Session, error) {
	normalizedEmail := auth.NormalizeEmail(email)

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return repo.User{}, repo.Session{}, err
	}
	// note: hashing password is intentionally before starting a transaction because it takes 100ms
	// we want to lessen the time the db connection is in a transaction as much as possible

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.User{}, repo.Session{}, err
	}
	defer tx.Rollback(ctx)
	qtx := s.queries.WithTx(tx)

	user, err := qtx.CreateUser(ctx, repo.CreateUserParams{
		Email: normalizedEmail,
		Name: pgtype.Text{
			String: name,
			Valid:  name != "",
		},
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return repo.User{}, repo.Session{}, auth.ErrEmailAlreadyTaken
		}
		return repo.User{}, repo.Session{}, err
	}

	_, err = qtx.CreateIdentity(ctx, repo.CreateIdentityParams{
		UserID:         user.ID,
		Provider:       auth.ProviderLocal,
		ProviderUserID: normalizedEmail,
		PasswordHash: pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		},
	})
	if err != nil {
		return repo.User{}, repo.Session{}, err
	}

	session, err := qtx.CreateSession(ctx, repo.CreateSessionParams{
		UserID: user.ID,
		ExpiresAt: pgtype.Timestamp{
			Time:  time.Now().Add(time.Hour * 24),
			Valid: true,
		},
	})
	if err != nil {
		return repo.User{}, repo.Session{}, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return repo.User{}, repo.Session{}, err
	}

	return user, session, nil
}

func (s *authService) Login(ctx context.Context, email string, password string) (repo.User, repo.Session, error) {
	normalizedEmail := auth.NormalizeEmail(email)

	identity, err := s.queries.GetIdentityByProvider(ctx, repo.GetIdentityByProviderParams{
		Provider:       auth.ProviderLocal,
		ProviderUserID: normalizedEmail,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {

			// this dummy password check is for security
			_ = auth.CheckPasswordHash(password, s.dummyHashedPassword)
			// note: If the email lookup fails (no identity), you could still run bcrypt.CompareHashAndPassword against a dummy hash
			// to keep the response time uniform. Without this, attackers can time the difference between
			// "no such user" (~1ms, just DB lookup) and "wrong password" (~100ms, bcrypt).

			return repo.User{}, repo.Session{}, auth.ErrInvalidCredentials
		}

		return repo.User{}, repo.Session{}, err
	}

	user, err := s.queries.GetUserByID(ctx, identity.UserID)
	if err != nil {
		return repo.User{}, repo.Session{}, err
	}

	err = auth.CheckPasswordHash(password, identity.PasswordHash.String)
	if err != nil {
		logger.Debug("Login() CheckPasswordHash(): ", err)

		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return repo.User{}, repo.Session{}, auth.ErrInvalidCredentials
		}

		return repo.User{}, repo.Session{}, err
	}

	session, err := s.queries.CreateSession(ctx, repo.CreateSessionParams{
		UserID: identity.UserID,
		ExpiresAt: pgtype.Timestamp{
			Time:  time.Now().Add(time.Hour * 24),
			Valid: true,
		},
	})
	if err != nil {
		return repo.User{}, repo.Session{}, err
	}

	return user, session, nil
}

func (s *authService) Logout(ctx context.Context, sessionID string) error {
	sessionIDUuid, err := uuid.Parse(sessionID)
	if err != nil {
		return err // malformed cookie - treat as "not logged in"
	}

	err = s.queries.DeleteSession(ctx, sessionIDUuid)
	if err != nil {
		return err
	}

	return nil
}

func (s *authService) LoginWithGoogle(ctx context.Context, googleSub string, email string, name string, emailVerified bool) (repo.User, repo.Session, error) {
	// TODO
	return repo.User{}, repo.Session{}, nil
}

func (s *authService) UserFromSession(ctx context.Context, sessionID string) (repo.User, error) {
	sessionIDUuid, err := uuid.Parse(sessionID)
	if err != nil {
		return repo.User{}, err
	}

	session, err := s.queries.GetSession(ctx, sessionIDUuid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// possible cause: logged out, expired-and-cleaned, or never existed
			return repo.User{}, auth.ErrNotAuthenticated
		}

		return repo.User{}, err
	}

	if time.Now().After(session.ExpiresAt.Time) {
		return repo.User{}, auth.ErrNotAuthenticated
	}

	user, err := s.queries.GetUserByID(ctx, session.UserID)
	if err != nil {
		return repo.User{}, err
	}

	return user, nil
}
