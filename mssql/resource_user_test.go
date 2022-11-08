package mssql

import (
  "fmt"
  "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
  "github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
  "os"
  "testing"
)

func TestAccUser_Local_Instance(t *testing.T) {
  resource.Test(t, resource.TestCase{
    PreCheck:          func() { testAccPreCheck(t) },
    IsUnitTest:        runLocalAccTests,
    ProviderFactories: testAccProviders,
    CheckDestroy:      func(state *terraform.State) error { return testAccCheckUserDestroy(state) },
    Steps: []resource.TestStep{
      {
        Config: testAccCheckUser(t, "instance", false, map[string]interface{}{"username": "instance", "login_name": "user_instance", "login_password": "valueIsH8kd$¡", "roles": "[\"db_owner\"]"}),
        Check: resource.ComposeTestCheckFunc(
          testAccCheckUserExists("mssql_user.instance"),
          testAccCheckDatabaseUserWorks("mssql_user.instance", "user_instance", "valueIsH8kd$¡"),
          resource.TestCheckResourceAttr("mssql_user.instance", "database", "master"),
          resource.TestCheckResourceAttr("mssql_user.instance", "username", "instance"),
          resource.TestCheckResourceAttr("mssql_user.instance", "login_name", "user_instance"),
          resource.TestCheckResourceAttr("mssql_user.instance", "authentication_type", "INSTANCE"),
          resource.TestCheckResourceAttr("mssql_user.instance", "default_schema", "dbo"),
          resource.TestCheckResourceAttr("mssql_user.instance", "default_language", ""),
          resource.TestCheckResourceAttr("mssql_user.instance", "roles.#", "1"),
          resource.TestCheckResourceAttr("mssql_user.instance", "roles.0", "db_owner"),
          resource.TestCheckResourceAttrSet("mssql_user.instance", "principal_id"),
          resource.TestCheckNoResourceAttr("mssql_user.instance", "password"),
        ),
      },
    },
  })
}

func TestAccUser_Azure_Instance(t *testing.T) {
  resource.Test(t, resource.TestCase{
    PreCheck:          func() { testAccPreCheck(t) },
    ProviderFactories: testAccProviders,
    CheckDestroy:      func(state *terraform.State) error { return testAccCheckUserDestroy(state) },
    Steps: []resource.TestStep{
      {
        Config: testAccCheckUser(t, "instance", true, map[string]interface{}{"database": "testdb", "username": "instance", "login_name": "user_instance", "login_password": "valueIsH8kd$¡", "roles": "[\"db_owner\"]"}),
        Check: resource.ComposeTestCheckFunc(
          testAccCheckUserExists("mssql_user.instance"),
          testAccCheckDatabaseUserWorks("mssql_user.instance", "user_instance", "valueIsH8kd$¡"),
          resource.TestCheckResourceAttr("mssql_user.instance", "database", "testdb"),
          resource.TestCheckResourceAttr("mssql_user.instance", "username", "instance"),
          resource.TestCheckResourceAttr("mssql_user.instance", "login_name", "user_instance"),
          resource.TestCheckResourceAttr("mssql_user.instance", "authentication_type", "INSTANCE"),
          resource.TestCheckResourceAttr("mssql_user.instance", "default_schema", "dbo"),
          resource.TestCheckResourceAttr("mssql_user.instance", "default_language", ""),
          resource.TestCheckResourceAttr("mssql_user.instance", "roles.#", "1"),
          resource.TestCheckResourceAttr("mssql_user.instance", "roles.0", "db_owner"),
          resource.TestCheckResourceAttrSet("mssql_user.instance", "principal_id"),
          resource.TestCheckNoResourceAttr("mssql_user.instance", "password"),
        ),
      },
    },
  })
}

func TestAccUser_Azure_Database(t *testing.T) {
  resource.Test(t, resource.TestCase{
    PreCheck:          func() { testAccPreCheck(t) },
    ProviderFactories: testAccProviders,
    CheckDestroy:      func(state *terraform.State) error { return testAccCheckUserDestroy(state) },
    Steps: []resource.TestStep{
      {
        Config: testAccCheckUser(t, "database", true, map[string]interface{}{"database": "testdb", "username": "database_user", "password": "valueIsH8kd$¡", "roles": "[\"db_owner\"]"}),
        Check: resource.ComposeTestCheckFunc(
          testAccCheckUserExists("mssql_user.database"),
          testAccCheckDatabaseUserWorks("mssql_user.database", "database_user", "valueIsH8kd$¡"),
          resource.TestCheckResourceAttr("mssql_user.database", "database", "testdb"),
          resource.TestCheckResourceAttr("mssql_user.database", "username", "database_user"),
          resource.TestCheckResourceAttr("mssql_user.database", "password", "valueIsH8kd$¡"),
          resource.TestCheckResourceAttr("mssql_user.database", "login_name", ""),
          resource.TestCheckResourceAttr("mssql_user.database", "authentication_type", "DATABASE"),
          resource.TestCheckResourceAttr("mssql_user.database", "default_schema", "dbo"),
          resource.TestCheckResourceAttr("mssql_user.database", "default_language", ""),
          resource.TestCheckResourceAttr("mssql_user.database", "roles.#", "1"),
          resource.TestCheckResourceAttr("mssql_user.database", "roles.0", "db_owner"),
          resource.TestCheckResourceAttrSet("mssql_user.database", "principal_id"),
        ),
      },
    },
  })
}

