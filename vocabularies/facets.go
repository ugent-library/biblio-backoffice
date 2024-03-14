package vocabularies

var Facets = map[string][][]string{
	"publication": {
		{

			"status",
			"locked",
			"faculty_id",
			"type",
			"wos_type",
			"classification",
			"year",
			"vabb_type",
			"has_files",
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
			"locked",
			"faculty_id",
			"created_since",
			"updated_since",
		},
	},
	"dataset_curation": {
		{
			"status",
			"locked",
			"faculty_id",
			"reviewer_tags",
			"has_message",
			"created_since",
			"updated_since",
			"identifier_type",
		},
	},
}
