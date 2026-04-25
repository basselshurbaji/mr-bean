package bean_test

import (
	"context"
	"testing"
	"time"

	"github.com/basselshurbaji/mr_bean/backend/internal/bean"
)

type mockBeanRepo struct {
	item *bean.Bean
	list []bean.Bean
	err  error
}

func (m *mockBeanRepo) ListBeans(_ context.Context, _ string) ([]bean.Bean, error) {
	return m.list, m.err
}
func (m *mockBeanRepo) CreateBean(_ context.Context, _ string, _ bean.BeanParams) (*bean.Bean, error) {
	return m.item, m.err
}
func (m *mockBeanRepo) UpdateBean(_ context.Context, _, _ string, _ bean.BeanParams) (*bean.Bean, error) {
	return m.item, m.err
}
func (m *mockBeanRepo) DeleteBean(_ context.Context, _, _ string) error { return m.err }

func makeBean(id string) bean.Bean {
	return bean.Bean{
		ID: id, UserID: "user-1", Name: "Ethiopia Yirgacheffe",
		CreatedAt: time.Now(), UpdatedAt: time.Now(),
	}
}

func TestBeanService_ListBeans_ReturnsList(t *testing.T) {
	items := []bean.Bean{makeBean("b1"), makeBean("b2")}
	svc := bean.NewBeanService(&mockBeanRepo{list: items})

	got, err := svc.ListBeans(context.Background(), "user-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 items, got %d", len(got))
	}
}

func TestBeanService_CreateBean_ReturnsBean(t *testing.T) {
	b := makeBean("b1")
	svc := bean.NewBeanService(&mockBeanRepo{item: &b})

	got, err := svc.CreateBean(context.Background(), "user-1", bean.BeanParams{Name: "Ethiopia Yirgacheffe"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.ID != "b1" {
		t.Fatalf("expected id %q, got %q", "b1", got.ID)
	}
}

func TestBeanService_UpdateBean_NotFound(t *testing.T) {
	svc := bean.NewBeanService(&mockBeanRepo{err: bean.ErrNotFound})

	_, err := svc.UpdateBean(context.Background(), "missing", "user-1", bean.BeanParams{Name: "X"})
	if err != bean.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestBeanService_DeleteBean_NotFound(t *testing.T) {
	svc := bean.NewBeanService(&mockBeanRepo{err: bean.ErrNotFound})

	err := svc.DeleteBean(context.Background(), "missing", "user-1")
	if err != bean.ErrNotFound {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestBeanService_DeleteBean_Success(t *testing.T) {
	svc := bean.NewBeanService(&mockBeanRepo{})

	if err := svc.DeleteBean(context.Background(), "b1", "user-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