func TestAccUser_Azure_External(t *testing.T) {
  tenantId := os.Getenv("MSSQL_TENANT_ID")
  clientId := os.Getenv("TF_ACC_AZURE_USER_CLIENT_ID")
  clientUser := os.Getenv("TF_ACC_AZURE_USER_CLIENT_USER")
  clientSecret := os.Getenv("TF_ACC_AZURE_USER_CLIENT_SECRET")
  resource.Test(t, resource.TestCase{
    PreCheck:          func() { testAccPreCheck(t) },
    ProviderFactories: testAccProviders,
    CheckDestroy:      func(state *terraform.State) error { return testAccCheckUserDestroy(state) },
    Steps: []resource.TestStep{
      {
        Config: testAccCheckUser(t, "database", true, map[string]interface{}{"database": "testdb", "username": clientUser, "roles": "[\"db_owner\"]"}),
        Check: resource.ComposeTestCheckFunc(
          testAccCheckUserExists("mssql_user.database"),
          testAccCheckExternalUserWorks("mssql_user.database", tenantId, clientId, clientSecret),
          resource.TestCheckResourceAttr("mssql_user.database", "database", "testdb"),
          resource.TestCheckResourceAttr("mssql_user.database", "username", clientUser),
          resource.TestCheckResourceAttr("mssql_user.database", "login_name", ""),
          resource.TestCheckResourceAttr("mssql_user.database", "authentication_type", "EXTERNAL"),
          resource.TestCheckResourceAttr("mssql_user.database", "default_schema", "dbo"),
          resource.TestCheckResourceAttr("mssql_user.database", "default_language", ""),
          resource.TestCheckResourceAttr("mssql_user.database", "roles.#", "1"),
          resource.TestCheckResourceAttr("mssql_user.database", "roles.0", "db_owner"),
          resource.TestCheckResourceAttrSet("mssql_user.database", "principal_id"),
          resource.TestCheckNoResourceAttr("mssql_user.database", "password"),
        ),
      },
    },
  })
}

func TestAccUser_Local_Update_DefaultSchema(t *testing.T) {
  resource.Test(t, resource.TestCase{
    PreCheck:          func() { testAccPreCheck(t) },
    IsUnitTest:        runLocalAccTests,
    ProviderFactories: testAccProviders,
    CheckDestroy:      func(state *terraform.State) error { return testAccCheckLoginDestroy(state) },
    Steps: []resource.TestStep{
      {
        Config: testAccCheckUser(t, "update", false, map[string]interface{}{"username": "test_update", "login_name": "user_update", "login_password": "valueIsH8kd$¡"}),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("mssql_user.update", "default_schema", "dbo"),
          testAccCheckUserExists("mssql_user.update", Check{"default_schema", "==", "dbo"}),
          testAccCheckDatabaseUserWorks("mssql_user.update", "user_update", "valueIsH8kd$¡"),
        ),
      },
      {
        Config: testAccCheckUser(t, "update", false, map[string]interface{}{"username": "test_update", "login_name": "user_update", "login_password": "valueIsH8kd$¡", "default_schema": "sys"}),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("mssql_user.update", "default_schema", "sys"),
          testAccCheckUserExists("mssql_user.update", Check{"default_schema", "==", "sys"}),
          testAccCheckDatabaseUserWorks("mssql_user.update", "user_update", "valueIsH8kd$¡"),
        ),
      },
    },
  })
}

