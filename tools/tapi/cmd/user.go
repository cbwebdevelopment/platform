package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	osUser "os/user"
	"path/filepath"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	sheets "google.golang.org/api/sheets/v4"

	"github.com/urfave/cli"

	"github.com/tidepool-org/platform/tools/tapi/api"
	"github.com/tidepool-org/platform/user"
)

const (
	UserIDFlag      = "user-id"
	EmailFlag       = "email"
	PasswordFlag    = "password"
	FullNameFlag    = "fullName"
	RoleFlag        = "role"
	GoogleSheetFlag = "sheetID"
)

func UserCommands() cli.Commands {
	return cli.Commands{
		{
			Name:  "user",
			Usage: "user management",
			Subcommands: []cli.Command{
				{
					Name:  "get",
					Usage: "get a user by id or email",
					Flags: CommandFlags(
						cli.StringFlag{
							Name:  UserIDFlag,
							Usage: "`USERID` of the user to get",
						},
						cli.StringFlag{
							Name:  EmailFlag,
							Usage: "`EMAIL` of the user to get",
						},
					),
					Before: ensureNoArgs,
					Action: userGet,
				},
				{
					Name:  "find",
					Usage: "find all users matching the specified search criteria",
					Flags: CommandFlags(
						cli.StringFlag{
							Name:  RoleFlag,
							Usage: "find users matching the specified `ROLE`",
						},
					),
					Before: ensureNoArgs,
					Action: userFind,
				},
				{
					Name:  "add-role",
					Usage: "add the specified role to the user specified by id",
					Flags: CommandFlags(
						cli.StringFlag{
							Name:  UserIDFlag,
							Usage: "`USERID` of the user to update",
						},
						cli.StringSliceFlag{
							Name:  RoleFlag,
							Usage: "`ROLE` to add to the user",
						},
					),
					Before: ensureNoArgs,
					Action: userAddRoles,
				},
				{
					Name:  "remove-role",
					Usage: "remove the specified role from the user specified by id",
					Flags: CommandFlags(
						cli.StringFlag{
							Name:  UserIDFlag,
							Usage: "`USERID` of the user to update",
						},
						cli.StringSliceFlag{
							Name:  RoleFlag,
							Usage: "`ROLE` to remove from the user",
						},
					),
					Before: ensureNoArgs,
					Action: userRemoveRoles,
				},
				{
					Name:  "delete",
					Usage: "delete a user",
					Flags: CommandFlags(
						cli.StringFlag{
							Name:  UserIDFlag,
							Usage: "`USERID` of the user to delete",
						},
						cli.StringFlag{
							Name:  PasswordFlag,
							Usage: "`PASSWORD` of the user to delete (required if authenticated as the user being deleted)",
						},
					),
					Before: ensureNoArgs,
					Action: userDelete,
				},
				{
					Name:  "create",
					Usage: "create a user",
					Flags: CommandFlags(
						cli.StringFlag{
							Name:  EmailFlag,
							Usage: "`EMAIL` of the user to create",
						},
						cli.StringFlag{
							Name:  PasswordFlag,
							Usage: "`PASSWORD` of the user to create",
						},
						cli.StringFlag{
							Name:  FullNameFlag,
							Usage: "`FULLNAME` of the user to create",
						},
						cli.StringFlag{
							Name:  GoogleSheetFlag,
							Usage: "`GOOGLE SHEET ID` to bulk create users from",
						},
					),
					Before: ensureNoArgs,
					Action: userCreate,
				},
			},
		},
	}
}

func userGet(c *cli.Context) error {
	var user *user.User
	var err error

	email := c.String(EmailFlag)
	if email != "" {
		if c.String(UserIDFlag) != "" {
			return errors.New("Must specified either EMAIL or USERID, but not both")
		}
		user, err = API(c).GetUserByEmail(email)
	} else {
		user, err = API(c).GetUserByID(c.String(UserIDFlag))
	}

	if err != nil {
		return err
	}

	return reportMessageWithJSON(c, user)
}

func userFind(c *cli.Context) error {
	if !c.IsSet(RoleFlag) {
		return errors.New("No search criteria specified")
	}

	role := c.String(RoleFlag)
	if role == "" {
		return errors.New("Role is missing")
	}

	query := &api.UsersQuery{}
	query.Role = &role
	users, err := API(c).FindUsers(query)
	if err != nil {
		return err
	}

	for _, user := range users {
		if err = reportMessageWithJSON(c, user); err != nil {
			return err
		}
	}

	return nil
}

