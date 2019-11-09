package manager

import (
	"fmt"
)

// get the property for app as a string, if property does not exist return err
func (ms *ManagerService) GetAppProperty(app string, property string) (string, error) {
	if ms.Config.Get(app+"."+property) != nil {
		return ms.Config.Get(app + "." + property).(string), nil
	} else {
		return "", fmt.Errorf("Configuration missing '%s' section under [%s] heading.\n", property, app)
	}
}
