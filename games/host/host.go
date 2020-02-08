package host

func (h Host) Plug(plugin Plugin) Host {
	h.handler = plugin.Plug(h.handler)
	return h
}

func (h Host) Handle(ctx Context) {
	h.handler.Handle(ctx)
}
