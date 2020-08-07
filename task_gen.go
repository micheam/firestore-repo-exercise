// THIS FILE IS A GENERATED CODE. DO NOT EDIT
// generated version: 0.4.0
package todolist

import (
	"context"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/go-utils/dedupe"
	"golang.org/x/xerrors"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//go:generate mockgen -source task_gen.go -destination mock/mock_task_gen/mock_task_gen.go

// TaskRepository Repository of Task
type TaskRepository interface {
	// Single
	Get(ctx context.Context, id string, options ...GetOption) (*Task, error)
	GetWithDoc(ctx context.Context, doc *firestore.DocumentRef, options ...GetOption) (*Task, error)
	Insert(ctx context.Context, subject *Task) (string, error)
	Update(ctx context.Context, subject *Task) error
	Delete(ctx context.Context, subject *Task, options ...DeleteOption) error
	DeleteByID(ctx context.Context, id string, options ...DeleteOption) error
	// Multiple
	GetMulti(ctx context.Context, ids []string, options ...GetOption) ([]*Task, error)
	InsertMulti(ctx context.Context, subjects []*Task) ([]string, error)
	UpdateMulti(ctx context.Context, subjects []*Task) error
	DeleteMulti(ctx context.Context, subjects []*Task, options ...DeleteOption) error
	DeleteMultiByIDs(ctx context.Context, ids []string, options ...DeleteOption) error
	// List
	List(ctx context.Context, req *TaskListReq, q *firestore.Query) ([]*Task, error)
	// misc
	GetCollection() *firestore.CollectionRef
	GetCollectionName() string
	GetDocRef(id string) *firestore.DocumentRef
}

// TaskRepositoryMiddleware middleware of TaskRepository
type TaskRepositoryMiddleware interface {
	BeforeInsert(ctx context.Context, subject *Task) (bool, error)
	BeforeUpdate(ctx context.Context, old, subject *Task) (bool, error)
	BeforeDelete(ctx context.Context, subject *Task, options ...DeleteOption) (bool, error)
	BeforeDeleteByID(ctx context.Context, ids []string, options ...DeleteOption) (bool, error)
}

type taskRepository struct {
	collectionName  string
	firestoreClient *firestore.Client
	middleware      []TaskRepositoryMiddleware
}

// NewTaskRepository constructor
func NewTaskRepository(firestoreClient *firestore.Client, middleware ...TaskRepositoryMiddleware) TaskRepository {
	return &taskRepository{
		collectionName:  "Task",
		firestoreClient: firestoreClient,
		middleware:      middleware,
	}
}

func (repo *taskRepository) setMeta(subject *Task, isInsert bool) {
	now := time.Now()

	if isInsert {
		subject.CreatedAt = time.Now()
	}
	subject.UpdatedAt = now
	subject.Version += 1
}

func (repo *taskRepository) beforeInsert(ctx context.Context, subject *Task) error {
	if subject.Version != 0 {
		return xerrors.Errorf("insert data must be Version == 0: %+v", subject)
	}
	if subject.DeletedAt != nil {
		return xerrors.Errorf("insert data must be DeletedAt == nil: %+v", subject)
	}
	repo.setMeta(subject, true)
	for _, m := range repo.middleware {
		c, err := m.BeforeInsert(ctx, subject)
		if err != nil {
			return xerrors.Errorf("beforeInsert.middleware error: %w", err)
		}
		if !c {
			continue
		}
	}
	return nil
}

func (repo *taskRepository) beforeUpdate(ctx context.Context, old, subject *Task) error {
	if old.Version > subject.Version {
		return xerrors.Errorf("The data in the database is newer: (db version: %d, target version: %d) %+v",
			old.Version, subject.Version, subject)
	}
	if subject.DeletedAt != nil {
		return xerrors.Errorf("update data must be DeletedAt == nil: %+v", subject)
	}
	repo.setMeta(subject, false)
	for _, m := range repo.middleware {
		c, err := m.BeforeUpdate(ctx, old, subject)
		if err != nil {
			return xerrors.Errorf("beforeUpdate.middleware error: %w", err)
		}
		if !c {
			continue
		}
	}
	return nil
}

func (repo *taskRepository) beforeDelete(ctx context.Context, subject *Task, options ...DeleteOption) error {
	repo.setMeta(subject, false)
	for _, m := range repo.middleware {
		c, err := m.BeforeDelete(ctx, subject, options...)
		if err != nil {
			return xerrors.Errorf("beforeDelete.middleware error: %w", err)
		}
		if !c {
			continue
		}
	}
	return nil
}

// GetCollection *firestore.CollectionRef getter
func (repo *taskRepository) GetCollection() *firestore.CollectionRef {
	return repo.firestoreClient.Collection(repo.collectionName)
}

// GetCollectionName CollectionName getter
func (repo *taskRepository) GetCollectionName() string {
	return repo.collectionName
}

// GetDocRef *firestore.DocumentRef getter
func (repo *taskRepository) GetDocRef(id string) *firestore.DocumentRef {
	return repo.GetCollection().Doc(id)
}

// TaskListReq List取得時に渡すリクエスト
// └─ bool/int(64)|float64 は stringの独自型で渡す(BoolCriteria | NumericCriteria)
type TaskListReq struct {
	Desc      *QueryChainer
	Done      *QueryChainer
	CreatedAt *QueryChainer
	CreatedBy *QueryChainer
	DeletedAt *QueryChainer
	DeletedBy *QueryChainer
	UpdatedAt *QueryChainer
	UpdatedBy *QueryChainer
	Version   *QueryChainer

	IncludeSoftDeleted bool
}

// List firestore.Queryを使用し条件抽出をする
//  └─ 第3引数はNOT/OR/IN/RANGEなど、より複雑な条件を適用したいときにつける
//      └─ 基本的にnilを渡せば良い
// BUG(54m) 潜在的なバグがあるかもしれない
func (repo *taskRepository) List(ctx context.Context, req *TaskListReq, q *firestore.Query) ([]*Task, error) {
	if (req == nil && q == nil) || (req != nil && q != nil) {
		return nil, xerrors.New("either one should be nil")
	}

	var query firestore.Query
	if q == nil {
		ref := repo.GetCollection()

		if req.CreatedAt != nil {
			for _, chain := range req.CreatedAt.QueryGroup {
				if chain.IsSlice() && chain.Operator == OpTypeIn {
					dedupe.Do(&chain.Value)
					ref.Query = ref.Query.Where("CreatedAt", chain.Operator, chain.Value)
					continue
				}
				ref.Query = ref.Query.Where("CreatedAt", chain.Operator, chain.Value)
			}
		}
		if req.CreatedBy != nil {
			for _, chain := range req.CreatedBy.QueryGroup {
				if chain.IsSlice() && chain.Operator == OpTypeIn {
					dedupe.Do(&chain.Value)
					ref.Query = ref.Query.Where("CreatedBy", chain.Operator, chain.Value)
					continue
				}
				ref.Query = ref.Query.Where("CreatedBy", chain.Operator, chain.Value)
			}
		}
		if req.DeletedAt != nil {
			for _, chain := range req.DeletedAt.QueryGroup {
				if chain.IsSlice() && chain.Operator == OpTypeIn {
					dedupe.Do(&chain.Value)
					ref.Query = ref.Query.Where("DeletedAt", chain.Operator, chain.Value)
					continue
				}
				ref.Query = ref.Query.Where("DeletedAt", chain.Operator, chain.Value)
			}
		}
		if req.DeletedBy != nil {
			for _, chain := range req.DeletedBy.QueryGroup {
				if chain.IsSlice() && chain.Operator == OpTypeIn {
					dedupe.Do(&chain.Value)
					ref.Query = ref.Query.Where("DeletedBy", chain.Operator, chain.Value)
					continue
				}
				ref.Query = ref.Query.Where("DeletedBy", chain.Operator, chain.Value)
			}
		}
		if req.UpdatedAt != nil {
			for _, chain := range req.UpdatedAt.QueryGroup {
				if chain.IsSlice() && chain.Operator == OpTypeIn {
					dedupe.Do(&chain.Value)
					ref.Query = ref.Query.Where("UpdatedAt", chain.Operator, chain.Value)
					continue
				}
				ref.Query = ref.Query.Where("UpdatedAt", chain.Operator, chain.Value)
			}
		}
		if req.UpdatedBy != nil {
			for _, chain := range req.UpdatedBy.QueryGroup {
				if chain.IsSlice() && chain.Operator == OpTypeIn {
					dedupe.Do(&chain.Value)
					ref.Query = ref.Query.Where("UpdatedBy", chain.Operator, chain.Value)
					continue
				}
				ref.Query = ref.Query.Where("UpdatedBy", chain.Operator, chain.Value)
			}
		}
		if req.Version != nil {
			for _, chain := range req.Version.QueryGroup {
				if chain.IsSlice() && chain.Operator == OpTypeIn {
					dedupe.Do(&chain.Value)
					ref.Query = ref.Query.Where("Version", chain.Operator, chain.Value)
					continue
				}
				ref.Query = ref.Query.Where("Version", chain.Operator, chain.Value)
			}
		}
		if req.Desc != nil {
			for _, chain := range req.Desc.QueryGroup {
				if chain.IsSlice() && chain.Operator == OpTypeIn {
					dedupe.Do(&chain.Value)
					ref.Query = ref.Query.Where("description", chain.Operator, chain.Value)
					continue
				}
				ref.Query = ref.Query.Where("description", chain.Operator, chain.Value)
			}
		}
		if req.Done != nil {
			for _, chain := range req.Done.QueryGroup {
				if chain.IsSlice() && chain.Operator == OpTypeIn {
					dedupe.Do(&chain.Value)
					ref.Query = ref.Query.Where("done", chain.Operator, chain.Value)
					continue
				}
				ref.Query = ref.Query.Where("done", chain.Operator, chain.Value)
			}
		}

		if !req.IncludeSoftDeleted {
			ref.Query = ref.Query.Where("DeletedAt", OpTypeEqual, nil)
		}

		query = ref.Query
	} else {
		query = *q
	}

	subjects := make([]*Task, 0)
	iter := query.Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, xerrors.Errorf("error in Next method: %w", err)
		}
		subject := new(Task)
		if err := doc.DataTo(&subject); err != nil {
			return nil, xerrors.Errorf("error in DataTo method: %w", err)
		}
		subject.ID = doc.Ref.ID
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

// Get 処理中の Task の取得処理一切の責任を持ち、これを行う
func (repo *taskRepository) Get(ctx context.Context, id string, options ...GetOption) (*Task, error) {
	doc := repo.GetCollection().Doc(id)
	return repo.GetWithDoc(ctx, doc, options...)
}

// GetWithDoc 処理中の Task の取得処理一切の責任を持ち、これを行う
func (repo *taskRepository) GetWithDoc(ctx context.Context, doc *firestore.DocumentRef, options ...GetOption) (*Task, error) {
	snapShot, err := doc.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, ErrNotFound
		}
		return nil, xerrors.Errorf("error in Get method: %w", err)
	}

	subject := new(Task)
	if err := snapShot.DataTo(&subject); err != nil {
		return nil, xerrors.Errorf("error in DataTo method: %w", err)
	}

	if len(options) == 0 || !options[0].IncludeSoftDeleted {
		if subject.DeletedAt != nil {
			return nil, ErrAlreadyDeleted
		}
	}
	subject.ID = snapShot.Ref.ID

	return subject, nil
}

