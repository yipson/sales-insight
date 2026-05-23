package tokencache

import (
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCache_SetAndGet(t *testing.T) {
	c := NewCache()
	id := uuid.New()
	td := TokenData{
		MerchantID:   id,
		AccessToken:  "access123",
		RefreshToken: "refresh456",
	}

	c.Set(id, td)

	got, ok := c.Get(id)
	if !ok {
		t.Fatal("expected token to be found")
	}
	if got.AccessToken != "access123" {
		t.Errorf("access token = %q, want %q", got.AccessToken, "access123")
	}
	if got.RefreshToken != "refresh456" {
		t.Errorf("refresh token = %q, want %q", got.RefreshToken, "refresh456")
	}
}

func TestCache_Get_NotFound(t *testing.T) {
	c := NewCache()
	id := uuid.New()

	_, ok := c.Get(id)
	if ok {
		t.Error("expected token to not be found")
	}
}

func TestCache_Delete(t *testing.T) {
	c := NewCache()
	id := uuid.New()
	c.Set(id, TokenData{MerchantID: id, AccessToken: "token"})

	c.Delete(id)

	_, ok := c.Get(id)
	if ok {
		t.Error("expected token to be deleted")
	}
}

func TestCache_IsAccessTokenExpired_NotExpired(t *testing.T) {
	c := NewCache()
	id := uuid.New()
	c.Set(id, TokenData{
		MerchantID:           id,
		AccessToken:          "token",
		AccessTokenExpiresAt: time.Now().Add(1 * time.Hour),
	})

	if c.IsAccessTokenExpired(id) {
		t.Error("expected token to NOT be expired (expires in 1 hour)")
	}
}

func TestCache_IsAccessTokenExpired_Expired(t *testing.T) {
	c := NewCache()
	id := uuid.New()
	c.Set(id, TokenData{
		MerchantID:           id,
		AccessToken:          "token",
		AccessTokenExpiresAt: time.Now().Add(-1 * time.Hour),
	})

	if !c.IsAccessTokenExpired(id) {
		t.Error("expected token to be expired (expired 1 hour ago)")
	}
}

func TestCache_IsAccessTokenExpired_NotFound(t *testing.T) {
	c := NewCache()
	id := uuid.New()

	if !c.IsAccessTokenExpired(id) {
		t.Error("expected not-found token to be treated as expired")
	}
}

func TestCache_IsAccessTokenExpired_FiveMinuteBuffer(t *testing.T) {
	c := NewCache()
	id := uuid.New()
	// Token expires in 3 minutes — within the 5-minute buffer
	c.Set(id, TokenData{
		MerchantID:           id,
		AccessToken:          "token",
		AccessTokenExpiresAt: time.Now().Add(3 * time.Minute),
	})

	if !c.IsAccessTokenExpired(id) {
		t.Error("expected token to be expired due to 5-minute buffer")
	}
}

func TestCache_ConcurrentAccess(t *testing.T) {
	c := NewCache()
	id := uuid.New()
	var wg sync.WaitGroup

	// Writers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			c.Set(id, TokenData{MerchantID: id, AccessToken: "token"})
		}(i)
	}

	// Readers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _ = c.Get(id)
		}()
	}

	// Expiry checkers
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = c.IsAccessTokenExpired(id)
		}()
	}

	wg.Wait()
	// If we get here without panic or data race, the test passes.
}
