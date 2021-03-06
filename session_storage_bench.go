package main

import (
	"flag"
	"fmt"
	"math"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"time"
)

var numConnections = flag.Int("no_connections", 100, "Number of parallel connections")
var storageType = flag.String("storage_type", "memcache", "Session backend: memcache(default) or redis")
var usePersistentConnections = flag.Bool("persistent", false, "Use persistent connections")
var numCpu = flag.Int("num_cpu", 1, "Number of used CPU")
var numRequests = flag.Int("num_requests", 100, "Number of requests")
var storageAddress = flag.String("storage_address", "127.0.0.1:11211", "Tested storage instance host and port")

var SESSION_CONTENT = []byte(`tracking|a:2:{s:6:"source";s:6:"direct";s:6:"medium";s:4:"none";}84c0dc92cba6c36e15f2cddb26232636__returnUrl|s:17:"/customer/account";SessionExpireTimestamp|i:1388720260;client_type|s:7:"desktop";cartParentSalesOrderItems|a:0:{}cartSimples|a:0:{}wt3_eid|s:19:"2137947614700087480";__utmz|s:69:"84101705.1379476148.1.1.utmcsr=(direct)|utmccn=(direct)|utmcmd=(none)";lastUrlkeyCategory|s:1:"/";catalog_placeholder_cache|a:61:{s:16:"NE771HLAJK34ANMY";i:5;s:16:"OE702EL52SMLANMY";i:1;s:16:"OE702ELAJDG0ANMY";i:4;s:16:"AP564ELAJORPANMY";i:5;s:16:"AP564ELAJORRANMY";i:5;s:16:"SA356ELAD85PANMY";i:4;s:16:"SA356ELAD85OANMY";i:5;s:16:"YO498ELAJHECANMY";i:3;s:16:"ID738EL56NAXANMY";i:3;s:16:"ID738EL55NAYANMY";i:4;s:16:"AI880EL40QXNANMY";i:4;s:16:"HT964ELAJL74ANMY";i:2;s:16:"SA356EL09CDEANMY";i:5;s:16:"SA356EL10CDDANMY";i:2;s:16:"TP481ELAJEF0ANMY";i:1;s:16:"WE539ELAJFM1ANMY";i:1;s:16:"WE539ELAJFM2ANMY";i:5;s:16:"WE539ELAJFM3ANMY";i:5;s:16:"WE539ELAJGSEANMY";i:5;s:16:"SA356EL03FPSANMY";i:1;s:16:"LE106ELAJA9UANMY";i:2;s:16:"LE106EL59UZOANMY";i:5;s:16:"KE060EL47PCEANMY";i:1;s:16:"BU826EL51QHSANMY";i:1;s:16:"IR017EL12LNFANMY";i:1;s:16:"EL637EL59LHOANMY";i:3;s:16:"XM625EL73AZUANMY";i:1;s:16:"CR734ELAJHDZANMY";i:1;s:16:"ED812ELAD8DOANMY";i:4;s:16:"ED812ELAD8DNANMY";i:1;s:16:"GA891ELAJK9GANMY";i:4;s:16:"SO406ELAJO8YANMY";i:4;s:16:"SA356ELAD7YUANMY";i:1;s:16:"SA356ELAD7YTANMY";i:2;s:16:"SA356ELAD7YVANMY";i:3;s:16:"FU879EL78RZXANMY";i:1;s:16:"FU879EL81RZUANMY";i:5;s:16:"FU879EL80RZVANMY";i:3;s:16:"FU879EL79RZWANMY";i:3;s:16:"FU879EL77RZYANMY";i:3;s:16:"HP961ELAJ3RKANMY";i:2;s:16:"CA673EL05TTAANMY";i:2;s:16:"CA673EL04TTBANMY";i:4;s:16:"SA356EL71QJWANMY";i:1;s:16:"SA356EL72QJVANMY";i:4;s:16:"NI220EL38FGRANMY";i:4;s:16:"PI935HL38JSBANMY";i:2;s:16:"AS673HL22PORANMY";i:3;s:16:"NE771HL84TXRANMY";i:4;s:16:"RO641HL70DJLANMY";i:2;s:16:"NE913HL75HOUANMY";i:4;s:16:"IB648HB02PLPANMY";i:1;s:16:"SP661HB12MILANMY";i:2;s:16:"LK013HBAJ3MOANMY";i:4;s:16:"SK969ELAD83FANMY";i:2;s:16:"BR924HB27JDCANMY";i:2;s:16:"KA299MEAJN51ANMY";i:2;s:16:"KA299MEAJN50ANMY";i:2;s:16:"AL989MEAJOZWANMY";i:2;s:16:"MP991MEAJOZXANMY";i:5;s:16:"MP991MEAJOZYANMY";i:2;}YII_CSRF_TOKEN|s:40:"124a79af22b922d9366bd8cfe92a6ad64291bc45";profiling_counts|a:0:{}URI_history|a:3:{i:0;s:1:"/";i:1;s:35:"/mega-deals/mega-deals-subcategory/";i:2;s:1:"/";}LastSearchterm|s:0:"";activeFacets|a:1:{s:14:"facet_category";a:16:{s:4:"name";s:8:"Category";s:5:"param";s:8:"category";s:4:"view";s:16:"category-segment";s:7:"display";i:1;s:13:"displayMobile";i:0;s:11:"displayZero";b:1;s:6:"expand";i:1;s:9:"showEmpty";i:0;s:14:"multipleSelect";i:2;s:9:"resolveId";i:1;s:12:"facet_search";s:14:"facet_category";s:5:"multi";i:0;s:7:"combine";s:3:"AND";s:6:"weight";i:100;s:10:"filterArgs";a:1:{s:6:"filter";i:513;}s:5:"value";s:4:"3895";}}LastViewedCategoryId|s:4:"3895";continueShoppingUrl|s:54:"http://alice2.local/mega-deals/mega-deals-subcategory/";persistentSessionId|s:26:"984ktf1mtvt7v25is0t9mh7rv1";last_searchresult_configids_sessionkey|a:0:{}laz10049@LAZHCM10049:/media/DATA/`)