// Insert 処理中の Task の登録処理一切の責任を持ち、これを行う
func (repo *taskRepository) Insert(ctx context.Context, subject *Task) (string, error) {
	err := repo.beforeInsert(ctx, subject)
	if err != nil {
		return "", xerrors.Errorf("before insert error: %w", err)
	}
	ref := repo.GetCollection().Doc(subject.ID)

	if _, err := ref.Get(ctx); err == nil {
		return "", ErrAlreadyExists
	}

	if _, err := ref.Set(ctx, subject); err != nil {
		return "", xerrors.Errorf("error in Set method: %w", err)
	}

	return ref.ID, nil
}

// Update 処理中の Task の更新処理一切の責任を持ち、これを行う
func (repo *taskRepository) Update(ctx context.Context, subject *Task) error {
	ref := repo.GetCollection().Doc(subject.ID)

	snapShot, err := ref.Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return ErrNotFound
		}
		return xerrors.Errorf("error in Get method: %w", err)
	}
	old := new(Task)
	if err := snapShot.DataTo(&old); err != nil {
		return xerrors.Errorf("error in DataTo method: %w", err)
	}

	err = repo.beforeUpdate(ctx, old, subject)
	if err != nil {
		return xerrors.Errorf("before update error: %w", err)
	}
	if _, err := ref.Set(ctx, subject); err != nil {
		return xerrors.Errorf("error in Set method: %w", err)
	}

	return nil
}

