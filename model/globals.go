package model

type Globals struct {
	Version VersionFlag `kong:"name='version',help='Version of the application'"`
}
