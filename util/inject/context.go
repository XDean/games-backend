package inject

type (
	Context interface {
		Register(obj interface{}) error
		Get(rec interface{}) error
		GetAll(recs interface{}) error
		Inject(obj interface{}) error
	}
)