func TestAccUser_Local_Update_DefaultLanguage(t *testing.T) {
  resource.Test(t, resource.TestCase{
    PreCheck:          func() { testAccPreCheck(t) },
    IsUnitTest:        runLocalAccTests,
    ProviderFactories: testAccProviders,
    CheckDestroy:      func(state *terraform.State) error { return testAccCheckLoginDestroy(state) },
    Steps: []resource.TestStep{
      {
        Config: testAccCheckUser(t, "update", false, map[string]interface{}{"username": "test_update", "login_name": "user_update", "login_password": "valueIsH8kd$¡"}),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("mssql_user.update", "default_language", ""),
          testAccCheckUserExists("mssql_user.update", Check{"default_language", "==", ""}),
          testAccCheckDatabaseUserWorks("mssql_user.update", "user_update", "valueIsH8kd$¡"),
        ),
      },
      {
        Config: testAccCheckUser(t, "update", false, map[string]interface{}{"username": "test_update", "login_name": "user_update", "login_password": "valueIsH8kd$¡", "default_language": "russian"}),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("mssql_user.update", "default_language", ""),
          testAccCheckUserExists("mssql_user.update", Check{"default_language", "==", ""}),
          testAccCheckDatabaseUserWorks("mssql_user.update", "user_update", "valueIsH8kd$¡"),
        ),
      },
    },
  })
}

func TestAccUser_Local_Update_Roles(t *testing.T) {
  resource.Test(t, resource.TestCase{
    PreCheck:          func() { testAccPreCheck(t) },
    IsUnitTest:        runLocalAccTests,
    ProviderFactories: testAccProviders,
    CheckDestroy:      func(state *terraform.State) error { return testAccCheckLoginDestroy(state) },
    Steps: []resource.TestStep{
      {
        Config: testAccCheckUser(t, "update", false, map[string]interface{}{"username": "test_update", "login_name": "user_update", "login_password": "valueIsH8kd$¡"}),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("mssql_user.update", "roles.#", "0"),
          testAccCheckUserExists("mssql_user.update", Check{"roles", "==", []string{}}),
          testAccCheckDatabaseUserWorks("mssql_user.update", "user_update", "valueIsH8kd$¡"),
        ),
      },
      {
        Config: testAccCheckUser(t, "update", false, map[string]interface{}{"username": "test_update", "login_name": "user_update", "login_password": "valueIsH8kd$¡", "roles": "[\"db_owner\",\"db_datawriter\"]"}),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("mssql_user.update", "roles.#", "2"),
          resource.TestCheckResourceAttr("mssql_user.update", "roles.0", "db_owner"),
          resource.TestCheckResourceAttr("mssql_user.update", "roles.1", "db_datawriter"),
          testAccCheckUserExists("mssql_user.update", Check{"roles", "==", []string{"db_owner","db_datawriter"}}),
          testAccCheckDatabaseUserWorks("mssql_user.update", "user_update", "valueIsH8kd$¡"),
        ),
      },
      {
        Config: testAccCheckUser(t, "update", false, map[string]interface{}{"username": "test_update", "login_name": "user_update", "login_password": "valueIsH8kd$¡", "roles": "[\"db_owner\"]"}),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("mssql_user.update", "roles.#", "1"),
          resource.TestCheckResourceAttr("mssql_user.update", "roles.0", "db_owner"),
          testAccCheckUserExists("mssql_user.update", Check{"roles", "==", []string{"db_owner"}}),
          testAccCheckDatabaseUserWorks("mssql_user.update", "user_update", "valueIsH8kd$¡"),
        ),
      },
    },
  })
}

func TestAccUser_Azure_Update_DefaultSchema(t *testing.T) {
  resource.Test(t, resource.TestCase{
    PreCheck:          func() { testAccPreCheck(t) },
    ProviderFactories: testAccProviders,
    CheckDestroy:      func(state *terraform.State) error { return testAccCheckLoginDestroy(state) },
    Steps: []resource.TestStep{
      {
        Config: testAccCheckUser(t, "update", true, map[string]interface{}{"username": "test_update", "login_name": "user_update", "login_password": "valueIsH8kd$¡"}),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("mssql_user.update", "default_schema", "dbo"),
          testAccCheckUserExists("mssql_user.update", Check{"default_schema", "==", "dbo"}),
          testAccCheckDatabaseUserWorks("mssql_user.update", "user_update", "valueIsH8kd$¡"),
        ),
      },
      {
        Config: testAccCheckUser(t, "update", true, map[string]interface{}{"username": "test_update", "login_name": "user_update", "login_password": "valueIsH8kd$¡", "default_schema": "sys"}),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("mssql_user.update", "default_schema", "sys"),
          testAccCheckUserExists("mssql_user.update", Check{"default_schema", "==", "sys"}),
          testAccCheckDatabaseUserWorks("mssql_user.update", "user_update", "valueIsH8kd$¡"),
        ),
      },
    },
  })
}

