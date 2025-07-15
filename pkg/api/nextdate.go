package api

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

const layout = "20060102"

func NextDate(now time.Time, dstart string, repeat string) (string, error) {
	//Преобразуем строку dstart в объект времени
	date, err := time.Parse(layout, dstart)
	if err != nil {
		return "", errors.New("Неправильный формат даты")
	}

	//Разбиваем параметр repeat на части
	sliceRepeat, err := split(repeat)
	if err != nil {
		return "", err
	}

	//Выбираем действие в зависимости от первого условия repeat
	switch sliceRepeat[0] {
	case "d":
		//Проверяем правильность интервала
		if len(sliceRepeat) != 2 {
			return "", errors.New("не указан интервал в днях")
		}

		interval, err := strconv.Atoi(sliceRepeat[1])
		if err != nil {
			return "", errors.New("не верный интервал")
		}

		if interval < 1 || interval > 400 {
			return "", errors.New("превышен максимально допустимый интервал")
		}

		//Переносим задачу на указанное количество дней
		for {
			date = date.AddDate(0, 0, interval)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(layout), nil
	case "y":
		//Переносим задачу на следующий год
		for {
			date = date.AddDate(1, 0, 0)
			if afterNow(date, now) {
				break
			}
		}
		return date.Format(layout), nil
	case "w":
		//Проверяем правильность формата при условии недели
		if len(sliceRepeat) != 2 {
			return "", errors.New("не указан день недели")
		}

		//В днях недели меняем 7 на 0, т.к. в date.Weekday() воскресенье будет 0, а в условиях задачи 7
		weekdays := strings.Replace(sliceRepeat[1], "7", "0", -1)
		//Проверяем правильность дней недели в параметре
		for _, chrday := range strings.Split(weekdays, ",") {
			i, err := strconv.Atoi(chrday)
			if err != nil {
				return "", errors.New("неверный день недели")
			}
			if i < 0 || i > 6 {
				return "", errors.New("неверный день недели")
			}
		}

		for {
			date = date.AddDate(0, 0, 1)
			//Проверяем что день недели даты совпадает с одним из указанных дней недели
			if afterNow(date, now) && strings.Contains(weekdays, strconv.Itoa(int(date.Weekday()))) {
				break
			}
		}
		return date.Format(layout), nil
	case "m":
		//Массив дней
		var day [32]bool
		//Массив месяцев
		var month [13]bool
		//Признак последнего дня месяца
		var last bool
		//Признак предпоследнего дня месяца
		var prelast bool

		//Проверяем что массив состоит из 2 или 3 элементов
		switch len(sliceRepeat) {
		case 2, 3:
			sliceday := strings.Split(sliceRepeat[1], ",")
			for _, daystr := range sliceday {
				//Проверяем правильность дня месяца
				dayint, err := strconv.Atoi(daystr)
				if err != nil {
					return "", errors.New("неверный день месяца")
				}
				if dayint < -2 || dayint > 31 || dayint == 0 {
					return "", errors.New("неверный день месяца")
				}
				//Заполняем массив дней днями задачи
				if dayint > 0 {
					day[dayint] = true
					//Устанавливаем признак последнего дня месяца
				} else if dayint == -1 {
					last = true
					//Устанавливаем признак предпоследнего дня месяца
				} else {
					prelast = true
				}
			}

			//Проверяем и заполняем массив месяцев месяцами задачи. Если месяцы не указаны, то заполняем массив всеми месяцами
			if len(sliceRepeat) == 3 {
				slicemonth := strings.Split(sliceRepeat[2], ",")
				for _, monthstr := range slicemonth {
					monthint, err := strconv.Atoi(monthstr)
					if err != nil {
						return "", errors.New("неверный месяц")
					}
					if monthint < 1 || monthint > 12 {
						return "", errors.New("неверный месяц")
					}
					month[monthint] = true
				}
			} else {
				for monthint := 1; monthint < 13; monthint++ {
					month[monthint] = true
				}
			}

			for {
				date = date.AddDate(0, 0, 1)
				if afterNow(date, now) && month[date.Month()] {
					//Проверяем что день совпадает с одним из дней задачи
					if day[date.Day()] {
						break
					}

					//Проверяем что дата совпадает с предпоследним днем месяца
					if prelast {
						y, m, _ := date.Date()
						lastDay := lastDayOfMonth(m, y)
						preLastDay := lastDay.AddDate(0, 0, -1)
						if EqualDays(date, preLastDay) {
							break
						}
					}

					//Проверяем что дата совпадает с последним днем месяца
					if last {
						y, m, _ := date.Date()
						lastDay := lastDayOfMonth(m, y)
						if EqualDays(date, lastDay) {
							break
						}
					}
				}
			}
			return date.Format(layout), nil
		default:
			return "", errors.New("не указан день месяца")
		}
	default:
		return "", errors.New("недопустимый символ")
	}
}

// Функция проверяет что дата после текущей даты
func afterNow(date time.Time, now time.Time) bool {
	return date.Truncate(24 * time.Hour).After(now.Truncate(24 * time.Hour))
}

// Функция разбивает параметр repeat на массив строк
func split(repeat string) ([]string, error) {
	if len(repeat) < 1 {
		return []string{}, errors.New("неверный формат правила повторения")
	}
	slice := strings.Split(repeat, " ")
	return slice, nil
}

// Функция возвращает последний день месяца
func lastDayOfMonth(m time.Month, y int) time.Time {
	return time.Date(y, m+1, 0, 0, 0, 0, 0, time.UTC)
}

// Функция проверяет равны ли даты днями, не учитывая время
func EqualDays(d1, d2 time.Time) bool {
	return d1.Truncate(24 * time.Hour).Equal(d2.Truncate(24 * time.Hour))
}
