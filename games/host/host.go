package host

func (h Host) Init() error {
	return h.Inject.Refresh()
}

func (h Host) Plug(plugin Plugin) Host {
	h.Inject.Register(plugin)
	h.Handler = plugin.Plug(h.Handler)
	return h
}

func (h Host) Handle(ctx Context) error {
	return h.Handler.Handle(ctx)
}
