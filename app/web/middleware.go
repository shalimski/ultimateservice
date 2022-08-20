package web

type Middleware func(Handler) Handler

func wrap(mw []Middleware, next Handler) Handler {
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			next = h(next)
		}
	}

	return next
}
