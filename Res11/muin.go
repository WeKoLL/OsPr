package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/fogleman/gg"
	"gopkg.in/telebot.v3"
)

// BotConfig содержит конфигурационные параметры бота
type BotConfig struct {
	TelegramAPIToken string
}

func main() {
	// Инициализация конфигурации
	config := BotConfig{
		TelegramAPIToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}
	
	if config.TelegramAPIToken == "" {
		log.Fatal("Требуется переменная окружения TELEGRAM_BOT_TOKEN")
	}

	// Настройки Telegram бота
	botSettings := telebot.Settings{
		Token:  config.TelegramAPIToken,
		Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
	}

	bot, err := telebot.NewBot(botSettings)
	if err != nil {
		log.Fatal("Ошибка инициализации бота: ", err)
	}

	// Приветственное сообщение
	bot.Handle("/start", func(ctx telebot.Context) error {
		welcomeMessage := `✨ <b>Графический калькулятор площадей</b> ✨

Я помогу вычислить площадь под графиком функции на заданном интервале.

📝 <b>Формат запроса:</b>
<code>функция нижняя_граница верхняя_граница</code>

📌 <b>Примеры:</b>
<code>sin(x) 0 3.14</code> - площадь синусоиды
<code>x^2 0 5</code> - площадь параболы
<code>sqrt(4-x^2) -2 2</code> - площадь полукруга

🔧 <b>Поддерживаемые операции:</b> +, -, *, /, ^
📐 <b>Функции:</b> sin, cos, tan, sqrt, log, exp, abs
📏 <b>Константы:</b> pi, e`

		return ctx.Send(welcomeMessage, telebot.ModeHTML)
	})

	// Обработчик математических запросов
	bot.Handle(telebot.OnText, func(ctx telebot.Context) error {
		userInput := ctx.Text()
		inputParts := strings.Fields(userInput)
		
		if len(inputParts) != 3 {
			return ctx.Reply("❌ Неверный формат. Используйте: <функция> <a> <b>\nПример: sin(x) 0 3.14")
		}

		funcExpression := inputParts[0]
		lowerBound, err1 := strconv.ParseFloat(inputParts[1], 64)
		upperBound, err2 := strconv.ParseFloat(inputParts[2], 64)
		
		if err1 != nil || err2 != nil {
			return ctx.Reply("🔢 Ошибка в границах интервала. Убедитесь, что это числа.\nПример: 0 3.14")
		}

		if lowerBound >= upperBound {
			return ctx.Reply("↗️ Левая граница должна быть меньше правой!")
		}

		computedArea, calcErr := computeIntegral(funcExpression, lowerBound, upperBound)
		if calcErr != nil {
			return ctx.Reply(fmt.Sprintf("🧮 Ошибка вычисления: %v", calcErr))
		}

		graphImage, plotErr := createFunctionPlot(funcExpression, lowerBound, upperBound, computedArea)
		if plotErr != nil {
			return ctx.Reply(fmt.Sprintf("🖼️ Ошибка построения графика: %v", plotErr))
		}

		resultCaption := fmt.Sprintf(
			"📊 <b>Результат</b>\n\nФункция: <code>%s</code>\nИнтервал: [%.2f, %.2f]\nПлощадь: <code>%.4f</code>",
			funcExpression, lowerBound, upperBound, computedArea,
		)
		
		return ctx.Send(
			&telebot.Photo{File: telebot.FromReader(graphImage), Caption: resultCaption},
			telebot.ModeHTML,
		)
	})

	log.Println("✅ Бот успешно запущен")
	bot.Start()
}

func computeIntegral(funcExpr string, start, end float64) (float64, error) {
	const segments = 1000
	stepSize := (end - start) / segments
	totalArea := 0.0

	for i := 0; i <= segments; i++ {
		currentX := start + float64(i)*stepSize
		funcValue, err := evaluateMathExpression(funcExpr, currentX)
		if err != nil {
			return 0, err
		}

		// Правило трапеций
		if i == 0 || i == segments {
			totalArea += funcValue
		} else {
			totalArea += 2 * funcValue
		}
	}

	return totalArea * stepSize / 2, nil
}

