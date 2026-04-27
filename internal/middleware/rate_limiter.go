package middleware

import (
	"net/http"
	"sync"
	"time"
)

var (
	visitors = make(map[string]*visitor)
	mu       sync.Mutex
)

type visitor struct {
	lastSeen time.Time
	count    int
}

func init() {
	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()
			for ip, v := range visitors {
				if time.Since(v.lastSeen) > 3*time.Minute {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()
}

// Limite à 50 requêtes par seconde
func RateLimiter(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		mu.Lock()
		v, exists := visitors[ip]
		if !exists {
			visitors[ip] = &visitor{lastSeen: time.Now(), count: 1}
		} else {
			if time.Since(v.lastSeen) > time.Second {
				v.count = 0
			}
			v.count++
			v.lastSeen = time.Now()

			if v.count > 50 {
				mu.Unlock()
				http.Error(w, "Rate Limit Exceeded: Trop de requêtes (Anti-DDoS activé)", http.StatusTooManyRequests)
				return
			}
		}
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
