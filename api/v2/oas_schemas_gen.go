// Code generated by ogen, DO NOT EDIT.

package api

import (
	"fmt"
	"time"
)

func (s *ErrorStatusCode) Error() string {
	return fmt.Sprintf("code %d: %+v", s.StatusCode, s.Response)
}

// AddPersonOK is response for AddPerson operation.
type AddPersonOK struct{}

// Ref: #/components/schemas/AddPersonRequest
type AddPersonRequest struct {
	Person AddPersonRequestPerson `json:"person"`
}

// GetPerson returns the value of Person.
func (s *AddPersonRequest) GetPerson() AddPersonRequestPerson {
	return s.Person
}

// SetPerson sets the value of Person.
func (s *AddPersonRequest) SetPerson(val AddPersonRequestPerson) {
	s.Person = val
}

type AddPersonRequestPerson struct {
	Identifiers         []Identifier `json:"identifiers"`
	Name                string       `json:"name"`
	PreferredName       OptString    `json:"preferredName"`
	GivenName           OptString    `json:"givenName"`
	PreferredGivenName  OptString    `json:"preferredGivenName"`
	FamilyName          OptString    `json:"familyName"`
	PreferredFamilyName OptString    `json:"preferredFamilyName"`
	HonorificPrefix     OptString    `json:"honorificPrefix"`
	Email               OptString    `json:"email"`
	Username            OptString    `json:"username"`
	Active              OptBool      `json:"active"`
	Attributes          []Attribute  `json:"attributes"`
}

// GetIdentifiers returns the value of Identifiers.
func (s *AddPersonRequestPerson) GetIdentifiers() []Identifier {
	return s.Identifiers
}

// GetName returns the value of Name.
func (s *AddPersonRequestPerson) GetName() string {
	return s.Name
}

// GetPreferredName returns the value of PreferredName.
func (s *AddPersonRequestPerson) GetPreferredName() OptString {
	return s.PreferredName
}

// GetGivenName returns the value of GivenName.
func (s *AddPersonRequestPerson) GetGivenName() OptString {
	return s.GivenName
}

// GetPreferredGivenName returns the value of PreferredGivenName.
func (s *AddPersonRequestPerson) GetPreferredGivenName() OptString {
	return s.PreferredGivenName
}

// GetFamilyName returns the value of FamilyName.
func (s *AddPersonRequestPerson) GetFamilyName() OptString {
	return s.FamilyName
}

// GetPreferredFamilyName returns the value of PreferredFamilyName.
func (s *AddPersonRequestPerson) GetPreferredFamilyName() OptString {
	return s.PreferredFamilyName
}

// GetHonorificPrefix returns the value of HonorificPrefix.
func (s *AddPersonRequestPerson) GetHonorificPrefix() OptString {
	return s.HonorificPrefix
}

// GetEmail returns the value of Email.
func (s *AddPersonRequestPerson) GetEmail() OptString {
	return s.Email
}

// GetUsername returns the value of Username.
func (s *AddPersonRequestPerson) GetUsername() OptString {
	return s.Username
}

// GetActive returns the value of Active.
func (s *AddPersonRequestPerson) GetActive() OptBool {
	return s.Active
}

// GetAttributes returns the value of Attributes.
func (s *AddPersonRequestPerson) GetAttributes() []Attribute {
	return s.Attributes
}

// SetIdentifiers sets the value of Identifiers.
func (s *AddPersonRequestPerson) SetIdentifiers(val []Identifier) {
	s.Identifiers = val
}

// SetName sets the value of Name.
func (s *AddPersonRequestPerson) SetName(val string) {
	s.Name = val
}

// SetPreferredName sets the value of PreferredName.
func (s *AddPersonRequestPerson) SetPreferredName(val OptString) {
	s.PreferredName = val
}

// SetGivenName sets the value of GivenName.
func (s *AddPersonRequestPerson) SetGivenName(val OptString) {
	s.GivenName = val
}

// SetPreferredGivenName sets the value of PreferredGivenName.
func (s *AddPersonRequestPerson) SetPreferredGivenName(val OptString) {
	s.PreferredGivenName = val
}

// SetFamilyName sets the value of FamilyName.
func (s *AddPersonRequestPerson) SetFamilyName(val OptString) {
	s.FamilyName = val
}

// SetPreferredFamilyName sets the value of PreferredFamilyName.
func (s *AddPersonRequestPerson) SetPreferredFamilyName(val OptString) {
	s.PreferredFamilyName = val
}

// SetHonorificPrefix sets the value of HonorificPrefix.
func (s *AddPersonRequestPerson) SetHonorificPrefix(val OptString) {
	s.HonorificPrefix = val
}

// SetEmail sets the value of Email.
func (s *AddPersonRequestPerson) SetEmail(val OptString) {
	s.Email = val
}

// SetUsername sets the value of Username.
func (s *AddPersonRequestPerson) SetUsername(val OptString) {
	s.Username = val
}

// SetActive sets the value of Active.
func (s *AddPersonRequestPerson) SetActive(val OptBool) {
	s.Active = val
}

// SetAttributes sets the value of Attributes.
func (s *AddPersonRequestPerson) SetAttributes(val []Attribute) {
	s.Attributes = val
}

