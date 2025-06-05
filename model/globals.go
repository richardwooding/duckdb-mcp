package model

type Globals struct {
	Version model.VersionFlag `kong:"name='version',help='Version of the application',default='0.1.1'"`
}
