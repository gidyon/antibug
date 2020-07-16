package antibiogram

import (
	"context"
	"fmt"
	"github.com/Pallinder/go-randomdata"
	"github.com/gidyon/antibug/internal/mocks"
	"github.com/gidyon/antibug/internal/modules/culture"
	antibiogram "github.com/gidyon/antibug/pkg/api/antibiogram"
	"github.com/gidyon/micros"
	"github.com/go-redis/redis"
	"math/rand"
	"time"

	// Imports mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"testing"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"
)

func TestAntibiogramService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Service Suite")
}

var (
	AntibiogramAPI    antibiogram.AntibiogramAPIServer
	AntibiogramServer *apiServer
)

const (
	dbName       = "antibug"
	dbAddress    = "localhost:3306"
	redisAddress = "localhost:6379"
)

func initDB() (*gorm.DB, error) {
	param := "charset=utf8&parseTime=true"
	dsn := fmt.Sprintf("root:hakty11@tcp(%s)/%s?%s", dbAddress, dbName, param)
	return gorm.Open("mysql", dsn)
}

var _ = BeforeSuite(func() {
	// Seed random number generator
	rand.Seed(time.Now().UnixNano())

	ctx := context.Background()

	db, err := initDB()
	Expect(err).ShouldNot(HaveOccurred())

	redisDB := redis.NewClient(&redis.Options{
		Addr: redisAddress,
	})

	opt := &Options{
		SQLDB:         db,
		Logger:        micros.NewLogger("antibiogram_app"),
		RedisDB:       redisDB,
		JWTSigningKey: randomdata.RandStringRunes(32),
	}

	AntibiogramAPI, err = NewAntibiogramAPIServer(ctx, opt)
	Expect(err).ShouldNot(HaveOccurred())

	var ok bool
	AntibiogramServer, ok = AntibiogramAPI.(*apiServer)
	Expect(ok).Should(BeTrue())

	// Mock auth
	AntibiogramServer.authAPI = mocks.AuthAPI

	_, err = NewAntibiogramAPIServer(nil, opt)
	Expect(err).Should(HaveOccurred())

	opt.SQLDB = nil
	_, err = NewAntibiogramAPIServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.SQLDB = db
	opt.RedisDB = nil
	_, err = NewAntibiogramAPIServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.RedisDB = redisDB
	opt.Logger = nil
	_, err = NewAntibiogramAPIServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.Logger = micros.NewLogger("antibiogram_app")
	opt.JWTSigningKey = ""
	_, err = NewAntibiogramAPIServer(ctx, opt)
	Expect(err).Should(HaveOccurred())

	opt.JWTSigningKey = randomdata.RandStringRunes(32)

	culture.InitCulture()

	// var count = 10000

	// opt.Logger.Infof("Creating %d cultures", count)

	// for i := 0; i < count; i++ {
	// 	culturePB := culture.FakeCulture()
	// 	cultureDB, err := culture.GetCultureDB(culturePB)
	// 	Expect(err).ShouldNot(HaveOccurred())

	// 	err = AntibiogramServer.sqlDB.Create(cultureDB).Error
	// 	Expect(err).ShouldNot(HaveOccurred())
	// }
})

var _ = AfterSuite(func() {
	Expect(AntibiogramServer.sqlDB.Close()).ShouldNot(HaveOccurred())
})

// Declarations for Ginkgo DSL
type Done ginkgo.Done
type Benchmarker ginkgo.Benchmarker

