package server

type HttpServer struct {
	ltpService LtpService
}

func NewHttpServer(ltpService LtpService) HttpServer {
	return HttpServer{
		ltpService: ltpService,
	}
}
