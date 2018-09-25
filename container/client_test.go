package container

// TODO: this tests break other tests on my machine
// func TestCreateSwarm(t *testing.T) {
// 	leaveSwarm()
// 	dockerClient, _ := godocker.NewClientFromEnv()
// 	ID, err := createSwarm(dockerClient)
// 	require.NoError(t, err)
// 	require.NotEqual(t, ID, "")
// }

// func TestClientWithCreateSwarm(t *testing.T) {
// 	leaveSwarm()
// 	client, err := Client()
// 	require.NoError(t, err)
// 	require.NotNil(t, client)
// }

// func leaveSwarm() {
// 	dockerClient, _ := godocker.NewClientFromEnv()
// 	dockerClient.LeaveSwarm(godocker.LeaveSwarmOptions{
// 		Context: context.Background(),
// 		Force:   true,
// 	})
// }
