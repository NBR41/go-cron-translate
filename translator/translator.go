package translator

import (
	"fmt"
	"regexp"
	"strings"
)

var (
	dOTW = map[string]string{
		"0": "Sunday",
		"1": "Monday",
		"2": "Tuesday",
		"3": "Wednesday",
		"4": "Thursday",
		"5": "Friday",
		"6": "Saturday",
		"7": "Sunday",
	}

	monthNames = map[string]string{
		"1":  "January",
		"2":  "February",
		"3":  "Marsh",
		"4":  "Aprim",
		"5":  "May",
		"6":  "June",
		"7":  "July",
		"8":  "August",
		"9":  "September",
		"10": "October",
		"11": "November",
		"12": "December",
	}
)

const (
	modeEvery  = iota // 0 *
	modeEveryN        // 1 */a
	modeAtN           // 2 a
	modeList          // 3 a,b
	modeRange         // 4 a-b
)

const (
	mm = iota
	hh
	dd
	mmm
	ddd
)

const (
	reg             = `^(?:([0-9]{1,2}-[0-9]{1,2}|\*)(/[0-9]{1,2})*)$|^([0-9]{1,2}((?:,[0-9]{1,2})*))$`
	fmtInvalidParam = "invalid crontab string [%s]"
	fmtInvalidDOTW  = "invalid day of the week value [%s]"
	fmtInvalidMonth = "invalid month value [%s]"
	fmtHour         = "%02sh%02s"
	fmtEvery        = "every %s"
	fmtEveryN       = "every %s minutes of %sh"
	fmtRange        = "from %s to %s"
)

type translator func([][]string, []int) (string, error)

var transFuncs = []translator{
	getMinuteHourTranslation,
	getDOTWTranslation,
	getDayTranslation,
	getMonthTranslation,
}

// GetTranslation return a translation for a crontab value
func GetTranslation(v string) (string, error) {
	var parts = strings.Split(v, " ")
	if len(parts) != 5 {
		return "", fmt.Errorf(fmtInvalidParam, v)
	}

	var reparts = make([][]string, 5, 5)
	var modes = make([]int, 5, 5)
	re := regexp.MustCompile(reg)
	for i := range parts {
		reparts[i] = re.FindStringSubmatch(parts[i])
		if len(reparts[i]) == 0 {
			return "", fmt.Errorf(fmtInvalidParam, v)
		}
		modes[i] = getMode(reparts[i])
	}

	return translate(reparts, modes)
}

func getMode(vals []string) int {
	switch vals[1] {
	case "*":
		if vals[2] == "" {
			return modeEvery
		}
		return modeEveryN
	case "":
		if vals[4] == "" {
			return modeAtN
		}
		return modeList
	default:
		return modeRange
	}
}

func translate(reparts [][]string, modes []int) (string, error) {
	var val string
	var ret []string
	var err error
	for i := range transFuncs {
		val, err = transFuncs[i](reparts, modes)
		if err != nil {
			return "", err
		}
		if val != "" {
			ret = append(ret, val)
		}
	}
	return strings.Join(ret, " "), nil
}

func getMinuteHourTranslation(reparts [][]string, modes []int) (string, error) {
	switch modes[hh] {
	case modeAtN:
		switch modes[mm] {
		case modeAtN:
			return "at " + fmt.Sprintf(fmtHour, reparts[hh][3], reparts[mm][3]), nil
		case modeEveryN:
			return fmt.Sprintf(fmtEveryN, strings.TrimPrefix(reparts[mm][2], "/"), reparts[hh][3]), nil
		}

	case modeEvery:
		switch modes[mm] {
		case modeEvery:
			return fmt.Sprintf(fmtEvery, "minutes"), nil
		case modeAtN:
			fallthrough
		case modeList:
			return fmt.Sprintf(fmtEvery, "hour past "+reparts[mm][3]+" minutes"), nil
		case modeEveryN:
			return fmt.Sprintf(fmtEvery, strings.TrimPrefix(reparts[mm][2], "/")+" minutes"), nil
		}
	}
	return "", nil
}

