package auth

import (
	"context"
	"errors"
	"fmt" // Added fmt import
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/securecookie"
	"github.com/hermesgen/hm"
)

const ( // TODO: Move to config
	sessionCookieName = "user_session"
	sessionMaxAge     = 3600 * 24 * 7 // 1 week
)

// SessionManager handles user sessions.
type SessionManager struct {
	*hm.BaseCore
	encoder *securecookie.SecureCookie
}

// NewSessionManager creates a new SessionManager.
func NewSessionManager(params hm.XParams) *SessionManager {
	core := hm.NewCore("session-manager", params)
	return &SessionManager{
		BaseCore: core,
	}
}

// Setup initializes the SessionManager.
func (sm *SessionManager) Setup(ctx context.Context) error {
	err := sm.BaseCore.Setup(ctx)
	if err != nil {
		return err
	}

	cfg := sm.Cfg()
	hashKey := cfg.ByteSliceVal(hm.Key.SecHashKey)
	blockKey := cfg.ByteSliceVal(hm.Key.SecBlockKey)

	if len(hashKey) == 0 || len(blockKey) == 0 {
		return errors.New("missing hashKey or blockKey in configuration for session manager")
	}

	sm.encoder = securecookie.New(hashKey, blockKey)
	return nil
}

// SetUserSession sets a user ID in a session cookie.
func (sm *SessionManager) SetUserSession(w http.ResponseWriter, userID uuid.UUID) error {
	if sm.encoder == nil {
		return errors.New("session manager not initialized")
	}

	value := map[string]string{
		"user_id": userID.String(),
	}

	encoded, err := sm.encoder.Encode(sessionCookieName, value)
	if err != nil {
		return fmt.Errorf("failed to encode session cookie: %w", err)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    encoded,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Should be true in production (HTTPS)
		MaxAge:   sessionMaxAge,
		Expires:  time.Now().Add(time.Duration(sessionMaxAge) * time.Second),
		SameSite: http.SameSiteLaxMode,
	})

	return nil
}

// GetUserSession retrieves a user ID from a session cookie.
func (sm *SessionManager) GetUserSession(r *http.Request) (uuid.UUID, error) {
	if sm.encoder == nil {
		return uuid.Nil, errors.New("session manager not initialized")
	}

	cookie, err := r.Cookie(sessionCookieName)
	if err != nil {
		return uuid.Nil, fmt.Errorf("session cookie not found: %w", err)
	}

	value := make(map[string]string)
	err = sm.encoder.Decode(sessionCookieName, cookie.Value, &value)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to decode session cookie: %w", err)
	}

	userIDStr, ok := value["user_id"]
	if !ok || userIDStr == "" {
		return uuid.Nil, errors.New("user ID not found in session")
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid user ID in session: %w", err)
	}

	return userID, nil
}

// ClearUserSession clears the user session cookie.
func (sm *SessionManager) ClearUserSession(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // Delete cookie now
		HttpOnly: true,
		Secure:   true, // Should be true in production (HTTPS)
		SameSite: http.SameSiteLaxMode,
	})
}
