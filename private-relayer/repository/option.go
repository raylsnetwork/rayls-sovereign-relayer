package repository

type Option func(*QueryOptions)

type QueryOptions struct {
	Limit int
}

func NewQueryOptions() *QueryOptions {
	return &QueryOptions{
		Limit: 0,
	}
}

func GetQueryOptions(opts ...Option) *QueryOptions {
	queryOptions := NewQueryOptions()
	for _, opt := range opts {
		opt(queryOptions)
	}
	return queryOptions
}

func WithLimit(limit int) Option {
	return func(q *QueryOptions) {
		q.Limit = limit
	}
}
