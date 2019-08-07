package kubernetes

// ClientAsync is a wrapper around the client-go package for Kubernetes
type ClientAsync interface {
	Sync() Client
}

// ClientAsyncImpl is the interface implementation of ClientAsync
type ClientAsyncImpl struct {
	syncClient Client
}

// MakeClient returns a ClientAsync
func MakeClient() (ClientAsync, error) {
	syncClient, err := makeClient()
	if err == nil {
		return ClientAsyncImpl{syncClient}, nil
	}
	return nil, err
}

// MakeFromClient returns a ClientAsync from the given sync client
func MakeFromClient(syncClient Client) ClientAsync {
	return ClientAsyncImpl{syncClient}
}

// Sync returns the synchronous client
func (c ClientAsyncImpl) Sync() Client {
	return c.syncClient
}
