# Groups Reader

Add all the groups reader in this package.


### GroupsReader Interface

Use the following interface to create your own custom groups reader.
```golang
type GroupReader interface {
	Members(grpName string) ([]string, error)
}
```

### Note

Use the following trick to enforce implementation of the GroupsReader

```golang
var _ GroupsReader = &EnterpriseGithubGroupReader{}
```
