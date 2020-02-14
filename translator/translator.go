// Package translator provide a function to translate crontab in natural language
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
		"4":  "April",
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

	xhour = "XX"

	fmtEveryNMonths = "every %s months"
	everyDay        = "every day"
	fmtEveryNDays   = "every %s days"

	everyMinute          = "every minute"
	fmtEveryMinuteOf     = "every minute of %s"
	fmtEveryMinuteFrom   = "every minute from %s to %s"
	fmtEveryMinuteFromOf = "every minute from %s to %s of %s"
	fmtEveryNMinutes     = "every %s minutes"
	fmtEveryNMinutesOf   = "every %s minutes of %s"

	everyHour        = "every hour"
	fmtEveryHourFrom = "every hour from %s to %s"
	fmtEveryNHours   = "every %s hours"

	fmtAtN   = "at %s"
	fmtAtNOf = "at %s of %s"
	fmtOf    = "of %s"
	fmtAnd   = "and %s"

	fmtHour     = "%02sh"
	fmtFullHour = "%02sh%02s"
	fmtEvery    = "every %s"
	fmtRange    = "from %s to %s"

	fmtModeError = "unexpected %s translation mode [%d]"
)

const partsCount = 5

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
	if len(parts) != partsCount {
		return "", fmt.Errorf(fmtInvalidParam, v)
	}

	var reparts = make([][]string, partsCount)
	var modes = make([]int, partsCount)
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

func getMinuteHourTranslationForModeEvery(reparts [][]string, mode int) (string, error) {
	switch mode {
	case modeEvery:
		return everyMinute, nil

	case modeEveryN:
		return fmt.Sprintf(fmtEveryNMinutes, getEveryValue(reparts[mm])), nil

	case modeAtN:
		fallthrough

	case modeList:
		mparts := strings.Split(reparts[mm][3], ",")
		ret := make([]string, len(mparts))
		for i := range mparts {
			ret[i] = fmt.Sprintf(fmtFullHour, xhour, mparts[i])
		}
		return fmt.Sprintf(fmtAtNOf, strings.Join(ret, ", "), everyHour), nil

	case modeRange:
		parts := strings.Split(reparts[mm][1], "-")
		return fmt.Sprintf(fmtEveryMinuteFromOf, fmt.Sprintf(fmtFullHour, xhour, parts[0]), fmt.Sprintf(fmtFullHour, xhour, parts[1]), everyHour), nil

	default:
		return "", fmt.Errorf(fmtModeError, "minute", mode)
	}
}

func getMinuteHourTranslationForModeEveryN(reparts [][]string, mode int) (string, error) {
	switch mode {
	case modeEvery:
		return fmt.Sprintf(fmtEveryMinuteOf, fmt.Sprintf(fmtEveryNHours, strings.TrimPrefix(reparts[hh][2], "/"))), nil

	case modeEveryN:
		return fmt.Sprintf(fmtEveryNMinutesOf, getEveryValue(reparts[mm]), fmt.Sprintf(fmtEveryNHours, strings.TrimPrefix(reparts[hh][2], "/"))), nil

	case modeAtN:
		fallthrough
	case modeList:
		parts := strings.Split(reparts[mm][3], ",")
		ret := make([]string, len(parts))
		for i := range parts {
			ret[i] = fmt.Sprintf(fmtFullHour, xhour, parts[i])
		}
		return fmt.Sprintf(fmtAtNOf, strings.Join(ret, ", "), fmt.Sprintf(fmtEveryNHours, strings.TrimPrefix(reparts[hh][2], "/"))), nil

	case modeRange:
		parts := strings.Split(reparts[mm][1], "-")
		return fmt.Sprintf(fmtEveryMinuteFromOf, fmt.Sprintf(fmtFullHour, xhour, parts[0]), fmt.Sprintf(fmtFullHour, xhour, parts[1]), fmt.Sprintf(fmtEveryNHours, strings.TrimPrefix(reparts[hh][2], "/"))), nil

	default:
		return "", fmt.Errorf(fmtModeError, "minute", mode)
	}
}

