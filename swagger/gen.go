package swagger

//go:generate mkdir -p ../server
//go:generate swagger generate server --target ../server --name experimental --spec swagger.yml
