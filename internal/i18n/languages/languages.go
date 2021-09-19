//
// @project GeniusRabbit 2016 - 2018
// @author Dmitry Ponomarev <demdxx@gmail.com> 2016 - 2018
//

package languages

import (
	"strings"
)

// Language item object
type Language struct {
	ID         uint    `json:"ID"`
	Code       [2]byte `json:"code"`
	Name       string  `json:"name"`
	NativeName string  `json:"native_name"`
}

// IntCode for ISO
func (l Language) IntCode() uint {
	return IntCode(l.Code)
}

// Languages array
var Languages = []Language{
	{ID: 0, Code: [2]byte{'*', '*'}, Name: "Undefined"},
	{ID: 1, Code: [2]byte{'A', 'B'}, Name: "Abkhaz", NativeName: "аҧсуа"},
	{ID: 2, Code: [2]byte{'A', 'A'}, Name: "Afar", NativeName: "Afaraf"},
	{ID: 3, Code: [2]byte{'A', 'F'}, Name: "Afrikaans", NativeName: "Afrikaans"},
	{ID: 4, Code: [2]byte{'A', 'K'}, Name: "Akan", NativeName: "Akan"},
	{ID: 5, Code: [2]byte{'S', 'Q'}, Name: "Albanian", NativeName: "Shqip"},
	{ID: 6, Code: [2]byte{'A', 'M'}, Name: "Amharic", NativeName: "አማርኛ"},
	{ID: 7, Code: [2]byte{'A', 'R'}, Name: "Arabic", NativeName: "العربية"},
	{ID: 8, Code: [2]byte{'A', 'N'}, Name: "Aragonese", NativeName: "Aragonés"},
	{ID: 9, Code: [2]byte{'H', 'Y'}, Name: "Armenian", NativeName: "Հայերեն"},
	{ID: 10, Code: [2]byte{'A', 'S'}, Name: "Assamese", NativeName: "অসমীয়া"},
	{ID: 11, Code: [2]byte{'A', 'V'}, Name: "Avaric", NativeName: "авар мацӀ,магӀарул мацӀ"},
	{ID: 12, Code: [2]byte{'A', 'E'}, Name: "Avestan", NativeName: "avesta"},
	{ID: 13, Code: [2]byte{'A', 'Y'}, Name: "Aymara", NativeName: "aymar aru"},
	{ID: 14, Code: [2]byte{'A', 'Z'}, Name: "Azerbaijani", NativeName: "azərbaycan dili"},
	{ID: 15, Code: [2]byte{'B', 'M'}, Name: "Bambara", NativeName: "bamanankan"},
	{ID: 16, Code: [2]byte{'B', 'A'}, Name: "Bashkir", NativeName: "башҡорт теле"},
	{ID: 17, Code: [2]byte{'E', 'U'}, Name: "Basque", NativeName: "euskara,euskera"},
	{ID: 18, Code: [2]byte{'B', 'E'}, Name: "Belarusian", NativeName: "Беларуская"},
	{ID: 19, Code: [2]byte{'B', 'N'}, Name: "Bengali", NativeName: "বাংলা"},
	{ID: 20, Code: [2]byte{'B', 'H'}, Name: "Bihari", NativeName: "भोजपुरी"},
	{ID: 21, Code: [2]byte{'B', 'I'}, Name: "Bislama", NativeName: "Bislama"},
	{ID: 22, Code: [2]byte{'B', 'S'}, Name: "Bosnian", NativeName: "bosanski jezik"},
	{ID: 23, Code: [2]byte{'B', 'R'}, Name: "Breton", NativeName: "brezhoneg"},
	{ID: 24, Code: [2]byte{'B', 'G'}, Name: "Bulgarian", NativeName: "български език"},
	{ID: 25, Code: [2]byte{'M', 'Y'}, Name: "Burmese", NativeName: "ဗမာစာ"},
	{ID: 26, Code: [2]byte{'C', 'A'}, Name: "Catalan; Valencian", NativeName: "Català"},
	{ID: 27, Code: [2]byte{'C', 'H'}, Name: "Chamorro", NativeName: "Chamoru"},
	{ID: 28, Code: [2]byte{'C', 'E'}, Name: "Chechen", NativeName: "нохчийн мотт"},
	{ID: 29, Code: [2]byte{'N', 'Y'}, Name: "Chichewa; Chewa; Nyanja", NativeName: "chiCheŵa,chinyanja"},
	{ID: 30, Code: [2]byte{'Z', 'H'}, Name: "Chinese", NativeName: "中文 (Zhōngwén),汉语,漢語"},
	{ID: 31, Code: [2]byte{'C', 'V'}, Name: "Chuvash", NativeName: "чӑваш чӗлхи"},
	{ID: 32, Code: [2]byte{'K', 'W'}, Name: "Cornish", NativeName: "Kernewek"},
	{ID: 33, Code: [2]byte{'C', 'O'}, Name: "Corsican", NativeName: "corsu,lingua corsa"},
	{ID: 34, Code: [2]byte{'C', 'R'}, Name: "Cree", NativeName: "ᓀᐦᐃᔭᐍᐏᐣ"},
	{ID: 35, Code: [2]byte{'H', 'R'}, Name: "Croatian", NativeName: "hrvatski"},
	{ID: 36, Code: [2]byte{'C', 'S'}, Name: "Czech", NativeName: "česky,čeština"},
	{ID: 37, Code: [2]byte{'D', 'A'}, Name: "Danish", NativeName: "dansk"},
	{ID: 38, Code: [2]byte{'D', 'V'}, Name: "Divehi; Dhivehi; Maldivian;", NativeName: "ދިވެހި"},
	{ID: 39, Code: [2]byte{'N', 'L'}, Name: "Dutch", NativeName: "Nederlands,Vlaams"},
	{ID: 40, Code: [2]byte{'E', 'N'}, Name: "English", NativeName: "English"},
	{ID: 41, Code: [2]byte{'E', 'O'}, Name: "Esperanto", NativeName: "Esperanto"},
	{ID: 42, Code: [2]byte{'E', 'T'}, Name: "Estonian", NativeName: "eesti,eesti keel"},
	{ID: 43, Code: [2]byte{'E', 'E'}, Name: "Ewe", NativeName: "Eʋegbe"},
	{ID: 44, Code: [2]byte{'F', 'O'}, Name: "Faroese", NativeName: "føroyskt"},
	{ID: 45, Code: [2]byte{'F', 'J'}, Name: "Fijian", NativeName: "vosa Vakaviti"},
	{ID: 46, Code: [2]byte{'F', 'I'}, Name: "Finnish", NativeName: "suomi,suomen kieli"},
	{ID: 47, Code: [2]byte{'F', 'R'}, Name: "French", NativeName: "français,langue française"},
	{ID: 48, Code: [2]byte{'F', 'F'}, Name: "Fula; Fulah; Pulaar; Pular", NativeName: "Fulfulde,Pulaar,Pular"},
	{ID: 49, Code: [2]byte{'G', 'L'}, Name: "Galician", NativeName: "Galego"},
	{ID: 50, Code: [2]byte{'K', 'A'}, Name: "Georgian", NativeName: "ქართული"},
	{ID: 51, Code: [2]byte{'D', 'E'}, Name: "German", NativeName: "Deutsch"},
	{ID: 52, Code: [2]byte{'E', 'L'}, Name: "Greek,Modern", NativeName: "Ελληνικά"},
	{ID: 53, Code: [2]byte{'G', 'N'}, Name: "Guaraní", NativeName: "Avañeẽ"},
	{ID: 54, Code: [2]byte{'G', 'U'}, Name: "Gujarati", NativeName: "ગુજરાતી"},
	{ID: 55, Code: [2]byte{'H', 'T'}, Name: "Haitian; Haitian Creole", NativeName: "Kreyòl ayisyen"},
	{ID: 56, Code: [2]byte{'H', 'A'}, Name: "Hausa", NativeName: "Hausa,هَوُسَ"},
	{ID: 57, Code: [2]byte{'H', 'E'}, Name: "Hebrew (modern)", NativeName: "עברית"},
	{ID: 58, Code: [2]byte{'H', 'Z'}, Name: "Herero", NativeName: "Otjiherero"},
	{ID: 59, Code: [2]byte{'H', 'I'}, Name: "Hindi", NativeName: "हिन्दी,हिंदी"},
	{ID: 60, Code: [2]byte{'H', 'O'}, Name: "Hiri Motu", NativeName: "Hiri Motu"},
	{ID: 61, Code: [2]byte{'H', 'U'}, Name: "Hungarian", NativeName: "Magyar"},
	{ID: 62, Code: [2]byte{'I', 'A'}, Name: "Interlingua", NativeName: "Interlingua"},
	{ID: 63, Code: [2]byte{'I', 'D'}, Name: "Indonesian", NativeName: "Bahasa Indonesia"},
	{ID: 64, Code: [2]byte{'I', 'E'}, Name: "Interlingue", NativeName: "Originally called OccIDental; then Interlingue after WWII"},
	{ID: 65, Code: [2]byte{'G', 'A'}, Name: "Irish", NativeName: "Gaeilge"},
	{ID: 66, Code: [2]byte{'I', 'G'}, Name: "Igbo", NativeName: "Asụsụ Igbo"},
	{ID: 67, Code: [2]byte{'I', 'K'}, Name: "Inupiaq", NativeName: "Iñupiaq,Iñupiatun"},
	{ID: 68, Code: [2]byte{'I', 'O'}, Name: "IDo", NativeName: "IDo"},
	{ID: 69, Code: [2]byte{'I', 'S'}, Name: "Icelandic", NativeName: "Íslenska"},
	{ID: 70, Code: [2]byte{'I', 'T'}, Name: "Italian", NativeName: "Italiano"},
	{ID: 71, Code: [2]byte{'I', 'U'}, Name: "Inuktitut", NativeName: "ᐃᓄᒃᑎᑐᑦ"},
	{ID: 72, Code: [2]byte{'J', 'A'}, Name: "Japanese", NativeName: "日本語 (にほんご／にっぽんご)"},
	{ID: 73, Code: [2]byte{'J', 'V'}, Name: "Javanese", NativeName: "basa Jawa"},
	{ID: 74, Code: [2]byte{'K', 'L'}, Name: "Kalaallisut,Greenlandic", NativeName: "kalaallisut,kalaallit oqaasii"},
	{ID: 75, Code: [2]byte{'K', 'N'}, Name: "Kannada", NativeName: "ಕನ್ನಡ"},
	{ID: 76, Code: [2]byte{'K', 'R'}, Name: "Kanuri", NativeName: "Kanuri"},
	{ID: 77, Code: [2]byte{'K', 'S'}, Name: "Kashmiri", NativeName: "कश्मीरी,كشميري‎"},
	{ID: 78, Code: [2]byte{'K', 'K'}, Name: "Kazakh", NativeName: "Қазақ тілі"},
	{ID: 79, Code: [2]byte{'K', 'M'}, Name: "Khmer", NativeName: "ភាសាខ្មែរ"},
	{ID: 80, Code: [2]byte{'K', 'I'}, Name: "Kikuyu,Gikuyu", NativeName: "Gĩkũyũ"},
	{ID: 81, Code: [2]byte{'R', 'W'}, Name: "Kinyarwanda", NativeName: "Ikinyarwanda"},
	{ID: 82, Code: [2]byte{'K', 'Y'}, Name: "Kirghiz,Kyrgyz", NativeName: "кыргыз тили"},
	{ID: 83, Code: [2]byte{'K', 'V'}, Name: "Komi", NativeName: "коми кыв"},
	{ID: 84, Code: [2]byte{'K', 'G'}, Name: "Kongo", NativeName: "KiKongo"},
	{ID: 85, Code: [2]byte{'K', 'O'}, Name: "Korean", NativeName: "한국어 (韓國語),조선말 (朝鮮語)"},
	{ID: 86, Code: [2]byte{'K', 'U'}, Name: "Kurdish", NativeName: "Kurdî,كوردی‎"},
	{ID: 87, Code: [2]byte{'K', 'J'}, Name: "Kwanyama,Kuanyama", NativeName: "Kuanyama"},
	{ID: 88, Code: [2]byte{'L', 'A'}, Name: "Latin", NativeName: "latine,lingua latina"},
	{ID: 89, Code: [2]byte{'L', 'B'}, Name: "Luxembourgish,Letzeburgesch", NativeName: "Lëtzebuergesch"},
	{ID: 90, Code: [2]byte{'L', 'G'}, Name: "Luganda", NativeName: "Luganda"},
	{ID: 91, Code: [2]byte{'L', 'I'}, Name: "Limburgish,Limburgan,Limburger", NativeName: "Limburgs"},
	{ID: 92, Code: [2]byte{'L', 'N'}, Name: "Lingala", NativeName: "Lingála"},
	{ID: 93, Code: [2]byte{'L', 'O'}, Name: "Lao", NativeName: "ພາສາລາວ"},
	{ID: 94, Code: [2]byte{'L', 'T'}, Name: "Lithuanian", NativeName: "lietuvių kalba"},
	{ID: 95, Code: [2]byte{'L', 'U'}, Name: "Luba-Katanga", NativeName: ""},
	{ID: 96, Code: [2]byte{'L', 'V'}, Name: "Latvian", NativeName: "latviešu valoda"},
	{ID: 97, Code: [2]byte{'G', 'V'}, Name: "Manx", NativeName: "Gaelg,Gailck"},
	{ID: 98, Code: [2]byte{'M', 'K'}, Name: "Macedonian", NativeName: "македонски јазик"},
	{ID: 99, Code: [2]byte{'M', 'G'}, Name: "Malagasy", NativeName: "Malagasy fiteny"},
	{ID: 100, Code: [2]byte{'M', 'S'}, Name: "Malay", NativeName: "bahasa Melayu,بهاس ملايو‎"},
	{ID: 101, Code: [2]byte{'M', 'L'}, Name: "Malayalam", NativeName: "മലയാളം"},
	{ID: 102, Code: [2]byte{'M', 'T'}, Name: "Maltese", NativeName: "Malti"},
	{ID: 103, Code: [2]byte{'M', 'I'}, Name: "Māori", NativeName: "te reo Māori"},
	{ID: 104, Code: [2]byte{'M', 'R'}, Name: "Marathi (Marāṭhī)", NativeName: "मराठी"},
	{ID: 105, Code: [2]byte{'M', 'H'}, Name: "Marshallese", NativeName: "Kajin M̧ajeļ"},
	{ID: 106, Code: [2]byte{'M', 'N'}, Name: "Mongolian", NativeName: "монгол"},
	{ID: 107, Code: [2]byte{'N', 'A'}, Name: "Nauru", NativeName: "Ekakairũ Naoero"},
	{ID: 108, Code: [2]byte{'N', 'V'}, Name: "Navajo,Navaho", NativeName: "Diné bizaad,Dinékʼehǰí"},
	{ID: 109, Code: [2]byte{'N', 'B'}, Name: "Norwegian Bokmål", NativeName: "Norsk bokmål"},
	{ID: 110, Code: [2]byte{'N', 'D'}, Name: "North Ndebele", NativeName: "isiNdebele"},
	{ID: 111, Code: [2]byte{'N', 'E'}, Name: "Nepali", NativeName: "नेपाली"},
	{ID: 112, Code: [2]byte{'N', 'G'}, Name: "Ndonga", NativeName: "Owambo"},
	{ID: 113, Code: [2]byte{'N', 'N'}, Name: "Norwegian Nynorsk", NativeName: "Norsk nynorsk"},
	{ID: 114, Code: [2]byte{'N', 'O'}, Name: "Norwegian", NativeName: "Norsk"},
	{ID: 115, Code: [2]byte{'I', 'I'}, Name: "Nuosu", NativeName: "ꆈꌠ꒿ Nuosuhxop"},
	{ID: 116, Code: [2]byte{'N', 'R'}, Name: "South Ndebele", NativeName: "isiNdebele"},
	{ID: 117, Code: [2]byte{'O', 'C'}, Name: "Occitan", NativeName: "Occitan"},
	{ID: 118, Code: [2]byte{'O', 'J'}, Name: "Ojibwe,Ojibwa", NativeName: "ᐊᓂᔑᓈᐯᒧᐎᓐ"},
	{ID: 119, Code: [2]byte{'C', 'U'}, Name: "Old Church Slavonic,Church Slavic,Church Slavonic,Old Bulgarian,Old Slavonic", NativeName: "ѩзыкъ словѣньскъ"},
	{ID: 120, Code: [2]byte{'O', 'M'}, Name: "Oromo", NativeName: "Afaan Oromoo"},
	{ID: 121, Code: [2]byte{'O', 'R'}, Name: "Oriya", NativeName: "ଓଡ଼ିଆ"},
	{ID: 122, Code: [2]byte{'O', 'S'}, Name: "Ossetian,Ossetic", NativeName: "ирон æвзаг"},
	{ID: 123, Code: [2]byte{'P', 'A'}, Name: "Panjabi,Punjabi", NativeName: "ਪੰਜਾਬੀ,پنجابی‎"},
	{ID: 124, Code: [2]byte{'P', 'I'}, Name: "Pāli", NativeName: "पाऴि"},
	{ID: 125, Code: [2]byte{'F', 'A'}, Name: "Persian", NativeName: "فارسی"},
	{ID: 126, Code: [2]byte{'P', 'L'}, Name: "Polish", NativeName: "polski"},
	{ID: 127, Code: [2]byte{'P', 'S'}, Name: "Pashto,Pushto", NativeName: "پښتو"},
	{ID: 128, Code: [2]byte{'P', 'T'}, Name: "Portuguese", NativeName: "Português"},
	{ID: 129, Code: [2]byte{'Q', 'U'}, Name: "Quechua", NativeName: "Runa Simi,Kichwa"},
	{ID: 130, Code: [2]byte{'R', 'M'}, Name: "Romansh", NativeName: "rumantsch grischun"},
	{ID: 131, Code: [2]byte{'R', 'N'}, Name: "Kirundi", NativeName: "kiRundi"},
	{ID: 132, Code: [2]byte{'R', 'O'}, Name: "Romanian,Moldavian,Moldovan", NativeName: "română"},
	{ID: 133, Code: [2]byte{'R', 'U'}, Name: "Russian", NativeName: "русский язык"},
	{ID: 134, Code: [2]byte{'S', 'A'}, Name: "Sanskrit (Saṁskṛta)", NativeName: "संस्कृतम्"},
	{ID: 135, Code: [2]byte{'S', 'C'}, Name: "Sardinian", NativeName: "sardu"},
	{ID: 136, Code: [2]byte{'S', 'D'}, Name: "Sindhi", NativeName: "सिन्धी,سنڌي، سندھی‎"},
	{ID: 137, Code: [2]byte{'S', 'E'}, Name: "Northern Sami", NativeName: "Davvisámegiella"},
	{ID: 138, Code: [2]byte{'S', 'M'}, Name: "Samoan", NativeName: "gagana faa Samoa"},
	{ID: 139, Code: [2]byte{'S', 'G'}, Name: "Sango", NativeName: "yângâ tî sängö"},
	{ID: 140, Code: [2]byte{'S', 'R'}, Name: "Serbian", NativeName: "српски језик"},
	{ID: 141, Code: [2]byte{'G', 'D'}, Name: "Scottish Gaelic; Gaelic", NativeName: "GàIDhlig"},
	{ID: 142, Code: [2]byte{'S', 'N'}, Name: "Shona", NativeName: "chiShona"},
	{ID: 143, Code: [2]byte{'S', 'I'}, Name: "Sinhala,Sinhalese", NativeName: "සිංහල"},
	{ID: 144, Code: [2]byte{'S', 'K'}, Name: "Slovak", NativeName: "slovenčina"},
	{ID: 145, Code: [2]byte{'S', 'L'}, Name: "Slovene", NativeName: "slovenščina"},
	{ID: 146, Code: [2]byte{'S', 'O'}, Name: "Somali", NativeName: "Soomaaliga,af Soomaali"},
	{ID: 147, Code: [2]byte{'S', 'T'}, Name: "Southern Sotho", NativeName: "Sesotho"},
	{ID: 148, Code: [2]byte{'E', 'S'}, Name: "Spanish; Castilian", NativeName: "español,castellano"},
	{ID: 149, Code: [2]byte{'S', 'U'}, Name: "Sundanese", NativeName: "Basa Sunda"},
	{ID: 150, Code: [2]byte{'S', 'W'}, Name: "Swahili", NativeName: "Kiswahili"},
	{ID: 151, Code: [2]byte{'S', 'S'}, Name: "Swati", NativeName: "SiSwati"},
	{ID: 152, Code: [2]byte{'S', 'V'}, Name: "Swedish", NativeName: "svenska"},
	{ID: 153, Code: [2]byte{'T', 'A'}, Name: "Tamil", NativeName: "தமிழ்"},
	{ID: 154, Code: [2]byte{'T', 'E'}, Name: "Telugu", NativeName: "తెలుగు"},
	{ID: 155, Code: [2]byte{'T', 'G'}, Name: "Tajik", NativeName: "тоҷикӣ,toğikī,تاجیکی‎"},
	{ID: 156, Code: [2]byte{'T', 'H'}, Name: "Thai", NativeName: "ไทย"},
	{ID: 157, Code: [2]byte{'T', 'I'}, Name: "Tigrinya", NativeName: "ትግርኛ"},
	{ID: 158, Code: [2]byte{'B', 'O'}, Name: "Tibetan Standard,Tibetan,Central", NativeName: "བོད་ཡིག"},
	{ID: 159, Code: [2]byte{'T', 'K'}, Name: "Turkmen", NativeName: "Türkmen,Түркмен"},
	{ID: 160, Code: [2]byte{'T', 'L'}, Name: "Tagalog", NativeName: "Wikang Tagalog,ᜏᜒᜃᜅ᜔ ᜆᜄᜎᜓᜄ᜔"},
	{ID: 161, Code: [2]byte{'T', 'N'}, Name: "Tswana", NativeName: "Setswana"},
	{ID: 162, Code: [2]byte{'T', 'O'}, Name: "Tonga (Tonga Islands)", NativeName: "faka Tonga"},
	{ID: 163, Code: [2]byte{'T', 'R'}, Name: "Turkish", NativeName: "Türkçe"},
	{ID: 164, Code: [2]byte{'T', 'S'}, Name: "Tsonga", NativeName: "Xitsonga"},
	{ID: 165, Code: [2]byte{'T', 'T'}, Name: "Tatar", NativeName: "татарча,tatarça,تاتارچا‎"},
	{ID: 166, Code: [2]byte{'T', 'W'}, Name: "Twi", NativeName: "Twi"},
	{ID: 167, Code: [2]byte{'T', 'Y'}, Name: "Tahitian", NativeName: "Reo Tahiti"},
	{ID: 168, Code: [2]byte{'U', 'G'}, Name: "Uighur,Uyghur", NativeName: "Uyƣurqə,ئۇيغۇرچە‎"},
	{ID: 169, Code: [2]byte{'U', 'K'}, Name: "Ukrainian", NativeName: "українська"},
	{ID: 170, Code: [2]byte{'U', 'R'}, Name: "Urdu", NativeName: "اردو"},
	{ID: 171, Code: [2]byte{'U', 'Z'}, Name: "Uzbek", NativeName: "zbek,Ўзбек,أۇزبېك‎"},
	{ID: 172, Code: [2]byte{'V', 'E'}, Name: "Venda", NativeName: "Tshivenḓa"},
	{ID: 173, Code: [2]byte{'V', 'I'}, Name: "Vietnamese", NativeName: "Tiếng Việt"},
	{ID: 174, Code: [2]byte{'V', 'O'}, Name: "Volapük", NativeName: "Volapük"},
	{ID: 175, Code: [2]byte{'W', 'A'}, Name: "Walloon", NativeName: "Walon"},
	{ID: 176, Code: [2]byte{'C', 'Y'}, Name: "Welsh", NativeName: "Cymraeg"},
	{ID: 177, Code: [2]byte{'W', 'O'}, Name: "Wolof", NativeName: "Wollof"},
	{ID: 178, Code: [2]byte{'F', 'Y'}, Name: "Western Frisian", NativeName: "Frysk"},
	{ID: 179, Code: [2]byte{'X', 'H'}, Name: "Xhosa", NativeName: "isiXhosa"},
	{ID: 180, Code: [2]byte{'Y', 'I'}, Name: "YIDdish", NativeName: "ייִדיש"},
	{ID: 181, Code: [2]byte{'Y', 'O'}, Name: "Yoruba", NativeName: "Yorùbá"},
	{ID: 182, Code: [2]byte{'Z', 'A'}, Name: "Zhuang,Chuang", NativeName: "Saɯ cueŋƅ,Saw cuengh"},
}

