package robot_service

type Response struct {
	topic  string
	MsgStr string
}

func  (s *robotService) ErrResponse(err error, sn, sessionid, title string) error {
	//robot:=s.getRobot(sn)
	return nil
}
