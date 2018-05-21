package dependency

type testDependency struct {
	Ports []string
}

func (t *testDependency) GetPorts() []string {
	return t.Ports
}
func (t *testDependency) GetVolumes() (volumes []string) {
	return
}
func (t *testDependency) GetVolumesfrom() (volumesFrom []string) {
	return
}
