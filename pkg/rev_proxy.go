package pkg

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"poc-gateway/pkg/interfaces"
	"poc-gateway/pkg/plugins"
)

type GatewayService struct {
	RegisteredPlugins []interfaces.GenericGatewayPlugin
	RevProxy          *httputil.ReverseProxy
	TargetURL         *url.URL
}

// func dealWithErrors(r *http.Response)

func (g GatewayService) ProxyDispatcher(proxy *httputil.ReverseProxy) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		for _, plugin := range g.RegisteredPlugins {
			if err := plugin.Process(r); err != nil {
				w.WriteHeader(err.StatusCode)
				w.Write(err.JSON())
				return
			}
		}
		proxy.Transport = &TransportLayer{}
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = g.TargetURL.Scheme
			req.URL.Host = g.TargetURL.Host
		}
		proxy.ServeHTTP(w, r)
	}
}

func (g GatewayService) Init() error {
	for _, plugin := range g.RegisteredPlugins {
		if err := plugin.Setup(); err != nil {
			return err
		}
	}
	return nil
}

func ReverseProxy() {

	svc := GatewayService{
		RegisteredPlugins: []interfaces.GenericGatewayPlugin{
			&plugins.OIDCPlugin{
				IDPUrl:      "https://id.magalu.com/oauth/certs",
				AllowedAuds: []string{"public"},
			},
			&plugins.HeaderCleanerPlugin{
				AllowedHeaders: []string{"Access-Control-Allow-Headers", "Access-Control-Allow-Methods",
					"Access-Control-Max-Age", "Access-Control-Request-Method",
					"Access-Control-Allow-Credentials", "Access-Control-Allow-Origin",
					"Origin", "content-type", "date", "x-api-key",
					"x-tenant-id", "authorization"},
			},
			plugins.TenantsPlugin{},
		},
	}

	url, err := url.Parse("https://api.github.com/")
	if err != nil {
		log.Fatal(err)
	}
	svc.TargetURL = url

	proxy, err := httputil.NewSingleHostReverseProxy(url), nil
	if err != nil {
		panic(err)
	}
	svc.RevProxy = proxy
	err = svc.Init()
	if err != nil {
		log.Fatal(err)

	}
	http.HandleFunc("/", svc.ProxyDispatcher(svc.RevProxy))
	log.Fatal(http.ListenAndServe(":8080", nil))

}