// AddProjectOK is response for AddProject operation.
type AddProjectOK struct{}

// Ref: #/components/schemas/AddProjectRequest
type AddProjectRequest struct {
	Project AddProjectRequestProject `json:"project"`
}

// GetProject returns the value of Project.
func (s *AddProjectRequest) GetProject() AddProjectRequestProject {
	return s.Project
}

// SetProject sets the value of Project.
func (s *AddProjectRequest) SetProject(val AddProjectRequestProject) {
	s.Project = val
}

type AddProjectRequestProject struct {
	Identifiers  []Identifier `json:"identifiers"`
	Names        []Text       `json:"names"`
	Descriptions []Text       `json:"descriptions"`
	StartDate    OptString    `json:"startDate"`
	EndDate      OptString    `json:"endDate"`
	Deleted      OptBool      `json:"deleted"`
	Attributes   []Attribute  `json:"attributes"`
}

// GetIdentifiers returns the value of Identifiers.
func (s *AddProjectRequestProject) GetIdentifiers() []Identifier {
	return s.Identifiers
}

// GetNames returns the value of Names.
func (s *AddProjectRequestProject) GetNames() []Text {
	return s.Names
}

// GetDescriptions returns the value of Descriptions.
func (s *AddProjectRequestProject) GetDescriptions() []Text {
	return s.Descriptions
}

// GetStartDate returns the value of StartDate.
func (s *AddProjectRequestProject) GetStartDate() OptString {
	return s.StartDate
}

// GetEndDate returns the value of EndDate.
func (s *AddProjectRequestProject) GetEndDate() OptString {
	return s.EndDate
}

// GetDeleted returns the value of Deleted.
func (s *AddProjectRequestProject) GetDeleted() OptBool {
	return s.Deleted
}

// GetAttributes returns the value of Attributes.
func (s *AddProjectRequestProject) GetAttributes() []Attribute {
	return s.Attributes
}

// SetIdentifiers sets the value of Identifiers.
func (s *AddProjectRequestProject) SetIdentifiers(val []Identifier) {
	s.Identifiers = val
}

// SetNames sets the value of Names.
func (s *AddProjectRequestProject) SetNames(val []Text) {
	s.Names = val
}

// SetDescriptions sets the value of Descriptions.
func (s *AddProjectRequestProject) SetDescriptions(val []Text) {
	s.Descriptions = val
}

// SetStartDate sets the value of StartDate.
func (s *AddProjectRequestProject) SetStartDate(val OptString) {
	s.StartDate = val
}

// SetEndDate sets the value of EndDate.
func (s *AddProjectRequestProject) SetEndDate(val OptString) {
	s.EndDate = val
}

// SetDeleted sets the value of Deleted.
func (s *AddProjectRequestProject) SetDeleted(val OptBool) {
	s.Deleted = val
}

// SetAttributes sets the value of Attributes.
func (s *AddProjectRequestProject) SetAttributes(val []Attribute) {
	s.Attributes = val
}

// Ref: #/components/schemas/Affiliation
type Affiliation struct {
	Organization Organization `json:"organization"`
}

// GetOrganization returns the value of Organization.
func (s *Affiliation) GetOrganization() Organization {
	return s.Organization
}

// SetOrganization sets the value of Organization.
func (s *Affiliation) SetOrganization(val Organization) {
	s.Organization = val
}

// Ref: #/components/schemas/AffiliationParams
type AffiliationParams struct {
	OrganizationIdentifier Identifier `json:"organizationIdentifier"`
}

// GetOrganizationIdentifier returns the value of OrganizationIdentifier.
func (s *AffiliationParams) GetOrganizationIdentifier() Identifier {
	return s.OrganizationIdentifier
}

// SetOrganizationIdentifier sets the value of OrganizationIdentifier.
func (s *AffiliationParams) SetOrganizationIdentifier(val Identifier) {
	s.OrganizationIdentifier = val
}

type ApiKey struct {
	APIKey string
}

// GetAPIKey returns the value of APIKey.
func (s *ApiKey) GetAPIKey() string {
	return s.APIKey
}

// SetAPIKey sets the value of APIKey.
func (s *ApiKey) SetAPIKey(val string) {
	s.APIKey = val
}

// Ref: #/components/schemas/Attribute
type Attribute struct {
	Scope string `json:"scope"`
	Key   string `json:"key"`
	Value string `json:"value"`
}

// GetScope returns the value of Scope.
func (s *Attribute) GetScope() string {
	return s.Scope
}

// GetKey returns the value of Key.
func (s *Attribute) GetKey() string {
	return s.Key
}

// GetValue returns the value of Value.
func (s *Attribute) GetValue() string {
	return s.Value
}

// SetScope sets the value of Scope.
func (s *Attribute) SetScope(val string) {
	s.Scope = val
}

// SetKey sets the value of Key.
func (s *Attribute) SetKey(val string) {
	s.Key = val
}

// SetValue sets the value of Value.
func (s *Attribute) SetValue(val string) {
	s.Value = val
}

// Ref: #/components/schemas/Error
type Error struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

// GetCode returns the value of Code.
func (s *Error) GetCode() int64 {
	return s.Code
}

