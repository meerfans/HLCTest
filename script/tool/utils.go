/**
	HLC FOUNDATION
	james
 */

package tool

import (
	"encoding/hex"
	"strings"
	"unicode"
	"os/user"
	"os"
	"path/filepath"
	"runtime"
	"log"
	"encoding/csv"
	"io/ioutil"
)

func WriteCsv(filename string,content [][]string){
	f, err := os.Create(filename)//创建文件
	if err != nil {
		panic(err)
	}
	defer f.Close()

	f.WriteString("\xEF\xBB\xBF") // 写入UTF-8 BOM

	w := csv.NewWriter(f)//创建一个新的写入文件流
	data := content
	w.WriteAll(data)//写入数据
	w.Flush()
}

func InArrayString(val string , vals []interface{}) bool {
	for _,v := range vals{
		if v.(string) == val{
			return true
		}
	}
	return false
}

func ReadCsv(filename string,start,length int) (result [][]string) {
	cntb,err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err.Error())
	}
	r2 := csv.NewReader(strings.NewReader(string(cntb)))
	ss,_ := r2.ReadAll()
	//fmt.Println(ss)
	//sz := len(ss)
	for i:=start;i<(start+length);i++{
		result = append(result,ss[i])
	}
	return
}

// Reverse reverses a byte array.
func Reverse(src []byte) []byte {
	dst := make([]byte, len(src))
	for i := len(src); i > 0; i-- {
		dst[len(src)-i] = src[i-1]
	}
	return dst
}


// appDataDir returns an operating system specific directory to be used for
// storing application data for an application.  See AppDataDir for more
// details.  This unexported version takes an operating system argument
// primarily to enable the testing package to properly test the function by
// forcing an operating system that is not the currently one.
func appDataDir(goos, appName string, roaming bool) string {
	if appName == "" || appName == "." {
		return "."
	}

	// The caller really shouldn't prepend the appName with a period, but
	// if they do, handle it gracefully by stripping it.
	appName = strings.TrimPrefix(appName, ".")
	appNameUpper := string(unicode.ToUpper(rune(appName[0]))) + appName[1:]
	appNameLower := string(unicode.ToLower(rune(appName[0]))) + appName[1:]

	// Get the OS specific home directory via the Go standard lib.
	var homeDir string
	usr, err := user.Current()
	if err == nil {
		homeDir = usr.HomeDir
	}

	// Fall back to standard HOME environment variable that works
	// for most POSIX OSes if the directory from the Go standard
	// lib failed.
	if err != nil || homeDir == "" {
		homeDir = os.Getenv("HOME")
	}

	switch goos {
	// Attempt to use the LOCALAPPDATA or APPDATA environment variable on
	// Windows.
	case "windows":
		// Windows XP and before didn't have a LOCALAPPDATA, so fallback
		// to regular APPDATA when LOCALAPPDATA is not set.
		appData := os.Getenv("LOCALAPPDATA")
		if roaming || appData == "" {
			appData = os.Getenv("APPDATA")
		}

		if appData != "" {
			return filepath.Join(appData, appNameUpper)
		}

	case "darwin":
		if homeDir != "" {
			return filepath.Join(homeDir, "Library",
				"Application Support", appNameUpper)
		}

	case "plan9":
		if homeDir != "" {
			return filepath.Join(homeDir, appNameLower)
		}

	default:
		if homeDir != "" {
			return filepath.Join(homeDir, "."+appNameLower)
		}
	}

	// Fall back to the current directory if all else fails.
	return "."
}

// AppDataDir returns an operating system specific directory to be used for
// storing application data for an application.
//
// The appName parameter is the name of the application the data directory is
// being requested for.  This function will prepend a period to the appName for
// POSIX style operating systems since that is standard practice.  An empty
// appName or one with a single dot is treated as requesting the current
// directory so only "." will be returned.  Further, the first character
// of appName will be made lowercase for POSIX style operating systems and
// uppercase for Mac and Windows since that is standard practice.
//
// The roaming parameter only applies to Windows where it specifies the roaming
// application data profile (%APPDATA%) should be used instead of the local one
// (%LOCALAPPDATA%) that is used by default.
//
// Example results:
//  dir := AppDataDir("myapp", false)
//   POSIX (Linux/BSD): ~/.myapp
//   Mac OS: $HOME/Library/Application Support/Myapp
//   Windows: %LOCALAPPDATA%\Myapp
//   Plan 9: $home/myapp
func AppDataDir(appName string, roaming bool) string {
	return appDataDir(runtime.GOOS, appName, roaming)
}

func Target2BlockBits(target string) []byte {
	// 8
	d,_ := hex.DecodeString(target[0:16])
	return Reverse(d)
}

func GetCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))  //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}