package response

type m map[string]interface{}

func InternalError() m {
	return m{"status": 500, "message": "Internal Server Error"}
}

func BadRequest() m {
	return m{"status": 400, "message": "Bad Request"}
}

func Failed() m {
	return m{"status": 400, "message": "Failed"}
}

func NotFound() m {
	return m{"status": 404, "message": "Not Found"}
}

func Success() m {
	return m{"status": 200, "message": "Success"}
}

func Custom(code int, message string) m {
	return m{"status": code, "message": message}
}

func DataArray(code int, message string, data []map[string]interface{}) m {
	return m{"status": code, "message": message, "data": data}
}
