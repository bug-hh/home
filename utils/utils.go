package utils

var (
	G_mysql_addr string
	G_mysql_port string
	G_fastdfs_addr string
	G_fastdfs_port string
)


func AddDomain2Url(url string) (domain_url string) {
	domain_url = "http://" + G_fastdfs_addr + ":" + G_fastdfs_port + "/" + url
	return domain_url
}