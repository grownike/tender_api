package handlers

type handler struct{
	storage storage
}

func New(storage storage) *handler{
	return &handler{
		storage: storage,
	}
}