// Delete 処理中の Task の削除処理一切の責任を持ち、これを行う
func (repo *taskRepository) Delete(ctx context.Context, subject *Task, options ...DeleteOption) error {
	ref := repo.GetCollection().Doc(subject.ID)
	if _, err := ref.Get(ctx); err != nil {
		if status.Code(err) == codes.NotFound {
			return ErrNotFound
		}
		return xerrors.Errorf("error in Get method: %w", err)
	}

	if err := repo.beforeDelete(ctx, subject, options...); err != nil {
		return xerrors.Errorf("before delete error: %w", err)
	}

	if len(options) > 0 && options[0].Mode == DeleteModeSoft {
		t := time.Now()
		subject.DeletedAt = &t
		if _, err := ref.Set(ctx, subject); err != nil {
			return xerrors.Errorf("error in Set method: %w", err)
		}
	} else {
		if _, err := ref.Delete(ctx); err != nil {
			return xerrors.Errorf("error in Delete method: %w", err)
		}
	}
	return nil
}

// DeleteByID 処理中の Task のIDから削除処理一切の責任を持ち、これを行う
func (repo *taskRepository) DeleteByID(ctx context.Context, id string, options ...DeleteOption) error {
	if err := repo.Delete(ctx, &Task{ID: id}, options...); err != nil {
		return xerrors.Errorf("error in repo.Delete method: %w", err)
	}
	return nil
}

