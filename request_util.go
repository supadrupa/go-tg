package tg

// helper function for add message identity to request
func addOptMessageIdentityToRequest(r *Request, k string, mi MessageIdentity) *Request {
	if mi != nil {
		r.AddInt(k, int(mi.GetMessageID()))
	}

	return r
}

// helper function for add ReplyMarkup to request
func addOptReplyMarkupToRequest(r *Request, k string, rm ReplyMarkup) (*Request, error) {
	if rm != nil {
		str, err := rm.EncodeReplyMarkup()
		if err != nil {
			return r, err
		}
		r.AddString(k, str)
		return r, nil
	}

	return r, nil
}
