package main

func (app *application) run() {
	app.logger.Infof("Running app with %v workers", app.processor.workers)
}
