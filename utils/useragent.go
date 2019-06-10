package utils

import (
	// "fmt"
	"github.com/ua-parser/uap-go/uaparser"
	"regexp"
	"strings"
)

/**
 * [UaParser 解析器初始化]
 * @type {[type]}
 */
var (
	UaParser  = uaparser.NewFromSaved()
	Devices   = []string{"IPHONE", "IPAD", "ANDROID", "TABLET"}
	TabletRgx = regexp.MustCompile("Mobile")
)

/**
 * user agent 解析结果
 */
type DInfoFUa struct {
	Device string
	Os     string
	Osv    string
}

/**
 * 解析 user-agent
 */
func UserAgentParse(ua string) DInfoFUa {
	info := DInfoFUa{}

	client := UaParser.Parse(ua)

	// fmt.Println("user agent.", ua)
	// fmt.Printf("ua parse: %+v\n, %+v\n, %+v\n", *client.Device, *client.UserAgent, *client.Os)

	info.Os = strings.ToUpper(client.Os.Family) //Android IOS
	info.Osv = strings.Join([]string{client.Os.Major, client.Os.Minor, client.Os.Patch}, ".")
	info.Device = strings.ToUpper(client.Device.Family)

	if "ANDROID" == info.Os {
		if TabletRgx.MatchString(ua) {
			info.Device = "ANDROID"
		} else {
			info.Device = "TABLET"
		}
	}

	return info
}
