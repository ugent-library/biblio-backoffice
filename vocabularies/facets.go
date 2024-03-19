package vocabularies

var Facets = map[string][][]string{
	"publication": {
		{
			"status",
			"classification",
			"faculty_id",
			"year",
			"type",
		},
		{
			"has_message",
			"locked",
			"has_files",
			"vabb_type",
			"created_since",
			"updated_since",
		},
	},
	"publication_curation": {
		{
			"status",
			"classification",
			"faculty_id",
			"year",
			"type",
		},
		{
			"publication_status",
			"reviewer_tags",
			"has_message",
			"locked",
			"extern",
		},
		{
			"wos_type",
			"vabb_type",
			"has_files",
			"file_relation",
			"created_since",
			"updated_since",
			"legacy",
		},
	},
	"dataset": {
		{
			"status",
			"faculty_id",
			"locked",
			"has_message",
			"created_since",
			"updated_since",
		},
	},
	"dataset_curation": {
		{
			"status",
			"faculty_id",
			"locked",
			"identifier_type",
		},
		{
			"reviewer_tags",
			"has_message",
			"created_since",
			"updated_since",
		},
	},
}
