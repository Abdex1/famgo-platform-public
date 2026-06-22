# 🧪 TESTING PATTERNS
## Extracted from uber-master

**Status:** Pattern 7/8

---

## Mock Pattern

```go
type MockService struct {
    GetByIDFunc func(ctx context.Context, id string) (*Entity, error)
    CreateFunc  func(ctx context.Context, entity *Entity) error
}

func (m *MockService) GetByID(ctx context.Context, id string) (*Entity, error) {
    return m.GetByIDFunc(ctx, id)
}

func (m *MockService) Create(ctx context.Context, entity *Entity) error {
    return m.CreateFunc(ctx, entity)
}
```

## Table-Driven Tests

```go
func TestGetDriver(t *testing.T) {
    tests := []struct {
        name    string
        id      string
        setup   func(*MockRepo)
        want    *Driver
        wantErr bool
    }{
        {
            name: "valid driver",
            id:   "driver-1",
            setup: func(m *MockRepo) {
                m.GetByIDFunc = func(ctx context.Context, id string) (*Driver, error) {
                    return &Driver{ID: "driver-1", Name: "John"}, nil
                }
            },
            want: &Driver{ID: "driver-1", Name: "John"},
        },
        {
            name:    "driver not found",
            id:      "invalid",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            repo := &MockRepo{}
            if tt.setup != nil {
                tt.setup(repo)
            }
            
            svc := NewService(repo)
            got, err := svc.GetDriver(context.Background(), tt.id)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("got error %v, want %v", err, tt.wantErr)
            }
            
            if !reflect.DeepEqual(got, tt.want) {
                t.Errorf("got %+v, want %+v", got, tt.want)
            }
        })
    }
}
```

## Integration Test Template

```go
func TestCreateDriverIntegration(t *testing.T) {
    db := setupTestDB(t)
    defer db.Close()
    
    repo := repository.New(db)
    svc := service.New(repo)
    
    driver := &Driver{Name: "Jane", Email: "jane@example.com"}
    err := svc.Create(context.Background(), driver)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, driver.ID)
    
    retrieved, _ := repo.GetByID(context.Background(), driver.ID)
    assert.Equal(t, "Jane", retrieved.Name)
}
```

**Pattern 7 Status:** READY FOR USE

---
