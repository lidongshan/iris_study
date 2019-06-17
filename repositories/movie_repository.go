package repositories

import (
	"errors"
	"sync"

	"awesomeProject/datamodels"
)

// Query代表一种“访客”和它的查询动作。
type Query func(datamodels.Movie) bool

// MovieRepository会处理一些关于movie实例的基本的操作 。
// 这是一个以测试为目的的接口，即是一个内存中的movie库
// 或是一个连接到数据库的实例。
type MovieRepository interface {
	Exec(query Query, action Query, limit int, mode int) (ok bool)

	Select(query Query) (movie datamodels.Movie, found bool)
	SelectMany(query Query, limit int) (results []datamodels.Movie)

	InsertOrUpdate(movie datamodels.Movie) (updatedMovie datamodels.Movie, err error)
	Delete(query Query, limit int) (deleted bool)
}

// NewMovieRepository返回一个新的基于内存的movie库。
// 库的类型在我们的例子中是唯一的。
func NewMovieRepository(source map[int64]datamodels.Movie) MovieRepository {
	return &movieMemoryRepository{source: source}
}

// movieMemoryRepository就是一个"MovieRepository"
// 它负责存储于内存中的实例数据(map)
type movieMemoryRepository struct {
	source map[int64]datamodels.Movie
	mu     sync.RWMutex
}

const (
	// ReadOnlyMode will RLock(read) the data .
	ReadOnlyMode = iota
	// ReadWriteMode will Lock(read/write) the data.
	ReadWriteMode
)

func (r *movieMemoryRepository) Exec(query Query, action Query, actionLimit int, mode int) (ok bool) {
	loops := 0

	if mode == ReadOnlyMode {
		r.mu.RLock()
		defer r.mu.RUnlock()
	} else {
		r.mu.Lock()
		defer r.mu.Unlock()
	}

	for _, movie := range r.source {
		ok = query(movie)
		if ok {
			if action(movie) {
				loops++
				if actionLimit >= loops {
					break // break
				}
			}
		}
	}

	return
}

// Select方法会收到一个查询方法
// 这个方法给出一个单独的movie实例
// 直到这个功能返回为true时停止迭代。
//
// 它返回最后一次查询成功所找到的结果的值
// 和最后的movie模型
// 以减少caller之间的通信
//
// 这是一个很简单但很聪明的雏形方法
// 我基本在所有会用到的地方使用自从我想到了它
// 也希望你们觉得好用
func (r *movieMemoryRepository) Select(query Query) (movie datamodels.Movie, found bool) {
	found = r.Exec(query, func(m datamodels.Movie) bool {
		movie = m
		return true
	}, 1, ReadOnlyMode)

	// set an empty datamodels.Movie if not found at all.
	if !found {
		movie = datamodels.Movie{}
	}

	return
}

// SelectMany作用相同于Select但是它返回一个切片
// 切片包含一个或多个实例
// 如果传入的参数limit<=0则返回所有
func (r *movieMemoryRepository) SelectMany(query Query, limit int) (results []datamodels.Movie) {
	r.Exec(query, func(m datamodels.Movie) bool {
		results = append(results, m)
		return true
	}, limit, ReadOnlyMode)

	return
}

// InsertOrUpdate添加或者更新一个movie实例到（内存）储存中。
//
// 返回最新操作成功的实例或抛出错误。
func (r *movieMemoryRepository) InsertOrUpdate(movie datamodels.Movie) (datamodels.Movie, error) {
	id := movie.ID

	if id == 0 { // 创建一个新的操作
		var lastID int64
		// 找到最大的ID，避免重复。
		// 在实际使用时您可以使用第三方库去生成
		// 一个string类型的UUID
		r.mu.RLock()
		for _, item := range r.source {
			if item.ID > lastID {
				lastID = item.ID
			}
		}
		r.mu.RUnlock()

		id = lastID + 1
		movie.ID = id

		// map-specific thing
		r.mu.Lock()
		r.source[id] = movie
		r.mu.Unlock()

		return movie, nil
	}

	// 更新操作是基于movie.ID的，
	// 在例子中我们允许了对poster和genre的更新（如果它们非空）。
	// 当然我们可以只是做单纯的数据替换操作:
	// r.source[id] = movie
	// 并注释掉下面的代码;
	current, exists := r.Select(func(m datamodels.Movie) bool {
		return m.ID == id
	})

	if !exists { // 当ID不存在时抛出一个error
		return datamodels.Movie{}, errors.New("failed to update a nonexistent movie")
	}

	// 或者注释下面这段然后用 r.source[id] = m 做单纯替换
	if movie.Poster != "" {
		current.Poster = movie.Poster
	}

	if movie.Genre != "" {
		current.Genre = movie.Genre
	}

	// map-specific thing
	r.mu.Lock()
	r.source[id] = current
	r.mu.Unlock()

	return movie, nil
}

func (r *movieMemoryRepository) Delete(query Query, limit int) bool {
	return r.Exec(query, func(m datamodels.Movie) bool {
		delete(r.source, m.ID)
		return true
	}, limit, ReadWriteMode)
}