var (
	codeLanguages map[uint]uint
	IDLanguages   map[uint]uint
)

func init() {
	codeLanguages = make(map[uint]uint)
	IDLanguages = make(map[uint]uint)

	for i, lg := range Languages {
		codeLanguages[lg.IntCode()] = uint(i)
		IDLanguages[lg.ID] = uint(i)
	}
}

///////////////////////////////////////////////////////////////////////////////
/// Global methods
///////////////////////////////////////////////////////////////////////////////

func GetLanguageByCodeString(code string) *Language {
	if len(code) < 2 {
		return nil
	}
	code = strings.ToUpper(code)
	return GetLanguageByCode([2]byte{code[0], code[1]})
}

func GetLanguageByCode(code [2]byte) *Language {
	if index, ok := codeLanguages[IntCode(code)]; ok {
		return &Languages[index]
	}
	return &Languages[0]
}

func GetLanguageIdByCode(code [2]byte) uint {
	if lang := GetLanguageByCode(code); nil != lang {
		return lang.ID
	}
	return 0
}

func GetLanguageIdByCodeString(code string) uint {
	if len(code) < 2 {
		return 0
	}
	code = strings.ToUpper(code)
	return GetLanguageIdByCode([2]byte{code[0], code[1]})
}

func GetLanguageByID(ID uint) *Language {
	if index, ok := IDLanguages[ID]; ok {
		return &Languages[index]
	}
	return nil
}

///////////////////////////////////////////////////////////////////////////////
/// Helpers
///////////////////////////////////////////////////////////////////////////////

func IntCode(geoCode [2]byte) uint {
	var code uint = (uint)(geoCode[0])
	code |= (uint)(geoCode[1]) << 8
	return code
}