// GetMessage returns the value of Message.
func (s *Error) GetMessage() string {
	return s.Message
}

// SetCode sets the value of Code.
func (s *Error) SetCode(val int64) {
	s.Code = val
}

// SetMessage sets the value of Message.
func (s *Error) SetMessage(val string) {
	s.Message = val
}

// ErrorStatusCode wraps Error with StatusCode.
type ErrorStatusCode struct {
	StatusCode int
	Response   Error
}

// GetStatusCode returns the value of StatusCode.
func (s *ErrorStatusCode) GetStatusCode() int {
	return s.StatusCode
}

// GetResponse returns the value of Response.
func (s *ErrorStatusCode) GetResponse() Error {
	return s.Response
}

// SetStatusCode sets the value of StatusCode.
func (s *ErrorStatusCode) SetStatusCode(val int) {
	s.StatusCode = val
}

// SetResponse sets the value of Response.
func (s *ErrorStatusCode) SetResponse(val Error) {
	s.Response = val
}

func (*ErrorStatusCode) getOrganizationRes()     {}
func (*ErrorStatusCode) getPersonRes()           {}
func (*ErrorStatusCode) importOrganizationsRes() {}
func (*ErrorStatusCode) importPersonRes()        {}
func (*ErrorStatusCode) importProjectRes()       {}

// Ref: #/components/schemas/GetOrganization
type GetOrganization struct {
	Organization Organization `json:"organization"`
}

// GetOrganization returns the value of Organization.
func (s *GetOrganization) GetOrganization() Organization {
	return s.Organization
}

// SetOrganization sets the value of Organization.
func (s *GetOrganization) SetOrganization(val Organization) {
	s.Organization = val
}

func (*GetOrganization) getOrganizationRes() {}

// Ref: #/components/schemas/GetOrganizationRequest
type GetOrganizationRequest struct {
	Identifier Identifier `json:"identifier"`
}

// GetIdentifier returns the value of Identifier.
func (s *GetOrganizationRequest) GetIdentifier() Identifier {
	return s.Identifier
}

// SetIdentifier sets the value of Identifier.
func (s *GetOrganizationRequest) SetIdentifier(val Identifier) {
	s.Identifier = val
}

// Ref: #/components/schemas/GetPerson
type GetPerson struct {
	Person Person `json:"person"`
}

// GetPerson returns the value of Person.
func (s *GetPerson) GetPerson() Person {
	return s.Person
}

// SetPerson sets the value of Person.
func (s *GetPerson) SetPerson(val Person) {
	s.Person = val
}

func (*GetPerson) getPersonRes() {}

// Ref: #/components/schemas/GetPersonRequest
type GetPersonRequest struct {
	Identifier Identifier `json:"identifier"`
}

// GetIdentifier returns the value of Identifier.
func (s *GetPersonRequest) GetIdentifier() Identifier {
	return s.Identifier
}

// SetIdentifier sets the value of Identifier.
func (s *GetPersonRequest) SetIdentifier(val Identifier) {
	s.Identifier = val
}

// Ref: #/components/schemas/Identifier
type Identifier struct {
	Kind  string `json:"kind"`
	Value string `json:"value"`
}

// GetKind returns the value of Kind.
func (s *Identifier) GetKind() string {
	return s.Kind
}

// GetValue returns the value of Value.
func (s *Identifier) GetValue() string {
	return s.Value
}

// SetKind sets the value of Kind.
func (s *Identifier) SetKind(val string) {
	s.Kind = val
}

// SetValue sets the value of Value.
func (s *Identifier) SetValue(val string) {
	s.Value = val
}

// Ref: #/components/schemas/ImportOrganizationParams
type ImportOrganizationParams struct {
	Identifiers      []Identifier  `json:"identifiers"`
	ParentIdentifier OptIdentifier `json:"parentIdentifier"`
	Names            []Text        `json:"names"`
	Ceased           OptBool       `json:"ceased"`
	CreatedAt        OptDateTime   `json:"createdAt"`
	UpdatedAt        OptDateTime   `json:"updatedAt"`
}

// GetIdentifiers returns the value of Identifiers.
func (s *ImportOrganizationParams) GetIdentifiers() []Identifier {
	return s.Identifiers
}

// GetParentIdentifier returns the value of ParentIdentifier.
func (s *ImportOrganizationParams) GetParentIdentifier() OptIdentifier {
	return s.ParentIdentifier
}

// GetNames returns the value of Names.
func (s *ImportOrganizationParams) GetNames() []Text {
	return s.Names
}

// GetCeased returns the value of Ceased.
func (s *ImportOrganizationParams) GetCeased() OptBool {
	return s.Ceased
}

// GetCreatedAt returns the value of CreatedAt.
func (s *ImportOrganizationParams) GetCreatedAt() OptDateTime {
	return s.CreatedAt
}

// GetUpdatedAt returns the value of UpdatedAt.
func (s *ImportOrganizationParams) GetUpdatedAt() OptDateTime {
	return s.UpdatedAt
}

// SetIdentifiers sets the value of Identifiers.
func (s *ImportOrganizationParams) SetIdentifiers(val []Identifier) {
	s.Identifiers = val
}

