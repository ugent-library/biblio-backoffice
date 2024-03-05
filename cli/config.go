package cli

// Version info
type Version struct {
	Branch string `env:"SOURCE_BRANCH"`
	Commit string `env:"SOURCE_COMMIT"`
	Image  string `env:"IMAGE_NAME"`
}

// Application config
type Config struct {
	// Env must be local, development, test or production
	// TODO rename MODE to ENV
	Env              string `env:"MODE" envDefault:"production"`
	BaseURL          string `env:"BASE_URL"`
	Timezone         string `env:"TIMEZONE" envDefault:"Europe/Brussels"`
	TokenSecret      string `env:"TOKEN_SECRET,notEmpty"`
	IndexRetention   int    `env:"INDEX_RETENTION" envDefault:"2"`
	PgConn           string `env:"PG_CONN,notEmpty"`
	Es6URL           string `env:"ES6_URL"`
	PublicationIndex string `env:"PUBLICATION_INDEX"`
	DatasetIndex     string `env:"DATASET_INDEX"`
	Host             string `env:"HOST"`
	Port             int    `env:"PORT" envDefault:"3000"`
	API              struct {
		Host string `env:"HOST"`
		Port int    `env:"PORT" envDefault:"30000"`
	} `envPrefix:"API_"`
	AdminUsername   string `env:"ADMIN_USERNAME"`
	AdminPassword   string `env:"ADMIN_PASSWORD"`
	CuratorUsername string `env:"CURATOR_USERNAME"`
	CuratorPassword string `env:"CURATOR_PASSWORD"`
	IPRanges        string `env:"IP_RANGES"`
	S3              struct {
		Endpoint   string `env:"ENDPOINT"`
		Region     string `env:"REGION" envDefault:"us-east-1"`
		ID         string `env:"ID"`
		Secret     string `env:"SECRET"`
		Bucket     string `env:"BUCKET"`
		TempBucket string `env:"TEMP_BUCKET"`
	} `envPrefix:"S3_"`
	Session struct {
		Name   string `env:"NAME" envDefault:"biblio-backoffice"`
		Secret string `env:"SECRET"`
		MaxAge int    `env:"MAX_AGE" envDefault:"2592000"` // default 30 days
	} `envPrefix:"SESSION_"`
	CSRF struct {
		Name   string `env:"NAME" envDefault:"biblio-backoffice.csrf-token"`
		Secret string `env:"SECRET"`
	} `envPrefix:"CSRF_"`
	MaxFileSize int    `env:"MAX_FILE_SIZE" envDefault:"2000000000"`
	FileDir     string `env:"FILE_DIR"`
	Frontend    struct {
		URL      string `env:"URL"`
		Username string `env:"USERNAME"`
		Password string `env:"PASSWORD"`
		Es6URL   string `env:"ES6_URL"`
	} `envPrefix:"FRONTEND_"`
	ORCID struct {
		ClientID     string `env:"CLIENT_ID"`
		ClientSecret string `env:"CLIENT_SECRET"`
		Sandbox      bool   `env:"SANDBOX"`
	} `envPrefix:"ORCID_"`
	OIDC struct {
		URL          string `env:"URL"`
		ClientID     string `env:"CLIENT_ID"`
		ClientSecret string `env:"CLIENT_SECRET"`
	} `envPrefix:"OIDC_"`
	CiteprocURL string `env:"CITEPROC_URL"`
	MongoDBURL  string `env:"MONGODB_URL"`
	APIKey      string `env:"API_KEY"`
	Handle      struct {
		Enabled  bool   `env:"ENABLED"`
		URL      string `env:"URL"`
		Prefix   string `env:"PREFIX"`
		Username string `env:"USERNAME"`
		Password string `env:"PASSWORD"`
	} `envPrefix:"HDL_SRV_"`
	OAI struct {
		APIURL string `env:"API_URL"`
		APIKey string `env:"API_KEY"`
	} `envPrefix:"OAI_"`
	// Feature flags
	FF struct {
		FilePath     string `env:"FILE_PATH"`
		GitHubToken  string `env:"GITHUB_TOKEN"`
		GitHubRepo   string `env:"GITHUB_REPO"`
		GitHubBranch string `env:"GITHUB_BRANCH" envDefault:"main"`
		GitHubPath   string `env:"GITHUB_PATH"`
	} `envPrefix:"FF_"`
}