func getMinuteHourTranslationForModeAtN(reparts [][]string, mode int) (string, error) {
	switch mode {
	case modeEvery:
		return fmt.Sprintf(fmtEveryMinuteOf, fmt.Sprintf(fmtHour, reparts[hh][3])), nil

	case modeEveryN:
		return fmt.Sprintf(fmtEveryNMinutesOf, getEveryValue(reparts[mm]), fmt.Sprintf(fmtHour, reparts[hh][3])), nil

	case modeAtN:
		return fmt.Sprintf(fmtAtN, fmt.Sprintf(fmtFullHour, reparts[hh][3], reparts[mm][3])), nil

	case modeList:
		parts := strings.Split(reparts[mm][3], ",")
		ret := make([]string, len(parts))
		for i := range parts {
			ret[i] = fmt.Sprintf(fmtFullHour, reparts[hh][3], parts[i])
		}
		return fmt.Sprintf(fmtAtN, strings.Join(ret, ", ")), nil

	case modeRange:
		parts := strings.Split(reparts[mm][1], "-")
		return fmt.Sprintf(fmtEveryMinuteFrom, fmt.Sprintf(fmtFullHour, reparts[hh][3], parts[0]), fmt.Sprintf(fmtFullHour, reparts[hh][3], parts[1])), nil

	default:
		return "", fmt.Errorf(fmtModeError, "minute", mode)
	}
}

func getMinuteHourTranslationForModeList(reparts [][]string, mode int) (string, error) {
	hparts := strings.Split(reparts[hh][3], ",")
	switch mode {
	case modeEvery:
		ret := make([]string, len(hparts))
		for i := range hparts {
			ret[i] = fmt.Sprintf(fmtHour, hparts[i])
		}
		return fmt.Sprintf(fmtEveryMinuteOf, strings.Join(ret, ", ")), nil

	case modeEveryN:
		ret := make([]string, len(hparts))
		for i := range hparts {
			ret[i] = fmt.Sprintf(fmtHour, hparts[i])
		}
		return fmt.Sprintf(fmtEveryNMinutesOf, getEveryValue(reparts[mm]), strings.Join(ret, ", ")), nil

	case modeAtN:
		fallthrough
	case modeList:
		mparts := strings.Split(reparts[mm][3], ",")
		ret := make([]string, len(hparts)*len(mparts))
		k := 0
		for i := range hparts {
			for j := range mparts {
				ret[k] = fmt.Sprintf(fmtFullHour, hparts[i], mparts[j])
				k++
			}
		}
		return fmt.Sprintf(fmtAtN, strings.Join(ret, ", ")), nil

	case modeRange:
		mparts := strings.Split(reparts[mm][1], "-")
		ret := make([]string, len(hparts))
		for i := range hparts {
			ret[i] = fmt.Sprintf(fmtRange, fmt.Sprintf(fmtFullHour, hparts[i], mparts[0]), fmt.Sprintf(fmtFullHour, hparts[i], mparts[1]))
		}
		return everyMinute + " " + strings.Join(ret, ", "), nil

	default:
		return "", fmt.Errorf(fmtModeError, "minute", mode)
	}
}

func getMinuteHourTranslationForModeRange(reparts [][]string, mode int) (string, error) {
	hparts := strings.Split(reparts[hh][1], "-")
	switch mode {
	case modeEvery:
		return fmt.Sprintf(fmtEveryMinuteOf, fmt.Sprintf(fmtEveryHourFrom, fmt.Sprintf(fmtHour, hparts[0]), fmt.Sprintf(fmtHour, hparts[1]))), nil

	case modeEveryN:
		return fmt.Sprintf(fmtEveryNMinutesOf, getEveryValue(reparts[mm]), fmt.Sprintf(fmtEveryHourFrom, fmt.Sprintf(fmtHour, hparts[0]), fmt.Sprintf(fmtHour, hparts[1]))), nil

	case modeAtN:
		return fmt.Sprintf(fmtAtNOf, fmt.Sprintf(fmtFullHour, xhour, reparts[mm][3]), fmt.Sprintf(fmtEveryHourFrom, fmt.Sprintf(fmtHour, hparts[0]), fmt.Sprintf(fmtHour, hparts[1]))), nil

	case modeList:
		parts := strings.Split(reparts[mm][3], ",")
		ret := make([]string, len(parts))
		for i := range parts {
			ret[i] = fmt.Sprintf(fmtFullHour, xhour, parts[i])
		}
		return fmt.Sprintf(fmtAtNOf, strings.Join(ret, ", "), fmt.Sprintf(fmtEveryHourFrom, fmt.Sprintf(fmtHour, hparts[0]), fmt.Sprintf(fmtHour, hparts[1]))), nil

	case modeRange:
		mparts := strings.Split(reparts[mm][1], "-")
		return fmt.Sprintf(fmtEveryMinuteFromOf, fmt.Sprintf(fmtFullHour, xhour, mparts[0]), fmt.Sprintf(fmtFullHour, xhour, mparts[1]), fmt.Sprintf(fmtEveryHourFrom, fmt.Sprintf(fmtHour, hparts[0]), fmt.Sprintf(fmtHour, hparts[1]))), nil

	default:
		return "", fmt.Errorf(fmtModeError, "minute", mode)
	}
}

