package cookies

import "net/http"

func SetCokies(w http.ResponseWriter, name, value string, sameSite http.SameSite, maxAge int) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		HttpOnly: true,
		SameSite: sameSite,
		MaxAge:   maxAge,
	})
}