// SetParentIdentifier sets the value of ParentIdentifier.
func (s *ImportOrganizationParams) SetParentIdentifier(val OptIdentifier) {
	s.ParentIdentifier = val
}

// SetNames sets the value of Names.
func (s *ImportOrganizationParams) SetNames(val []Text) {
	s.Names = val
}

// SetCeased sets the value of Ceased.
func (s *ImportOrganizationParams) SetCeased(val OptBool) {
	s.Ceased = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *ImportOrganizationParams) SetCreatedAt(val OptDateTime) {
	s.CreatedAt = val
}

// SetUpdatedAt sets the value of UpdatedAt.
func (s *ImportOrganizationParams) SetUpdatedAt(val OptDateTime) {
	s.UpdatedAt = val
}

// ImportOrganizationsOK is response for ImportOrganizations operation.
type ImportOrganizationsOK struct{}

func (*ImportOrganizationsOK) importOrganizationsRes() {}

// Ref: #/components/schemas/ImportOrganizationsRequest
type ImportOrganizationsRequest struct {
	Organizations []ImportOrganizationParams `json:"organizations"`
}

// GetOrganizations returns the value of Organizations.
func (s *ImportOrganizationsRequest) GetOrganizations() []ImportOrganizationParams {
	return s.Organizations
}

// SetOrganizations sets the value of Organizations.
func (s *ImportOrganizationsRequest) SetOrganizations(val []ImportOrganizationParams) {
	s.Organizations = val
}

// ImportPersonOK is response for ImportPerson operation.
type ImportPersonOK struct{}

func (*ImportPersonOK) importPersonRes() {}

// Ref: #/components/schemas/ImportPersonParams
type ImportPersonParams struct {
	Identifiers         []Identifier        `json:"identifiers"`
	Name                string              `json:"name"`
	PreferredName       OptString           `json:"preferredName"`
	GivenName           OptString           `json:"givenName"`
	PreferredGivenName  OptString           `json:"preferredGivenName"`
	FamilyName          OptString           `json:"familyName"`
	PreferredFamilyName OptString           `json:"preferredFamilyName"`
	HonorificPrefix     OptString           `json:"honorificPrefix"`
	Email               OptString           `json:"email"`
	Active              OptBool             `json:"active"`
	Role                OptString           `json:"role"`
	Username            OptString           `json:"username"`
	Attributes          []Attribute         `json:"attributes"`
	Tokens              []Token             `json:"tokens"`
	Affiliations        []AffiliationParams `json:"affiliations"`
	CreatedAt           OptDateTime         `json:"createdAt"`
	UpdatedAt           OptDateTime         `json:"updatedAt"`
}

// GetIdentifiers returns the value of Identifiers.
func (s *ImportPersonParams) GetIdentifiers() []Identifier {
	return s.Identifiers
}

// GetName returns the value of Name.
func (s *ImportPersonParams) GetName() string {
	return s.Name
}

// GetPreferredName returns the value of PreferredName.
func (s *ImportPersonParams) GetPreferredName() OptString {
	return s.PreferredName
}

// GetGivenName returns the value of GivenName.
func (s *ImportPersonParams) GetGivenName() OptString {
	return s.GivenName
}

// GetPreferredGivenName returns the value of PreferredGivenName.
func (s *ImportPersonParams) GetPreferredGivenName() OptString {
	return s.PreferredGivenName
}

// GetFamilyName returns the value of FamilyName.
func (s *ImportPersonParams) GetFamilyName() OptString {
	return s.FamilyName
}

// GetPreferredFamilyName returns the value of PreferredFamilyName.
func (s *ImportPersonParams) GetPreferredFamilyName() OptString {
	return s.PreferredFamilyName
}

// GetHonorificPrefix returns the value of HonorificPrefix.
func (s *ImportPersonParams) GetHonorificPrefix() OptString {
	return s.HonorificPrefix
}

// GetEmail returns the value of Email.
func (s *ImportPersonParams) GetEmail() OptString {
	return s.Email
}

// GetActive returns the value of Active.
func (s *ImportPersonParams) GetActive() OptBool {
	return s.Active
}

// GetRole returns the value of Role.
func (s *ImportPersonParams) GetRole() OptString {
	return s.Role
}

// GetUsername returns the value of Username.
func (s *ImportPersonParams) GetUsername() OptString {
	return s.Username
}

// GetAttributes returns the value of Attributes.
func (s *ImportPersonParams) GetAttributes() []Attribute {
	return s.Attributes
}

// GetTokens returns the value of Tokens.
func (s *ImportPersonParams) GetTokens() []Token {
	return s.Tokens
}

// GetAffiliations returns the value of Affiliations.
func (s *ImportPersonParams) GetAffiliations() []AffiliationParams {
	return s.Affiliations
}

// GetCreatedAt returns the value of CreatedAt.
func (s *ImportPersonParams) GetCreatedAt() OptDateTime {
	return s.CreatedAt
}

// GetUpdatedAt returns the value of UpdatedAt.
func (s *ImportPersonParams) GetUpdatedAt() OptDateTime {
	return s.UpdatedAt
}

// SetIdentifiers sets the value of Identifiers.
func (s *ImportPersonParams) SetIdentifiers(val []Identifier) {
	s.Identifiers = val
}