func getDayTranslation(reparts [][]string, modes []int) (string, error) {
	var pref string
	if modes[ddd] != modeEvery {
		pref = "and "
	}

	switch modes[dd] {
	case modeEvery:
		if modes[ddd] != modeEvery || modes[hh] == modeEvery {
			return "", nil
		}
		return pref + fmt.Sprintf(fmtEvery, "day"), nil
	case modeEveryN:
		return pref + fmt.Sprintf(fmtEvery, strings.TrimPrefix(reparts[dd][2], "/")+" days"), nil
	case modeAtN:
		fallthrough
	case modeList:
		return pref + fmt.Sprintf(fmtEvery, reparts[dd][3]), nil
	case modeRange:
		parts := strings.Split(reparts[dd][1], "-")
		return pref + fmt.Sprintf(fmtRange, parts[0], parts[1]), nil
	default:
		return "", nil
	}
}

func getMonthTranslation(reparts [][]string, modes []int) (string, error) {
	switch modes[mmm] {
	case modeEvery:
		if modes[dd] == modeEvery {
			return "", nil
		}
		return "of the month", nil
	case modeEveryN:
		return fmt.Sprintf(fmtEvery, strings.TrimPrefix(reparts[mmm][2], "/")+" months"), nil
	case modeAtN:
		fallthrough
	case modeList:
		var parts = strings.Split(reparts[mmm][3], ",")
		var ret = make([]string, len(parts))
		for i := range parts {
			if _, present := monthNames[parts[i]]; !present {
				return "", fmt.Errorf(fmtInvalidMonth, reparts[mmm][3])
			}
			ret[i] = monthNames[parts[i]]
		}
		return "of " + strings.Join(ret, ","), nil
	case modeRange:
		parts := strings.Split(reparts[mmm][1], "-")
		var ret = make([]string, len(parts))
		for i := range parts {
			if _, present := monthNames[parts[i]]; !present {
				return "", fmt.Errorf(fmtInvalidMonth, parts[i])
			}
			ret[i] = monthNames[parts[i]]
		}
		return fmt.Sprintf(fmtRange, ret[0], ret[1]), nil
	default:
		return "", nil
	}
}

func getDOTWTranslation(reparts [][]string, modes []int) (string, error) {
	switch modes[ddd] {
	case modeEvery:
		return "", nil
	case modeRange:
		parts := strings.Split(reparts[ddd][1], "-")
		var ret = make([]string, len(parts))
		for i := range parts {

			if _, present := dOTW[parts[i]]; !present {
				return "", fmt.Errorf(fmtInvalidDOTW, parts[i])
			}
			ret[i] = dOTW[parts[i]]
		}
		return fmt.Sprintf(fmtRange, ret[0], ret[1]), nil
	case modeAtN:
		fallthrough
	case modeList:
		var parts = strings.Split(reparts[ddd][3], ",")
		var ret = make([]string, len(parts))
		for i := range parts {
			if _, present := dOTW[parts[i]]; !present {
				return "", fmt.Errorf(fmtInvalidDOTW, parts[i])
			}
			ret[i] = dOTW[parts[i]]
		}
		return fmt.Sprintf(fmtEvery, strings.Join(ret, ",")), nil
	// TODO missing modes
	default:
		return "", nil
	}
}

/*
0 5 5 {0:0 1: 2: 3:0 4: } 2
* 5 5 {0:* 1:* 2: 3: 4: } 0
5,10,25 5 5 {0:5,10,25 1: 2: 3:5,10,25 4:,10,25 } 3
2-5 5 5 {0:2-5 1:2-5 2: 3: 4: } 4
*\/2 5 5 {0:*\/2 1:* 2:/2 3: 4: } 1
22 5 5 {0:22 1: 2: 3:22 4: } 2
*/
