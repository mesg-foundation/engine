package service

// // Resume run the resume command for a service
// var Resume = &cobra.Command{
// 	Use:   "resume SERVICE_ID",
// 	Short: "Resume a service",
// 	Long: `Resume a previously paused service.

// To pause a service, see the [pause command](mesg-core_service_pause.md)`,
// 	Args:              cobra.MinimumNArgs(1),
// 	Example:           "mesg-core service resume SERVICE_ID --account ACCOUNT --confirm",
// 	Run:               resumeHandler,
// 	DisableAutoGenTag: true,
// }

// func resumeHandler(cmd *cobra.Command, args []string) {
// 	account := utils.AccountFromFlagOrAsk(cmd, "Select an account:")
// 	if !utils.Confirm(cmd, "Are you sure?") {
// 		return
// 	}
// 	// TODO start and when ready resume (onchan) the service
// 	fmt.Println("Service resumed with success", args, account)
// }

// func init() {
// 	utils.Confirmable(Resume)
// 	utils.Accountable(Resume)
// }