func TestAccUser_Azure_Update_DefaultLanguage(t *testing.T) {
  resource.Test(t, resource.TestCase{
    PreCheck:          func() { testAccPreCheck(t) },
    ProviderFactories: testAccProviders,
    CheckDestroy:      func(state *terraform.State) error { return testAccCheckLoginDestroy(state) },
    Steps: []resource.TestStep{
      {
        Config: testAccCheckUser(t, "update", true, map[string]interface{}{"username": "test_update", "login_name": "user_update", "login_password": "valueIsH8kd$¡"}),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("mssql_user.update", "default_language", ""),
          testAccCheckUserExists("mssql_user.update", Check{"default_language", "==", ""}),
          testAccCheckDatabaseUserWorks("mssql_user.update", "user_update", "valueIsH8kd$¡"),
        ),
      },
      {
        Config: testAccCheckUser(t, "update", true, map[string]interface{}{"username": "test_update", "login_name": "user_update", "login_password": "valueIsH8kd$¡", "default_language": "russian"}),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("mssql_user.update", "default_language", ""),
          testAccCheckUserExists("mssql_user.update", Check{"default_language", "==", ""}),
          testAccCheckDatabaseUserWorks("mssql_user.update", "user_update", "valueIsH8kd$¡"),
        ),
      },
    },
  })
}

func TestAccUser_Azure_Update_Roles(t *testing.T) {
  resource.Test(t, resource.TestCase{
    PreCheck:          func() { testAccPreCheck(t) },
    ProviderFactories: testAccProviders,
    CheckDestroy:      func(state *terraform.State) error { return testAccCheckLoginDestroy(state) },
    Steps: []resource.TestStep{
      {
        Config: testAccCheckUser(t, "update", true, map[string]interface{}{"database": "testdb", "username": "test_update", "login_name": "user_update", "login_password": "valueIsH8kd$¡"}),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("mssql_user.update", "roles.#", "0"),
          testAccCheckUserExists("mssql_user.update", Check{"roles", "==", []string{}}),
          testAccCheckDatabaseUserWorks("mssql_user.update", "user_update", "valueIsH8kd$¡"),
        ),
      },
      {
        Config: testAccCheckUser(t, "update", true, map[string]interface{}{"database": "testdb", "username": "test_update", "login_name": "user_update", "login_password": "valueIsH8kd$¡", "roles": "[\"db_owner\",\"db_datawriter\"]"}),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("mssql_user.update", "roles.#", "2"),
          resource.TestCheckResourceAttr("mssql_user.update", "roles.0", "db_owner"),
          resource.TestCheckResourceAttr("mssql_user.update", "roles.1", "db_datawriter"),
          testAccCheckUserExists("mssql_user.update", Check{"roles", "==", []string{"db_owner","db_datawriter"}}),
          testAccCheckDatabaseUserWorks("mssql_user.update", "user_update", "valueIsH8kd$¡"),
        ),
      },
      {
        Config: testAccCheckUser(t, "update", true, map[string]interface{}{"database": "testdb", "username": "test_update", "login_name": "user_update", "login_password": "valueIsH8kd$¡", "roles": "[\"db_owner\"]"}),
        Check: resource.ComposeTestCheckFunc(
          resource.TestCheckResourceAttr("mssql_user.update", "roles.#", "1"),
          resource.TestCheckResourceAttr("mssql_user.update", "roles.0", "db_owner"),
          testAccCheckUserExists("mssql_user.update", Check{"roles", "==", []string{"db_owner"}}),
          testAccCheckDatabaseUserWorks("mssql_user.update", "user_update", "valueIsH8kd$¡"),
        ),
      },
    },
  })
}

func testAccCheckUser(t *testing.T, name string, azure bool, data map[string]interface{}) string {
  text := `{{ if .login_name }}
           resource "mssql_login" "{{ .name }}" {
             login_name = "{{ .login_name }}"
             password   = "{{ .login_password }}"
           }
           {{ end }}
           resource "mssql_user" "{{ .name }}" {
             {{ with .database }}database = "{{ . }}"{{ end }}
             username = "{{ .username }}"
             {{ with .password }}password = "{{ . }}"{{ end }}
             {{ with .login_name }}login_name = "{{ . }}"{{ end }}
             {{ with .default_schema }}default_schema = "{{ . }}"{{ end }}
             {{ with .default_language }}default_language = "{{ . }}"{{ end }}
             {{ with .roles }}roles = {{ . }}{{ end }}
           }`
  data["name"] = name
  res, err := templateToString(name, text, data)
  if err != nil {
    t.Fatalf("%s", err)
  }
  return res
}

