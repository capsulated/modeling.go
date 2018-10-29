package main

import (
	"github.com/gonum/stat/distuv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"gonum.org/v1/plot/plotter"
	"log"
	"math"
	"math/rand"
	"net/http"
	"time"
)

var apiPath = "/modeling"

type Trials struct {
	Trials int `json:"trials" form:"trials" query:"trials"`
}

type TrialsDices struct {
	Trials int `json:"trials" form:"trials" query:"trials"`
	Dices int `json:"dices" form:"dices" query:"dices"`
}

type Response struct {
	Bins   []float64
	Values []float64
}

type Value struct {
	Num int
	Val float64
}

// todo генерить нормальные массивы для вывода
func main() {
	// todo
	// var cfgArg string = "/etc/binatex/quoteBroadcaster.yaml"
	// flag.StringVar(&cfgArg, "cfg", cfgArg, "config file location")
	// flag.Parse()

	rand.Seed(time.Now().UnixNano())

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:8080", "http://localhost", "http://sci.logiq.one"},
		//AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, "X-Token"},
		AllowMethods: []string{echo.GET},
	}))
	e.POST(apiPath+"/onedice", oneDiceHandler)
	e.POST(apiPath+"/twodice", twoDiceHandler)
	e.POST(apiPath+"/normal", normalDistributionHandler)
	e.POST(apiPath+"/exponential", exponentialDistributionHandler)
	e.POST(apiPath+"/advanced", advancedNormalDistributionHandler)
	e.POST(apiPath+"/goadvanced", goAdvancedNormalDistributionHandler)
	e.Logger.Fatal(e.Start(":8000"))
}

func oneDiceHandler(c echo.Context) error {

	t := new(Trials)

	if err := c.Bind(t); err != nil {
		return err
	}

	rand.Seed(time.Now().UnixNano())

	values := make([]int, t.Trials)

	for i := 0; i < t.Trials; i++ {
		values[i] = rand.Intn(6) + 1
	}

	//fmt.Printf("%#v\n", values)

	response := &Response{
		Bins:   []float64{0, 1, 2, 3, 4, 5, 6},
		Values: make([]float64, 7),
	}

	for i := 0; i < t.Trials; i++ {
		response.Values[int(values[i])]++
	}

	//fmt.Printf("%#v\n", response)

	return c.JSON(http.StatusOK, response)
}

func twoDiceHandler(c echo.Context) error {
	rand.Seed(time.Now().UnixNano())
	t := new(Trials)
	if err := c.Bind(t); err != nil {
		return err
	}

	response := &Response{
		Bins:   []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14},
		Values: make([]float64, 14, 14),
	}

	for i := 0; i < t.Trials; i++ {
		firstDice := rand.Intn(6) + 1
		secondDice := rand.Intn(6) + 1
		//fmt.Printf("fi: %d se: %d sum %d \n", firstDice, secondDice, firstDice+secondDice)

		response.Values[firstDice+secondDice - 1]++
	}

	return c.JSON(http.StatusOK, response)
}

func normalDistributionHandler(c echo.Context) error {
	var step = 13
	var mu = 7.0
	var sigma = 2.2

	t := &Trials{
		Trials: 100,
	}
	if err := c.Bind(t); err != nil {
		return err
	}

	// Create a normal distribution
	dist := distuv.Normal{
		Mu:    mu,
		Sigma: sigma,
	}

	data := make(plotter.Values, t.Trials)

	// Draw some random values from the standard normal distribution
	for i := range data {
		data[i] = dist.Rand()
	}

	//mean, std := stat.MeanStdDev(data, nil)
	//meanErr := stat.StdErr(std, float64(len(data)))

	//fmt.Printf("%#v\n", data)

	//fmt.Printf("mean= %1.1f ± %0.1v\n", mean, meanErr)

	hist, err := plotter.NewHist(data, step)
	if err != nil {
		log.Panic(err)
	}

	response := &Response{
		Bins:   make([]float64, step),
		Values: make([]float64, step),
	}

	for i, bin := range hist.Bins {
		mean := (bin.Max-bin.Min)/2 + bin.Min
		response.Bins[i] = math.Round(mean*100) / 100
		//response.Bins[i] = float64(int(mean))
		response.Values[i] = bin.Weight
	}

	return c.JSON(http.StatusOK, response)
}

