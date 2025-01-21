package docs

// @title Template API
// @version 1.0
// @description A template API with multitenancy support
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey OrganizationID
// @in header
// @name organization_id
// @description Organization ID for multitenancy

// @securityDefinitions.apikey BranchID
// @in header
// @name branch_id
// @description Branch ID for multitenancy

// @SecurityRequirement BearerAuth
// @SecurityRequirement OrganizationID
// @SecurityRequirement BranchID

func SwaggerDocs() {
	// This function exists only for Swagger annotations
	// The actual Swagger docs are generated from these annotations
}
