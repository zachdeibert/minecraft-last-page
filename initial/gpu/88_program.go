package gpu

//	#cgo pkg-config: OpenCL
//	#cgo CFLAGS: -Wall
//	#include <errno.h>
//	#include <stdlib.h>
//	#include <CL/cl.h>
//
//	static void *create_program(void *_ctx, char *src, int src_len, void *_device_id) {
//		cl_context *ctx = (cl_context *) _ctx;
//		cl_device_id *device_id = (cl_device_id *) _device_id;
//		cl_program *program = (cl_program *) malloc(sizeof(cl_program));
//		if (program == NULL) {
//			errno = ENOBUFS;
//			return NULL;
//		}
//		cl_int ret;
//		*program = clCreateProgramWithSource(*ctx, 1, (const char **) &src, (const size_t *) &src_len, &ret);
//		if (ret != CL_SUCCESS) {
//			free(program);
//			errno = ret;
//			return NULL;
//		}
//		if ((ret = clBuildProgram(*program, 1, device_id, NULL, NULL, NULL)) != CL_SUCCESS) {
//			errno = ret;
//			return program;
//		}
//		errno = 0;
//		return program;
//	}
//
//	static void *create_kernel(void *_program, char *name) {
//		cl_program *program = (cl_program *) _program;
//		cl_kernel *kernel = (cl_kernel *) malloc(sizeof(cl_kernel));
//		if (kernel == NULL) {
//			errno = ENOBUFS;
//			return NULL;
//		}
//		cl_int ret;
//		*kernel = clCreateKernel(*program, name, &ret);
//		if (ret != CL_SUCCESS) {
//			free(kernel);
//			errno = ret;
//			return NULL;
//		}
//		errno = 0;
//		return kernel;
//	}
//
//	static char *compile_log(void *_program, void *_device_id) {
//		cl_program *program = (cl_program *) _program;
//		cl_device_id *device_id = (cl_device_id *) _device_id;
//		size_t len;
//		cl_int ret;
//		if ((ret = clGetProgramBuildInfo(*program, *device_id, CL_PROGRAM_BUILD_LOG, 0, NULL, &len)) != CL_SUCCESS) {
//			errno = ret;
//			return NULL;
//		}
//		if (len == 0) {
//			errno = 0;
//			return NULL;
//		}
//		char *log = (char *) malloc(len);
//		if (log == NULL) {
//			errno = ENOBUFS;
//			return NULL;
//		}
//		if ((ret = clGetProgramBuildInfo(*program, *device_id, CL_PROGRAM_BUILD_LOG, len, log, NULL)) != CL_SUCCESS) {
//			errno = ret;
//			return NULL;
//		}
//		errno = 0;
//		return log;
//	}
//
//	static void set_arg(void *_kernel, int i, void *_mem) {
//		cl_kernel *kernel = (cl_kernel *) _kernel;
//		cl_mem *mem = (cl_mem *) _mem;
//		cl_int ret;
//		if ((ret = clSetKernelArg(*kernel, i, sizeof(cl_mem), mem)) != CL_SUCCESS) {
//			errno = ret;
//			return;
//		}
//		errno = 0;
//	}
//
//	static void cleanup_kernel(void *_kernel) {
//		cl_kernel *kernel = (cl_kernel *) _kernel;
//		clReleaseKernel(*kernel);
//		free(kernel);
//	}
//
//	static void cleanup_program(void *_program) {
//		cl_program *program = (cl_program *) _program;
//		clReleaseProgram(*program);
//		free(program);
//	}
import "C"
import (
	"runtime"
	"unsafe"
)

type kernel struct {
	ptr unsafe.Pointer
}
type program struct {
	ptr unsafe.Pointer
}

const (
	clBuildProgramFailure = -11
)

func createProgram(src string) (*program, string, error) {
	str := C.CString(src)
	defer C.free(unsafe.Pointer(str))
	runtime.GC() // CL_OUT_OF_HOST_MEMORY without this here
	p, e := C.create_program(context, str, C.int(len(src)), deviceID)
	if e != nil && p == nil {
		return nil, "", convertError(e)
	}
	l, ce := C.compile_log(p, deviceID)
	if ce != nil {
		C.cleanup_program(p)
		return nil, "", convertError(ce)
	}
	log := ""
	if l != nil {
		log = C.GoString(l)
		C.free(unsafe.Pointer(l))
	}
	if e != nil {
		C.cleanup_program(p)
		return nil, log, convertError(e)
	}
	return &program{
		ptr: p,
	}, log, nil
}

func createKernel(program program, name string) (kernel, error) {
	str := C.CString(name)
	defer C.free(unsafe.Pointer(str))
	k, e := C.create_kernel(program.ptr, str)
	return kernel{
		ptr: k,
	}, convertError(e)
}

func (k kernel) setArg(i int, buf buffer) error {
	_, e := C.set_arg(k.ptr, C.int(i), buf.ptr)
	return convertError(e)
}

func (p program) Close() {
	C.cleanup_program(p.ptr)
}

func (k kernel) Close() {
	C.cleanup_kernel(k.ptr)
}
