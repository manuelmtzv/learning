package main

func (app *application) logError(err error) {
	app.logger.Error(err)
}
