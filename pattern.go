package reply

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"
)

var patterns = map[string][]*regexp2.Regexp{
	"REMOVE_PGP_MARKERS_REGEX": {
		regexp2.MustCompile(`\A-----BEGIN PGP SIGNED MESSAGE-----\n(?:Hash: \w+)?\s+`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		regexp2.MustCompile(
			`^-----BEGIN PGP SIGNATURE-----$[\s\S]+^-----END PGP SIGNATURE-----`,
			regexp2.RE2,
		),
	},

	"REMOVE_UNSUBSCRIBE_REGEX": {
		regexp2.MustCompile(`^Unsubscribe: .+@.+(\n.+http:.+)?\s*\z`, regexp2.IgnoreCase|regexp2.RE2),
	},

	"REMOVE_ALIAS_REGEX": {
		regexp2.MustCompile(`^.*>{5} "[^"\n]+" == .+ writes:`, regexp2.RE2),
	},

	"CHANGE_ENCLOSED_QUOTE_ONE_REGEX": {
		regexp2.MustCompile(`^>>> ?(.+) ?>>>$\n([\s\S]+?)\n^<<< ?1 ?<<<$`, regexp2.RE2),
	},

	"CHANGE_ENCLOSED_QUOTE_TWO_REGEX": {
		regexp2.MustCompile(`^>{4,}[[:blank:]]*$\n([\s\S]+?)\n^<{4,}[[:blank:]]*$`, regexp2.RE2|regexp2.Multiline),
	},

	"FIX_QUOTES_FORMAT_REGEX": {
		regexp2.MustCompile(`^((?:[[:blank:]]*[[:alpha:]]*[>|])+)`, regexp2.RE2|regexp2.Multiline),
	},

	// On init
	"FIX_EMBEDDED_REGEX": {},

	// Envoyé depuis mon iPhone
	// Von meinem Mobilgerät gesendet
	// Diese Nachricht wurde von meinem Android-Mobiltelefon mit K-9 Mail gesendet.
	// Nik from mobile
	// From My Iphone 6
	// Sent via mobile
	// Sent with Airmail
	// Sent from Windows Mail
	// Sent from my TI-85
	// <<sent by galaxy>>
	// (sent from a phone)
	// (Sent from mobile device)
	// 從我的 iPhone 傳送
	"SIGNATURE_REGEXES": {
		// Chinese
		regexp2.MustCompile(
			`^[[:blank:]]*從我的 iPhone 傳送`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		// English
		regexp2.MustCompile(
			`^[[:blank:]]*[[:word:]]+ from mobile`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		regexp2.MustCompile(
			`^[[:blank:]]*[(<]*Sent (from|via|with|by) .+[)>]*`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		regexp2.MustCompile(
			`^[[:blank:]]*From my .{1,20}`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		regexp2.MustCompile(
			`^[[:blank:]]*Get Outlook for `,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		// French
		regexp2.MustCompile(
			`^[[:blank:]]*Envoyé depuis (mon|Yahoo Mail)`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		// German
		regexp2.MustCompile(
			`^[[:blank:]]*Von meinem .+ gesendet`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		regexp2.MustCompile(
			`^[[:blank:]]*Diese Nachricht wurde von .+ gesendet`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		// Italian
		regexp2.MustCompile(
			`^[[:blank:]]*Inviato da `,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		// Norwegian
		regexp2.MustCompile(
			`^[[:blank:]]*Sendt fra min `,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		// Portuguese
		regexp2.MustCompile(
			`^[[:blank:]]*Enviado do meu `,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		// Spanish
		regexp2.MustCompile(
			`^[[:blank:]]*Enviado desde mi `,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		// Dutch
		regexp2.MustCompile(
			`^[[:blank:]]*Verzonden met `,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		regexp2.MustCompile(
			`^[[:blank:]]*Verstuurd vanaf mijn `,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		// Swedish
		regexp2.MustCompile(
			`^[[:blank:]]*från min `,
			regexp2.IgnoreCase|regexp2.RE2,
		),
	},

	// On init
	"EMAIL_HEADERS_WITH_DATE_REGEXES": {},
	"EMAIL_HEADERS_WITH_TEXT_REGEXES": {},
	"EMAIL_HEADER_REGEXES":            {},

	// On Wed, Sep 25, 2013, at 03:57 PM, jorge_castro wrote:
	// On Thursday, June 27, 2013, knwang via Discourse Meta wrote:
	// On Wed, 2015-12-02 at 13:58 +0000, Tom Newsom wrote:
	// On 10/12/15 12:30, Jeff Atwood wrote:
	// ---- On Tue, 22 Dec 2015 14:17:36 +0530 Sam Saffron&lt;info@discourse.org&gt; wrote ----
	// Op 24 aug. 2013 om 16:48 heeft ven88 via Discourse Meta <info@discourse.org> het volgende geschreven:
	// Le 4 janv. 2016 19:03, "Neil Lalonde" <info@discourse.org> a écrit :
	// Dnia 14 lip 2015 o godz. 00:25 Michael Downey <info@discourse.org> napisał(a):
	// Em seg, 27 de jul de 2015 17:13, Neil Lalonde <info@discourse.org> escreveu:
	// El jueves, 21 de noviembre de 2013, codinghorror escribió:
	// At 6/16/2016 08:32 PM, you wrote:
	"ON_DATE_SOMEONE_WROTE_REGEXES": {
		// Chinese
		regexp2.MustCompile(
			`^[[:blank:]<>-]*在 (?:(?!\b(?>在|写道)\b).)+?写道[[:blank:].:>-]*$`,
			regexp2.IgnoreCase|regexp2.Singleline|regexp2.Multiline|regexp2.RE2,
		),
		// Dutch
		regexp2.MustCompile(
			`^[[:blank:]<>-]*Op (?:(?!\b(?>Op|het\svolgende\sgeschreven|schreef)\b).)+?(het\svolgende\sgeschreven|schreef[^:]+)[[:blank:].:>-]*$`,
			regexp2.IgnoreCase|regexp2.Singleline|regexp2.Multiline|regexp2.RE2,
		),
		// English
		regexp2.MustCompile(
			`^[[:blank:]<>-]*In message (?:(?!\b(?>In message|writes)\b).)+?writes[[:blank:].:>-]*$`,
			regexp2.IgnoreCase|regexp2.Singleline|regexp2.Multiline|regexp2.RE2,
		),
		regexp2.MustCompile(
			`^[[:blank:]<>-]*(On|At) (?:(?!\b(?>On|wrote|writes|says|said)\b).)+?(wrote|writes|says|said)[[:blank:].:>-]*$`,
			regexp2.IgnoreCase|regexp2.Singleline|regexp2.Multiline|regexp2.RE2,
		),
		// French
		regexp2.MustCompile(
			`^[[:blank:]<>-]*Le (?:(?!\b(?>Le|nous\sa\sdit|a\s+écrit)\b).)+?(nous\sa\sdit|a\s+écrit)[[:blank:].:>- ]*$`,
			regexp2.IgnoreCase|regexp2.Singleline|regexp2.Multiline|regexp2.RE2,
		),
		// German
		regexp2.MustCompile(
			`^[[:blank:]<>-]*Am (?:(?!\b(?>Am|schrieben\sSie)\b).)+?schrieben\sSie[[:blank:].:>-]*$`,
			regexp2.IgnoreCase|regexp2.Singleline|regexp2.Multiline|regexp2.RE2,
		),
		regexp2.MustCompile(
			`^[[:blank:]<>-]*Am (?:(?!\b(?>Am|geschrieben)\b).)+?(geschrieben|schrieb[^:]+)[[:blank:].:>-]*$`,
			regexp2.IgnoreCase|regexp2.Singleline|regexp2.Multiline|regexp2.RE2,
		),
		// Italian
		regexp2.MustCompile(
			`^[[:blank:]<>-]*Il (?:(?!\b(?>Il|ha\sscritto)\b).)+?ha\sscritto[[:blank:].:>-]*$`,
			regexp2.IgnoreCase|regexp2.Singleline|regexp2.Multiline|regexp2.RE2,
		),
		// Polish
		regexp2.MustCompile(
			`^[[:blank:]<>-]*(Dnia|Dňa) (?:(?!\b(?>Dnia|Dňa|napisał)\b).)+?napisał(\(a\))?[[:blank:].:>-]*$`,
			regexp2.IgnoreCase|regexp2.Singleline|regexp2.Multiline|regexp2.RE2,
		),
		// Portuguese
		regexp2.MustCompile(
			`^[[:blank:]<>-]*Em (?:(?!\b(?>Em|escreveu)\b).)+?escreveu[[:blank:].:>-]*$`,
			regexp2.IgnoreCase|regexp2.Singleline|regexp2.Multiline|regexp2.RE2,
		),
		// Spanish
		regexp2.MustCompile(
			`^[[:blank:]<>-]*El (?:(?!\b(?>El|escribió)\b).)+?escribió[[:blank:].:>-]*$`,
			regexp2.IgnoreCase|regexp2.Singleline|regexp2.Multiline|regexp2.RE2,
		),
	},

	// On init
	"ON_DATE_WROTE_SOMEONE_REGEXES": {},
	"DATE_SOMEONE_WROTE_REGEXES":    {},

	// 2015-10-18 0:17 GMT+03:00 Matt Palmer <info@discourse.org>:
	// 2013/10/2 camilohollanda <info@discourse.org>
	// вт, 5 янв. 2016 г. в 23:39, Erlend Sogge Heggen <info@discourse.org>:
	// ср, 1 апр. 2015, 18:29, Denis Didkovsky <info@discourse.org>:
	"DATE_SOMEONE_EMAIL_REGEX": {
		regexp2.MustCompile(
			`\d{4}.{1,80}\s?<[^@<>]+@[^@<>.]+\.[^@<>]+>:?$`,
			regexp2.RE2|regexp2.Multiline,
		),
	},

	// Max Mustermann <try_discourse@discoursemail.com> schrieb am Fr., 28. Apr. 2017 um 11:53 Uhr:
	"SOMEONE_WROTE_ON_DATE_REGEXES": {
		// English
		regexp2.MustCompile(
			`^.+\bwrote\b[[:space:]]+\bon\b.+[^:]+:`,
			regexp2.RE2|regexp2.Multiline,
		),
		// German
		regexp2.MustCompile(
			`^.+\bschrieb\b[[:space:]]+\bam\b.+[^:]+:`,
			regexp2.RE2|regexp2.Multiline,
		),
	},

	// 2016-03-03 17:21 GMT+01:00 Some One
	"ISO_DATE_SOMEONE_REGEX": {
		regexp2.MustCompile(
			`^[[:blank:]>]*20\d\d-\d\d-\d\d \d\d:\d\d GMT\+\d\d:\d\d [\w[:blank:]]+$`,
			regexp2.RE2,
		),
	},

	// Some One <info@discourse.org> wrote:
	// Gavin Sinclair (gsinclair@soyabean.com.au) wrote:
	"SOMEONE_EMAIL_WROTE_REGEX": {
		regexp2.MustCompile(
			`^.+\b[\w.+-]+@[\w.-]+\.\w{2,}\b.+wrote:?$`,
			regexp2.RE2,
		),
	},

	"SOMEONE_VIA_SOMETHING_WROTE_REGEXES": {},

	// Posted by mpalmer on 01/21/2016
	"POSTED_BY_SOMEONE_ON_DATE_REGEX": {
		regexp2.MustCompile(
			`^[[:blank:]>]*Posted by .+ on \d{2}\/\d{2}\/\d{4}$`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
	},

	// Begin forwarded message:
	// Reply Message
	// ----- Forwarded Message -----
	// ----- Original Message -----
	// -----Original Message-----
	// ----- Mensagem Original -----
	// -----Mensagem Original-----
	// *----- Original Message -----*
	// ----- Reply message -----
	// ------------------ 原始邮件 ------------------
	"FORWARDED_EMAIL_REGEXES": {
		// English
		regexp2.MustCompile(
			`^[[:blank:]>]*Begin forwarded message:`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		regexp2.MustCompile(
			`^[[:blank:]>*]*-{2,}[[:blank:]]*(Forwarded|Original|Reply) Message[[:blank:]]*-{2,}`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		// French
		regexp2.MustCompile(
			`^[[:blank:]>]*Début du message transféré :`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		regexp2.MustCompile(
			`^[[:blank:]>*]*-{2,}[[:blank:]]*Message transféré[[:blank:]]*-{2,}`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		// German
		regexp2.MustCompile(
			`^[[:blank:]>*]*-{2,}[[:blank:]]*Ursprüngliche Nachricht[[:blank:]]*-{2,}`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		// Spanish
		regexp2.MustCompile(
			`^[[:blank:]>*]*-{2,}[[:blank:]]*Mensaje original[[:blank:]]*-{2,}`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		// Portuguese
		regexp2.MustCompile(
			`^[[:blank:]>*]*-{2,}[[:blank:]]*Mensagem original[[:blank:]]*-{2,}`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
		// Chinese
		regexp2.MustCompile(
			`^[[:blank:]>*]*-{2,}[[:blank:]]*原始邮件[[:blank:]]*-{2,}`,
			regexp2.IgnoreCase|regexp2.RE2,
		),
	},

	// on init
	"EMBEDDED_REGEXES": {},
}

// init ON_DATE_WROTE_SOMEONE_REGEXES
func init() {
	dateMarkers := [][]string{
		// Norwegian
		{"Sendt"},
		// English
		{"Sent", "Date"},
		// French
		{"Date", "Le"},
		// German
		{"Gesendet"},
		// Portuguese
		{"Enviada em"},
		// Spanish
		{"Enviado"},
		// Spanish (Mexican)
		{"Fecha"},
		// Italian
		{"Data"},
		// Dutch
		{"Datum"},
		// Swedish
		{"Skickat"},
		// Chinese
		{"发送时间"},
	}

	textMarkers := [][]string{
		// Norwegian
		{"Fra", "Til", "Emne"},
		// English
		{"From", "To", "Cc", "Reply-To", "Subject"},
		// French
		{"De", "Expéditeur", "À", "Destinataire", "Répondre à", "Objet"},
		// German
		{"Von", "An", "Betreff"},
		// Portuguese
		{"De", "Para", "Assunto"},
		// Spanish
		{"De", "Para", "Asunto"},
		// Italian
		{"Da", "Risposta", "A", "Oggetto"},
		// Dutch
		{"Van", "Beantwoorden - Aan", "Aan", "Onderwerp"},
		// Swedish
		{"Från", "Till", "Ämne"},
		// Chinese
		{"发件人", "收件人", "主题"},
	}

	// Op 10 dec. 2015 18:35 schreef "Arpit Jalan" <info@discourse.org>:
	// Am 18.09.2013 um 16:24 schrieb codinghorror <info@discourse.org>:
	// Den 15. jun. 2016 kl. 20.42 skrev Jeff Atwood <info@discourse.org>:
	// søn. 30. apr. 2017 kl. 00.26 skrev David Taylor <meta@discoursemail.com>:
	onDateWroteSomeoneMarkers := [][]string{
		// Dutch
		{"Op", "schreef"},
		// German
		{"Am", "schrieb"},
		// Norwegian
		{"Den", "skrev"},
		// Dutch
		{`søn\.`, "skrev"},
	}

	// суббота, 14 марта 2015 г. пользователь etewiah написал:
	// 23 mar 2017 21:25 "Neil Lalonde" <meta@discoursemail.com> napisał(a):
	// 30 серп. 2016 р. 20:45 "Arpit" no-reply@example.com пише:
	dateSomeoneWroteMarkers := [][]string{
		// Russian
		{"пользователь", "написал"},
		// Polish
		{"", "napisał\\(a\\)"},
		// Ukrainian
		{"", "пише"},
	}

	// codinghorror via Discourse Meta wrote:
	// codinghorror via Discourse Meta <info@discourse.org> schrieb:
	someoneViaSomethingWroteMarkers := []string{
		// English
		"wrote",
		// German
		"schrieb",
	}

	// date
	for _, markers := range dateMarkers {
		pattern := regexp2.MustCompile(
			fmt.Sprintf(`^[[:blank:]*]*(?:%s)[[:blank:]*]*:.*\d+`, strings.Join(markers, "|")),
			regexp2.RE2|regexp2.Multiline,
		)
		patterns["EMAIL_HEADERS_WITH_DATE_REGEXES"] = append(
			patterns["EMAIL_HEADERS_WITH_DATE_REGEXES"],
			pattern,
		)

		patterns["EMAIL_HEADER_REGEXES"] = append(patterns["EMAIL_HEADER_REGEXES"], pattern)
	}

	// text
	for _, markers := range textMarkers {
		pattern := regexp2.MustCompile(
			fmt.Sprintf(`^[[:blank:]*]*(?:%s)[[:blank:]*]*:.*[[:word:]]+`, strings.Join(markers, "|")),
			regexp2.IgnoreCase|regexp2.Multiline|regexp2.RE2,
		)
		patterns["EMAIL_HEADERS_WITH_TEXT_REGEXES"] = append(
			patterns["EMAIL_HEADERS_WITH_TEXT_REGEXES"],
			pattern,
		)

		patterns["EMAIL_HEADER_REGEXES"] = append(patterns["EMAIL_HEADER_REGEXES"], pattern)
	}

	for _, marker := range onDateWroteSomeoneMarkers {
		patterns["ON_DATE_WROTE_SOMEONE_REGEXES"] = append(
			patterns["ON_DATE_WROTE_SOMEONE_REGEXES"],
			regexp2.MustCompile(fmt.Sprintf(`^[[:blank:]>]*%s\s.+\s%s\s[^:]+:`, marker[0], marker[1]), regexp2.RE2),
		)
	}

	for _, marker := range dateSomeoneWroteMarkers {
		if marker[0] == "" {
			patterns["DATE_SOMEONE_WROTE_REGEXES"] = append(
				patterns["DATE_SOMEONE_WROTE_REGEXES"],
				regexp2.MustCompile(fmt.Sprintf(`\d{4}.{1,80}\n?.{0,80}?%s:`, marker[1]), regexp2.RE2),
			)
		} else {
			patterns["DATE_SOMEONE_WROTE_REGEXES"] = append(
				patterns["DATE_SOMEONE_WROTE_REGEXES"],
				regexp2.MustCompile(fmt.Sprintf(`\d{4}.{1,80}%s.{0,80}\n?.{0,80}?%s:`, marker[0], marker[1]), regexp2.RE2),
			)
		}
	}

	for _, marker := range someoneViaSomethingWroteMarkers {
		patterns["SOMEONE_VIA_SOMETHING_WROTE_REGEXES"] = append(
			patterns["SOMEONE_VIA_SOMETHING_WROTE_REGEXES"],
			regexp2.MustCompile(fmt.Sprintf(`^.+ via .+ %s:?[[:blank:]]*$`, marker), regexp2.RE2),
		)
	}

	patterns["EMBEDDED_REGEXES"] = append(patterns["EMBEDDED_REGEXES"], patterns["ON_DATE_SOMEONE_WROTE_REGEXES"]...)
	patterns["EMBEDDED_REGEXES"] = append(patterns["EMBEDDED_REGEXES"], patterns["ON_DATE_WROTE_SOMEONE_REGEXES"]...)
	patterns["EMBEDDED_REGEXES"] = append(patterns["EMBEDDED_REGEXES"], patterns["DATE_SOMEONE_WROTE_REGEXES"]...)
	patterns["EMBEDDED_REGEXES"] = append(patterns["EMBEDDED_REGEXES"], patterns["DATE_SOMEONE_EMAIL_REGEX"]...)
	patterns["EMBEDDED_REGEXES"] = append(patterns["EMBEDDED_REGEXES"], patterns["SOMEONE_WROTE_ON_DATE_REGEXES"]...)
	patterns["EMBEDDED_REGEXES"] = append(patterns["EMBEDDED_REGEXES"], patterns["ISO_DATE_SOMEONE_REGEX"]...)
	patterns["EMBEDDED_REGEXES"] = append(patterns["EMBEDDED_REGEXES"], patterns["SOMEONE_VIA_SOMETHING_WROTE_REGEXES"]...)
	patterns["EMBEDDED_REGEXES"] = append(patterns["EMBEDDED_REGEXES"], patterns["SOMEONE_EMAIL_WROTE_REGEX"]...)
	patterns["EMBEDDED_REGEXES"] = append(patterns["EMBEDDED_REGEXES"], patterns["POSTED_BY_SOMEONE_ON_DATE_REGEX"]...)
	patterns["EMBEDDED_REGEXES"] = append(patterns["EMBEDDED_REGEXES"], patterns["FORWARDED_EMAIL_REGEXES"]...)

	patterns["FIX_EMBEDDED_REGEX"] = append(patterns["FIX_EMBEDDED_REGEX"], patterns["ON_DATE_SOMEONE_WROTE_REGEXES"]...)
	patterns["FIX_EMBEDDED_REGEX"] = append(patterns["FIX_EMBEDDED_REGEX"], patterns["SOMEONE_WROTE_ON_DATE_REGEXES"]...)
	patterns["FIX_EMBEDDED_REGEX"] = append(patterns["FIX_EMBEDDED_REGEX"], patterns["DATE_SOMEONE_WROTE_REGEXES"]...)
	patterns["FIX_EMBEDDED_REGEX"] = append(patterns["FIX_EMBEDDED_REGEX"], patterns["DATE_SOMEONE_EMAIL_REGEX"]...)

}
