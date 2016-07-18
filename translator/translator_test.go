package translator

import (
	"testing"
)

// TODO: need more test
func TestGetTranslation(t *testing.T) {

	var test = [][]string{
		{`0 0 13 * 5`, `at 00h00 of every Friday and every 13 of the month`},
		{`30 23 * * *`, `at 23h30 of every day`},
		{`5 * * * *`, `every hour past 5 minutes of every day`},
		{`5,10,25 * * * *`, `every hour past 5,10,25 minutes of every day`},
		{`30 23 1 * *`, `at 23h30 of every 1 of the month`}, //5
		{`28 22 * * 1`, `at 22h28 of every Monday`},
		{`22 11 13 * 5`, `at 11h22 of every Friday and every 13 of the month`},
		{`12 10 2-5 * *`, `at 10h12 of from 2 to 5 of the month`},
		{`59 23 */2 * *`, `at 23h59 of every 2 days of the month`},
		{`0 22 * * 1-5`, `at 22h00 from Monday to Friday`}, //10
		{`*/5 22 * * 1-5`, `every 5 minutes of 22h from Monday to Friday`},
		{`*/5 * * * *`, `every 5 minutes of every day`},
		{`* * * * *`, `every minute of every day`},
		{`0 */12 * * *`, `every 12 hours past 0 minutes of every day`},
		{`5 */12 * * *`, `every 12 hours past 5 minutes of every day`}, //15
		{`59 12-23 * * *`, `every hour past 59 minutes from 12h to 23h of every day`},
		{`0 4,16 * * *`, `at 04h00, 16h00 of every day`},
		{
			`0 0,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23 * * *`,
			`at 00h00, 07h00, 08h00, 09h00, 10h00, 11h00, 12h00, 13h00, 14h00, 15h00, 16h00, 17h00, 18h00, 19h00, 20h00, 21h00, 22h00, 23h00 of every day`,
		},
		{`*/15 */2 * * *`, `every 15 minutes every 2 hours of every day`},
		{`* * * * */2`, `every minute of every 2 days`}, //20
		{`* * * */3 *`, `every minute of every day of every 3 months`},
		{`* * * 3 *`, `every minute of every day of Marsh`},
		{`* * * 1-6 *`, `every minute of every day from January to June`},
		{`* * * 1,2 *`, `every minute of every day of January,February`},
		{`*/2 5 * * *`, `every 2 minutes of 05h of every day`}, //25
		{`10,15,20 5 * * *`, `at 05h10, 05h15, 05h20 of every day`},
		{`10-20 5 * * *`, `every minute from 05h10 to 05h20 of every day`},
		{`* 5 * * *`, `every minute of 05h of every day`},
		{`10-20 * * * *`, `every minute from XXh10 to XXh20 of every hour of every day`},
		{`* */2 * * *`, `every minute of every 2 hours of every day`},
		{`10,15,20 */2 * * *`, `at XXh10, XXh15, XXh20 of every 2 hours of every day`},
		{`25-30 */2 * * *`, `every minute from XXh25 to XXh30 of every 2 hours of every day`},
		{`* 10,15,20 * * *`, `every minute of 10h, 15h, 20h of every day`},
		{`*/2 10,15,20 * * *`, `every 2 minutes of 10h, 15h, 20h of every day`},
		{`25,30,35 10,15,20 * * *`, `at 10h25, 10h30, 10h35, 15h25, 15h30, 15h35, 20h25, 20h30, 20h35 of every day`},
		{`25-30 10,15,20 * * *`, `every minute from 10h25 to 10h30, from 15h25 to 15h30, from 20h25 to 20h30 of every day`},
	}

	var out string
	var err error
	for i := range test {
		out, err = GetTranslation(test[i][0])
		switch {
		case err != nil:
			t.Logf("case %d [%s]:", i+1, test[i][0])
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

func TestGetTranslationWithErrors(t *testing.T) {
	var test = [][]string{
		{`* * * * 10`, `invalid day of the week value [10]`},
		{`* * * * 1-10`, `invalid day of the week value [1-10]`},
		{`* * * * 1,10`, `invalid day of the week value [1,10]`},
		{`* * * 13 *`, `invalid month value [13]`},
		{`* * * 1-13 *`, `invalid month value [1-13]`},
		{`* * * 1,13 *`, `invalid month value [1,13]`},
		{``, `invalid crontab string []`},
		{`a * * * b`, `invalid crontab string [a * * * b]`},
		{`* * * * * a`, `invalid crontab string [* * * * * a]`},
	}

	var err error
	for i := range test {
		_, err = GetTranslation(test[i][0])
		switch {
		case err == nil:
			t.Logf("case %d", i+1)
			t.Logf("    expected error not return")
			t.Fail()
		case test[i][1] != err.Error():
			t.Logf("case %d [%s]: unexpected error message", i+1, test[i][0])
			t.Logf("    expected: [%s]", test[i][1])
			t.Logf("    having:   [%s]", err.Error())
			t.Fail()
		}
	}
}
