package filter

import (
	"regexp"
	"strings"
)

var badWords []string = []string{
	"qwerty",
	"йцукен",
	"zxvbnm",
}

type Result struct {
	Filtered bool `json:"filtered"`
}

/*
 * Тут у нас как в известном анекдоте: два путя:
 * 1. Собрать из списка плохих слов регулярку и мэтчить сроку по ней
 * 2. Разобрать строку на слова и каждое прогонять по массиву
 *
 * В первом случае в список должен содержать слова только из букв
 * и цифр. Иначе всё пропало
 *
 * Во втором нам всё равно накладывать регулярку на строку, чтобы
 * разбить последнюю на слова, но потом ещё M*N итераций цикла для
 * проверки.
 *
 * Так что пойдём первым путём и при пополнении списка плохих
 * слов будем держать в голове наложенные на него ограничения
 */

func BadWords(s string) bool {
	re_line := "(?:^|[^\\pL\\pN])("
	re_line += strings.Join(badWords, "|")
	re_line += ")(?:$|[^\\pL\\pN])"
	re := regexp.MustCompile(re_line)
	return re.MatchString(strings.ToLower(s))
}