// GetMulti 処理中の Task の一括取得処理一切の責任を持ち、これを行う
func (repo *taskRepository) GetMulti(ctx context.Context, ids []string, options ...GetOption) ([]*Task, error) {
	collect := repo.GetCollection()
	docRefs := make([]*firestore.DocumentRef, 0, len(ids))
	for _, id := range ids {
		ref := collect.Doc(id)
		docRefs = append(docRefs, ref)
	}

	snapShots, err := repo.firestoreClient.GetAll(ctx, docRefs)
	if err != nil {
		return nil, xerrors.Errorf("error in GetAll method: %w", err)
	}

	subjects := make([]*Task, 0, len(ids))
	for _, snapShot := range snapShots {
		subject := new(Task)
		if err := snapShot.DataTo(&subject); err != nil {
			return nil, xerrors.Errorf("error in DataTo method: %w", err)
		}

		if len(options) == 0 || !options[0].IncludeSoftDeleted {
			if subject.DeletedAt != nil {
				continue
			}
		}
		subject.ID = snapShot.Ref.ID
		subjects = append(subjects, subject)
	}

	return subjects, nil
}

// InsertMulti 処理中の Task の一括挿入処理一切の責任を持ち、これを行う
func (repo *taskRepository) InsertMulti(ctx context.Context, subjects []*Task) ([]string, error) {
	ids := make([]string, 0, len(subjects))
	batches := make([]*firestore.WriteBatch, 0)
	batch := repo.firestoreClient.Batch()
	collect := repo.GetCollection()

	for i, subject := range subjects {
		err := repo.beforeInsert(ctx, subject)
		if err != nil {
			return nil, xerrors.Errorf("before insert error: %w", err)
		}
		var ref *firestore.DocumentRef
		if subject.ID == "" {
			ref = collect.NewDoc()
			subject.ID = ref.ID
		} else {
			ref = collect.Doc(subject.ID)
			if s, err := ref.Get(ctx); err == nil {
				return nil, xerrors.Errorf("already exists [%v]: %#v", subject.ID, s)
			}
		}
		batch.Set(ref, subject)
		ids = append(ids, ref.ID)
		i++
		if (i%500) == 0 && len(subjects) != i {
			batches = append(batches, batch)
			batch = repo.firestoreClient.Batch()
		}
	}
	batches = append(batches, batch)

	for _, b := range batches {
		if _, err := b.Commit(ctx); err != nil {
			return nil, xerrors.Errorf("error in Commit method: %w", err)
		}
	}

	return ids, nil
}