var GinkgoWriter = ginkgo.GinkgoWriter
var GinkgoRandomSeed = ginkgo.GinkgoRandomSeed
var GinkgoParallelNode = ginkgo.GinkgoParallelNode
var GinkgoT = ginkgo.GinkgoT
var CurrentGinkgoTestDescription = ginkgo.CurrentGinkgoTestDescription
var RunSpecs = ginkgo.RunSpecs
var RunSpecsWithDefaultAndCustomReporters = ginkgo.RunSpecsWithDefaultAndCustomReporters
var RunSpecsWithCustomReporters = ginkgo.RunSpecsWithCustomReporters
var Skip = ginkgo.Skip
var Fail = ginkgo.Fail
var GinkgoRecover = ginkgo.GinkgoRecover
var Describe = ginkgo.Describe
var FDescribe = ginkgo.FDescribe
var PDescribe = ginkgo.PDescribe
var XDescribe = ginkgo.XDescribe
var Context = ginkgo.Context
var FContext = ginkgo.FContext
var PContext = ginkgo.PContext
var XContext = ginkgo.XContext
var When = ginkgo.When
var FWhen = ginkgo.FWhen
var PWhen = ginkgo.PWhen
var XWhen = ginkgo.XWhen
var It = ginkgo.It
var FIt = ginkgo.FIt
var PIt = ginkgo.PIt
var XIt = ginkgo.XIt
var Specify = ginkgo.Specify
var FSpecify = ginkgo.FSpecify
var PSpecify = ginkgo.PSpecify
var XSpecify = ginkgo.XSpecify
var By = ginkgo.By
var Measure = ginkgo.Measure
var FMeasure = ginkgo.FMeasure
var PMeasure = ginkgo.PMeasure
var XMeasure = ginkgo.XMeasure
var BeforeSuite = ginkgo.BeforeSuite
var AfterSuite = ginkgo.AfterSuite
var SynchronizedBeforeSuite = ginkgo.SynchronizedBeforeSuite
var SynchronizedAfterSuite = ginkgo.SynchronizedAfterSuite
var BeforeEach = ginkgo.BeforeEach
var JustBeforeEach = ginkgo.JustBeforeEach
var JustAfterEach = ginkgo.JustAfterEach
var AfterEach = ginkgo.AfterEach

// Declarations for Gomega DSL
var RegisterFailHandler = gomega.RegisterFailHandler
var RegisterFailHandlerWithT = gomega.RegisterFailHandlerWithT
var RegisterTestingT = gomega.RegisterTestingT
var InterceptGomegaFailures = gomega.InterceptGomegaFailures
var Ω = gomega.Ω
var Expect = gomega.Expect
var ExpectWithOffset = gomega.ExpectWithOffset
var Eventually = gomega.Eventually
var EventuallyWithOffset = gomega.EventuallyWithOffset
var Consistently = gomega.Consistently
var ConsistentlyWithOffset = gomega.ConsistentlyWithOffset
var SetDefaultEventuallyTimeout = gomega.SetDefaultEventuallyTimeout
var SetDefaultEventuallyPollingInterval = gomega.SetDefaultEventuallyPollingInterval
var SetDefaultConsistentlyDuration = gomega.SetDefaultConsistentlyDuration
var SetDefaultConsistentlyPollingInterval = gomega.SetDefaultConsistentlyPollingInterval
var NewWithT = gomega.NewWithT
var NewGomegaWithT = gomega.NewGomegaWithT

// Declarations for Gomega Matchers
var Equal = gomega.Equal
var BeEquivalentTo = gomega.BeEquivalentTo
var BeIdenticalTo = gomega.BeIdenticalTo
var BeNil = gomega.BeNil
var BeTrue = gomega.BeTrue
var BeFalse = gomega.BeFalse
var HaveOccurred = gomega.HaveOccurred
var Succeed = gomega.Succeed
var MatchError = gomega.MatchError
var BeClosed = gomega.BeClosed
var Receive = gomega.Receive
var BeSent = gomega.BeSent
var MatchRegexp = gomega.MatchRegexp
var ContainSubstring = gomega.ContainSubstring
var HavePrefix = gomega.HavePrefix
var HaveSuffix = gomega.HaveSuffix
var MatchJSON = gomega.MatchJSON
var MatchXML = gomega.MatchXML
var MatchYAML = gomega.MatchYAML
var BeEmpty = gomega.BeEmpty
var HaveLen = gomega.HaveLen
var HaveCap = gomega.HaveCap
var BeZero = gomega.BeZero
var ContainElement = gomega.ContainElement
var BeElementOf = gomega.BeElementOf
var ConsistOf = gomega.ConsistOf
var HaveKey = gomega.HaveKey
var HaveKeyWithValue = gomega.HaveKeyWithValue
var BeNumerically = gomega.BeNumerically
var BeTemporally = gomega.BeTemporally
var BeAssignableToTypeOf = gomega.BeAssignableToTypeOf
var Panic = gomega.Panic
var BeAnExistingFile = gomega.BeAnExistingFile
var BeARegularFile = gomega.BeARegularFile
var BeADirectory = gomega.BeADirectory
var And = gomega.And
var SatisfyAll = gomega.SatisfyAll
var Or = gomega.Or
var SatisfyAny = gomega.SatisfyAny
var Not = gomega.Not
var WithTransform = gomega.WithTransform