func exponentialDistributionHandler(c echo.Context) error {
	// var arrival = godes.NewExpDistr(true)
	var step = 20
	var lambda = 7.0

	t := &Trials{
		Trials: 100,
	}

	if err := c.Bind(t); err != nil {
		return err
	}

	// Create a normal distribution
	dist := distuv.Exponential{
		Rate:    lambda,
	}

	values := make(plotter.Values, t.Trials)
	for i := 1; i < t.Trials; i++ {
		// values[i] = arrival.Get(lambda)
		values[i] = dist.Rand()
	}

	//log.Printf("%#v", values)

	hist, err := plotter.NewHist(values, step)
	if err != nil {
		log.Panic(err)
	}

	//log.Printf("Historgam: %v", hist)

	response := &Response{
		Bins:   make([]float64, step),
		Values: make([]float64, step),
	}

	for i, bin := range hist.Bins {
		mean := (bin.Max-bin.Min)/2 + bin.Min
		response.Bins[i] = math.Round(mean*100) / 100
		response.Values[i] = bin.Weight
	}

	//log.Printf("Resp Bins: %v", response.Bins)

	return c.JSON(http.StatusOK, response)
}

func advancedNormalDistributionHandler(c echo.Context) error {
	var step = 20

	t := &TrialsDices{
		Trials: 150000,
		Dices: 200,
	}

	if err := c.Bind(t); err != nil {
		return err
	}

	data := make(plotter.Values, t.Trials)

	// Draw some random values from the standard normal distribution
	for i := range data {
		for j := 1; j <= t.Dices; j++ {
			data[i] += float64(rand.Intn(6) + 1)
		}
	}

	//mean, std := stat.MeanStdDev(data, nil)
	//meanErr := stat.StdErr(std, float64(len(data)))

	//fmt.Printf("%#v\n", data)

	//fmt.Printf("mean= %1.1f ± %0.1v\n", mean, meanErr)

	hist, err := plotter.NewHist(data, step)
	if err != nil {
		log.Panic(err)
	}

	response := &Response{
		Bins:   make([]float64, step),
		Values: make([]float64, step),
	}

	for i, bin := range hist.Bins {
		mean := (bin.Max-bin.Min)/2 + bin.Min
		//response.Bins[i] = math.Round(mean*100) / 100
		response.Bins[i] = float64(int(mean))
		response.Values[i] = bin.Weight
	}

	return c.JSON(http.StatusOK, response)
}

func goAdvancedNormalDistributionHandler(c echo.Context) error {
	var step = 20

	t := &TrialsDices{
		Trials: 150000,
		Dices: 200,
	}

	if err := c.Bind(t); err != nil {
		return err
	}

	data := make(plotter.Values, t.Trials)

	ch := make(chan Value)

	for i := range data {
		go func(in chan<- Value, it int, dices int) {
			var sum = 0.0
			for j := 1; j <= dices; j++ {
				sum += float64(rand.Intn(6) + 1)
			}
			in <- Value{
				Num: it,
				Val: sum,
			}
		}(ch, i, t.Dices)
	}

	for v := range ch {
		data[v.Num] = v.Val
	}

	//mean, std := stat.MeanStdDev(data, nil)
	//meanErr := stat.StdErr(std, float64(len(data)))

	//fmt.Printf("%#v\n", data)

	//fmt.Printf("mean= %1.1f ± %0.1v\n", mean, meanErr)

	hist, err := plotter.NewHist(data, step)
	if err != nil {
		log.Panic(err)
	}

	response := &Response{
		Bins:   make([]float64, step),
		Values: make([]float64, step),
	}

	for i, bin := range hist.Bins {
		mean := (bin.Max-bin.Min)/2 + bin.Min
		//response.Bins[i] = math.Round(mean*100) / 100
		response.Bins[i] = float64(int(mean))
		response.Values[i] = bin.Weight
	}

	return c.JSON(http.StatusOK, response)
}