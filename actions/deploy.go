package actions

func (act *Action) DeployService() error {
	act.ExecComposeCommandAndWait("down")
	act.ExecComposeCommandAndWait("up", "-d", "--force-recreate", "--no-deps")
	return nil
}
