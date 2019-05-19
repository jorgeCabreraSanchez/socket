package Api

import "net/http"

func ReturnHttpError(err error, w http.ResponseWriter, httpErrType int) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "limit, includes, skip, sort, filter, select, destination, self, Accept, Accept-Language, Content-Type, Content-Length, Accept-Encoding, Authorization,X-CSRF-Token	")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(httpErrType)
	w.Write([]byte("{\"success\":false,\"error\": \"" + err.Error() + "\"}"))
}
