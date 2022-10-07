package mongo

type spanKindTagName string

const (
	HttpUrl                = spanKindTagName("http.url")
	NodeId                 = spanKindTagName("node_id")
	RequestSize            = spanKindTagName("request_size")
	PeerAddress            = spanKindTagName("peer.address")
	HttpStatusCode         = spanKindTagName("http.status_code")
	HttpMethod             = spanKindTagName("http.method")
	IstioCanonicalService  = spanKindTagName("istio.canonical_service")
	UpstreamCluster        = spanKindTagName("upstream_cluster")
	HttpProtocol           = spanKindTagName("http.protocol")
	IstioCanonicalRevision = spanKindTagName("istio.canonical_revision")
	DownstreamCluster      = spanKindTagName("downstream_cluster")
	UpstreamClusterName    = spanKindTagName("upstream_cluster.name")
	IstioNamespace         = spanKindTagName("istio.namespace")
	ResponseFlags          = spanKindTagName("response_flags")
	IstioMeshId            = spanKindTagName("istio.mesh_id")
	UserAgent              = spanKindTagName("user_agent")
	ResponseSize           = spanKindTagName("response_size")
	Component              = spanKindTagName("component")
	SpanKind               = spanKindTagName("span.kind")
	InternalSpanFormat     = spanKindTagName("internal.span.format")
)
