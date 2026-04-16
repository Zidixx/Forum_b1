package utils

import (
	"fmt"
	"math"
	"time"
)

func TimeAgo(t time.Time) string {
	diff := time.Since(t)
	seconds := diff.Seconds()

	switch {
	case seconds < 60:
		return "à l'instant"
	case seconds < 3600:
		m := int(math.Floor(seconds / 60))
		if m == 1 {
			return "il y a 1 minute"
		}
		return fmt.Sprintf("il y a %d minutes", m)
	case seconds < 86400:
		h := int(math.Floor(seconds / 3600))
		if h == 1 {
			return "il y a 1 heure"
		}
		return fmt.Sprintf("il y a %d heures", h)
	default:
		d := int(math.Floor(seconds / 86400))
		if d == 1 {
			return "il y a 1 jour"
		}
		return fmt.Sprintf("il y a %d jours", d)
	}
}

func FormatDate(t time.Time) string {
	return t.Format("02/01/2006 15:04")
}
