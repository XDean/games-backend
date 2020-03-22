package inject

type (
	Context interface {
		Register(obj interface{})
		Get(rec interface{}) error
		GetAll(recs interface{}) error
		Inject(obj interface{}) error
		Refresh() error
	}
)

func NewContext() Context {
	return &impl{}
}
