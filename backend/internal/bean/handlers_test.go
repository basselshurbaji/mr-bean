package bean_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/bean"
	"github.com/basselshurbaji/mr_bean/backend/internal/handler"
	"github.com/basselshurbaji/mr_bean/backend/internal/principal"
)

type mockBeanService struct {
	item *bean.Bean
	list []bean.Bean
	err  error
}

func (m *mockBeanService) ListBeans(_ context.Context, _ string) ([]bean.Bean, error) {
	return m.list, m.err
}
func (m *mockBeanService) CreateBean(_ context.Context, _ string, _ bean.BeanParams) (*bean.Bean, error) {
	return m.item, m.err
}
func (m *mockBeanService) UpdateBean(_ context.Context, _, _ string, _ bean.BeanParams) (*bean.Bean, error) {
	return m.item, m.err
}
func (m *mockBeanService) DeleteBean(_ context.Context, _, _ string) error { return m.err }

func authedCtx() context.Context {
	return principal.WithUserID(context.Background(), "user-1")
}

var sampleBean = bean.Bean{
	ID:        "bean-1",
	UserID:    "user-1",
	Name:      "Ethiopia Yirgacheffe",
	CreatedAt: time.Now(),
	UpdatedAt: time.Now(),
}

func strPtr(s string) *string { return &s }

// --- CreateBeanHandler ---

func TestCreateBeanHandler_Validate(t *testing.T) {
	h := bean.NewCreateBeanHandler(nil)

	cases := []struct {
		name    string
		req     bean.CreateBeanRequest
		wantErr bool
	}{
		{"valid minimal", bean.CreateBeanRequest{Name: "Ethiopia"}, false},
		{"missing name", bean.CreateBeanRequest{}, true},
		{"valid process", bean.CreateBeanRequest{Name: "X", Process: strPtr("washed")}, false},
		{"invalid process", bean.CreateBeanRequest{Name: "X", Process: strPtr("unknown")}, true},
		{"valid roast_level", bean.CreateBeanRequest{Name: "X", RoastLevel: strPtr("light")}, false},
		{"invalid roast_level", bean.CreateBeanRequest{Name: "X", RoastLevel: strPtr("extra_dark")}, true},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := h.Validate(tc.req)
			if tc.wantErr && err == nil {
				t.Fatal("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
		})
	}
}

func TestCreateBeanHandler_Serve_Success(t *testing.T) {
	svc := &mockBeanService{item: &sampleBean}
	h := bean.NewCreateBeanHandler(svc)

	res, err := h.Serve(authedCtx(), bean.CreateBeanRequest{Name: "Ethiopia"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ID != sampleBean.ID {
		t.Fatalf("expected id %q, got %q", sampleBean.ID, res.ID)
	}
}

// --- UpdateBeanHandler ---

func TestUpdateBeanHandler_Validate(t *testing.T) {
	h := bean.NewUpdateBeanHandler(nil)

	if err := h.Validate(bean.UpdateBeanRequest{Name: "Colombia"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := h.Validate(bean.UpdateBeanRequest{}); err == nil {
		t.Fatal("expected error for missing name")
	}
	if err := h.Validate(bean.UpdateBeanRequest{Name: "X", Process: strPtr("bad")}); err == nil {
		t.Fatal("expected error for invalid process")
	}
}

func TestUpdateBeanHandler_Serve_NotFound(t *testing.T) {
	svc := &mockBeanService{err: bean.ErrNotFound}
	h := bean.NewUpdateBeanHandler(svc)

	_, err := h.Serve(authedCtx(), bean.UpdateBeanRequest{ID: "missing", Name: "X"})
	var appErr *handler.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T: %v", err, err)
	}
	if appErr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", appErr.Code)
	}
}

// --- DeleteBeanHandler ---

func TestDeleteBeanHandler_Validate(t *testing.T) {
	h := bean.NewDeleteBeanHandler(nil)

	if err := h.Validate(bean.DeleteBeanRequest{ID: "x"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := h.Validate(bean.DeleteBeanRequest{}); err == nil {
		t.Fatal("expected error for missing id")
	}
}

func TestDeleteBeanHandler_Serve_NotFound(t *testing.T) {
	svc := &mockBeanService{err: bean.ErrNotFound}
	h := bean.NewDeleteBeanHandler(svc)

	_, err := h.Serve(authedCtx(), bean.DeleteBeanRequest{ID: "missing"})
	var appErr *handler.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected AppError, got %T", err)
	}
	if appErr.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", appErr.Code)
	}
}

// --- ListBeansHandler ---

func TestListBeansHandler_Serve_EmptyList(t *testing.T) {
	svc := &mockBeanService{list: []bean.Bean{}}
	h := bean.NewListBeansHandler(svc)

	res, err := h.Serve(authedCtx(), bean.ListBeansRequest{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res) != 0 {
		t.Fatalf("expected empty list, got %d items", len(res))
	}
}
