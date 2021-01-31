package vmess

// Name 协议名称
const Name = "vmess"

// Request Options
const (
	OptBasicFormat byte = 0 // 不加密传输
	OptChunkStream byte = 1 // 分块传输，每个分块使用如下Security方法加密
	// OptReuseTCPConnection byte = 2
	// OptMetadataObfuscate  byte = 4
)

// Security types
const (
	SecurityAES128GCM        byte = 3
	SecurityChacha20Poly1305 byte = 4
	SecurityNone             byte = 5
)

// CMD types
const (
	CmdTCP byte = 1
	CmdUDP byte = 2
)

// Atyp
const (
	AddrTypeIPv4   byte = 1
	AddrTypeDomain byte = 2
	AddrTypeIPv6   byte = 3
)
