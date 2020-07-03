package main

//os specific close window function
func snip_close() error {
	logDebug("F:snip_close:start")

	w.Terminate()
	return nil
}
