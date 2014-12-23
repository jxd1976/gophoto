package main

import (
	"crypto/md5"
	"fmt"
	"github.com/jxd1976/photo/util"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type ExifInfo struct {
	takendate string
	takentime string
	camModel  string
	md5       string
}

func main() {
	args := len(os.Args)
	if args != 4 {
		fmt.Printf("program description:\n")
		fmt.Printf("  This program for organize your photographs only,\n")
		fmt.Printf("the source files in the directory in accordance with \n")
		fmt.Printf("the taken date or last modified date is transferred to \n")
		fmt.Printf("the target directory.\n\n")
		fmt.Printf("usage:\n   %s <srcdir> <dstdir> <filetype[jpg|avi]>\n", os.Args[0])
	} else {
		srcdir := os.Args[1]
		dstdir := os.Args[2]
		filetype := os.Args[3]

		files, _ := WalkDir(srcdir, filetype)
		//fmt.Println(files)
		t := time.Now()
		outfile := util.GetCurrPath() + "log_" + strings.Replace(t.String()[:19], ":", "_", 3) + ".txt"

		fout, err := os.Create(outfile)
		if err != nil {
			fmt.Println(outfile, err)
			return
		}
		defer fout.Close()

		for i := 0; i < len(files); i++ {
			exifInfo, err := ExampleDecode(files[i])
			if err != nil {
				fmt.Println("图片", i, ":", files[i], err)
				strdate := util.Lastmodified(files[i])

				dir := dstdir + util.Fileseprater + util.Substr(strdate, 0, 4) + util.Fileseprater + util.Substr(strdate, 4, 2)
				if !util.Exist(dir) {
					util.MakeDirAll(dir)
				}
				destfile := dir + util.Fileseprater + util.ExtactFileName(files[i])
				fout.WriteString(fmt.Sprintf("%s,%s,%s\r\n", files[i], "--", destfile))
				_, err = util.CopyFile(files[i], destfile)
				if err != nil {
					fmt.Println("图片", i, ":", files[i], err)
				}
				continue
			}
			//fmt.Println((*exifInfo).camModel)
			//fmt.Println((*exifInfo).takendate)
			//fmt.Println((*exifInfo).md5)
			strdate := (*exifInfo).takendate
			dir := dstdir + util.Fileseprater + util.Substr(strdate, 0, 4) + util.Fileseprater + util.Substr(strdate, 4, 2)
			if !util.Exist(dir) {
				util.MakeDirAll(dir)
			}
			destfile := dir + util.Fileseprater + util.ExtactFileName(files[i])
			//fmt.Println(files[i], "--", destfile)
			fout.WriteString(fmt.Sprintf("%s,%s,%s\r\n", files[i], "--", destfile))
			_, err = util.CopyFile(files[i], destfile)
			if err != nil {
				fmt.Println("图片", i, ":", files[i], err)
				continue
			}
			//fout.WriteString(fmt.Sprintf("%s,%s,%s,%s\r\n",
			//	(*exifInfo).md5, (*exifInfo).takendate, (*exifInfo).camModel, files[i]))

		}
	}
}

func ExampleDecode(fname string) (exifInfo *ExifInfo, err error) {

	pr := new(ExifInfo)
	f, err := os.Open(fname)
	if err != nil {
		return pr, err
	}
	defer f.Close()

	exif.RegisterParsers(mknote.All...)

	x, err := exif.Decode(f)
	if err != nil {
		return pr, err
	}

	// normally, don't ignore errors!
	camModel, err := x.Get(exif.Model)
	if err != nil {
		return pr, err
	}
	//fmt.Println(camModel.StringVal())
	e1 := camModel.String()

	//focal, err := x.Get(exif.FocalLength)
	//if err != nil {
	//	return pr, err
	//}

	//// retrieve first (only) rat. value
	//numer, denom, err := focal.Rat2(0)
	//if err != nil {
	//return pr, err
	//}
	//fmt.Printf("%v/%v", numer, denom)

	// Two convenience functions exist for date/time taken and GPS coords:
	tm, err := x.DateTime()
	if err != nil {
		return pr, err
	}
	//fmt.Println("Taken: ", tm.Format("20060102150405"))
	e2 := tm.Format("20060102")
	e3 := tm.Format("150405")
	//lat, long, _ := x.LatLong()
	//fmt.Println("lat, long: ", lat, ", ", long)
	(*pr).camModel = e1
	(*pr).takendate = e2
	(*pr).takentime = e3

	md5h := md5.New()
	io.Copy(md5h, f)

	(*pr).md5 = fmt.Sprintf("%x", md5h.Sum([]byte("")))

	return pr, nil

}

func WalkDir(dirPth, suffix string) (files []string, err error) {
	files = make([]string, 0, 30)
	//忽略后缀匹配的大小写
	suffix = strings.ToUpper(suffix)
	err = filepath.Walk(dirPth, func(filename string, fi os.FileInfo, err error) error { //遍历目录
		//if err != nil { //忽略错误
		// return err
		//}
		if fi.IsDir() { // 忽略目录
			return nil
		}
		if strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) {
			files = append(files, filename)
		}
		return nil
	})
	return files, err
}