// SetName sets the value of Name.
func (s *ImportPersonParams) SetName(val string) {
	s.Name = val
}

// SetPreferredName sets the value of PreferredName.
func (s *ImportPersonParams) SetPreferredName(val OptString) {
	s.PreferredName = val
}

// SetGivenName sets the value of GivenName.
func (s *ImportPersonParams) SetGivenName(val OptString) {
	s.GivenName = val
}

// SetPreferredGivenName sets the value of PreferredGivenName.
func (s *ImportPersonParams) SetPreferredGivenName(val OptString) {
	s.PreferredGivenName = val
}

// SetFamilyName sets the value of FamilyName.
func (s *ImportPersonParams) SetFamilyName(val OptString) {
	s.FamilyName = val
}

// SetPreferredFamilyName sets the value of PreferredFamilyName.
func (s *ImportPersonParams) SetPreferredFamilyName(val OptString) {
	s.PreferredFamilyName = val
}

// SetHonorificPrefix sets the value of HonorificPrefix.
func (s *ImportPersonParams) SetHonorificPrefix(val OptString) {
	s.HonorificPrefix = val
}

// SetEmail sets the value of Email.
func (s *ImportPersonParams) SetEmail(val OptString) {
	s.Email = val
}

// SetActive sets the value of Active.
func (s *ImportPersonParams) SetActive(val OptBool) {
	s.Active = val
}

// SetRole sets the value of Role.
func (s *ImportPersonParams) SetRole(val OptString) {
	s.Role = val
}

// SetUsername sets the value of Username.
func (s *ImportPersonParams) SetUsername(val OptString) {
	s.Username = val
}

// SetAttributes sets the value of Attributes.
func (s *ImportPersonParams) SetAttributes(val []Attribute) {
	s.Attributes = val
}

// SetTokens sets the value of Tokens.
func (s *ImportPersonParams) SetTokens(val []Token) {
	s.Tokens = val
}

// SetAffiliations sets the value of Affiliations.
func (s *ImportPersonParams) SetAffiliations(val []AffiliationParams) {
	s.Affiliations = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *ImportPersonParams) SetCreatedAt(val OptDateTime) {
	s.CreatedAt = val
}

// SetUpdatedAt sets the value of UpdatedAt.
func (s *ImportPersonParams) SetUpdatedAt(val OptDateTime) {
	s.UpdatedAt = val
}

// Ref: #/components/schemas/ImportPersonRequest
type ImportPersonRequest struct {
	Person ImportPersonParams `json:"person"`
}

// GetPerson returns the value of Person.
func (s *ImportPersonRequest) GetPerson() ImportPersonParams {
	return s.Person
}

// SetPerson sets the value of Person.
func (s *ImportPersonRequest) SetPerson(val ImportPersonParams) {
	s.Person = val
}

// ImportProjectOK is response for ImportProject operation.
type ImportProjectOK struct{}

func (*ImportProjectOK) importProjectRes() {}

// Ref: #/components/schemas/ImportProjectParams
type ImportProjectParams struct {
	Identifiers  []Identifier `json:"identifiers"`
	Names        []Text       `json:"names"`
	Descriptions []Text       `json:"descriptions"`
	StartDate    OptString    `json:"startDate"`
	EndDate      OptString    `json:"endDate"`
	Deleted      OptBool      `json:"deleted"`
	Attributes   []Attribute  `json:"attributes"`
	CreatedAt    OptDateTime  `json:"createdAt"`
	UpdatedAt    OptDateTime  `json:"updatedAt"`
}

// GetIdentifiers returns the value of Identifiers.
func (s *ImportProjectParams) GetIdentifiers() []Identifier {
	return s.Identifiers
}

// GetNames returns the value of Names.
func (s *ImportProjectParams) GetNames() []Text {
	return s.Names
}

// GetDescriptions returns the value of Descriptions.
func (s *ImportProjectParams) GetDescriptions() []Text {
	return s.Descriptions
}

// GetStartDate returns the value of StartDate.
func (s *ImportProjectParams) GetStartDate() OptString {
	return s.StartDate
}

// GetEndDate returns the value of EndDate.
func (s *ImportProjectParams) GetEndDate() OptString {
	return s.EndDate
}

// GetDeleted returns the value of Deleted.
func (s *ImportProjectParams) GetDeleted() OptBool {
	return s.Deleted
}

// GetAttributes returns the value of Attributes.
func (s *ImportProjectParams) GetAttributes() []Attribute {
	return s.Attributes
}

// GetCreatedAt returns the value of CreatedAt.
func (s *ImportProjectParams) GetCreatedAt() OptDateTime {
	return s.CreatedAt
}

// GetUpdatedAt returns the value of UpdatedAt.
func (s *ImportProjectParams) GetUpdatedAt() OptDateTime {
	return s.UpdatedAt
}

// SetIdentifiers sets the value of Identifiers.
func (s *ImportProjectParams) SetIdentifiers(val []Identifier) {
	s.Identifiers = val
}

// SetNames sets the value of Names.
func (s *ImportProjectParams) SetNames(val []Text) {
	s.Names = val
}

// SetDescriptions sets the value of Descriptions.
func (s *ImportProjectParams) SetDescriptions(val []Text) {
	s.Descriptions = val
}

