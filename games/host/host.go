package host

func (h Host) Plug(plugin Plugin) Host {
	h.Handler = plugin.Plug(h.Handler)
	return h
}

func (h Host) Handle(ctx Context) error {
	return h.Handler.Handle(ctx)
}