func userAddRoles(c *cli.Context) error {
	updater, err := api.NewAddRolesUserUpdater(c.StringSlice(RoleFlag))
	if err != nil {
		return err
	}

	updateUser, err := API(c).ApplyUpdatersToUserByID(c.String(UserIDFlag), []api.UserUpdater{updater})
	if err != nil {
		return err
	}

	return reportMessageWithJSON(c, updateUser)
}

func userRemoveRoles(c *cli.Context) error {
	updater, err := api.NewRemoveRolesUserUpdater(c.StringSlice(RoleFlag))
	if err != nil {
		return err
	}

	updateUser, err := API(c).ApplyUpdatersToUserByID(c.String(UserIDFlag), []api.UserUpdater{updater})
	if err != nil {
		return err
	}

	return reportMessageWithJSON(c, updateUser)
}

func userDelete(c *cli.Context) error {
	userID := c.String(UserIDFlag)
	password := c.String(PasswordFlag)
	if password == "" && API(c).IsSessionUserID(userID) {
		var err error
		if password, err = readFromConsoleNoEcho("Password: "); err != nil {
			return err
		}
	}

	err := API(c).DeleteUserByID(userID, password)
	if err != nil {
		return err
	}

	return reportMessage(c, "User deleted.")
}

func userCreate(c *cli.Context) error {
	email := c.String(EmailFlag)
	password := c.String(PasswordFlag)
	fullName := c.String(FullNameFlag)
	sheetID := c.String(GoogleSheetFlag)

	if sheetID != "" {
		// Fetch user definitions from google sheet
		return getUsersFromSheet(c, sheetID)
	}

	if err := API(c).CreateUser(email, password, fullName); err != nil {
		return err
	}

	return reportMessage(c, "User Created.")
}

// getClient uses a Context and Config to retrieve a Token
// then generate a Client. It returns the generated Client.
func getClient(ctx context.Context, config *oauth2.Config) *http.Client {
	cacheFile, err := tokenCacheFile()
	if err != nil {
		log.Fatalf("Unable to get path to cached credential file. %v", err)
	}
	tok, err := tokenFromFile(cacheFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(cacheFile, tok)
	}
	return config.Client(ctx, tok)
}

// getTokenFromWeb uses Config to request a Token.
// It returns the retrieved Token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// tokenCacheFile generates credential file path/filename.
// It returns the generated credential path/filename.
func tokenCacheFile() (string, error) {
	usr, err := osUser.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".credentials")
	os.MkdirAll(tokenCacheDir, 0700)
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("sheets.googleapis.com-go-quickstart.json")), err
}

// tokenFromFile retrieves a Token from a given file path.
// It returns the retrieved Token and any read error encountered.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	t := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(t)
	defer f.Close()
	return t, err
}

// saveToken uses a file path to create a file and store the
// token in it.
func saveToken(file string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", file)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}

func getUsersFromSheet(c *cli.Context, sheetID string) error {
	ctx := context.Background()

	b, err := ioutil.ReadFile("client_secret.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved credentials
	// at ~/.credentials/sheets.googleapis.com-go-quickstart.json
	config, err := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(ctx, config)

	srv, err := sheets.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Sheets Client %v", err)
	}

	// Prints the names and majors of students in a sample spreadsheet:
	// https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit
	// sheetID := "17ON4Sxmf-2njEFODCRZv2-lFl_RjwR_AI65dqrbq6jo"
	readRange := "Hackweek Test!A2:F"
	resp, err := srv.Spreadsheets.Values.Get(sheetID, readRange).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve data from sheet. %v", err)
	}

	// headings := resp.Values[0]

	if len(resp.Values) > 0 {
		// fmt.Println(headings[1], headings[2], headings[3])
		for _, row := range resp.Values {
			var email = row[2].(string)
			var password = row[3].(string)
			var fullName = row[0].(string)
			// var err error
			if err := API(c).CreateUser(email, password, fullName); err != nil {
				// return nil, err
			}
		}
		return reportMessage(c, "Users Created.")
	} else {
		fmt.Print("No data found.")
	}

	return nil
}
