/**
https://github.com/openatx/uiautomator2#app-management
*/
package uiautomator

/*
Install an app
TODO: api "/install" not work
*/
func (ua *UIAutomator) AppInstall(url string) error {
	return nil
}

/*
Launch app
*/
func (ua *UIAutomator) AppStart(packageName string) error {
	_, err := ua.Shell(
		[]string{
			"monkey", "-p", packageName, "-c",
			"android.intent.category.LAUNCHER", "1",
		},
		10,
	)

	return err
}

/*
Stop app
*/
func (ua *UIAutomator) AppStop(packageName string) error {
	_, err := ua.Shell(
		[]string{
			"am", "force-stop", packageName,
		},
		10,
	)

	return err
}
