package translator

import (
	"testing"
)

// TODO: need more test
func TestGetTranslation(t *testing.T) {

	var test = [][]string{
		{`0 0 13 * 5`, `at 00h00 every Friday and every 13 of the month`},
		{`30 23 * * *`, `at 23h30 every day`},
		{`5 * * * *`, `every hour past 5 minutes`},
		{`5,10,25 * * * *`, `every hour past 5,10,25 minutes`},
		{`30 23 1 * *`, `at 23h30 every 1 of the month`},
		{`28 22 * * 1`, `at 22h28 every Monday`},
		{`22 11 13 * 5`, `at 11h22 every Friday and every 13 of the month`},
		{`12 10 2-5 * *`, `at 10h12 from 2 to 5 of the month`},
		{`59 23 */2 * *`, `at 23h59 every 2 days of the month`},
		{`0 22 * * 1-5`, `at 22h00 from Monday to Friday`},
		{`*/5 22 * * 1-5`, `every 5 minutes of 22h from Monday to Friday`},
		{`*/5 * * * *`, `every 5 minutes`},
		{`* * * * *`, `every minute`},
		{`0 */12 * * *`, `every 12 hours past 0 minutes`},
		{`5 */12 * * *`, `every 12 hours past 5 minutes`},
		{`59 12-23 * * *`, `every hour past 59 minutes from 12h to 23h`},
		{`0 4,16 * * *`, `at 04h00,16h00`},
		{
			`0 0,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23 * * *`,
			`at 00h00,07h00,08h00,09h00,10h00,11h00,12h00,13h00,14h00,15h00,16h00,17h00,18h00,19h00,20h00,21h00,22h00,23h00`,
		},
		{`*/15 */2 * * *`, `every 15 minutes every 2 hours`},
	}

	var out string
	var err error
	for i := range test {
		out, err = GetTranslation(test[i][0])
		switch {
		case err != nil:
			t.Logf("case %d", i+1)
			t.Logf("    unexpected error: %s", err)
			t.Fail()
		case test[i][1] != out:
			t.Logf("case %d [%s]: unexpected output", i+1, test[i][0])
			t.Logf("    expected: [%s]", test[i][1])
			t.Logf("    having:   [%s]", out)
			t.Fail()
		}
	}
}
