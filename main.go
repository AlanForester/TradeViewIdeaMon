package main

import (
	"github.com/AlexCollin/TradeViewIdeaMon/libs"
	"github.com/AlexCollin/TradeViewIdeaMon/sql"
	"log"
	"time"
)

//var (
//	originalWindow = new(opencv.Window)
//	image          = new(opencv.IplImage)
//	seq            *opencv.Seq
//	redColor       = opencv.NewScalar(0, 0, 255, 255) // (blue, green, red, alpha)
//	blackColor     = opencv.NewScalar(0, 0, 0, 255)   // (blue, green, red, alpha)
//	blueColor      = opencv.NewScalar(255, 0, 0, 255) // (blue, green, red, alpha)
//	slider         = 15
//)

func main() {

	db := new(sql.Postgres)
	sql.DB = db.Connect()
	sql.DB.AutoMigrate(&sql.Author{})
	sql.DB.AutoMigrate(&sql.Post{})
	sql.DB.AutoMigrate(&sql.User{})

	newPostCh := make(chan sql.Post)

	tick := time.NewTicker(10 * time.Second)
	go func() {
		for range tick.C {
			scraper := libs.Scraper{}
			log.Printf("Tick ")
			results := scraper.GetLastIdeas(newPostCh)

			for _, post := range results {
				db := sql.DB.First(&post)
				if db.RowsAffected == 0 {
					db.Save(&post)
				}
			}
			log.Printf("Scraper res: %v\n", results)
		}
	}()

	tbot := new(libs.Telebot)

	go tbot.Start()

	tbot.Sender(newPostCh)
	//urlIdea := "https://www.tradingview.com/chart/AUDCAD/Yr6omBQt-bay-D/"
	//err := libs.Screenshot(urlIdea)
	//if err != nil {
	//	panic(err)
	//}

	//image = opencv.LoadImage(urlIdea)
	//if image == nil {
	//	panic("LoadImage failed")
	//}
	//defer image.Release()
	//
	//originalWindow = opencv.NewWindow("Find contours in image")
	//defer originalWindow.Destroy()
	//
	//originalWindow.CreateTrackbar("Levels : ", 0, slider, trackBar)
	//
	//// initialize by drawing the top level contours(level = 0)
	//trackBar(0, 0)
	//
	//for {
	//	key := opencv.WaitKey(20)
	//	if key == 27 {
	//		os.Exit(0)
	//	}
	//}
	//
	//os.Exit(0)
}

//
//func trackBar(position int, param ...interface{}) {
//
//	width := image.Width()
//	height := image.Height()
//
//	// Convert to grayscale
//	gray := opencv.CreateImage(width, height, opencv.IPL_DEPTH_8U, 1)
//	defer gray.Release()
//
//	opencv.CvtColor(image, gray, opencv.CV_BGR2GRAY)
//
//	// for edge detection
//	cannyImage := opencv.CreateImage(width, height, opencv.IPL_DEPTH_8U, 1)
//	defer cannyImage.Release()
//
//	// Run the edge detector on grayscale
//	//opencv.Canny(gray, cannyImage, float64(position), float64(position)*2, 3)
//
//	// ** For better result, use 50 for the canny threshold instead of tying the value
//	//    to the track bar position
//	opencv.Canny(gray, cannyImage, float64(50), float64(50)*2, 3)
//
//	// Find contours sequence from canny edge processed image
//	// see http://docs.opencv.org/2.4/modules/imgproc/doc/structural_analysis_and_shape_descriptors.html
//	// for mode and method
//
//	seq = cannyImage.FindContours(opencv.CV_RETR_TREE, opencv.CV_CHAIN_APPROX_NONE, opencv.Point{0, 0})
//	defer seq.Release()
//
//	// based on the sequence, draw the contours
//	// back on the original image
//	finalImage := image.Clone()
//
//	maxLevel := position
//	opencv.DrawContours(finalImage, seq, redColor, blueColor, maxLevel, 2, 8, opencv.Point{0, 0})
//	originalWindow.ShowImage(finalImage)
//
//	fmt.Printf("Levels  = %d\n", position)
//}
