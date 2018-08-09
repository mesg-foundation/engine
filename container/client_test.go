package container

// TODO: this tests break other tests on my machine
// func TestCreateSwarm(t *testing.T) {
// 	leaveSwarm()
// 	dockerClient, _ := godocker.NewClientFromEnv()
// 	ID, err := createSwarm(dockerClient)
// 	assert.Nil(t, err)
// 	assert.NotEqual(t, ID, "")
// }

// func TestClientWithCreateSwarm(t *testing.T) {
// 	leaveSwarm()
// 	client, err := Client()
// 	assert.Nil(t, err)
// 	assert.NotNil(t, client)
// }

// func leaveSwarm() {
// 	dockerClient, _ := godocker.NewClientFromEnv()
// 	dockerClient.LeaveSwarm(godocker.LeaveSwarmOptions{
// 		Context: context.Background(),
// 		Force:   true,
// 	})
// }
