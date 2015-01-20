package main

import (
	"fmt"
	"github.com/sap/gorfc/gorfc"
	"time"
)

// BAPI calls log
func printLog(bapi_return interface{}) {
	for _, line := range bapi_return.([]interface{}) {
		fmt.Printf("%s: %s\n",
			line.(map[string]interface{})["TYPE"],
			line.(map[string]interface{})["MESSAGE"])
	}
}

func main() {
	// Connect to ABAP system
	SAPROUTER := "/H/123.12.123.12/E/yt6ntx/H/123.14.131.111/H/"

	EC4 := gorfc.ConnectionParameter{
		User:      "abapuser",
		Passwd:    "abappass",
		Ashost:    "10.11.12.13",
		Saprouter: SAPROUTER,
		Sysnr:     "200",
		Client:    "300",
		Trace:     "3",
		Lang:      "EN",
	}

	c, _ := gorfc.Connection(EC4)

	// The source user, to be copied
	unameFrom := "UNAMEFROM"

	// Defaults if source user validity not maintained (undefined)
	validFrom := time.Date(2015, time.January, 19, 0, 0, 0, 0, time.UTC)
	validTo := time.Date(2015, time.December, 31, 0, 0, 0, 0, time.UTC)

	// New users" password. For automatic generation check CREATE BAPI
	initpwd := "InitPa$$21"

	// Users to be created
	users := []string{"UNAMETO1", "UNAMETO2"}

	// Get source user details
	r, _ := c.Call("BAPI_USER_GET_DETAIL", map[string]interface{}{"USERNAME": unameFrom, "CACHE_RESULTS": " "})

	// Set new users" defaults
	logonData := r["LOGONDATA"].(map[string]interface{})

	if logonData["GLTGV"] == nil {
		logonData["GLTGV"] = validFrom
	}
	if logonData["GLTGB"] == nil {
		logonData["GLTGB"] = validTo
	}
	password := map[string]string{"BAPIPWD": initpwd}

	// Create new users
	address := r["ADDRESS"].(map[string]interface{})

	for _, unameTo := range users {
		fmt.Println(unameTo)

		address["LASTNAME"] = unameTo
		address["FULLNAME"] = unameTo

		x, _ := c.Call("BAPI_USER_CREATE1", map[string]interface{}{
			"USERNAME":  unameTo,
			"LOGONDATA": logonData,
			"PASSWORD":  password,
			"DEFAULTS":  r["DEFAULTS"],
			"ADDRESS":   address,
			"COMPANY":   r["COMPANY"],
			"REF_USER":  r["REF_USER"],
			"PARAMETER": r["PARAMETER"],
			"GROUPS":    r["GROUPS"],
		})

		printLog(x["RETURN"])

		x, _ = c.Call("BAPI_USER_PROFILES_ASSIGN", map[string]interface{}{
			"USERNAME": unameTo,
			"PROFILES": r["PROFILES"],
		})

		printLog(x["RETURN"])

		x, _ = c.Call("BAPI_USER_ACTGROUPS_ASSIGN", map[string]interface{}{
			"USERNAME":       unameTo,
			"ACTIVITYGROUPS": r["ACTIVITYGROUPS"],
		})

		printLog(x["RETURN"])
	}

	// Finished
	fmt.Printf("%s copied to %d new users.\nBye!", unameFrom, len(users))
}