// SetStartDate sets the value of StartDate.
func (s *ImportProjectParams) SetStartDate(val OptString) {
	s.StartDate = val
}

// SetEndDate sets the value of EndDate.
func (s *ImportProjectParams) SetEndDate(val OptString) {
	s.EndDate = val
}

// SetDeleted sets the value of Deleted.
func (s *ImportProjectParams) SetDeleted(val OptBool) {
	s.Deleted = val
}

// SetAttributes sets the value of Attributes.
func (s *ImportProjectParams) SetAttributes(val []Attribute) {
	s.Attributes = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *ImportProjectParams) SetCreatedAt(val OptDateTime) {
	s.CreatedAt = val
}

// SetUpdatedAt sets the value of UpdatedAt.
func (s *ImportProjectParams) SetUpdatedAt(val OptDateTime) {
	s.UpdatedAt = val
}

// Ref: #/components/schemas/ImportProjectRequest
type ImportProjectRequest struct {
	Project ImportProjectParams `json:"project"`
}

// GetProject returns the value of Project.
func (s *ImportProjectRequest) GetProject() ImportProjectParams {
	return s.Project
}

// SetProject sets the value of Project.
func (s *ImportProjectRequest) SetProject(val ImportProjectParams) {
	s.Project = val
}

// NewOptBool returns new OptBool with value set to v.
func NewOptBool(v bool) OptBool {
	return OptBool{
		Value: v,
		Set:   true,
	}
}

// OptBool is optional bool.
type OptBool struct {
	Value bool
	Set   bool
}