func generateRequests(inRequestsChanel chan RequestInfo, osSignal chan os.Signal, numRequests int) error {
	for i := 0; i < numRequests; i++ {
		select {
		case sign := <-osSignal:
			fmt.Printf("\nCatch signal %#v\n", sign)
			close(inRequestsChanel)
			return nil
		default:
			request := NewBaseRequestInfo()
			if *storageType == "memcache" {
				inRequestsChanel <- NewMemcacheRequestInfo(request)
			} else if *storageType == "redis" {
				inRequestsChanel <- NewRedisRequestInfo(request)
			} else {
				panic("Invalid backend type")
			}
		}
	}
	close(inRequestsChanel)

	return nil
}
func makeRequest(this RequestInfo) {
	if !this.initConnect() {
		return
	}

	if !this.addLock() {
		return
	}

	if !this.getData() {
		return
	}

	time.Sleep(200 * time.Millisecond)

	if !this.setData() {
		return
	}

	if !this.deleteLock() {
		return
	}

	if !this.closeConnect() {
		return
	}

	this.SetIsFailed(false)
}

func makeRequests(inRequestsChanel chan RequestInfo, outRequestsChanel chan RequestInfo, noParallelRoutines int) {
	routines := make(chan int, noParallelRoutines)
	numRoutines := 0

	defer close(outRequestsChanel)

	for request := range inRequestsChanel {
		if numRoutines >= noParallelRoutines {
			<-routines
			numRoutines--
		}

		go func(routines chan int, request RequestInfo, outRequestsChanel chan RequestInfo) {
			makeRequest(request)
			outRequestsChanel <- request
			routines <- 1

		}(routines, request, outRequestsChanel)

		numRoutines++
	}

	for i := 0; i < numRoutines; i++ {
		<-routines
	}
}

type Stats struct {
	TotalTime              int64
	LongestTransactionTime float64
	ShortesTransactionTime float64
	TotalFailed            int
	TotalSuccess           int
	TotalHttpErrors        int
	TotalContentLength     int64
}

func printStat(stats *Stats) {
	total := stats.TotalFailed + stats.TotalSuccess + stats.TotalHttpErrors
	var availability float64
	if total == 0 {
		return
	}

	title := "Total stats"
	fmt.Println("===================================================================")
	fmt.Printf("|| %s\n", title)
	fmt.Println("===================================================================")

	availability = 100.0 - float64(stats.TotalFailed+stats.TotalHttpErrors)/float64(total)*100.0
	fmt.Printf("Transactions: %d hits\n", total)
	fmt.Printf("Availability: %s %%\n", strconv.FormatFloat(availability, 'f', 2, 64))
	//fmt.Printf("Elapsed time(Lock,Get,Set,Unlock): %s \n", time.Duration(stats.TotalTime).String())
	fmt.Printf("Response time(Lock,Get,Set,Unlock): %s\n", time.Duration((stats.TotalTime / int64(total))).String())
	//fmt.Printf("Transaction rate: %s\n", strconv.FormatFloat(float64(total)/time.Duration(stats.TotalTime).Seconds(), 'f', 2, 64))
	fmt.Printf("Successful transactions: %d\n", stats.TotalSuccess)
	fmt.Printf("Failed transactions: %d\n", stats.TotalFailed)
	//fmt.Printf("HTTP error transactions: %d\n", stats.TotalHttpErrors)
	fmt.Printf("Longest transaction: %s \n", time.Duration(int64(stats.LongestTransactionTime)).String())
	fmt.Printf("Shortest transaction: %s \n", time.Duration(int64(stats.ShortesTransactionTime)).String())
}

func calculateStat(outRequestsChanel chan RequestInfo) {

	stats := new(Stats)
	stats.ShortesTransactionTime = math.MaxFloat64
	for request := range outRequestsChanel {
		if request == nil {
			break
		}
		if request.IsFailed() {
			stats.TotalFailed++
		} else {
			stats.TotalSuccess++
		}
		stats.TotalTime += request.GetDuration().Nanoseconds()

		stats.LongestTransactionTime = math.Max(stats.LongestTransactionTime, float64(request.GetDuration().Nanoseconds()))
		stats.ShortesTransactionTime = math.Min(stats.ShortesTransactionTime, float64(request.GetDuration().Nanoseconds()))
	}

	printStat(stats)
}

func main() {
	flag.Parse()
	runtime.GOMAXPROCS(*numCpu)

	inRequestsChanel := make(chan RequestInfo)
	outRequestsChanel := make(chan RequestInfo)
	osSignal := make(chan os.Signal, 1)

	signal.Notify(osSignal)

	go generateRequests(inRequestsChanel, osSignal, *numRequests)
	go makeRequests(inRequestsChanel, outRequestsChanel, *numConnections)
	calculateStat(outRequestsChanel)
}