func testAccCheckUserDestroy(state *terraform.State) error {
  for _, rs := range state.RootModule().Resources {
    if rs.Type != "mssql_user" {
      continue
    }

    connector, err := getTestConnector(rs.Primary.Attributes)
    if err != nil {
      return err
    }

    database := rs.Primary.Attributes["database"]
    username := rs.Primary.Attributes["username"]
    login, err := connector.GetUser(database, username)
    if login != nil {
      return fmt.Errorf("user still exists")
    }
    if err != nil {
      return fmt.Errorf("expected no error, got %s", err)
    }
  }
  return nil
}

func testAccCheckUserExists(resource string, checks ...Check) resource.TestCheckFunc {
  return func(state *terraform.State) error {
    rs, ok := state.RootModule().Resources[resource]
    if !ok {
      return fmt.Errorf("not found: %s", resource)
    }
    if rs.Type != "mssql_user" {
      return fmt.Errorf("expected resource of type %s, got %s", "mssql_user", rs.Type)
    }
    if rs.Primary.ID == "" {
      return fmt.Errorf("no record ID is set")
    }
    connector, err := getTestConnector(rs.Primary.Attributes)
    if err != nil {
      return err
    }

    database := rs.Primary.Attributes["database"]
    username := rs.Primary.Attributes["username"]
    user, err := connector.GetUser(database, username)
    if user == nil {
      return fmt.Errorf("user does not exist")
    }
    if err != nil {
      return fmt.Errorf("expected no error, got %s", err)
    }

    var actual interface{}
    for _, check := range checks {
      switch check.name {
      case "password":
        actual = user.Password
      case "login_name":
        actual = user.LoginName
      case "default_schema":
        actual = user.DefaultSchema
      case "default_language":
        actual = user.DefaultLanguage
      case "roles":
        actual = user.Roles
      case "authentication_type":
        actual = user.AuthType
      default:
        return fmt.Errorf("unknown property %s", check.name)
      }
      if (check.op == "" || check.op == "==") && !equal(check.expected, actual) {
        return fmt.Errorf("expected %s == %s, got %s", check.name, check.expected, actual)
      }
      if check.op == "!=" && equal(check.expected, actual) {
        return fmt.Errorf("expected %s != %s, got %s", check.name, check.expected, actual)
      }
    }
    return nil
  }
}

func equal(a, b interface{}) bool {
  switch a.(type) {
  case []string:
    aa := a.([]string)
    bb := b.([]string)
    if len(aa) != len(bb) {
      return false
    }
    for i, v := range aa {
      if v != bb[i] {
        return false
      }
    }
    return true
  default:
    return a == b
  }
}

func testAccCheckDatabaseUserWorks(resource string, username, password string) resource.TestCheckFunc {
  return func(state *terraform.State) error {
    rs, ok := state.RootModule().Resources[resource]
    if !ok {
      return fmt.Errorf("not found: %s", resource)
    }
    if rs.Type != "mssql_user" {
      return fmt.Errorf("expected resource of type %s, got %s", "mssql_user", rs.Type)
    }
    if rs.Primary.ID == "" {
      return fmt.Errorf("no record ID is set")
    }
    connector, err := getTestUserConnector(rs.Primary.Attributes, username, password)
    if err != nil {
      return err
    }
    current, system, err := connector.GetCurrentUser(rs.Primary.Attributes[databaseProp])
    if err != nil {
      return fmt.Errorf("error: %s", err)
    }
    if current != rs.Primary.Attributes[usernameProp] {
      return fmt.Errorf("expected to be user %s, got %s (%s)", rs.Primary.Attributes[usernameProp], current, system)
    }
    return nil
  }
}

func testAccCheckExternalUserWorks(resource string, tenantId, clientId, clientSecret string) resource.TestCheckFunc {
  return func(state *terraform.State) error {
    rs, ok := state.RootModule().Resources[resource]
    if !ok {
      return fmt.Errorf("not found: %s", resource)
    }
    if rs.Type != "mssql_user" {
      return fmt.Errorf("expected resource of type %s, got %s", "mssql_user", rs.Type)
    }
    if rs.Primary.ID == "" {
      return fmt.Errorf("no record ID is set")
    }
    connector, err := getTestExternalConnector(rs.Primary.Attributes, tenantId, clientId, clientSecret)
    if err != nil {
      return err
    }
    current, system, err := connector.GetCurrentUser(rs.Primary.Attributes[databaseProp])
    if err != nil {
      return fmt.Errorf("error: %s", err)
    }
    if current != rs.Primary.Attributes[usernameProp] {
      return fmt.Errorf("expected to be user %s, got %s (%s)", rs.Primary.Attributes[usernameProp], current, system)
    }
    return nil
  }
}
