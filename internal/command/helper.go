package command

import (
	"strings"
)

// func escapeArgs(args []string) []string {
// 	var processedArgs []string

// 	for _, arg := range args {
// 		_, errInt := strconv.ParseInt(arg, 10, 64)
// 		_, errFloat := strconv.ParseFloat(arg, 64)
// 		_, errBool := strconv.ParseBool(arg)

// 		if errInt == nil || errFloat == nil || errBool == nil {

// 			processedArgs = append(processedArgs, arg)
// 		} else {

// 			processedArgs = append(processedArgs, "\""+arg+"\"")
// 		}
// 	}

// 	return processedArgs
// }

func parseBaseCommand(baseCommand string) (string, string) {
	x := strings.Split(baseCommand, ":")
	if len(x) > 1 {
		return x[0], x[1]
	}
	return x[0], "latest"
}
