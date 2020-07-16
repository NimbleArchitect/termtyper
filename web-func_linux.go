package main

//os specific close window function
func snipClose() error {
	logDebug("F:snip_close:start")

	w.Dispatch(func() {
		w.Terminate()
	})
	return nil
}