func getMinuteHourTranslation(reparts [][]string, modes []int) (string, error) {
	switch modes[hh] {
	case modeEvery:
		return getMinuteHourTranslationForModeEvery(reparts, modes[mm])

	case modeEveryN:
		return getMinuteHourTranslationForModeEveryN(reparts, modes[mm])

	case modeAtN:
		return getMinuteHourTranslationForModeAtN(reparts, modes[mm])

	case modeList:
		return getMinuteHourTranslationForModeList(reparts, modes[mm])

	case modeRange:
		return getMinuteHourTranslationForModeRange(reparts, modes[mm])
	}

	return "", fmt.Errorf(fmtModeError, "hour", modes[hh])
}

func getDayTranslation(reparts [][]string, modes []int) (string, error) {
	var fmtRet string
	if modes[ddd] != modeEvery {
		fmtRet = fmtAnd
	} else {
		fmtRet = fmtOf
	}

	switch modes[dd] {
	case modeEvery:
		if modes[ddd] != modeEvery {
			return "", nil
		}
		return fmt.Sprintf(fmtRet, everyDay), nil

	case modeEveryN:
		return fmt.Sprintf(fmtRet, fmt.Sprintf(fmtEveryNDays, getEveryValue(reparts[dd]))), nil

	case modeAtN:
		fallthrough

	case modeList:
		return fmt.Sprintf(fmtRet, fmt.Sprintf(fmtEvery, reparts[dd][3])), nil

	case modeRange:
		parts := strings.Split(reparts[dd][1], "-")
		return fmt.Sprintf(fmtRet, fmt.Sprintf(fmtRange, parts[0], parts[1])), nil
	}

	return "", fmt.Errorf(fmtModeError, "day", modes[dd])
}

func getMonthTranslation(reparts [][]string, modes []int) (string, error) {
	switch modes[mmm] {
	case modeEvery:
		if modes[dd] == modeEvery {
			return "", nil
		}
		return fmt.Sprintf(fmtOf, "the month"), nil

	case modeEveryN:
		return fmt.Sprintf(fmtOf, fmt.Sprintf(fmtEveryNMonths, getEveryValue(reparts[mmm]))), nil

	case modeAtN:
		fallthrough

	case modeList:
		parts := strings.Split(reparts[mmm][3], ",")
		ret := make([]string, len(parts))
		for i := range parts {
			if _, present := monthNames[parts[i]]; !present {
				return "", fmt.Errorf(fmtInvalidMonth, reparts[mmm][3])
			}
			ret[i] = monthNames[parts[i]]
		}
		return fmt.Sprintf(fmtOf, strings.Join(ret, ",")), nil

	case modeRange:
		parts := strings.Split(reparts[mmm][1], "-")
		ret := make([]string, len(parts))
		for i := range parts {
			if _, present := monthNames[parts[i]]; !present {
				return "", fmt.Errorf(fmtInvalidMonth, reparts[mmm][1])
			}
			ret[i] = monthNames[parts[i]]
		}
		return fmt.Sprintf(fmtRange, ret[0], ret[1]), nil
	}

	return "", fmt.Errorf(fmtModeError, "month", modes[mmm])
}

func getDOTWTranslation(reparts [][]string, modes []int) (string, error) {
	switch modes[ddd] {
	case modeEvery:
		return "", nil

	case modeEveryN:
		return fmt.Sprintf(fmtOf, fmt.Sprintf(fmtEveryNDays, getEveryValue(reparts[ddd]))), nil

	case modeAtN:
		fallthrough

	case modeList:
		parts := strings.Split(reparts[ddd][3], ",")
		ret := make([]string, len(parts))
		for i := range parts {
			if _, present := dOTW[parts[i]]; !present {
				return "", fmt.Errorf(fmtInvalidDOTW, reparts[ddd][3])
			}
			ret[i] = dOTW[parts[i]]
		}
		return fmt.Sprintf(fmtOf, fmt.Sprintf(fmtEvery, strings.Join(ret, ","))), nil

	case modeRange:
		parts := strings.Split(reparts[ddd][1], "-")
		ret := make([]string, len(parts))
		for i := range parts {
			if _, present := dOTW[parts[i]]; !present {
				return "", fmt.Errorf(fmtInvalidDOTW, reparts[ddd][1])
			}
			ret[i] = dOTW[parts[i]]
		}
		return fmt.Sprintf(fmtRange, ret[0], ret[1]), nil
	}

	return "", fmt.Errorf(fmtModeError, "DOTW", modes[ddd])
}

func getEveryValue(reparts []string) string {
	return strings.TrimPrefix(reparts[2], "/")
}
