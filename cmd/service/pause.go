package cmdService

// // Pause run the pause command for a service
// var Pause = &cobra.Command{
// 	Use:   "pause SERVICE_ID",
// 	Short: "Pause a service",
// 	Long: `Pause a service. The service should have been previously started.

// You should always pause services before quitting the CLI, otherwise you may loss your stake.

// You will **NOT** get your stake back with this command. The goal of this command is to give you an opportunity to stop running a service for a short period of time without losing your stake.

// When a service is paused, the stake duration count is also paused.`,
// 	Args: cobra.MinimumNArgs(1),
// 	Example: `mesg-core service pause SERVICE_ID
// mesg-core service pause SERVICE_ID --account ACCOUNT --confirm`,
// 	Run:               pauseHandler,
// 	DisableAutoGenTag: true,
// }

// func pauseHandler(cmd *cobra.Command, args []string) {
// 	account := cmdUtils.AccountFromFlagOrAsk(cmd, "Select the account to use:")
// 	if !cmdUtils.Confirm(cmd, "Are you sure?") {
// 		return
// 	}
// 	// TODO pause (onchain) and then stop the service
// 	fmt.Println("Service paused", args, account)
// }

// func init() {
// 	cmdUtils.Confirmable(Pause)
// 	cmdUtils.Accountable(Pause)
// }