func evaluateMathExpression(expr string, xValue float64) (float64, error) {
	// Заменяем синтаксис степеней для совместимости
	expr = strings.ReplaceAll(expr, "^", "**")

	functions := map[string]govaluate.ExpressionFunction{
		"sin": func(args ...interface{}) (interface{}, error) {
			return math.Sin(args[0].(float64)), nil
		},
		"cos": func(args ...interface{}) (interface{}, error) {
			return math.Cos(args[0].(float64)), nil
		},
		"tan": func(args ...interface{}) (interface{}, error) {
			return math.Tan(args[0].(float64)), nil
		},
		"sqrt": func(args ...interface{}) (interface{}, error) {
			return math.Sqrt(args[0].(float64)), nil
		},
		"log": func(args ...interface{}) (interface{}, error) {
			return math.Log(args[0].(float64)), nil
		},
		"exp": func(args ...interface{}) (interface{}, error) {
			return math.Exp(args[0].(float64)), nil
		},
		"abs": func(args ...interface{}) (interface{}, error) {
			return math.Abs(args[0].(float64)), nil
		},
	}

	parameters := map[string]interface{}{
		"x":  xValue,
		"pi": math.Pi,
		"e":  math.E,
	}

	expression, err := govaluate.NewEvaluableExpressionWithFunctions(expr, functions)
	if err != nil {
		return 0, fmt.Errorf("ошибка в выражении: %v", err)
	}

	result, err := expression.Evaluate(parameters)
	if err != nil {
		return 0, fmt.Errorf("ошибка вычисления: %v", err)
	}

	switch v := result.(type) {
	case float64:
		return v, nil
	case int:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("неподдерживаемый тип результата: %T", result)
	}
}

func createFunctionPlot(funcExpr string, start, end, area float64) (*bytes.Buffer, error) {
	const (
		imageWidth  = 800
		imageHeight = 600
		margin      = 50.0
	)

	canvas := gg.NewContext(imageWidth, imageHeight)

	// Оформление фона
	canvas.SetColor(color.RGBA{245, 245, 245, 255})
	canvas.Clear()
	canvas.SetColor(color.Black)

	// Оси координат
	canvas.SetLineWidth(2)
	canvas.DrawLine(margin, imageHeight-margin, imageWidth-margin, imageHeight-margin)
	canvas.DrawLine(margin, margin, margin, imageHeight-margin)
	canvas.Stroke()

	// Вычисление точек графика
	pointCount := imageWidth - int(2*margin)
	values := make([]float64, pointCount)
	
	for i := range values {
		x := start + (end-start)*float64(i)/float64(pointCount-1)
		y, err := evaluateMathExpression(funcExpr, x)
		if err != nil {
			return nil, err
		}
		values[i] = y
	}

	// Определение масштаба
	minVal, maxVal := math.Inf(1), math.Inf(-1)
	for _, y := range values {
		if y < minVal {
			minVal = y
		}
		if y > maxVal {
			maxVal = y
		}
	}

	verticalScale := (imageHeight - 2*margin) / (maxVal - minVal)

	// Заливка области под кривой
	canvas.SetColor(color.RGBA{70, 130, 180, 150})
	canvas.MoveTo(margin, imageHeight-margin)
	for i, y := range values {
		x := margin + float64(i)
		scaledY := imageHeight - margin - (y-minVal)*verticalScale
		canvas.LineTo(x, scaledY)
	}
	canvas.LineTo(imageWidth-margin, imageHeight-margin)
	canvas.ClosePath()
	canvas.Fill()

	// Рисование графика
	canvas.SetColor(color.RGBA{25, 25, 112, 255})
	canvas.SetLineWidth(3)
	for i, y := range values {
		x := margin + float64(i)
		scaledY := imageHeight - margin - (y-minVal)*verticalScale
		if i == 0 {
			canvas.MoveTo(x, scaledY)
		} else {
			canvas.LineTo(x, scaledY)
		}
	}
	canvas.Stroke()

	// Добавление текста
	canvas.SetColor(color.Black)
	if err := canvas.LoadFontFace("Arial.ttf", 20); err == nil {
		canvas.DrawStringAnchored(fmt.Sprintf("y = %s", funcExpr), imageWidth/2, margin+30, 0.5, 0.5)
		canvas.DrawStringAnchored(fmt.Sprintf("Площадь: %.4f", area), imageWidth/2, imageHeight-margin-30, 0.5, 0.5)
	} else {
		canvas.DrawStringAnchored(fmt.Sprintf("y = %s", funcExpr), imageWidth/2, margin+20, 0.5, 0.5)
		canvas.DrawStringAnchored(fmt.Sprintf("Площадь: %.4f", area), imageWidth/2, imageHeight-margin-20, 0.5, 0.5)
	}

	// Сохранение в буфер
	imageBuffer := new(bytes.Buffer)
	if err := canvas.EncodePNG(imageBuffer); err != nil {
		return nil, fmt.Errorf("ошибка генерации изображения: %v", err)
	}

	return imageBuffer, nil
}