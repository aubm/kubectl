package test

// ControlPlane is a struct that knows how to start your test control plane.
//
// Right now, that means Etcd and your APIServer. This is likely to increase in future.
type ControlPlane struct {
	APIServer ControlPlaneProcess
}

// ControlPlaneProcess knows how to start and stop a ControlPlane process.
// This interface is potentially going to be expanded to e.g. allow access to the processes StdOut/StdErr
// and other internals.
type ControlPlaneProcess interface {
	Start() error
	Stop() error
	URL() (string, error)
}

//go:generate counterfeiter . ControlPlaneProcess

// NewControlPlane will give you a ControlPlane struct that's properly wired together.
func NewControlPlane() (*ControlPlane, error) {
	apiServer, err := NewAPIServer()
	if err != nil {
		return nil, err
	}

	controlPlane := &ControlPlane{
		APIServer: apiServer,
	}

	return controlPlane, nil
}

// Start will start your control plane. To stop it, call Stop().
func (f *ControlPlane) Start() error {
	started := make(chan error)
	starter := func(process ControlPlaneProcess) {
		started <- process.Start()
	}
	processes := []ControlPlaneProcess{
		f.APIServer,
	}

	for _, process := range processes {
		go starter(process)
	}

	for range processes {
		if err := <-started; err != nil {
			return err
		}
	}

	return nil
}

// Stop will stop your control plane, and clean up their data.
func (f *ControlPlane) Stop() error {
	return f.APIServer.Stop()
}

// APIServerURL returns the URL to the APIServer. Clients can use this URL to connect to the APIServer.
func (f *ControlPlane) APIServerURL() (string, error) {
	return f.APIServer.URL()
}
