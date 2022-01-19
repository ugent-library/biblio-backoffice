package engine

func (e *Engine) Vocabularies() map[string][]string {
	return vocabularies
}

var vocabularies = map[string][]string{
	"cc_licenses": {
		"cc_0",
		"cc_by",
		"cc_by_sa",
		"cc_by_nd",
		"cc_by_nc",
		"cc_by_nc_sa",
		"cc_by_nc_nd",
	},
	"confirmations": {
		"yes",
		"no",
		"dontknow",
	},
	"credit_roles": {
		"first_author",
		"last_author",
		"conceptualization",
		"data_curation",
		"formal_analysis",
		"funding_acquisition",
		"investigation",
		"methodology",
		"project_administration",
		"resources",
		"software",
		"supervision",
		"validation",
		"visualization",
		"writing_original_draft",
		"writing_review_editing",
	},
	"access_levels": {
		"info:eu-repo/semantics/openAccess",
		"info:eu-repo/semantics/restrictedAccess",
		"info:eu-repo/semantics/embargoedAccess",
		"info:eu-repo/semantics/closedAccess",
	},
	"publication_types": {
		"journal_article",
		"book",
		"book_chapter",
		"book_editor",
		"issue_editor",
		"conference",
		"dissertation",
		"miscellaneous",
	},
	"publication_statuses": {
		"new",
		"private",
		"returned",
		"public",
		"deleted",
	},
	"publication_publishing_statuses": {
		"unpublished",
		"accepted",
		"published",
	},
	"publication_versions": {
		"publishedVersion",
		"authorVersion",
		"acceptedVersion",
		"updatedVersion",
	},
	"publication_sorts": {
		"date-updated-desc",
		"date-created-desc",
		"year-desc",
	},
	"publication_file_access_levels": {
		"open_access",
		"local",
		"closed",
	},
	"publication_file_relations": {
		"main_file",
		"colophon",
		"data_fact_sheet",
		"peer_review_report",
		"table_of_contents",
		"agreement",
	},
	"publication_link_relations": {
		"data_management_plan",
		"homepage",
		"peer_review_report",
		"related_information",
		"software",
		"table_of_contents",
		"main_file",
	},
	"conference_types": {
		"proceedingsPaper",
		"abstract",
		"poster",
		"other",
	},
	"journal_article_types": {
		"original",
		"review",
		"letter_note",
		"proceedingsPaper",
	},
	"miscellaneous_types": {
		"artReview",
		"artisticWork",
		"bibliography",
		"biography",
		"blogPost",
		"bookReview",
		"correction",
		"dictionaryEntry",
		"editorialMaterial",
		"encyclopediaEntry",
		"exhibitionReview",
		"filmReview",
		"lectureSpeech",
		"lemma",
		"magazinePiece",
		"manual",
		"musicEdition",
		"musicReview",
		"newsArticle",
		"newspaperPiece",
		"other",
		"preprint",
		"productReview",
		"report",
		"technicalStandard",
		"textEdition",
		"textTranslation",
		"theatreReview",
		"workingPaper",
	},
	"research_fields": {
		"Agriculture and Food Sciences",
		"Arts and Architecture",
		"Biology and Life Sciences",
		"Business and Economics",
		"Chemistry",
		"Cultural Sciences",
		"Earth and Environmental Sciences",
		"General Works",
		"History and Archaeology",
		"Languages and Literatures",
		"Law and Political Science",
		"Mathematics and Statistics",
		"Medicine and Health Sciences",
		"Performing Arts",
		"Philosophy and Religion",
		"Physics and Astronomy",
		"Science General",
		"Social Sciences",
		"Technology and Engineering",
		"Veterinary Sciences",
	},
	"language_codes": {
		"und",
		"eng",
		"dut",
		"fre",
		"ger",
		"aar",
		"abk",
		"ace",
		"ach",
		"ada",
		"ady",
		"afa",
		"afh",
		"afr",
		"ain",
		"aka",
		"akk",
		"alb",
		"ale",
		"alg",
		"alt",
		"amh",
		"ang",
		"anp",
		"apa",
		"ara",
		"arc",
		"arg",
		"arm",
		"arn",
		"arp",
		"art",
		"arw",
		"asm",
		"ast",
		"ath",
		"aus",
		"ava",
		"ave",
		"awa",
		"aym",
		"aze",
		"bad",
		"bai",
		"bak",
		"bal",
		"bam",
		"ban",
		"baq",
		"bas",
		"bat",
		"bej",
		"bel",
		"bem",
		"ben",
		"ber",
		"bho",
		"bih",
		"bik",
		"bin",
		"bis",
		"bla",
		"bnt",
		"tib",
		"bos",
		"bra",
		"bre",
		"btk",
		"bua",
		"bug",
		"bul",
		"bur",
		"byn",
		"cad",
		"cai",
		"car",
		"cat",
		"cau",
		"ceb",
		"cel",
		"cze",
		"cha",
		"chb",
		"che",
		"chg",
		"chi",
		"chk",
		"chm",
		"chn",
		"cho",
		"chp",
		"chr",
		"chu",
		"chv",
		"chy",
		"cmc",
		"cop",
		"cor",
		"cos",
		"cpe",
		"cpf",
		"cpp",
		"cre",
		"crh",
		"crp",
		"csb",
		"cus",
		"wel",
		"cze",
		"dak",
		"dan",
		"dar",
		"day",
		"del",
		"den",
		"ger",
		"dgr",
		"din",
		"div",
		"doi",
		"dra",
		"dsb",
		"dua",
		"dum",
		"dyu",
		"dzo",
		"efi",
		"egy",
		"eka",
		"gre",
		"elx",
		"enm",
		"epo",
		"est",
		"baq",
		"ewe",
		"ewo",
		"fan",
		"fao",
		"per",
		"fat",
		"fij",
		"fil",
		"fin",
		"fiu",
		"fon",
		"frm",
		"fro",
		"frr",
		"frs",
		"fry",
		"ful",
		"fur",
		"gaa",
		"gay",
		"gba",
		"gem",
		"geo",
		"gez",
		"gil",
		"gla",
		"gle",
		"glg",
		"glv",
		"gmh",
		"goh",
		"gon",
		"gor",
		"got",
		"grb",
		"grc",
		"gre",
		"grn",
		"gsw",
		"guj",
		"gwi",
		"hai",
		"hat",
		"hau",
		"haw",
		"heb",
		"her",
		"hil",
		"him",
		"hin",
		"hit",
		"hmn",
		"hmo",
		"hrv",
		"hsb",
		"hun",
		"hup",
		"arm",
		"iba",
		"ibo",
		"ice",
		"ido",
		"iii",
		"ijo",
		"iku",
		"ile",
		"ilo",
		"ina",
		"inc",
		"ind",
		"ine",
		"inh",
		"ipk",
		"ira",
		"iro",
		"ice",
		"ita",
		"jav",
		"jbo",
		"jpn",
		"jpr",
		"jrb",
		"kaa",
		"kab",
		"kac",
		"kal",
		"kam",
		"kan",
		"kar",
		"kas",
		"geo",
		"kau",
		"kaw",
		"kaz",
		"kbd",
		"kha",
		"khi",
		"khm",
		"kho",
		"kik",
		"kin",
		"kir",
		"kmb",
		"kok",
		"kom",
		"kon",
		"kor",
		"kos",
		"kpe",
		"krc",
		"krl",
		"kro",
		"kru",
		"kua",
		"kum",
		"kur",
		"kut",
		"lad",
		"lah",
		"lam",
		"lao",
		"lat",
		"lav",
		"lez",
		"lim",
		"lin",
		"lit",
		"lol",
		"loz",
		"ltz",
		"lua",
		"lub",
		"lug",
		"lui",
		"lun",
		"luo",
		"lus",
		"mac",
		"mad",
		"mag",
		"mah",
		"mai",
		"mak",
		"mal",
		"man",
		"mao",
		"map",
		"mar",
		"mas",
		"may",
		"mdf",
		"mdr",
		"men",
		"mga",
		"mic",
		"min",
		"mis",
		"mac",
		"mkh",
		"mlg",
		"mlt",
		"mnc",
		"mni",
		"mno",
		"moh",
		"mon",
		"mos",
		"mao",
		"may",
		"mul",
		"mun",
		"mus",
		"mwl",
		"mwr",
		"bur",
		"myn",
		"myv",
		"nah",
		"nai",
		"nap",
		"nau",
		"nav",
		"nbl",
		"nde",
		"ndo",
		"nds",
		"nep",
		"new",
		"nia",
		"nic",
		"niu",
		"dut",
		"nno",
		"nob",
		"nog",
		"non",
		"nor",
		"nqo",
		"nso",
		"nub",
		"nwc",
		"nya",
		"nym",
		"nyn",
		"nyo",
		"nzi",
		"oci",
		"oji",
		"ori",
		"orm",
		"osa",
		"oss",
		"ota",
		"oto",
		"paa",
		"pag",
		"pal",
		"pam",
		"pan",
		"pap",
		"pau",
		"peo",
		"per",
		"phi",
		"phn",
		"pli",
		"pol",
		"pon",
		"por",
		"pra",
		"pro",
		"pus",
		"qaa",
		"que",
		"raj",
		"rap",
		"rar",
		"roa",
		"roh",
		"rom",
		"rum",
		"rum",
		"run",
		"rup",
		"rus",
		"sad",
		"sag",
		"sah",
		"sai",
		"sal",
		"sam",
		"san",
		"sas",
		"sat",
		"scn",
		"sco",
		"sel",
		"sem",
		"sga",
		"sgn",
		"shn",
		"sid",
		"sin",
		"sio",
		"sit",
		"sla",
		"slo",
		"slo",
		"slv",
		"sma",
		"sme",
		"smi",
		"smj",
		"smn",
		"smo",
		"sms",
		"sna",
		"snd",
		"snk",
		"sog",
		"som",
		"son",
		"sot",
		"spa",
		"alb",
		"srd",
		"srn",
		"srp",
		"srr",
		"ssa",
		"ssw",
		"suk",
		"sun",
		"sus",
		"sux",
		"swa",
		"swe",
		"syc",
		"syr",
		"tah",
		"tai",
		"tam",
		"tat",
		"tel",
		"tem",
		"ter",
		"tet",
		"tgk",
		"tgl",
		"tha",
		"tib",
		"tig",
		"tir",
		"tiv",
		"tkl",
		"tlh",
		"tli",
		"tmh",
		"tog",
		"ton",
		"tpi",
		"tsi",
		"tsn",
		"tso",
		"tuk",
		"tum",
		"tup",
		"tur",
		"tut",
		"tvl",
		"twi",
		"tyv",
		"udm",
		"uga",
		"uig",
		"ukr",
		"umb",
		"urd",
		"uzb",
		"vai",
		"ven",
		"vie",
		"vol",
		"vot",
		"wak",
		"wal",
		"war",
		"was",
		"wel",
		"wen",
		"wln",
		"wol",
		"xal",
		"xho",
		"yao",
		"yap",
		"yid",
		"yor",
		"ypk",
		"zap",
		"zbl",
		"zen",
		"zha",
		"chi",
		"znd",
		"zul",
		"zun",
		"zxx",
		"zza",
	},
}
