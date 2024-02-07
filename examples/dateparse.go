package main

import (
	"fmt"

	"github.com/scylladb/termtables"

	"github.com/oarkflow/pkg/dateparse"
)

var examples = []string{
	// mon day year (time)
	"May 8, 2009 5:57:51 PM",
	"oct 7, 1970",
	"oct 7, '70",
	"oct. 7, 1970",
	"oct. 7, 70",
	"October 7, 1970",
	"October 7th, 1970",
	"Sept. 7, 1970 11:15:26pm",
	"Sep 7 2009 11:15:26.123 PM PST",
	"September 3rd, 2009 11:15:26.123456789pm",
	"September 17 2012 10:09am",
	"September 17, 2012, 10:10:09",
	"Sep 17, 2012 at 10:02am (EST)",
	// (PST-08 will have an offset of -0800, and a zone name of "PST")
	"September 17, 2012 at 10:09am PST-08",
	// (UTC-0700 has the same offset as -0700, and the returned zone name will be empty)
	"September 17 2012 5:00pm UTC-0700",
	"September 17 2012 5:00pm GMT-0700",
	// (weekday) day mon year (time)
	"7 oct 70",
	"7 Oct 1970",
	"7 September 1970 23:15",
	"7 September 1970 11:15:26pm",
	"03 February 2013",
	"12 Feb 2006, 19:17",
	"12 Feb 2006 19:17",
	"14 May 2019 19:11:40.164",
	"4th Sep 2012",
	"1st February 2018 13:58:24",
	"Mon, 02 Jan 2006 15:04:05 MST", // RFC1123
	"Mon, 02 Jan 2006 15:04:05 -0700",
	"Tue, 11 Jul 2017 16:28:13 +0200 (CEST)",
	"Mon 30 Sep 2018 09:09:09 PM UTC",
	"Sun, 07 Jun 2020 00:00:00 +0100",
	"Wed,  8 Feb 2023 19:00:46 +1100 (AEDT)",
	// ANSIC and UnixDate - weekday month day time year
	"Mon Jan  2 15:04:05 2006",
	"Mon Jan  2 15:04:05 MST 2006",
	"Monday Jan 02 15:04:05 -0700 2006",
	"Mon Jan 2 15:04:05.103786 2006",
	// RubyDate - weekday month day time offset year
	"Mon Jan 02 15:04:05 -0700 2006",
	// ANSIC_GLIBC - weekday day month year time
	"Mon 02 Jan 2006 03:04:05 PM UTC",
	"Monday 02 Jan 2006 03:04:05 PM MST",
	// weekday month day time timezone-offset year
	"Mon Aug 10 15:44:11 UTC+0000 2015",
	// git log default date format
	"Thu Apr 7 15:13:13 2005 -0700",
	// variants of git log default date format
	"Thu Apr 7 15:13:13 2005 -07:00",
	"Thu Apr 7 15:13:13 2005 -07:00 PST",
	"Thu Apr 7 15:13:13 2005 -07:00 PST (Pacific Standard Time)",
	"Thu Apr 7 15:13:13 -0700 2005",
	"Thu Apr 7 15:13:13 -07:00 2005",
	"Thu Apr 7 15:13:13 -0700 PST 2005",
	"Thu Apr 7 15:13:13 -07:00 PST 2005",
	"Thu Apr 7 15:13:13 PST 2005",
	// Variants of the above with a (full time zone description)
	"Fri Jul 3 2015 06:04:07 PST-0700 (Pacific Daylight Time)",
	"Fri Jul 03 2015 18:04:07 GMT+0100 (GMT Daylight Time)",
	"Sun, 3 Jan 2021 00:12:23 +0800 (GMT+08:00)",
	// year month day
	"2013 May 2",
	"2013 May 02 11:37:55",
	// dd/Mon/year  alpha Months
	"06/Jan/2008 15:04:05 -0700",
	"06/January/2008 15:04:05 -0700",
	"06/Jan/2008:15:04:05 -0700", // ngnix-log
	"06/January/2008:08:11:17 -0700",
	// mm/dd/year (see also PreferMonthFirst and RetryAmbiguousDateWithSwap options)
	"3/31/2014",
	"03/31/2014",
	"08/21/71",
	"8/1/71",
	"4/8/2014 22:05",
	"04/08/2014 22:05",
	"04/08/2014, 22:05",
	"4/8/14 22:05",
	"04/2/2014 03:00:51",
	"8/8/1965 1:00 PM",
	"8/8/1965 01:00 PM",
	"8/8/1965 12:00 AM",
	"8/8/1965 12:00:00AM",
	"8/8/1965 01:00:01 PM",
	"8/8/1965 01:00:01PM -0700",
	"8/8/1965 13:00:01 -0700 PST",
	"8/8/1965 01:00:01 PM -0700 PST",
	"8/8/1965 01:00:01 PM -07:00 PST (Pacific Standard Time)",
	"4/02/2014 03:00:51",
	"03/19/2012 10:11:59",
	"03/19/2012 10:11:59.3186369",
	// mon/dd/year
	"Oct/ 7/1970",
	"Oct/03/1970 22:33:44",
	"February/03/1970 11:33:44.555 PM PST",
	// yyyy/mm/dd
	"2014/3/31",
	"2014/03/31",
	"2014/4/8 22:05",
	"2014/04/08 22:05",
	"2014/04/2 03:00:51",
	"2014/4/02 03:00:51",
	"2012/03/19 10:11:59",
	"2012/03/19 10:11:59.3186369",
	// weekday, day-mon-yy time
	"Fri, 03-Jul-15 08:08:08 CEST",
	"Monday, 02-Jan-06 15:04:05 MST", // RFC850
	"Monday, 02 Jan 2006 15:04:05 -0600",
	"02-Jan-06 15:04:05 MST",
	// RFC3339 - yyyy-mm-ddThh
	"2006-01-02T15:04:05+0000",
	"2009-08-12T22:15:09-07:00",
	"2009-08-12T22:15:09",
	"2009-08-12T22:15:09.988",
	"2009-08-12T22:15:09Z",
	"2009-08-12T22:15:09.52Z",
	"2017-07-19T03:21:51:897+0100",
	"2019-05-29T08:41-04", // no seconds, 2 digit TZ offset
	// yyyy-mm-dd hh:mm:ss
	"2014-04-26 17:24:37.3186369",
	"2012-08-03 18:31:59.257000000",
	"2014-04-26 17:24:37.123",
	"2014-04-01 12:01am",
	"2014-04-01 12:01:59.765 AM",
	"2014-04-01 12:01:59,765",
	"2014-04-01 22:43",
	"2014-04-01 22:43:22",
	"2014-12-16 06:20:00 UTC",
	"2014-12-16 06:20:00 GMT",
	"2014-04-26 05:24:37 PM",
	"2014-04-26 13:13:43 +0800",
	"2014-04-26 13:13:43 +0800 +08",
	"2014-04-26 13:13:44 +09:00",
	"2012-08-03 18:31:59.257000000 +0000 UTC",
	"2015-09-30 18:48:56.35272715 +0000 UTC",
	"2015-02-18 00:12:00 +0000 GMT", // golang native format
	"2015-02-18 00:12:00 +0000 UTC",
	"2015-02-08 03:02:00 +0300 MSK m=+0.000000001",
	"2015-02-08 03:02:00.001 +0300 MSK m=+0.000000001",
	"2017-07-19 03:21:51+00:00",
	"2017-04-03 22:32:14.322 CET",
	"2017-04-03 22:32:14,322 CET",
	"2017-04-03 22:32:14:322 CET",
	"2018-09-30 08:09:13.123PM PMDT", // PMDT time zone
	"2018-09-30 08:09:13.123 am AMT", // AMT time zone
	"2014-04-26",
	"2014-04",
	"2014",
	// yyyy-mm-dd(offset)
	"2020-07-20+08:00",
	"2020-07-20+0800",
	// year-mon-dd
	"2013-Feb-03",
	"2013-February-03 09:07:08.123",
	// dd-mon-year
	"03-Feb-13",
	"03-Feb-2013",
	"07-Feb-2004 09:07:07 +0200",
	"07-February-2004 09:07:07 +0200",
	// dd-mm-year (this format (common in Europe) always puts the day first, regardless of PreferMonthFirst)
	"28-02-02",
	"28-02-02 15:16:17",
	"28-02-2002",
	"28-02-2002 15:16:17",
	// mm.dd.yy (see also PreferMonthFirst and RetryAmbiguousDateWithSwap options)
	"3.31.2014",
	"03.31.14",
	"03.31.2014",
	"03.31.2014 10:11:59 MST",
	"03.31.2014 10:11:59.3186369Z",
	// year.mm.dd
	"2014.03",
	"2014.03.30",
	"2014.03.30 08:33pm",
	"2014.03.30T08:33:44.555 PM -0700 MST",
	"2014.03.30-0600",
	// yyyy:mm:dd
	"2014:3:31",
	"2014:03:31",
	"2014:4:8 22:05",
	"2014:04:08 22:05",
	"2014:04:2 03:00:51",
	"2014:4:02 03:00:51",
	"2012:03:19 10:11:59",
	"2012:03:19 10:11:59.3186369",
	// mm:dd:yyyy (see also PreferMonthFirst and RetryAmbiguousDateWithSwap options)
	"08:03:2012",
	"08:04:2012 18:31:59+00:00",
	// yyyymmdd and similar
	"20140601",
	"20140722105203",
	"20140722105203.364",
	// Chinese
	"2014年4月25日",
	"2014年04月08日",
	"2014年04月08日 19:17:22 -0700",
	// RabbitMQ log format
	"8-Mar-2018::14:09:27",
	"08-03-2018::02:09:29 PM",
	// yymmdd hh:mm:yy mysql log
	// 080313 05:21:55 mysqld started
	"171113 14:14:20",
	"190910 11:51:49",
	// unix seconds, ms, micro, nano
	"1332151919",
	"1384216367189",
	"1384216367111222",
	"1384216367111222333",
}

var (
	timezone = ""
)

func main() {
	table := termtables.CreateTable()

	table.AddHeaders("Input", "Parsed, and Output as %v")
	for _, dateExample := range examples {
		t, err := dateparse.ParseAny(dateExample)
		if err != nil {
			panic(err.Error())
		}
		table.AddRow(dateExample, fmt.Sprintf("%v", t))
	}
	fmt.Println(table.Render())
}
