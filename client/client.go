package client

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// Client represents an instance of a connection to AWS.
type Client struct {
	awsAccessKey    string   `json:"awsAccessKey"` // IAM access key to AWS.
	awsAccessSecret string   `json:"-"`            // IAM access secret to AWS.
	awsSession      *ssm.SSM `json:"-"`            // AWS connection.
}

// New is a factory function that returns a new client instance.
func New(key string, secret string, region string) *Client {
	s := session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(key, secret, ""),
	}))

	return &Client{
		awsAccessKey:    key,
		awsAccessSecret: secret,
		awsSession:      ssm.New(s),
	}
}

// Version prints the version of the client then exits.
func (c *Client) Version() string {
	return fmt.Sprintf("%s ver. %s", applicationName, version)
}

// A parameter value for exporting.
type Parameter struct {
	Name    string `json:"name"`    // Name of the key.
	Type    string `json:"type"`    // Data type.
	Value   string `json:"value"`   // Value.
	Version int64  `json:"version"` // Version of the value.
}

// A list of parameters.
type Parameters []*Parameter

// Get returns a single value from AWS Parameter Store for a given key and version.
func (c *Client) Query(key string, version string) (*Parameter, error) {
	var ver64 int64
	var err error
	if version != "" {
		ver64, err = strconv.ParseInt(version, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	i := &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	}
	o, err := c.awsSession.GetParameter(i)
	if err != nil {
		return nil, err
	}

	// Prep record.
	name := *o.Parameter.Name
	dataType := *o.Parameter.Type
	value := *o.Parameter.Value
	v := *o.Parameter.Version

	// Filter by version?
	if version != "" && ver64 != v {
		found := false
		historyPage, err := c.getParameterHistory(*o.Parameter.Name)
		if err != nil {
			return nil, err
		}
		for _, hp := range historyPage {
			for _, h := range hp.Parameters {
				if *h.Version == ver64 {
					value = *h.Value
					v = *h.Version
					found = true
				}
			}
		}
		if !found {
			e := errors.New(fmt.Sprintf("Version %s not found for key: %s", version, name))
			return nil, e
		}
	}

	return &Parameter{
		Name:    name,
		Type:    dataType,
		Value:   value,
		Version: v,
	}, nil
}

// Get returns a value from AWS Parameter Store for a given directory and version.
func (c *Client) QueryPath(path string, recursive bool, version string) (Parameters, error) {
	var ver64 int64
	var err error
	if version != "" {
		ver64, err = strconv.ParseInt(version, 10, 64)
		if err != nil {
			return nil, err
		}
	}
	result, err := c.getParametersByPath(path, recursive)
	if err != nil {
		return nil, err
	}

	var parms Parameters
	for _, r := range result {
		for _, p := range r.Parameters {
			// Prep record.
			name := *p.Name
			dataType := *p.Type
			value := *p.Value
			v := *p.Version

			// Filter by version?
			if version != "" && ver64 != v {
				found := false
				historyPage, err := c.getParameterHistory(*p.Name)
				if err != nil {
					return nil, err
				}
				for _, hp := range historyPage {
					for _, h := range hp.Parameters {
						if *h.Version == ver64 {
							value = *h.Value
							v = *h.Version
							found = true
						}
					}
				}
				if !found {
					e := errors.New(fmt.Sprintf("Version %s not found for key: %s", version, name))
					return nil, e
				}
			}

			// Append to results
			parms = append(parms, &Parameter{
				Name:    name,
				Type:    dataType,
				Value:   value,
				Version: v,
			})
		}
	}
	return parms, nil
}

// Given a directory path, return the k/v and version. Optionally perform
// recursively on directory path. Loops through the result sets of multiple
// calls to AWS.
func (c *Client) getParametersByPath(prefix string,
	recursive bool) ([]*ssm.GetParametersByPathOutput, error) {
	var err error
	var parameters []*ssm.GetParametersByPathOutput
	i := &ssm.GetParametersByPathInput{
		Path:           aws.String(prefix),
		Recursive:      aws.Bool(recursive),
		WithDecryption: aws.Bool(true),
	}
	err = c.awsSession.GetParametersByPathPages(i,
		func(page *ssm.GetParametersByPathOutput, lastPage bool) bool {
			parameters = append(parameters, page)
			return !lastPage
		})
	return parameters, err
}

// Given a key retrieve the history of a parameter, or error if not found
func (c *Client) getParameterHistory(key string) ([]*ssm.GetParameterHistoryOutput, error) {
	var err error
	var history []*ssm.GetParameterHistoryOutput
	i := &ssm.GetParameterHistoryInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	}
	err = c.awsSession.GetParameterHistoryPages(i,
		func(page *ssm.GetParameterHistoryOutput, lastPage bool) bool {
			history = append(history, page)
			return !lastPage
		})
	return history, err
}