// UpdateMulti 処理中の Task の一括更新処理一切の責任を持ち、これを行う
func (repo *taskRepository) UpdateMulti(ctx context.Context, subjects []*Task) error {
	batches := make([]*firestore.WriteBatch, 0)
	batch := repo.firestoreClient.Batch()
	collect := repo.GetCollection()

	for i, subject := range subjects {
		ref := collect.Doc(subject.ID)
		snapShot, err := ref.Get(ctx)
		if err != nil {
			if status.Code(err) == codes.NotFound {
				return xerrors.Errorf("not found [%v]: %w", subject.ID, err)
			}
			return xerrors.Errorf("error in Get method [%v]: %w", subject.ID, err)
		}
		old := new(Task)
		if err := snapShot.DataTo(&old); err != nil {
			return xerrors.Errorf("error in DataTo method: %w", err)
		}

		err = repo.beforeUpdate(ctx, old, subject)
		if err != nil {
			return xerrors.Errorf("before update error: %w", err)
		}

		batch.Set(ref, subject)
		i++
		if (i%500) == 0 && len(subjects) != i {
			batches = append(batches, batch)
			batch = repo.firestoreClient.Batch()
		}
	}
	batches = append(batches, batch)

	for _, b := range batches {
		if _, err := b.Commit(ctx); err != nil {
			return xerrors.Errorf("error in Commit method: %w", err)
		}
	}

	return nil
}

// DeleteMulti 処理中の Task の一括削除処理一切の責任を持ち、これを行う
func (repo *taskRepository) DeleteMulti(ctx context.Context, subjects []*Task, options ...DeleteOption) error {
	batches := make([]*firestore.WriteBatch, 0)
	batch := repo.firestoreClient.Batch()
	collect := repo.GetCollection()

	for i, subject := range subjects {
		ref := collect.Doc(subject.ID)
		if _, err := ref.Get(ctx); err != nil {
			if status.Code(err) == codes.NotFound {
				return xerrors.Errorf("not found [%v]: %w", subject.ID, err)
			}
			return xerrors.Errorf("error in Get method [%v]: %w", subject.ID, err)
		}

		if err := repo.beforeDelete(ctx, subject, options...); err != nil {
			return xerrors.Errorf("before delete error: %w", err)
		}

		if len(options) > 0 && options[0].Mode == DeleteModeSoft {
			t := time.Now()
			subject.DeletedAt = &t
			batch.Set(ref, subject)
		} else {
			batch.Delete(ref)
		}
		i++
		if (i%500) == 0 && len(subjects) != i {
			batches = append(batches, batch)
			batch = repo.firestoreClient.Batch()
		}
	}
	batches = append(batches, batch)

	for _, b := range batches {
		if _, err := b.Commit(ctx); err != nil {
			return xerrors.Errorf("error in Commit method: %w", err)
		}
	}

	return nil
}

// DeleteMultiByIDs 処理中の Task のID群を元に一括削除処理一切の責任を持ち、これを行う
func (repo *taskRepository) DeleteMultiByIDs(ctx context.Context, ids []string, options ...DeleteOption) error {
	subjects := make([]*Task, 0, len(ids))

	for _, id := range ids {
		subjects = append(subjects, &Task{ID: id})
	}

	return repo.DeleteMulti(ctx, subjects, options...)
}
