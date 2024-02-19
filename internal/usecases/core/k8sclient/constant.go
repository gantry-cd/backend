package k8sclient

// gantryで扱うラベルの定義
const (
	// AppLabel は ns やリソースに付与するラベルのキー。
	AppLabel = "app"
	// CreatedByLabel は ns やリソースに付与するラベルのキー。
	// このラベルが付与されたリソースは gantry によって作成されたものとみなす
	CreatedByLabel = "created"
	// RepositoryLabel はリソースに付与するラベルのキー。
	// このラベルが付与されたリソースはどのリポジトリに関連するかを示す
	RepositoryLabel = "repository"
	// PrIDLabel はリソースに付与するラベルのキー。
	// このラベルが付与されたリソースはどの PR に関連するかを示す
	PullRequestID = "pr-id"
	// BaseBranchLabel はリソースに付与するラベルのキー。
	// このラベルが付与されたリソースはどのブランチに関連するかを示す
	BaseBranchLabel = "base-branch"
	// EnvirionmentLabel はリソースに付与するラベルのキー。
	// このラベルが付与されたリソースはどの環境に関連するかを示す
	// 基本Botから生成されるデプロイにはpreviewが付与される
	// ex: staging , production , preview
	EnvirionmentLabel = "envirionment"
)

const (
	// AppIdentifier は gantry が扱うリソースの識別子
	AppIdentifier = "gantry"

	// EnvProduction は本番環境を示す
	EnvProduction = "production"
	// EnvStaging はステージング環境を示す
	EnvStaging = "staging"
	// EnvPreview はプレビュー環境を示す
	EnvPreview = "preview"
)
