package dllUtil

import (
	"golang.org/x/text/encoding/simplifiedchinese"
	"log"
	"syscall"
	"unsafe"
)

// Base2Ptr convert basic types to uintptr
func Base2Ptr[T int | byte | float64 | float32](n T) uintptr {
	return uintptr(n)
}

// Ptr2Base convert uintptr to basic types
func Ptr2Base[T int | byte | float64 | float32](n uintptr) T {
	return T(n)
}

// this does not work... a compilation error
//func BasePtr2Ptr[T *int] (n T) uintptr {
//	return uintptr(unsafe.Pointer(T))
//}

// BasePtr2Ptr convert the Pointer of basic types to uintptr
func BasePtr2Ptr[T int | byte | float64 | float32](n *T) uintptr {
	return uintptr(unsafe.Pointer(n))
}

// Bool2Ptr convert bool to uintptr
func Bool2Ptr(b bool) uintptr {
	if b {
		return uintptr(1)
	}
	return uintptr(0)
}

// Float2Ptr if @see Base2Ptr does not work, try this
func Float2Ptr[T float32 | float64](f T) uintptr {
	return uintptr(*(*uint32)(unsafe.Pointer(&f)))
}

// Ptr2Float if @see Ptr2Base does not work, try this
func Ptr2Float[T float64 | float32](u uintptr) T {
	p := uint32(u)
	return *(*T)(unsafe.Pointer(&p))
}

// Str2Ptr convert string to uintptr
func Str2Ptr(s string) (uintptr, error) {
	if b, err := syscall.BytePtrFromString(s); err != nil {
		return 0, err
	} else {
		return uintptr(unsafe.Pointer(b)), nil
	}
}

// Utf82Gbk for chinese, convert encoding first @see Str2Ptr
func Utf82Gbk(s string) string {
	decoder := simplifiedchinese.GBK.NewEncoder()
	res, err := decoder.String(s)
	if err != nil {
		log.Println(err.Error())
		return ""
	}
	return res
}

// Ptr2Str convert uintptr to string
func Ptr2Str(u uintptr) string {
	p := (*byte)(unsafe.Pointer(u)) // this first byte of the returned string
	data := make([]byte, 0)         // this slice of []byte store the result
	for *p != 0 {                   // the string of C end with '\0'
		data = append(data, *p)        // append the current byte to the result
		u += unsafe.Sizeof(byte(0))    // move the pointer(u) to the next byte of the returned string
		p = (*byte)(unsafe.Pointer(u)) // get the value(p) of the pointer(u), now p is point the next byte
	}
	return string(data)
}

// EqualDllString this func compare the go string and the result returned by Dll
// note that dllStr is end with '\0', we need find '\0' to get the real dllStr
func EqualDllString(dllStr, str string) bool {
	lenS := len([]byte(str))
	dllByte := []byte(dllStr)
	if len(dllByte) < lenS {
		return false
	}
	realDllByte := dllByte[:lenS]
	return string(realDllByte) == str
}

/*
CallBack
consider now we have a method in DLL that takes one arg is a callback function
. typedef int (*Callback, int)(int);
. int Callback(int x);

we convert Func1 first
*/

// Callback the arg and return of Callback needs uintptr type
func Callback(s uintptr) uintptr {
	log.Printf(Ptr2Str(s)) // your operation
	return Base2Ptr(0)
}

// UseFunc1 we reference Func1 like this
func UseFunc1() {
	dll := syscall.MustLoadDLL("test.dll")
	needCallback := dll.MustFindProc("needCallback")
	c := syscall.NewCallback(Callback)
	r, _, err := needCallback.Call(c, Base2Ptr(0))
	log.Println(r, err)
}