// IsSet returns true if OptBool was set.
func (o OptBool) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptBool) Reset() {
	var v bool
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptBool) SetTo(v bool) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptBool) Get() (v bool, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptBool) Or(d bool) bool {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptDateTime returns new OptDateTime with value set to v.
func NewOptDateTime(v time.Time) OptDateTime {
	return OptDateTime{
		Value: v,
		Set:   true,
	}
}

// OptDateTime is optional time.Time.
type OptDateTime struct {
	Value time.Time
	Set   bool
}

// IsSet returns true if OptDateTime was set.
func (o OptDateTime) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptDateTime) Reset() {
	var v time.Time
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptDateTime) SetTo(v time.Time) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptDateTime) Get() (v time.Time, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptDateTime) Or(d time.Time) time.Time {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptIdentifier returns new OptIdentifier with value set to v.
func NewOptIdentifier(v Identifier) OptIdentifier {
	return OptIdentifier{
		Value: v,
		Set:   true,
	}
}

// OptIdentifier is optional Identifier.
type OptIdentifier struct {
	Value Identifier
	Set   bool
}

// IsSet returns true if OptIdentifier was set.
func (o OptIdentifier) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptIdentifier) Reset() {
	var v Identifier
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptIdentifier) SetTo(v Identifier) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptIdentifier) Get() (v Identifier, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptIdentifier) Or(d Identifier) Identifier {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptInt returns new OptInt with value set to v.
func NewOptInt(v int) OptInt {
	return OptInt{
		Value: v,
		Set:   true,
	}
}

// OptInt is optional int.
type OptInt struct {
	Value int
	Set   bool
}

// IsSet returns true if OptInt was set.
func (o OptInt) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptInt) Reset() {
	var v int
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptInt) SetTo(v int) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptInt) Get() (v int, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptInt) Or(d int) int {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// NewOptString returns new OptString with value set to v.
func NewOptString(v string) OptString {
	return OptString{
		Value: v,
		Set:   true,
	}
}

// OptString is optional string.
type OptString struct {
	Value string
	Set   bool
}

// IsSet returns true if OptString was set.
func (o OptString) IsSet() bool { return o.Set }

// Reset unsets value.
func (o *OptString) Reset() {
	var v string
	o.Value = v
	o.Set = false
}

// SetTo sets value to v.
func (o *OptString) SetTo(v string) {
	o.Set = true
	o.Value = v
}

// Get returns value and boolean that denotes whether value was set.
func (o OptString) Get() (v string, ok bool) {
	if !o.Set {
		return v, false
	}
	return o.Value, true
}

// Or returns value if set, or given parameter if does not.
func (o OptString) Or(d string) string {
	if v, ok := o.Get(); ok {
		return v
	}
	return d
}

// Ref: #/components/schemas/Organization
type Organization struct {
	Identifiers []Identifier         `json:"identifiers"`
	Names       []Text               `json:"names"`
	Ceased      bool                 `json:"ceased"`
	CreatedAt   time.Time            `json:"createdAt"`
	Position    OptInt               `json:"position"`
	UpdatedAt   time.Time            `json:"updatedAt"`
	Parents     []ParentOrganization `json:"parents"`
}

// GetIdentifiers returns the value of Identifiers.
func (s *Organization) GetIdentifiers() []Identifier {
	return s.Identifiers
}

// GetNames returns the value of Names.
func (s *Organization) GetNames() []Text {
	return s.Names
}

// GetCeased returns the value of Ceased.
func (s *Organization) GetCeased() bool {
	return s.Ceased
}

// GetCreatedAt returns the value of CreatedAt.
func (s *Organization) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetPosition returns the value of Position.
func (s *Organization) GetPosition() OptInt {
	return s.Position
}

// GetUpdatedAt returns the value of UpdatedAt.
func (s *Organization) GetUpdatedAt() time.Time {
	return s.UpdatedAt
}

// GetParents returns the value of Parents.
func (s *Organization) GetParents() []ParentOrganization {
	return s.Parents
}

// SetIdentifiers sets the value of Identifiers.
func (s *Organization) SetIdentifiers(val []Identifier) {
	s.Identifiers = val
}

// SetNames sets the value of Names.
func (s *Organization) SetNames(val []Text) {
	s.Names = val
}

// SetCeased sets the value of Ceased.
func (s *Organization) SetCeased(val bool) {
	s.Ceased = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *Organization) SetCreatedAt(val time.Time) {
	s.CreatedAt = val
}

// SetPosition sets the value of Position.
func (s *Organization) SetPosition(val OptInt) {
	s.Position = val
}

// SetUpdatedAt sets the value of UpdatedAt.
func (s *Organization) SetUpdatedAt(val time.Time) {
	s.UpdatedAt = val
}

// SetParents sets the value of Parents.
func (s *Organization) SetParents(val []ParentOrganization) {
	s.Parents = val
}

// Ref: #/components/schemas/ParentOrganization
type ParentOrganization struct {
	Identifiers []Identifier `json:"identifiers"`
	Names       []Text       `json:"names"`
	Ceased      bool         `json:"ceased"`
}

// GetIdentifiers returns the value of Identifiers.
func (s *ParentOrganization) GetIdentifiers() []Identifier {
	return s.Identifiers
}

// GetNames returns the value of Names.
func (s *ParentOrganization) GetNames() []Text {
	return s.Names
}

// GetCeased returns the value of Ceased.
func (s *ParentOrganization) GetCeased() bool {
	return s.Ceased
}

// SetIdentifiers sets the value of Identifiers.
func (s *ParentOrganization) SetIdentifiers(val []Identifier) {
	s.Identifiers = val
}

// SetNames sets the value of Names.
func (s *ParentOrganization) SetNames(val []Text) {
	s.Names = val
}

// SetCeased sets the value of Ceased.
func (s *ParentOrganization) SetCeased(val bool) {
	s.Ceased = val
}

// Ref: #/components/schemas/Person
type Person struct {
	Identifiers         []Identifier  `json:"identifiers"`
	Name                string        `json:"name"`
	PreferredName       OptString     `json:"preferredName"`
	GivenName           OptString     `json:"givenName"`
	PreferredGivenName  OptString     `json:"preferredGivenName"`
	FamilyName          OptString     `json:"familyName"`
	PreferredFamilyName OptString     `json:"preferredFamilyName"`
	HonorificPrefix     OptString     `json:"honorificPrefix"`
	Email               OptString     `json:"email"`
	Username            OptString     `json:"username"`
	Active              bool          `json:"active"`
	Role                OptString     `json:"role"`
	Attributes          []Attribute   `json:"attributes"`
	Tokens              []Token       `json:"tokens"`
	Affiliations        []Affiliation `json:"affiliations"`
	CreatedAt           time.Time     `json:"createdAt"`
	UpdatedAt           time.Time     `json:"updatedAt"`
}

// GetIdentifiers returns the value of Identifiers.
func (s *Person) GetIdentifiers() []Identifier {
	return s.Identifiers
}

// GetName returns the value of Name.
func (s *Person) GetName() string {
	return s.Name
}

// GetPreferredName returns the value of PreferredName.
func (s *Person) GetPreferredName() OptString {
	return s.PreferredName
}

// GetGivenName returns the value of GivenName.
func (s *Person) GetGivenName() OptString {
	return s.GivenName
}

// GetPreferredGivenName returns the value of PreferredGivenName.
func (s *Person) GetPreferredGivenName() OptString {
	return s.PreferredGivenName
}

// GetFamilyName returns the value of FamilyName.
func (s *Person) GetFamilyName() OptString {
	return s.FamilyName
}

// GetPreferredFamilyName returns the value of PreferredFamilyName.
func (s *Person) GetPreferredFamilyName() OptString {
	return s.PreferredFamilyName
}

// GetHonorificPrefix returns the value of HonorificPrefix.
func (s *Person) GetHonorificPrefix() OptString {
	return s.HonorificPrefix
}

// GetEmail returns the value of Email.
func (s *Person) GetEmail() OptString {
	return s.Email
}

// GetUsername returns the value of Username.
func (s *Person) GetUsername() OptString {
	return s.Username
}

// GetActive returns the value of Active.
func (s *Person) GetActive() bool {
	return s.Active
}

// GetRole returns the value of Role.
func (s *Person) GetRole() OptString {
	return s.Role
}

// GetAttributes returns the value of Attributes.
func (s *Person) GetAttributes() []Attribute {
	return s.Attributes
}

// GetTokens returns the value of Tokens.
func (s *Person) GetTokens() []Token {
	return s.Tokens
}

// GetAffiliations returns the value of Affiliations.
func (s *Person) GetAffiliations() []Affiliation {
	return s.Affiliations
}

// GetCreatedAt returns the value of CreatedAt.
func (s *Person) GetCreatedAt() time.Time {
	return s.CreatedAt
}

// GetUpdatedAt returns the value of UpdatedAt.
func (s *Person) GetUpdatedAt() time.Time {
	return s.UpdatedAt
}

// SetIdentifiers sets the value of Identifiers.
func (s *Person) SetIdentifiers(val []Identifier) {
	s.Identifiers = val
}

// SetName sets the value of Name.
func (s *Person) SetName(val string) {
	s.Name = val
}

// SetPreferredName sets the value of PreferredName.
func (s *Person) SetPreferredName(val OptString) {
	s.PreferredName = val
}

// SetGivenName sets the value of GivenName.
func (s *Person) SetGivenName(val OptString) {
	s.GivenName = val
}

// SetPreferredGivenName sets the value of PreferredGivenName.
func (s *Person) SetPreferredGivenName(val OptString) {
	s.PreferredGivenName = val
}

// SetFamilyName sets the value of FamilyName.
func (s *Person) SetFamilyName(val OptString) {
	s.FamilyName = val
}

// SetPreferredFamilyName sets the value of PreferredFamilyName.
func (s *Person) SetPreferredFamilyName(val OptString) {
	s.PreferredFamilyName = val
}

// SetHonorificPrefix sets the value of HonorificPrefix.
func (s *Person) SetHonorificPrefix(val OptString) {
	s.HonorificPrefix = val
}

// SetEmail sets the value of Email.
func (s *Person) SetEmail(val OptString) {
	s.Email = val
}

// SetUsername sets the value of Username.
func (s *Person) SetUsername(val OptString) {
	s.Username = val
}

// SetActive sets the value of Active.
func (s *Person) SetActive(val bool) {
	s.Active = val
}

// SetRole sets the value of Role.
func (s *Person) SetRole(val OptString) {
	s.Role = val
}

// SetAttributes sets the value of Attributes.
func (s *Person) SetAttributes(val []Attribute) {
	s.Attributes = val
}

// SetTokens sets the value of Tokens.
func (s *Person) SetTokens(val []Token) {
	s.Tokens = val
}

// SetAffiliations sets the value of Affiliations.
func (s *Person) SetAffiliations(val []Affiliation) {
	s.Affiliations = val
}

// SetCreatedAt sets the value of CreatedAt.
func (s *Person) SetCreatedAt(val time.Time) {
	s.CreatedAt = val
}

// SetUpdatedAt sets the value of UpdatedAt.
func (s *Person) SetUpdatedAt(val time.Time) {
	s.UpdatedAt = val
}

// Ref: #/components/schemas/SearchOrganizations
type SearchOrganizations struct {
	Hits []Organization `json:"hits"`
}

// GetHits returns the value of Hits.
func (s *SearchOrganizations) GetHits() []Organization {
	return s.Hits
}

// SetHits sets the value of Hits.
func (s *SearchOrganizations) SetHits(val []Organization) {
	s.Hits = val
}

// Ref: #/components/schemas/SearchOrganizationsRequest
type SearchOrganizationsRequest struct {
	Query OptString `json:"query"`
}

// GetQuery returns the value of Query.
func (s *SearchOrganizationsRequest) GetQuery() OptString {
	return s.Query
}

// SetQuery sets the value of Query.
func (s *SearchOrganizationsRequest) SetQuery(val OptString) {
	s.Query = val
}

// Ref: #/components/schemas/SearchPeople
type SearchPeople struct {
	Hits []Person `json:"hits"`
}

// GetHits returns the value of Hits.
func (s *SearchPeople) GetHits() []Person {
	return s.Hits
}

// SetHits sets the value of Hits.
func (s *SearchPeople) SetHits(val []Person) {
	s.Hits = val
}

// Ref: #/components/schemas/SearchPeopleRequest
type SearchPeopleRequest struct {
	Query OptString `json:"query"`
}

// GetQuery returns the value of Query.
func (s *SearchPeopleRequest) GetQuery() OptString {
	return s.Query
}

// SetQuery sets the value of Query.
func (s *SearchPeopleRequest) SetQuery(val OptString) {
	s.Query = val
}

// Ref: #/components/schemas/Text
type Text struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

// GetLang returns the value of Lang.
func (s *Text) GetLang() string {
	return s.Lang
}

// GetValue returns the value of Value.
func (s *Text) GetValue() string {
	return s.Value
}

// SetLang sets the value of Lang.
func (s *Text) SetLang(val string) {
	s.Lang = val
}

// SetValue sets the value of Value.
func (s *Text) SetValue(val string) {
	s.Value = val
}

// Ref: #/components/schemas/Token
type Token struct {
	Kind  string `json:"kind"`
	Value []byte `json:"value"`
}

// GetKind returns the value of Kind.
func (s *Token) GetKind() string {
	return s.Kind
}

// GetValue returns the value of Value.
func (s *Token) GetValue() []byte {
	return s.Value
}

// SetKind sets the value of Kind.
func (s *Token) SetKind(val string) {
	s.Kind = val
}

// SetValue sets the value of Value.
func (s *Token) SetValue(val []byte) {
	s.Value = val
}
