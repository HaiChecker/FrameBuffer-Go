package lib

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"os/exec"
	"sync"
	"syscall"
	"unsafe"
)

var (
	FBIOGET_VSCREENINFO = 0x4600
)

type fbBitfield struct {
	offset    uint32 /* beginning of bitfield */
	length    uint32 /* length of bitfield */
	msb_right uint32 /* != 0 : Most significant bit is */
}

type fbVarInfo struct {
	xres         uint32 /* visible resolution */
	yres         uint32
	xres_virtual uint32 /* virtual resolution */
	yres_virtual uint32
	xoffset      uint32 /* offset from virtual to visible */
	yoffset      uint32 /* resolution */

	bits_per_pixel uint32 /* guess what */
	grayscale      uint32 /* 0 = color, 1 = grayscale, */
	/* >1 = FOURCC          */
	red    fbBitfield /* bitfield in fb mem if true color, */
	green  fbBitfield /* else only length is significant */
	blue   fbBitfield
	transp fbBitfield /* transparency         */

	nonstd uint32 /* != 0 Non standard pixel format */

	activate uint32 /* see FB_ACTIVATE_* */

	height uint32 /* height of picture in mm */
	width  uint32 /* width of picture in mm */

	accel_flags uint32 /* (OBSOLETE) see fb_info.flags */

	/* Timing: All values in pixclocks, except pixclock (of course) */
	pixclock     uint32 /* pixel clock in ps (pico seconds) */
	left_margin  uint32 /* time from sync to picture */
	right_margin uint32 /* time from picture to sync */
	upper_margin uint32 /* time from sync to picture */
	lower_margin uint32
	hsync_len    uint32    /* length of horizontal sync */
	vsync_len    uint32    /* length of vertical sync */
	sync         uint32    /* see FB_SYNC_* */
	vmode        uint32    /* see FB_VMODE_* */
	rotate       uint32    /* angle we rotate counter clockwise */
	colorspace   uint32    /* colorspace for FOURCC-based modes */
	reserved     [4]uint32 /* Reserved for future compatibility */
}

type AndroidFramebuffer struct {
	Image    image.Image
	filePath string
	lock     sync.RWMutex
}

func GetAndroidFramebuffer(devFile string) AndroidFramebuffer {
	if len(devFile) == 0 {
		devFile = "/dev/graphics/fb0"
	}
	return AndroidFramebuffer{filePath: devFile}
}

func (f *AndroidFramebuffer) getPermission() {
	out, err := exec.Command("/system/bin/su", "-c", "/system/bin/sh", "chmod 777 "+f.filePath).CombinedOutput()

	if err != nil {
		log.Printf("改变FB权限失败:%v error:%v", string(out), err)
	} else {
		log.Printf("获取FB权限成功:%v", string(out))
	}
}
func (f *AndroidFramebuffer) releasePermission() {
	out, err := exec.Command("/system/bin/su", "-c", "/system/bin/sh", "chmod 660 "+f.filePath).CombinedOutput()
	if err != nil {
		log.Printf("释放FB权限失败:%v error:%v", string(out), err)
	} else {
		log.Printf("释放FB权限成功:%v", string(out))
	}
}

func (f *AndroidFramebuffer) Load(img chan image.Image) (err error) {
	f.lock.RLock()
	defer f.lock.RUnlock()
	f.getPermission()
	defer f.releasePermission()
	var fb *os.File
	fb, err = os.Open(f.filePath)
	if err != nil {
		return
	}
	defer fb.Close()
	fmt.Printf("file open %v \n", fb.Fd())

	var varInfo fbVarInfo
	_, _, errno := syscall.RawSyscall(syscall.SYS_IOCTL, fb.Fd(), uintptr(FBIOGET_VSCREENINFO), uintptr(unsafe.Pointer(&varInfo)))
	if errno != 0 {
		fmt.Println(errno)
		err = errors.New("can't ioctl ... ")
	}
	fmt.Printf("width = %d height = %d\n", varInfo.xres, varInfo.yres)
	fmt.Printf("xoffset = %d yoffset = %d\n", varInfo.xoffset, varInfo.yoffset)
	fmt.Printf("depth = %d\n", varInfo.bits_per_pixel/8)

	// func Mmap(fd int, offset int64, length int, prot int, flags int) (data []byte, err error)
	bpp := varInfo.bits_per_pixel / 8
	size := varInfo.xres * varInfo.yres * bpp
	offset := varInfo.xoffset*bpp + varInfo.xres*bpp*varInfo.yoffset
	mapsize := size + offset
	fmt.Printf("mapsize = %d\n", mapsize)
	var data []byte
	data, err = syscall.Mmap(int(fb.Fd()), 0, int(mapsize), syscall.PROT_READ, syscall.MAP_PRIVATE)
	if err != nil {
		return
	}
	defer syscall.Munmap(data)

	// save as png image.
	//
	// func Encode(w io.Writer, m image.Image) error
	// m := image.NewRGBA(image.Rect(int(varInfo.xoffset), int(varInfo.yoffset), int(varInfo.xres), int(varInfo.yres)))
	// m.Pix = data[offset:]
	// f.image = m
	// log.Printf(hex.EncodeToString())
	screen := data[offset:]

	f.Image, err = png.Decode(bytes.NewBuffer(screen))
	if img != nil {
		img <- f.Image
	}
	return
}

func (f *AndroidFramebuffer) SaveImage(filePath string) (err error) {
	var outputFile *os.File
	outputFile, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	defer outputFile.Close()
	err = png.Encode(outputFile, f.Image)
	return
}
